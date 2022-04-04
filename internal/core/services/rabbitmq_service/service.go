package rabbitmq_service

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"user-service/internal/core/domain"
	"user-service/pkg/rabbitmq"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel}
}

func (rmq *rabbitmqPublisher) CreateUser(user domain.User) error {
	js, err := json.Marshal(user)

	if err != nil {
		return err
	}

	err = rmq.publishMessage("user.create", js)

	return err
}

func (rmq *rabbitmqPublisher) UpdateUserDetails(user domain.User) error {
	js, err := json.Marshal(user)

	if err != nil {
		return err
	}

	err = rmq.publishMessage("user.update", js)

	return err
}

func (rmq *rabbitmqPublisher) publishMessage(key string, body []byte) error {
	err := rmq.Channel.Publish(
		"topics",
		key,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)

	return err
}
