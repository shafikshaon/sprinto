package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() Config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default for local development
		dsn = "host=localhost port=5432 dbname=sprinto user=postgres password=postgres sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{DatabaseURL: dsn, Port: port}
}
