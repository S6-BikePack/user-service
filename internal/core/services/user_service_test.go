package services

import (
	"context"
	"errors"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"user-service/internal/core/domain"
	"user-service/internal/core/interfaces"
	"user-service/internal/mock"
)

type UserServiceTestSuite struct {
	suite.Suite
	MockRepository *mock.UserRepository
	MockPublisher  *mock.MessageBusPublisher
	TestService    interfaces.UserService
	TestData       struct {
		User domain.User
	}
}

func (suite *UserServiceTestSuite) SetupSuite() {
	repository := new(mock.UserRepository)
	publisher := new(mock.MessageBusPublisher)

	srv := NewUserService(repository, publisher)

	suite.MockRepository = repository
	suite.MockPublisher = publisher
	suite.TestService = srv
	suite.TestData = struct {
		User domain.User
	}{
		User: domain.User{
			ID:       "test-id",
			Name:     "test-name",
			LastName: "test-lastname",
			Email:    "test@email.com",
		},
	}
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.MockPublisher.ExpectedCalls = nil
	suite.MockRepository.ExpectedCalls = nil
}

func (suite *UserServiceTestSuite) TestUserService_GetAll() {
	suite.MockRepository.On("GetAll").Return([]domain.User{suite.TestData.User}, nil)

	result, err := suite.TestService.GetAll(context.Background())

	suite.NoError(err)

	suite.MockRepository.AssertCalled(suite.T(), "GetAll")
	suite.Equal(1, len(result))
	suite.EqualValues(suite.TestData.User, result[0])
}

func (suite *UserServiceTestSuite) TestUserService_Get() {
	suite.MockRepository.On("Get", suite.TestData.User.ID).Return(suite.TestData.User, nil)

	result, err := suite.TestService.Get(context.Background(), suite.TestData.User.ID)

	suite.NoError(err)

	suite.MockRepository.AssertCalled(suite.T(), "Get", suite.TestData.User.ID)
	suite.EqualValues(suite.TestData.User, result)
}

func (suite *UserServiceTestSuite) TestUserService_Get_NotFound() {
	suite.MockRepository.On("Get", suite.TestData.User.ID).Return(domain.User{}, errors.New("could not find user"))

	result, err := suite.TestService.Get(context.Background(), suite.TestData.User.ID)

	suite.Error(err)
	suite.EqualValues(domain.User{}, result)

	suite.MockRepository.AssertCalled(suite.T(), "Get", suite.TestData.User.ID)
}

func (suite *UserServiceTestSuite) TestUserService_Create() {
	suite.MockRepository.On("GetUser", suite.TestData.User.ID).Return(suite.TestData.User, nil)
	suite.MockRepository.On("Save", mock2.Anything).Return(suite.TestData.User, nil)
	suite.MockPublisher.On("CreateUser", suite.TestData.User).Return(nil)

	result, err := suite.TestService.Create(context.Background(), suite.TestData.User.ID, suite.TestData.User.Name, suite.TestData.User.LastName, suite.TestData.User.Email)

	suite.NoError(err)

	suite.MockPublisher.AssertCalled(suite.T(), "CreateUser", suite.TestData.User)
	suite.EqualValues(suite.TestData.User, result)
}

func (suite *UserServiceTestSuite) TestUserService_Create_MissingData() {
	_, err := suite.TestService.Create(context.Background(), suite.TestData.User.ID, suite.TestData.User.Name, suite.TestData.User.LastName, "")

	suite.MockRepository.AssertNotCalled(suite.T(), "Save")
	suite.Error(err)
}

func (suite *UserServiceTestSuite) TestUserService_Create_CouldNotSave() {
	suite.MockRepository.On("GetUser", suite.TestData.User.ID).Return(suite.TestData.User, nil)
	suite.MockRepository.On("Save", mock2.Anything).Return(domain.User{}, errors.New("could not save user"))
	suite.MockPublisher.On("CreateUser", suite.TestData.User).Return(nil)

	_, err := suite.TestService.Create(context.Background(), suite.TestData.User.ID, suite.TestData.User.Name, suite.TestData.User.LastName, suite.TestData.User.Email)

	suite.Error(err)

	suite.MockPublisher.AssertNotCalled(suite.T(), "CreateUser")
}

func (suite *UserServiceTestSuite) TestUserService_UpdateUserDetails() {
	updated := suite.TestData.User
	updated.Name = "new-name"
	updated.LastName = "new-last-name"

	suite.MockRepository.On("Get", suite.TestData.User.ID).Return(suite.TestData.User, nil)
	suite.MockRepository.On("Update", updated).Return(updated, nil)
	suite.MockPublisher.On("UpdateUserDetails", updated).Return(nil)

	result, err := suite.TestService.UpdateUserDetails(context.Background(), updated.ID, updated.Name, updated.LastName, updated.Email)

	suite.NoError(err)

	suite.EqualValues(updated, result)
}

func (suite *UserServiceTestSuite) TestUserService_UpdateServiceArea_UserNotFound() {
	updated := suite.TestData.User
	updated.Name = "new-name"
	updated.LastName = "new-last-name"

	suite.MockRepository.On("Get", suite.TestData.User.ID).Return(domain.User{}, errors.New("user not found"))

	_, err := suite.TestService.UpdateUserDetails(context.Background(), updated.ID, updated.Name, updated.LastName, updated.Email)

	suite.Error(err)
}

func (suite *UserServiceTestSuite) TestUserService_UpdateServiceArea_CouldNotUpdate() {
	updated := suite.TestData.User
	updated.Name = "new-name"
	updated.LastName = "new-last-name"

	suite.MockRepository.On("Get", suite.TestData.User.ID).Return(suite.TestData.User, nil)
	suite.MockRepository.On("Update", updated).Return(suite.TestData.User, errors.New("could not update user"))

	result, err := suite.TestService.UpdateUserDetails(context.Background(), updated.ID, updated.Name, updated.LastName, updated.Email)

	suite.Error(err)
	suite.EqualValues(suite.TestData.User, result)
}

func TestUnit_UserServiceTestSuite(t *testing.T) {
	testSuite := new(UserServiceTestSuite)
	suite.Run(t, testSuite)
}
