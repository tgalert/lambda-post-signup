package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"io/ioutil"
	appsV1 "k8s.io/api/apps/v1"
	//coreV1 "k8s.io/api/core/v1"
	//metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type RabbitMqCreds struct {
	username string
	password string
	vhost    string
}

var clientset *kubernetes.Clientset

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() {
	// TODO: get userId from event source (Cognito sub)
	userId := fmt.Sprint(uuid.New())

	initKubernetesClient()

	rabbitMqCreds := generateRabbitMqCreds(userId)
	createRabbitMqUser(rabbitMqCreds)
	saveRabbitMqCreds(rabbitMqCreds)

	createBackend(userId)
}

// Initialises 'clientset' variable
func initKubernetesClient() {
	kubeconfig := "./kubeconfig"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

// Generate RabbitMQ credentials (username, password, vhost) for the given user
func generateRabbitMqCreds(userId string) RabbitMqCreds {
	return RabbitMqCreds{userId, userId, userId}
}

// Create a RabbitMQ user with the given credentials, as well as a vhost, and
// grant the new RabbitMQ user permissions for only this vhost
func createRabbitMqUser(creds RabbitMqCreds) {
}

// Save the credentials for this RabbitMQ user in a secret in the cluster
func saveRabbitMqCreds(creds RabbitMqCreds) {
}

// Create a backend pod (controlled by a deployment) for the new user
func createBackend(userId string) {

	// Read deployment specification from YAML file
	bytes, err := ioutil.ReadFile("deployment.yml")
	if err != nil {
		panic(err.Error())
	}
	var spec appsV1.Deployment
	err = yaml.Unmarshal(bytes, &spec)
	if err != nil {
		panic(err.Error())
	}

	// Set user-specific values in deployment specification
	spec.ObjectMeta.Name = userId
	env := spec.Spec.Template.Spec.Containers[0].Env
	for _, e := range env {
		if e.Name == "RABBITMQ_CREDS" {
			e.ValueFrom.SecretKeyRef.Key = userId
		}
	}

	// Create deployment
	_, err = clientset.AppsV1().Deployments("default").Create(&spec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %s\n", spec.ObjectMeta.Name)
}
