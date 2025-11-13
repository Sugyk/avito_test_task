package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/Sugyk/avito_test_task/internal/api"
	"github.com/Sugyk/avito_test_task/internal/api/handlers"
	"github.com/Sugyk/avito_test_task/internal/repository"
	"github.com/Sugyk/avito_test_task/internal/service"
	"github.com/Sugyk/avito_test_task/pkg/database"
)

type Application struct {
	db      *sql.DB
	logger  *slog.Logger
	repo    *repository.Repository
	service *service.Service
	router  *api.Router

	wg      sync.WaitGroup
	errChan chan error

	listening_port string
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("init config: %w", err)
	}
	if err := a.initLogger(); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	if err := a.initDatabase(); err != nil {
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

func (a *Application) initDatabase() error {
	db, err := database.NewDbConnection()
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	a.db = db

	a.logger.Info("database connection established")
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
