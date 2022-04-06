package dto

import "user-service/internal/core/domain"

type BodyCreateUser struct {
	ID       string
	Name     string
	LastName string
	Email    string
}

type ResponseCreateUser domain.User

func BuildResponseCreateUser(model domain.User) ResponseCreateUser {
	return ResponseCreateUser(model)
}
