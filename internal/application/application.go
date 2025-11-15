package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Sugyk/avito_test_task/internal/api"
	"github.com/Sugyk/avito_test_task/internal/api/handlers"
	"github.com/Sugyk/avito_test_task/internal/repository"
	"github.com/Sugyk/avito_test_task/internal/service"
	"github.com/Sugyk/avito_test_task/pkg/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

type Application struct {
	db      *sqlx.DB
	logger  *slog.Logger
	repo    *repository.Repository
	service *service.Service
	router  *api.Router

	wg      sync.WaitGroup
	errChan chan error

	listening_port string
}

func NewApplication() *Application {
	return &Application{
		errChan: make(chan error),
	}
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("init config: %w", err)
	}
	if err := a.initLogger(); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	if err := a.initDatabase(ctx); err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	if err := a.initRepository(); err != nil {
		return fmt.Errorf("init repository: %w", err)
	}
	if err := a.initService(); err != nil {
		return fmt.Errorf("init service: %w", err)
	}
	if err := a.initRouter(); err != nil {
		return fmt.Errorf("init router: %w", err)
	}

	if err := a.migrate(); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	a.startHTTPServer()

	a.logger.Info("application started successfully")
	return nil

}

func (a *Application) initConfig() error {
	a.listening_port = os.Getenv("LISTEN_PORT")
	if a.listening_port == "" {
		return fmt.Errorf("parsing LISTEN_PORT error: listen port is missing")
	}
	return nil
}

func (a *Application) initLogger() error {
	a.logger = slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)
	return nil
}

func (a *Application) initDatabase(ctx context.Context) error {
	db, err := database.NewDbConnection(ctx)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	a.db = db

	a.logger.Info("database connection established")
	return nil
}

func (a *Application) migrate() error {
	if err := database.RunMigrations(a.db); err != nil {
		if err == migrate.ErrNoChange {
			a.logger.Info("no changes to migrate")
			return nil
		}
		return fmt.Errorf("database migration: %w", err)
	}
	a.logger.Info("migrated successfully")
	return nil
}

func (a *Application) initRepository() error {
	a.repo = repository.NewRepository(
		a.db,
		a.logger,
	)
	return nil
}

func (a *Application) initService() error {
	a.service = service.NewService(
		a.repo,
		a.logger,
	)
	return nil
}

func (a *Application) initRouter() error {
	handler := handlers.NewHandler(
		a.service,
		a.logger,
	)

	a.router = api.NewRouter(
		a.listening_port,
		handler,
	)
	return nil
}

func (a *Application) startHTTPServer() {
	a.wg.Add(1)

	go func() {
		defer a.wg.Done()

		a.logger.Info("Starting HTTP server on " + a.listening_port)
		if err := a.router.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()
}

func (a *Application) Wait(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received, starting graceful shutdown...")
	case err := <-a.errChan:
		a.logger.Error("error received, initiating shutdown", "error", err)
		return err
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := a.router.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", "error", err)
	}

	if err := a.db.Close(); err != nil {
		a.logger.Error("database closed with error", "error", err)
	} else {
		a.logger.Info("database connections closed")
	}

	a.wg.Wait()

	a.logger.Info("graceful shutdown completed")

	return nil
}
