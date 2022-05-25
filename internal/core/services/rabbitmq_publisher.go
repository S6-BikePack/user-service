package services

import (
	"context"
	"encoding/json"
	"fmt"
	"user-service/config"
	"user-service/internal/core/domain"
	"user-service/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type rabbitmqPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
	tracer   trace.Tracer
	config   *config.Config
}

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ, tracerProvider trace.TracerProvider, cfg *config.Config) *rabbitmqPublisher {
	return &rabbitmqPublisher{rabbitmq: rabbitmq, tracer: tracerProvider.Tracer("RabbitMQ.Publisher"), config: cfg}
}

func (rmq *rabbitmqPublisher) CreateUser(ctx context.Context, user domain.User) error {
	return rmq.publishJson(ctx, "create", user)
}

func (rmq *rabbitmqPublisher) UpdateUserDetails(ctx context.Context, user domain.User) error {
	return rmq.publishJson(ctx, "update", user)
}

func (rmq *rabbitmqPublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	if rmq.tracer != nil {
		_, span := rmq.tracer.Start(ctx, "publish")

		span.AddEvent(
			"Published message to rabbitmq",
			trace.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("body", string(js))))
		span.End()
	}

	err = rmq.rabbitmq.Channel.Publish(
		rmq.config.RabbitMQ.Exchange,
		fmt.Sprintf("user.%s", topic),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         js,
		},
	)

	return err
}
