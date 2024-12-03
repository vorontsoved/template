package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PgHost     string `env:"PG_HOST" env-required:"true" env-description:"PostgreSQL host"`
	PgPort     string `env:"PG_PORT" env-required:"true" env-description:"PostgreSQL port"`
	PgDB       string `env:"PG_DB" env-required:"true" env-description:"PostgreSQL database name"`
	PgUser     string `env:"PG_USER" env-required:"true" env-description:"PostgreSQL user"`
	PgPassword string `env:"PG_PASSWORD" env-required:"true" env-description:"PostgreSQL password"`
	HTTPPort   string `env:"HTTP_PORT" env-default:"4000" env-description:"HTTP server port"`
	Migration  bool   `env:"MIGRATION" env-default:"false" env-description:"Run migrations at startup"`
	LogLVL     int    `env:"LOG_LVL" env-default:"-1" env-description:"Log level"`
}

func Parse() (Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
