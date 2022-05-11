package dto

import "user-service/internal/core/domain"

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func CreateUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
	}
}
