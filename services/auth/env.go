package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				return
			}
		}
	}
	slog.Warn("Aviso: nenhum arquivo .env encontrado")
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("environment variable %s is required", "key", key)
		os.Exit(1)
	}
	return value
}
