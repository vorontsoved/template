package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"simple/internal/pkg/logging"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type Config struct {
	Port     string
	Host     string
	User     string
	Password string
	DB       string
	LogLevel int
}

type Postgres struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

type CustomLogger struct {
	logger logging.Logger
}

func (cl *CustomLogger) Log(_ context.Context, level tracelog.LogLevel, msg string, _ map[string]interface{}) {
	var file string

	var line int

	var ok bool

	for i := 2; i < 10; i++ {
		_, file, line, ok = runtime.Caller(i)
		if ok && !(strings.Contains(file, "pgx") || strings.Contains(file, "tracelog")) {
			break
		}
	}

	if !strings.Contains(msg, "Query") {
		return
	}

	file = filepath.Base(file)
	logEntryStr := fmt.Sprintf("Query\nFile: %s:%d\n", file, line)

	switch level {
	case tracelog.LogLevelTrace:
		cl.logger.DebugWithoutContext(logEntryStr)
	case tracelog.LogLevelDebug:
		cl.logger.DebugWithoutContext(logEntryStr)
	case tracelog.LogLevelInfo:
		cl.logger.InfoWithoutContext(logEntryStr)
	case tracelog.LogLevelWarn:
		cl.logger.WarnWithoutContext(logEntryStr)
	case tracelog.LogLevelError:
		cl.logger.ErrorWithoutContext(logEntryStr, nil)
	case tracelog.LogLevelNone:
	}
}

func InitDatabase(cfg Config, logger logging.Logger) (Postgres, error) {
	connConfig := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	config, err := pgxpool.ParseConfig(connConfig)
	if err != nil {
		logger.ErrorWithoutContext("Error parsing connection config: ", err)

		return Postgres{}, err
	}

	customLogger := &CustomLogger{logger: logger}
	tracer := &tracelog.TraceLog{
		Logger:   customLogger,
		LogLevel: tracelog.LogLevelDebug,
	}
	config.ConnConfig.Tracer = tracer

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.ErrorWithoutContext("Error connecting to the database: ", err)

		return Postgres{}, err
	}

	logger.InfoWithoutContext("Successfully connected to the database")

	return Postgres{db: conn, logger: logger}, nil
}

func (p Postgres) Close() error {
	if p.db != nil {
		p.db.Close()
	}

	return nil
}
