package migration

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"{{cookiecutter.project_slug}}/internal/pkg/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

const versionTable = "db_version"

type DatabaseMigrator struct {
	migrator *migrate.Migrator
}

//go:embed data/*.sql
var migrationFiles embed.FS

func NewDatabaseMigrator(dbDNS string) (DatabaseMigrator, error) {
	conn, err := pgx.Connect(context.Background(), dbDNS)
	if err != nil {
		return DatabaseMigrator{}, err
	}

	migrator, err := migrate.NewMigratorEx(
		context.Background(), conn, versionTable,
		&migrate.MigratorOptions{
			DisableTx: false,
		})
	if err != nil {
		return DatabaseMigrator{}, err
	}

	migrationRoot, _ := fs.Sub(migrationFiles, "data")

	err = migrator.LoadMigrations(migrationRoot)
	if err != nil {
		return DatabaseMigrator{}, err
	}

	return DatabaseMigrator{
		migrator: migrator,
	}, nil
}

// GetMigrationStatus returns the current migration version and the latest available migration,
// along with a textual representation of the migration state.
func (dm DatabaseMigrator) GetMigrationStatus() (int32, int32, string, error) {
	version, err := dm.migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return 0, 0, "", err
	}

	var info string
	var last int32

	for _, migration := range dm.migrator.Migrations {
		last = migration.Sequence

		isCurrent := version == migration.Sequence
		indicator := "  "

		if isCurrent {
			indicator = "->"
		}

		info += fmt.Sprintf(
			"%2s %3d %s\n",
			indicator,
			migration.Sequence, migration.Name)
	}

	return version, last, info, nil
}

// ApplyMigrations migrates the database to the latest schema version.
func (dm DatabaseMigrator) ApplyMigrations() error {
	return dm.migrator.Migrate(context.Background())
}

// ApplyMigrationsToVersion migrates the database to a specific schema version.
// Use version '0' to roll back all migrations.
func (dm DatabaseMigrator) ApplyMigrationsToVersion(version int32) error {
	return dm.migrator.MigrateTo(context.Background(), version)
}

type DatabaseConfig struct {
	Port     string
	Host     string
	User     string
	Password string
	DB       string
}

// RunDatabaseMigrations initializes and applies database migrations.
func RunDatabaseMigrations(cfg DatabaseConfig, logger logging.Logger) error {
	connConfig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	logger.InfoWithoutContext("Starting database migration logic")

	migrator, err := NewDatabaseMigrator(connConfig)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Retrieve migration status
	currentVersion, latestVersion, migrationInfo, err := migrator.GetMigrationStatus()
	if err != nil {
		return fmt.Errorf("failed to get migration info: %w", err)
	}

	if currentVersion < latestVersion {
		logger.InfoWithoutContext("Migrations required. Current state: ", "info", migrationInfo)

		err = migrator.ApplyMigrations()
		if err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}

		logger.InfoWithoutContext("Migrations applied successfully!")
	} else {
		logger.InfoWithoutContext("no database migration needed")
	}

	return nil
}
