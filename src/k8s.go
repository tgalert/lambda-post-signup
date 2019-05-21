package main

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	appsV1 "k8s.io/api/apps/v1"
	//coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var clientset *kubernetes.Clientset

// Save RabbitMQ credentials for new user in Kubernetes secret in cluster
func saveRabbitMqCredsInCluster(creds RabbitMqCreds, userId string) {
	if clientset == nil {
		initKubernetesClient()
	}
	secretsClient := clientset.CoreV1().Secrets("default")
	// Read RabbitMQ credentials secret
	secretName := "tgalert/rabbitmq-creds"
	secret, err := secretsClient.Get(secretName, metaV1.GetOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	// Add a new key/value pair to secret
	if secret.StringData == nil {
		secret.StringData = map[string]string{}
	}
	secret.StringData[userId] = creds.serialize()
	_, err = secretsClient.Update(secret)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Create a backend pod (controlled by a deployment) for the new user
func launchBackend(userId string) {
	if clientset == nil {
		initKubernetesClient()
	}
	// Read deployment specification from YAML file
	file := "./assets/backend.yml"
	log.Printf("Reading Kubernetes deployment specification from file %s", file)
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	var spec appsV1.Deployment
	err = yaml.Unmarshal(bytes, &spec)
	if err != nil {
		log.Fatal(err.Error())
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
	log.Printf("Creating Kubernetes deployment %s", spec.ObjectMeta.Name)
	_, err = clientset.AppsV1().Deployments("default").Create(&spec)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deployment %s created successfully", spec.ObjectMeta.Name)
}

// Helper: initialise a Kubernetes client
func initKubernetesClient() {
	kubeconfig := "./assets/kubeconfig"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
