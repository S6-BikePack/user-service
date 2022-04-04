package ports

import (
	"user-service/internal/core/domain"
)

type UserRepository interface {
	GetAll() ([]domain.User, error)
	Get(id string) (domain.User, error)
	Save(customer domain.User) (domain.User, error)
	Update(customer domain.User) (domain.User, error)
}
