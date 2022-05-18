package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"user-service/internal/core/domain"
)

type MessageBusPublisher struct {
	mock.Mock
}

func (m *MessageBusPublisher) CreateUser(ctx context.Context, user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MessageBusPublisher) UpdateUserDetails(ctx context.Context, user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
