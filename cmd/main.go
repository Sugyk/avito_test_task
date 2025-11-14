package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/Sugyk/avito_test_task/internal/application"
)

func main() {
	if err := run(); err != nil {
		slog.Info("error while starting server", "error", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	app := application.NewApplication()
	if err := app.Start(ctx); err != nil {
		log.Fatalln("can not start application:", err)
	}

	if err := app.Wait(ctx, cancel); err != nil {
		log.Fatalln("All systems closed with errors. LastError:", err)
	}
	return nil
}
