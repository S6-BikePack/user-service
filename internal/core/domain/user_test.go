package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	user *User
}

func (s *Suite) SetupSuite() {
	s.user = &User{
		ID:       "test-id",
		Name:     "test-name",
		LastName: "test-lastname",
		Email:    "test@test.com",
	}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestUser_NewUser() {
	res, err := NewUser(s.user.ID, s.user.Name, s.user.LastName, s.user.Email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *s.user, res)
}

func (s *Suite) TestUser_NewUserMissingInfo() {
	res, err := NewUser(s.user.ID, "", s.user.LastName, s.user.Email)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), User{}, res)
}

func (s *Suite) TestUser_NewUserInvalidName() {
	res, err := NewUser(s.user.ID, "2222", s.user.LastName, s.user.Email)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), User{}, res)
}

func (s *Suite) TestUser_NewUserInvalidLastName() {
	res, err := NewUser(s.user.ID, s.user.Name, "2222", s.user.Email)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), User{}, res)
}

func (s *Suite) TestUser_NewUserInvalidEmail() {
	res, err := NewUser(s.user.ID, s.user.Name, s.user.LastName, "test")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), User{}, res)
}
