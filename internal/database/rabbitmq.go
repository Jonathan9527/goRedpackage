package database

import (
	"encoding/json"
	"fmt"

	"learnGO/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn *amqp091.Connection
}

func NewRabbitMQ(cfg config.RabbitMQConfig) (*RabbitMQPublisher, error) {
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

	return &RabbitMQPublisher{conn: conn}, nil
}

func (p *RabbitMQPublisher) Close() error {
	if p == nil || p.conn == nil {
		return nil
	}
	return p.conn.Close()
}

func (p *RabbitMQPublisher) PublishJSON(queue string, payload any) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	if _, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return channel.Publish(
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *RabbitMQPublisher) Consume(queue string) (<-chan amqp091.Delivery, error) {
	channel, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}

	if _, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		channel.Close()
		return nil, err
	}
	channel.Qos(2, 0, false)
	deliveries, err := channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		return nil, err
	}

	return deliveries, nil
}
