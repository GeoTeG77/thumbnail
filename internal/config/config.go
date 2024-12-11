package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	slog.Info(".env file load successfully!")
	return nil

}
