package config

import (
	"fmt"
	"os"
)

type Config struct {
	RedisAddr  string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Load() *Config {
	cfg := &Config{
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "nabeel"),
		DBPassword: getEnv("DB_PASSWORD", "nabeel"),
		DBName:     getEnv("DB_NAME", "microservice_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}

	fmt.Printf("Config loaded: RedisAddr=%s DBHost=%s\n", cfg.RedisAddr, cfg.DBHost)

	return cfg
}
