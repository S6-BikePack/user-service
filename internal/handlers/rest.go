package handlers

import (
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag/example/basic/docs"
	"user-service/internal/core/ports"
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
	api.GET("/user/all", handler.GetAll)
	api.GET("/user", handler.Get)
	api.POST("/user", handler.Create)
	api.PUT("/user/:id", handler.Update)
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
// @Router       /api/users/all [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	users, err := handler.userService.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, users)
}

// Get godoc
// @Summary  get user
// @Schemes
// @Param        id     path  string           true  "User id"
// @Description  gets a user from the system by its ID
// @Produce      json
// @Success      200  {object}  domain.User
// @Router       /api/user [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	user, err := handler.userService.Get(c.GetHeader("X-User-Id"))

	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	c.JSON(200, user)
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
		c.AbortWithStatus(500)
	}

	user, err := handler.userService.Create(c.GetHeader("X-User-Id"), body.Name, body.LastName, c.GetHeader("X-User-Email"))

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreateUser(user))
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
	body := dto.BodyCreateUser{}
	err := c.BindJSON(&body)

	if err != nil {
		c.AbortWithStatus(500)
	}

	user, err := handler.userService.UpdateUserDetails(c.Param("id"), body.Name, body.LastName, "")

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, dto.BuildResponseCreateUser(user))
}
