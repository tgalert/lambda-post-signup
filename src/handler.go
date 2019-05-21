package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	//"log"
)

type RabbitMqCreds struct {
	username string
	password string
	vhost    string
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() {
	userId := fmt.Sprint(uuid.New())
	rabbitMqCreds := generateRabbitMqCreds(userId)
	createRabbitMqUserVhost(rabbitMqCreds)
	//saveRabbitMqCredsInCluster(rabbitMqCreds, userId)
	//launchBackend(userId)
}

// Generate RabbitMQ credentials (username, password, vhost) for new user
func generateRabbitMqCreds(userId string) RabbitMqCreds {
	return RabbitMqCreds{userId, userId, userId}
}

// Serialise RabbitMQ credentials struct to JSON
func (c RabbitMqCreds) serialize() string {
	return fmt.Sprintf(`{"username":"%s","password":"%s","vhost":"%s"`, c.username, c.password, c.vhost)
}
