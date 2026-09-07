package config

import "os"

type Config struct {
	Database_URL string
	JWT_SECRET   string
}

func EnvConfig() Config {
	return Config{
		Database_URL: os.Getenv("Database_URL"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
	}
}
