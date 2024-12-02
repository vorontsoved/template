package main

import (
	"log/slog"
	"os"
	"os/signal"
	"simple/cmd"
	"simple/internal/pkg/logging"
	"sync"
	"syscall"
)

func main() {
	baseLogger := logging.NewBaseLogger(slog.Level(-1))
	logger := logging.NewSlogWrapper(baseLogger)

	app := cmd.New(logger)

	var closeOnce sync.Once
	closeApp := func() {
		closeOnce.Do(func() {
			logger.InfoWithoutContext("App: Closing...")
			if err := app.Close(); err != nil {
				logger.ErrorWithoutContext("App: Failed to close properly", err)
			} else {
				logger.InfoWithoutContext("App: Closed successfully")
			}
		})
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	runErrChan := make(chan error, 1)

	go func() {
		runErrChan <- app.Run()
	}()

	logger.InfoWithoutContext("App: Starting...")

	select {
	case sig := <-shutdownChan:
		logger.InfoWithoutContext("App: Received shutdown signal", slog.String("signal", sig.String()))
	case err := <-runErrChan:
		if err != nil {
			logger.ErrorWithoutContext("App: Error occurred during execution", err)
		}
	}

	closeApp()

	logger.InfoWithoutContext("App: Stopped successfully")
}
