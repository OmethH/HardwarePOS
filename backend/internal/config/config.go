package config

import "os"

type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	Migrations string
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "data/hardwarepos.db"),
		JWTSecret:  getEnv("JWT_SECRET", "hardwarepos-dev-secret-change-in-production"),
		Migrations: getEnv("MIGRATIONS_PATH", "migrations"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
