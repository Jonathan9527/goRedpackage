package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort  string
	DB       DatabaseConfig
	RabbitMQ RabbitMQConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string
}
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load() Config {
	return Config{
		AppPort: getEnv("APP_PORT", "8080"),
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "learngo"),
			Password: getEnv("DB_PASSWORD", "learngo_password"),
			Name:     getEnv("DB_NAME", "learngo"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "learngo"),
			Password: getEnv("RABBITMQ_PASSWORD", "learngo_password"),
			VHost:    getEnv("RABBITMQ_VHOST", "/"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	intVal := 0
	_, err := fmt.Sscanf(value, "%d", &intVal)
	if err != nil {
		return fallback
	}
	return intVal
}
