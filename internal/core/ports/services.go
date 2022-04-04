package ports

import "user-service/internal/core/domain"

type UserService interface {
	GetAll() ([]domain.User, error)
	Get(id string) (domain.User, error)
	Create(id, name, lastName, email string) (domain.User, error)
	UpdateUserDetails(id, name, lastName, email string) (domain.User, error)
}
