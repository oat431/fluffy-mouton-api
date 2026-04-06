package config

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

func init() {
	log.Info("loading environment configuration...")
}

func LoadEnvConfig() {
	failedToLoadEnv := godotenv.Load(".env")
	if failedToLoadEnv != nil {
		log.Fatal("Failed to load environment variables from .env.development file")
	}
}
