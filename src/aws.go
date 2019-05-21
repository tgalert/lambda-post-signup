package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
)

var sess *session.Session

// Get credentials of RabbitMQ administrator user from AWS Secrets Manager
func getRabbitMqAdminCreds() (username, password string) {
	if sess == nil {
		initAwsSession()
	}
	secretName := "tgalert/rabbitmq-admin"
	log.Printf("Getting RabbitMQ admin credentials from secret %s in AWS Secrets Manager", secretName)
	secret, err := secretsmanager.New(sess).GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// Parse secret of format {"username":"xxx","password":"yyy"}
	var creds map[string]string
	json.Unmarshal([]byte(*secret.SecretString), &creds)
	log.Printf("RabbitMQ admin credentials: username=%s, password=********", creds["username"])
	return creds["username"], creds["password"]
}

// Get RabbitMQ host from AWS Systems Manager Parameter Store
func getRabbitMqHost() string {
	if sess == nil {
		initAwsSession()
	}
	paramName := "/tgalert/rabbitmq-host"
	log.Printf("Getting RabbitMQ host from parameter %s in AWS Systems Manager Parameter Store", paramName)
	param, err := ssm.New(sess).GetParameter(&ssm.GetParameterInput{
		Name: aws.String(paramName),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	host := *param.Parameter.Value
	log.Printf("RabbitMQ host: %s", host)
	return host
}

// Helper: initialise an AWS SDK session struct
func initAwsSession() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		log.Printf("AWS API request: %s %v", r.ClientInfo.ServiceName, *r.Operation)
	})
}
