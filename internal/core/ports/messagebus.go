package ports

import (
	"user-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateUser(user domain.User) error
	UpdateUserDetails(user domain.User) error
}
