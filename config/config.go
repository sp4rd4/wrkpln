package config

import "time"

type Config struct {
	Port            int           `env:"HTTP_PORT" envDefault:"8080"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"INFO"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	MaxHeaderBytes  int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"` //1MB
	DBPath          string        `env:"SQLITE_DB" envDefault:"work_planning.db"`
	DBSchemaPath    string        `env:"SQLITE_SCHEMA" envDefault:"db/schema.sql"`
}
