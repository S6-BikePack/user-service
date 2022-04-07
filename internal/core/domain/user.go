package domain

import (
	"errors"
	"regexp"
	"strings"
)

type User struct {
	ID string

	Name     string
	LastName string
	Email    string
}

func NewUser(id, name, lastName, email string) (User, error) {
	if id == "" || name == "" || lastName == "" || email == "" {
		return User{}, errors.New("missing data")
	}

	emailReg := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !emailReg.MatchString(email) {
		return User{}, errors.New("email is not valid")
	}

	nameReg := regexp.MustCompile(`^[A-Za-z]+(([,.] |[ '-])[A-Za-z]+)*([.,'-]?)$`)

	if !nameReg.MatchString(name) {
		return User{}, errors.New("name not allowed")
	}

	if !nameReg.MatchString(lastName) {
		return User{}, errors.New("last name not allowed")
	}

	return User{
		ID:       id,
		Name:     strings.ToLower(name),
		LastName: strings.ToLower(lastName),
		Email:    strings.ToLower(email),
	}, nil
}
