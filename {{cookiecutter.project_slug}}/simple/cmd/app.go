package cmd

import (
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

func (a App) Run() error {
	cfg, err := config.Parse()
	if err != nil {
		a.logger.ErrorWithoutContext("failed parse config", err)
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
			a.logger.ErrorWithoutContext("failed migrate", err)
			return err
		}

	}

	// Package service (or repository) represents the infrastructure layer responsible
	// for data persistence, external system interactions, and message publishing.
	// This layer implements the interfaces defined in the business logic layer (usecase/domain)
	// and contains the actual data access logic, such as database queries, API calls,
	// and producers for publishing messages to Kafka or other message brokers.
	// It acts as a bridge between the core business logic and external dependencies.
	a.pg, err = postgres.InitDatabase(postgres.Config{
		Port:     cfg.PgPort,
		Host:     cfg.PgHost,
		User:     cfg.PgUser,
		Password: cfg.PgPassword,
		DB:       cfg.PgDB,
		LogLevel: cfg.LogLVL,
	}, a.logger)
	if err != nil {
		a.logger.ErrorWithoutContext("failed", err)
		return err
	}

	// Package usecase (or domain) represents the core business logic layer.
	// This layer contains the application's primary business rules and use cases,
	// decoupled from infrastructure and external frameworks. It interacts with
	// repositories and other external systems through defined interfaces,
	// ensuring the business logic remains independent and testable.
	uc := usecase.New(a.pg, a.logger)

	// Package transport represents the layer responsible for communication between
	// the application and external clients. This layer handles incoming requests
	// (e.g., HTTP, gRPC) and outgoing responses, acting as the entry point for the system.
	// It validates and parses input, invokes the appropriate business logic in the
	// usecase/domain layer, and formats the output to be returned to the client.
	// Transport may also handle protocols like REST, gRPC, WebSocket, or other
	// communication mechanisms, ensuring a clear boundary between external clients
	// and the internal application logic.
	tr := transport.NewHandler(uc, a.logger)

	go func() {
		tr.Run()
	}()

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
