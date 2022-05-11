package dto

import "user-service/internal/core/domain"

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func CreateUserResponse(customer domain.User) UserResponse {
	return UserResponse{
		ID:       customer.ID,
		Name:     customer.Name,
		LastName: customer.LastName,
		Email:    customer.Email,
	}
}
