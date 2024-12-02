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
			logger.InfoWithOutContext("App: Closing...")
			if err := app.Close(); err != nil {
				logger.ErrorWithOutContext("App: Failed to close properly", err)
			} else {
				logger.InfoWithOutContext("App: Closed successfully")
			}
		})
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	runErrChan := make(chan error, 1)

	go func() {
		runErrChan <- app.Run()
	}()

	logger.InfoWithOutContext("App: Starting...")

	select {
	case sig := <-shutdownChan:
		logger.InfoWithOutContext("App: Received shutdown signal", slog.String("signal", sig.String()))
	case err := <-runErrChan:
		if err != nil {
			logger.ErrorWithOutContext("App: Error occurred during execution", err)
		}
	}

	closeApp()

	logger.InfoWithOutContext("App: Stopped successfully")
}
