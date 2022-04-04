package main

import (
	"log"
	"os"
	"user-service/internal/core/services/rabbitmq_service"
	"user-service/internal/core/services/user_service"
	"user-service/internal/handlers"
	"user-service/internal/repositories"
	"user-service/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/user"

func main() {
	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	userRepository, err := repositories.NewCockroachDB(dbConn)

	if err != nil {
		panic(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	userService := user_service.New(userRepository, rmqPublisher)

	router := gin.New()

	userHandler := handlers.NewRest(userService, router)
	userHandler.SetupEndpoints()
	userHandler.SetupSwagger()

	port := GetEnvOrDefault("PORT", defaultPort)

	log.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
