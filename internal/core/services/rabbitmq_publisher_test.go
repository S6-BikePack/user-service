package services

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/sdk/trace"
	"testing"
	"user-service/config"
	"user-service/internal/core/domain"
	"user-service/internal/core/interfaces"
	"user-service/pkg/rabbitmq"
)

type RabbitMQPublisherTestSuite struct {
	suite.Suite
	TestRabbitMQ  *rabbitmq.RabbitMQ
	TestPublisher interfaces.MessageBusPublisher
	Cfg           *config.Config
	TestData      struct {
		User domain.User
	}
}

func (suite *RabbitMQPublisherTestSuite) SetupSuite() {
	cfgPath := "../../../test/user.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(errors.WithStack(err))
	}

	tracer := trace.NewTracerProvider()

	rmqPublisher := NewRabbitMQPublisher(rmqServer, tracer, cfg)

	suite.Cfg = cfg
	suite.TestRabbitMQ = rmqServer
	suite.TestPublisher = rmqPublisher
	suite.TestData = struct {
		User domain.User
	}{
		User: domain.User{
			ID:       "test-id",
			Name:     "test-name",
			LastName: "test-lastname",
		},
	}
}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_CreateUser() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"user.create",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.CreateUser(context.Background(), suite.TestData.User)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("user.create", msg.RoutingKey)

		var user domain.User

		err = json.Unmarshal(msg.Body, &user)
		suite.NoError(err)

		suite.Equal(suite.TestData.User, user)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}
}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_UpdateUserDetails() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"user.update",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.UpdateUserDetails(context.Background(), suite.TestData.User)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("user.update", msg.RoutingKey)

		var user domain.User

		err = json.Unmarshal(msg.Body, &user)
		suite.NoError(err)

		suite.Equal(suite.TestData.User, user)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}
}

func TestIntegration_RabbitMQPublisherTestSuite(t *testing.T) {
	testSuite := new(RabbitMQPublisherTestSuite)
	suite.Run(t, testSuite)
}
