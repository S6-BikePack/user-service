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

	srv.messagePublisher.CreateUser(user)
	return user, nil
}

func (srv *service) UpdateUserDetails(id string, name, lastName, email string) (domain.User, error) {
	user, err := srv.Get(id)

	if err != nil {
		return domain.User{}, errors.New("could not find user with id")
	}

	if name != "" {
		user.Name = name
	}

	if lastName != "" {
		user.LastName = lastName
	}

	if email != "" {
		user.Email = email
	}

	user, err = srv.userRepository.Update(user)

	if err != nil {
		return domain.User{}, errors.New("saving new user failed")
	}

	srv.messagePublisher.UpdateUserDetails(user)
	return user, nil
}
