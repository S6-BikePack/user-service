package interfaces

import (
	"context"
	"user-service/internal/core/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Get(ctx context.Context, id string) (domain.User, error)
	Save(ctx context.Context, user domain.User) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
}
