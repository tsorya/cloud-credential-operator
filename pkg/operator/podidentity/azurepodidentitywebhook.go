package podidentity

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

type AzurePodIdentity struct {
}

func (a AzurePodIdentity) ShouldBeDeployed(clientSet *kubernetes.Clientset, namespace string) bool {
	secret, err := clientSet.CoreV1().Secrets(namespace).Get(context.TODO(),
		"azure-credentials", metav1.GetOptions{})
	if err != nil {
		return false
	}
	_, ok := secret.Data["azure_tenant_id"]
	return ok
}

func (a AzurePodIdentity) Deployment() string {
	return "v4.1.0/azure-pod-identity-webhook/deployment.yaml"
}

func (a AzurePodIdentity) GetImagePullSpec() string {
	return os.Getenv("AZURE_POD_IDENTITY_WEBHOOK_IMAGE")
}
