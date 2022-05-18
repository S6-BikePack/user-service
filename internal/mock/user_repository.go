package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"user-service/internal/core/domain"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *UserRepository) Get(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) Save(ctx context.Context, user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}
