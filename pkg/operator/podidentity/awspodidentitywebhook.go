package podidentity

import (
	"k8s.io/client-go/kubernetes"
	"os"
)

type PodIdentityInterface interface {
	Deployment() string
	GetImagePullSpec() string
	ShouldBeDeployed(clientSet *kubernetes.Clientset, namespace string) bool
}

type AwsPodIdentity struct {
}

func (a AwsPodIdentity) ShouldBeDeployed(clientSet *kubernetes.Clientset, namespace string) bool {
	return true
}

func (a AwsPodIdentity) Deployment() string {
	return "v4.1.0/aws-pod-identity-webhook/deployment.yaml"
}

func (a AwsPodIdentity) GetImagePullSpec() string {
	return os.Getenv("AWS_POD_IDENTITY_WEBHOOK_IMAGE")
}
