package dto

import "user-service/internal/core/domain"

type BodyCreateUser struct {
	Name     string
	LastName string
}

type ResponseCreateUser domain.User

func BuildResponseCreateUser(model domain.User) ResponseCreateUser {
	return ResponseCreateUser(model)
}
