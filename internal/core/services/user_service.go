package services

import (
	"context"
	"errors"
	"user-service/internal/core/domain"
	"user-service/internal/core/interfaces"
)

type userService struct {
	userRepository   interfaces.UserRepository
	messagePublisher interfaces.MessageBusPublisher
}

func NewUserService(riderRepository interfaces.UserRepository, messagePublisher interfaces.MessageBusPublisher) *userService {
	return &userService{
		userRepository:   riderRepository,
		messagePublisher: messagePublisher,
	}
}

func (srv *userService) GetAll(ctx context.Context) ([]domain.User, error) {
	return srv.userRepository.GetAll(ctx)
}

func (srv *userService) Get(ctx context.Context, id string) (domain.User, error) {
	return srv.userRepository.Get(ctx, id)
}

func (srv *userService) Create(ctx context.Context, id, name, lastName, email string) (domain.User, error) {
	user, err := domain.NewUser(id, name, lastName, email)

	if err != nil || user.ID == "" {
		return domain.User{}, err
	}

	user, err = srv.userRepository.Save(ctx, user)

	if err != nil {
		return domain.User{}, errors.New("saving new user failed")
	}

	err = srv.messagePublisher.CreateUser(ctx, user)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (srv *userService) UpdateUserDetails(ctx context.Context, id string, name, lastName, email string) (domain.User, error) {
	existing, err := srv.Get(ctx, id)
	updated := existing

	if err != nil {
		return domain.User{}, errors.New("could not find user with id")
	}

	if name != "" {
		updated.Name = name
	}

	if lastName != "" {
		updated.LastName = lastName
	}

	if email != "" {
		updated.Email = email
	}

	updated, err = srv.userRepository.Update(ctx, updated)

	if err != nil {
		return existing, errors.New("saving new user failed")
	}

	err = srv.messagePublisher.UpdateUserDetails(ctx, updated)

	if err != nil {
		return updated, err
	}

	return updated, nil
}
