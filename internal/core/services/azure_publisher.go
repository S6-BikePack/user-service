package services

import (
	"context"
	"encoding/json"
	"fmt"
	"user-service/config"
	"user-service/internal/core/domain"
	"user-service/pkg/azure"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type azurePublisher struct {
	serviceBus *azure.ServiceBus
	sender     *azservicebus.Sender
	config     *config.Config
}

func NewAzurePublisher(serviceBus *azure.ServiceBus, cfg *config.Config) *azurePublisher {
	return &azurePublisher{serviceBus: serviceBus, config: cfg}
}

func (rmq *azurePublisher) CreateUser(ctx context.Context, user domain.User) error {
	return rmq.publishJson(ctx, "create", user)
}

func (rmq *azurePublisher) UpdateUserDetails(ctx context.Context, user domain.User) error {
	return rmq.publishJson(ctx, "update", user)
}

func (az *azurePublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	sender, err := az.serviceBus.Client.NewSender(fmt.Sprintf("user.%s", topic), nil)

	defer func(sender *azservicebus.Sender, ctx context.Context) {
		_ = sender.Close(ctx)
	}(sender, ctx)

	if err != nil {
		return err
	}

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body: js,
	}, nil)

	if err != nil {
		return err
	}

	return err
}
