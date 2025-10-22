package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int
	DatabasePath string
}

func Load() *Config {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	databasePath := "tasks.db"
	if path := os.Getenv("DATABASE_PATH"); path != "" {
		databasePath = path
	}

	return &Config{
		Port:         port,
		DatabasePath: databasePath,
	}
}
