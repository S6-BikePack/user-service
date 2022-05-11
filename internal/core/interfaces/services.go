package interfaces

import (
	"context"
	"user-service/internal/core/domain"
)

type UserService interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Get(ctx context.Context, id string) (domain.User, error)
	Create(ctx context.Context, id, name, lastName, email string) (domain.User, error)
	UpdateUserDetails(ctx context.Context, id, name, lastName, email string) (domain.User, error)
}
