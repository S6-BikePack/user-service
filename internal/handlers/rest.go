package handlers

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag/example/basic/docs"
	"net/http"
	"user-service/internal/core/ports"
	"user-service/pkg/authorization"
	"user-service/pkg/dto"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	userService ports.UserService
	router      *gin.Engine
}

func NewRest(userService ports.UserService, router *gin.Engine) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
		router:      router,
	}
}

func (handler *HTTPHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/users", handler.GetAll)
	api.GET("/users/:id", handler.Get)
	api.POST("/users", handler.Create)
	api.PUT("/users/:id", handler.Update)
}

func (handler *HTTPHandler) SetupSwagger() {
	docs.SwaggerInfo.Title = "User service API"
	docs.SwaggerInfo.Description = "The user service manages all users for the BikePack system."

	handler.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all users
// @Schemes
// @Description  gets all users in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []domain.User
// @Router       /api/users [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	if authorization.NewRest(c).AuthorizeAdmin() {

		users, err := handler.userService.GetAll()

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, users)
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Get godoc
// @Summary  get user
// @Schemes
// @Param        id     path  string           true  "User id"
// @Description  gets a user from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.User
// @Router       /api/users/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		user, err := handler.userService.Get(c.Param("id"))

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, user)
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Create godoc
// @Summary  create user
// @Schemes
// @Description  creates a new user
// @Accept       json
// @Param        user  body  dto.BodyCreateUser  true  "Add user"
// @Produce      json
// @Success      200  {object}  dto.ResponseCreateUser
// @Router       /api/users [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	body := dto.BodyCreateUser{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(body.ID) {

		user, err := handler.userService.Create(body.ID, body.Name, body.LastName, body.Email)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, dto.BuildResponseCreateUser(user))
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Update godoc
// @Summary  update user
// @Schemes
// @Description  updates a users name, last name
// @Accept       json
// @Param        user  body  dto.BodyUpdateUser  true  "Update user"
// @Param        id  path  string  true  "User id"
// @Produce      json
// @Success      200  {object}  dto.ResponseUpdateUser
// @Router       /api/users/{id} [put]
func (handler *HTTPHandler) Update(c *gin.Context) {
	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		body := dto.BodyCreateUser{}
		err := c.BindJSON(&body)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		user, err := handler.userService.UpdateUserDetails(c.Param("id"), body.Name, body.LastName, body.Email)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, dto.BuildResponseCreateUser(user))

	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
