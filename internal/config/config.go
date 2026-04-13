package config

import "os"

type Config struct {
	AppPort  string
	DB       DatabaseConfig
	RabbitMQ RabbitMQConfig
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
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
