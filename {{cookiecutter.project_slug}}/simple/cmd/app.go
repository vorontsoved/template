package cmd

import (
	"context"

	"{{cookiecutter.project_slug}}/cmd/config"
	"{{cookiecutter.project_slug}}/internal/pkg/logging"
	"{{cookiecutter.project_slug}}/internal/service/postgres"
	"{{cookiecutter.project_slug}}/internal/transport"
	"{{cookiecutter.project_slug}}/internal/usecase"
	"{{cookiecutter.project_slug}}/internal/utils/migration"
)

// Two naming approaches:
// transport - usecase - service
// transport - service - repository
// Difference? The first approach focuses on building applications with clean architecture in mind.
// All calls to external systems (databases, message brokers, API calls) are treated as services
// and belong to the service layer.
// In the second approach, calls to brokers or APIs can reside in the business logic layer
// and are considered part of the service layer. The repository layer, in this case, is limited
// to handling only data-related operations (e.g., databases, Elasticsearch).

type App struct {
	logger logging.Logger
	pg     postgres.Postgres
}

func New(logger logging.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	cfg, err := config.Parse()
	if err != nil {
		a.logger.ErrorWithoutContext("failed to parse config", err)
		return err
	}

	if cfg.Migration {
		if err = migration.RunDatabaseMigrations(migration.DatabaseConfig{
			Port:     cfg.PgPort,
			Host:     cfg.PgHost,
			User:     cfg.PgUser,
			Password: cfg.PgPassword,
			DB:       cfg.PgDB,
		}, a.logger); err != nil {
			a.logger.ErrorWithoutContext("failed to migrate database", err)
			return err
		}
	}

	a.pg, err = postgres.InitDatabase(postgres.Config{
		Port:     cfg.PgPort,
		Host:     cfg.PgHost,
		User:     cfg.PgUser,
		Password: cfg.PgPassword,
		DB:       cfg.PgDB,
		LogLevel: cfg.LogLVL,
	}, a.logger)
	if err != nil {
		a.logger.ErrorWithoutContext("failed to initialize postgres", err)
		return err
	}

	uc := usecase.New(a.pg, a.logger)
	tr := transport.NewHandler(uc, a.logger)

	// Горутин для запуска транспорта
	go func() {
		if err := tr.Run(ctx); err != nil {
			a.logger.ErrorWithoutContext("transport stopped with error", err)
		}
	}()

	// Блокируем выполнение до отмены контекста
	<-ctx.Done()
	a.logger.InfoWithoutContext("App: Context cancelled, shutting down...")
	return nil
}

func (a App) Close() error {
	defer func() {
		a.logger.InfoWithoutContext("close app")
	}()

	err := a.pg.Close()
	if err != nil {
		a.logger.ErrorWithoutContext("failed close postgres", err)
		return err
	}

	return nil
}
