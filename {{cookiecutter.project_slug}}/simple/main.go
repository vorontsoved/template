package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"{{cookiecutter.project_slug}}/cmd"
	"{{cookiecutter.project_slug}}/internal/pkg/logging"
)

func main() {
	baseLogger := logging.NewBaseLogger(slog.Level(-1))
	logger := logging.NewSlogWrapper(baseLogger)

	app := cmd.New(logger)

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Обрабатываем сигналы завершения
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Канал для ошибок выполнения приложения
	runErrChan := make(chan error, 1)

	// Горутин для выполнения приложения
	go func() {
		runErrChan <- app.Run(ctx)
	}()

	logger.InfoWithoutContext("App: Starting...")

	select {
	case sig := <-shutdownChan:
		logger.InfoWithoutContext("App: Received shutdown signal", slog.String("signal", sig.String()))
		cancel() // Останавливаем приложение через контекст
	case err := <-runErrChan:
		if err != nil {
			logger.ErrorWithoutContext("App: Error occurred during execution", err)
			cancel() // Завершаем приложение в случае ошибки
		}
	}

	// Закрываем приложение
	if err := app.Close(); err != nil {
		logger.ErrorWithoutContext("App: Failed to close", err)
	} else {
		logger.InfoWithoutContext("App: Stopped successfully")
	}
}
