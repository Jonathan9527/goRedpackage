package database

import (
	"context"
	"fmt"
	"time"

	"learnGO/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(cfg config.RabbitMQConfig) (*amqp091.Connection, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- conn.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		conn.Close()
		return nil, fmt.Errorf("ping rabbitmq timeout")
	}

	return conn, nil
}
