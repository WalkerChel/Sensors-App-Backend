package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sensors-app/internal/entities"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func New() (entities.Config, error) {
	if err := godotenv.Load(); err != nil {
		if pathErr, ok := errors.AsType[*os.PathError](err); ok {
			if os.IsNotExist(err) {
				log.Println("Info: No .env file found; relying on system environment variables")
			} else {
				log.Printf("Warning: Could not read .env file (%v); falling back to system environment variables", pathErr)
			}
		} else {
			return entities.Config{}, err
		}
	}

	cnf := entities.Config{}
	if err := cleanenv.ReadEnv(&cnf); err != nil {
		return entities.Config{}, fmt.Errorf("config err: %w", err)
	}
	return cnf, nil
}
