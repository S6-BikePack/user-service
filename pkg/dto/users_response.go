package dto

import "user-service/internal/core/domain"

type userResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

func createUserResponse(user domain.User) userResponse {
	return userResponse{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
	}
}

type UserListResponse []*userResponse

func CreateUserListResponse(users []domain.User) UserListResponse {
	response := UserListResponse{}
	for _, s := range users {
		user := createUserResponse(s)
		response = append(response, &user)
	}
	return response
}
