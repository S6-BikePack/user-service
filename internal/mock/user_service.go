package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"user-service/internal/core/domain"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *UserService) Get(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserService) Create(ctx context.Context, id, name, lastName, email string) (domain.User, error) {
	args := m.Called(id, name, lastName, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserService) UpdateUserDetails(ctx context.Context, id, name, lastName, email string) (domain.User, error) {
	args := m.Called(id, name, lastName, email)
	return args.Get(0).(domain.User), args.Error(1)
}
