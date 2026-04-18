package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("Aviso: não foi possível carregar .env (talvez ele não exista?)", "err", err)
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("environment variable %s is required", "key", key)
		os.Exit(1)
	}
	return value
}
