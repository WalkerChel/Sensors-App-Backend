package configs

import (
	"fmt"
	"sensors-app/internal/entities"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func New() (entities.Config, error) {
	if err := godotenv.Load(); err != nil {
		return entities.Config{}, err
	}

	cnf := entities.Config{}
	if err := cleanenv.ReadEnv(&cnf); err != nil {
		return entities.Config{}, fmt.Errorf("config err: %w", err)
	}
	return cnf, nil
}
