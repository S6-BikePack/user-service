package main

import (
	"context"
	"fmt"
	"os"
	"user-service/config"
	"user-service/internal/core/services"
	"user-service/internal/handlers"
	"user-service/internal/repositories"
	"user-service/pkg/logging"
	"user-service/pkg/rabbitmq"
	"user-service/pkg/tracing"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Logging and Tracing
	//--------------------------------------------------------------------------------------

	logger, err := logging.NewSimpleLogger(cfg)

	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewOpenTracing(cfg.Server.Service, cfg.Tracing.Host, cfg.Tracing.Port)

	if err != nil {
		logger.Warning(context.Background(), "Failed to setup tracing: %v", err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	if tracer != nil {
		if err = db.Use(otelgorm.NewPlugin(otelgorm.WithTracerProvider(tracer))); err != nil {
			logger.Fatal(context.Background(), err)
		}
	}

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	userRepository, err := repositories.NewUserRepository(db)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	//--------------------------------------------------------------------------------------
	// Setup RabbitMQ
	//--------------------------------------------------------------------------------------

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	rmqPublisher := services.NewRabbitMQPublisher(rmqServer, tracer, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	userService := services.NewUserService(userRepository, rmqPublisher)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()

	if tracer != nil {
		router.Use(otelgin.Middleware(cfg.Server.Service, otelgin.WithTracerProvider(tracer)))
	}

	deliveryHandler := handlers.NewRest(userService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()

	logger.Fatal(context.Background(), router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
