package domain

type User struct {
	ID string

	Name     string
	LastName string
	Email    string
}

func NewUser(id, name, lastName, email string) User {
	return User{
		ID:       id,
		Name:     name,
		LastName: lastName,
		Email:    email,
	}
}
