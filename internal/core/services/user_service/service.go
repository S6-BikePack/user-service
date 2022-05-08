package user_service

import (
	"errors"
	"user-service/internal/core/domain"
	"user-service/internal/core/ports"
)

type service struct {
	userRepository   ports.UserRepository
	messagePublisher ports.MessageBusPublisher
}

func New(riderRepository ports.UserRepository, messagePublisher ports.MessageBusPublisher) *service {
	return &service{
		userRepository:   riderRepository,
		messagePublisher: messagePublisher,
	}
}

func (srv *service) GetAll() ([]domain.User, error) {
	return srv.userRepository.GetAll()
}

func (srv *service) Get(id string) (domain.User, error) {
	return srv.userRepository.Get(id)
}

func (srv *service) Create(id, name, lastName, email string) (domain.User, error) {
	user, err := domain.NewUser(id, name, lastName, email)

	if err != nil {
		return domain.User{}, err
	}

	user, err = srv.userRepository.Save(user)

	if err != nil {
		return domain.User{}, errors.New("saving new user failed")
	}

	err = srv.messagePublisher.CreateUser(user)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (srv *service) UpdateUserDetails(id string, name, lastName, email string) (domain.User, error) {
	existing, err := srv.Get(id)
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

	updated, err = srv.userRepository.Update(updated)

	if err != nil {
		return existing, errors.New("saving new user failed")
	}

	err = srv.messagePublisher.UpdateUserDetails(updated)

	if err != nil {
		return updated, err
	}

	return updated, nil
}
