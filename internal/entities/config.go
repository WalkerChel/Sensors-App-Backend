package entities

import "time"

type Config struct {
	Postgres Postgres
	Redis    Redis
	JWT      JWT
	Server   Server
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Password string `env:"POSTGRES_PASSWORD"`
	User     string `env:"POSTGRES_USER"`
	Name     string `env:"POSTGRES_DB"`
	SSLMode  string `env:"POSTGRES_SSL_MODE"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type JWT struct {
	TTL          time.Duration `env:"TOKEN_TTL"`
	SignatureKey string        `env:"TOKEN_SIGN_KEY"`
}

type Server struct {
	Host            string        `env:"HTTP_SERVER_HOST"`
	Port            string        `env:"HTTP_SERVER_PORT"`
	ReadTimeout     time.Duration `env:"HTTP_SERVER_READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"HTTP_SERVER_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT"`
}
