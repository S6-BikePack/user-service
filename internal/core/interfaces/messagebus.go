package interfaces

import (
	"context"
	"user-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateUser(ctx context.Context, user domain.User) error
	UpdateUserDetails(ctx context.Context, user domain.User) error
}
