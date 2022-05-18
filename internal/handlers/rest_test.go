package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-service/config"
	"user-service/internal/core/domain"
	"user-service/internal/mock"
	"user-service/pkg/dto"
	"user-service/pkg/logging"
)

type RestHandlerTestSuite struct {
	suite.Suite
	MockService *mock.UserService
	TestHandler *HTTPHandler
	TestRouter  *gin.Engine
	Cfg         *config.Config
	TestData    struct {
		User domain.User
	}
}

func (suite *RestHandlerTestSuite) SetupSuite() {
	cfgPath := "../../test/user.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	logger := logging.MockLogger{}

	mockService := new(mock.UserService)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	deliveryHandler := NewRest(mockService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()

	suite.Cfg = cfg
	suite.MockService = mockService
	suite.TestRouter = router
	suite.TestHandler = deliveryHandler
	suite.TestData = struct {
		User domain.User
	}{
		User: domain.User{
			ID:       "test-id-2",
			Name:     "test-name-2",
			LastName: "test-lastname-2",
			Email:    "test@email.com",
		},
	}
}

func (suite *RestHandlerTestSuite) SetupTest() {
	suite.MockService.ExpectedCalls = nil
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll() {
	suite.MockService.On("GetAll").Return([]domain.User{suite.TestData.User}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.UserListResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.Len(responseObject, 1)

	suite.EqualValues(suite.TestData.User.ID, responseObject[0].ID)
	suite.EqualValues(suite.TestData.User.Name, responseObject[0].Name)
	suite.EqualValues(suite.TestData.User.LastName, responseObject[0].LastName)
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll_NoneFound() {
	suite.MockService.On("GetAll").Return([]domain.User{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get() {
	suite.MockService.On("Get", suite.TestData.User.ID).Return(suite.TestData.User, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%s", suite.TestData.User.ID), nil)
	request.Header.Set("X-User-Id", suite.TestData.User.ID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.UserResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.User.Name, responseObject.Name)
	suite.EqualValues(suite.TestData.User.ID, responseObject.ID)
	suite.EqualValues(suite.TestData.User.Email, responseObject.Email)
	suite.EqualValues(suite.TestData.User.LastName, responseObject.LastName)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_BadID() {
	suite.MockService.On("Get", "test").Return(domain.User{}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%s", "test"), nil)
	request.Header.Set("X-User-Id", suite.TestData.User.ID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_NotFound() {
	suite.MockService.On("Get", suite.TestData.User.ID).Return(domain.User{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%s", suite.TestData.User.ID), nil)
	request.Header.Set("X-User-Id", suite.TestData.User.ID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create() {
	suite.MockService.On("Create", suite.TestData.User.ID, suite.TestData.User.Name, suite.TestData.User.LastName, suite.TestData.User.Email).Return(suite.TestData.User, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateUser(suite.TestData.User))

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusCreated, rr.Code)

	var responseObject dto.UserResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.User.Name, responseObject.Name)
	suite.EqualValues(suite.TestData.User.ID, responseObject.ID)
	suite.EqualValues(suite.TestData.User.Email, responseObject.Email)
	suite.EqualValues(suite.TestData.User.LastName, responseObject.LastName)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_BadInput() {
	rr := httptest.NewRecorder()

	data, err := json.Marshal(struct {
		Test string
	}{
		Test: suite.TestData.User.ID,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_CouldNotCreate() {
	suite.MockService.On("Create", suite.TestData.User.ID, suite.TestData.User.Name, suite.TestData.User.LastName, suite.TestData.User.Email).Return(domain.User{}, errors.New("could not create"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateUser(suite.TestData.User))

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Update() {
	updated := suite.TestData.User
	updated.Name = "new-name"

	suite.MockService.On("UpdateUserDetails", suite.TestData.User.ID, updated.Name, suite.TestData.User.LastName, suite.TestData.User.Email).Return(updated, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateUser(updated))

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/%s", suite.TestData.User.ID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.UserResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.User.ID, responseObject.ID)
	suite.EqualValues(updated.Name, responseObject.Name)
}

func (suite *RestHandlerTestSuite) TestHandler_Update_CouldNotCreate() {
	updated := suite.TestData.User
	updated.Name = "new-name"

	suite.MockService.On("UpdateUserDetails", suite.TestData.User.ID, updated.Name, suite.TestData.User.LastName, suite.TestData.User.Email).Return(domain.User{}, errors.New("could not update"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateUser(updated))

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/%s", suite.TestData.User.ID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func TestIntegration_RestHandlerTestSuite(t *testing.T) {
	testSuite := new(RestHandlerTestSuite)
	suite.Run(t, testSuite)
}
