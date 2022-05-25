package handlers

import (
	"net/http"
	"user-service/config"
	"user-service/internal/core/interfaces"
	"user-service/pkg/authorization"
	"user-service/pkg/dto"
	"user-service/pkg/logging"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag/example/basic/docs"
	"go.opentelemetry.io/otel/trace"
)

type HTTPHandler struct {
	userService interfaces.UserService
	router      *gin.Engine
	logger      logging.Logger
	config      *config.Config
}

func NewRest(userService interfaces.UserService, router *gin.Engine, logger logging.Logger, config *config.Config) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
		router:      router,
		config:      config,
		logger:      logger,
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
	docs.SwaggerInfo.Title = handler.config.Server.Service + " API"
	docs.SwaggerInfo.Description = handler.config.Server.Description

	handler.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (handler *HTTPHandler) SetupHealthprobe() {
	handler.router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

// GetAll godoc
// @Summary  get all users
// @Schemes
// @Description  gets all users in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.UserListResponse
// @Router       /api/users [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if authorization.NewRest(c).AuthorizeAdmin() {

		users, err := handler.userService.GetAll(ctx)

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, dto.CreateUserListResponse(users))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Get godoc
// @Summary  get user
// @Schemes
// @Param        id     path  string           true  "User id"
// @Description  gets a user from the system by its ID
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Router       /api/users/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		user, err := handler.userService.Get(ctx, c.Param("id"))

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, dto.CreateUserResponse(user))
		return
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
// @Success      200  {object}  dto.UserResponse
// @Router       /api/users [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := dto.BodyCreateUser{}
	err := c.BindJSON(&body)

	if err != nil || body.ID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(body.ID) {

		user, err := handler.userService.Create(ctx, body.ID, body.Name, body.LastName, body.Email)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusCreated, dto.CreateUserResponse(user))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// Update godoc
// @Summary  update user
// @Schemes
// @Description  updates a users name, last name
// @Accept       json
// @Param        user  body  dto.BodyCreateUser  true  "Update user"
// @Param        id  path  string  true  "User id"
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Router       /api/users/{id} [put]
func (handler *HTTPHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || auth.AuthorizeMatchingId(c.Param("id")) {

		body := dto.BodyCreateUser{}
		err := c.BindJSON(&body)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		user, err := handler.userService.UpdateUserDetails(ctx, c.Param("id"), body.Name, body.LastName, body.Email)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			handler.logger.Error(ctx, err.Error())
			return
		}

		c.JSON(http.StatusOK, dto.CreateUserResponse(user))

	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
