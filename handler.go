package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func HandleRequest() {
	// TODO: get userId from event source (Cognito sub)
	userId := fmt.Sprint(uuid.New())
	rabbitMq := getRabbitMqCredentials(userId)
	createDeployment(userId, rabbitMq)
}

func main() {
	lambda.Start(HandleRequest)
}

type RabbitMqCredentials struct {
	user     string
	password string
	vhost    string
}

func (r RabbitMqCredentials) getUri() string {
	return fmt.Sprintf("amqp://%s:%s@rabbitmq/%s", r.user, r.password, r.vhost)
}

func getRabbitMqCredentials(userId string) RabbitMqCredentials {
	// TODO: generate different values for user and password
	return RabbitMqCredentials{userId, userId, userId}
}

func createDeployment(name string, rabbitMqCredentials RabbitMqCredentials) {

	// Path to kubeconfig file
	kubeconfig := "./kubeconfig"

	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an API Clientset (k8s.io/client-go/kubernetes)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create an AppsV1Client (k8s.io/client-go/kubernetes/typed/apps/v1)
	appsV1Client := clientset.AppsV1()

	// Specification of the Deployment (k8s.io/api/apps/v1)
	deploymentSpec := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: func() *int32 { i := int32(1); return &i }(),
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "tgalert",
				},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						"app": "tgalert",
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:  "tgalert",
							Image: "weibeld/tg-monitor:core-0.0.1",
							//Image: "weibeld/tmp",
							Env: []coreV1.EnvVar{
								// Non-sensitive vars
								{
									Name:  "MAILGUN_SENDING_ADDRESS",
									Value: "tgalert@quantumsense.ai",
								},
								{
									Name:  "MAILGUN_DOMAIN",
									Value: "quantumsense.ai",
								},
								// Sensitive vars (dynamically generated or from secret)
								{
									Name:  "AMQP_URI",
									Value: rabbitMqCredentials.getUri(),
								},
								{
									Name: "TG_API_ID",
									ValueFrom: &coreV1.EnvVarSource{
										SecretKeyRef: &coreV1.SecretKeySelector{
											LocalObjectReference: coreV1.LocalObjectReference{
												Name: "secrets",
											},
											Key: "TG_API_ID",
										},
									},
								},
								{
									Name: "TG_API_HASH",
									ValueFrom: &coreV1.EnvVarSource{
										SecretKeyRef: &coreV1.SecretKeySelector{
											LocalObjectReference: coreV1.LocalObjectReference{
												Name: "secrets",
											},
											Key: "TG_API_HASH",
										},
									},
								},
								{
									Name: "MAILGUN_API_KEY",
									ValueFrom: &coreV1.EnvVarSource{
										SecretKeyRef: &coreV1.SecretKeySelector{
											LocalObjectReference: coreV1.LocalObjectReference{
												Name: "secrets",
											},
											Key: "MAILGUN_API_KEY",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	deployment, err := appsV1Client.Deployments("default").Create(deploymentSpec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %s\n", deployment.ObjectMeta.Name)

}
