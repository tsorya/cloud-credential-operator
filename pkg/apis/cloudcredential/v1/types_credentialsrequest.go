/*
Copyright 2018 The OpenShift Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// FinalizerDeprovision is used on CredentialsRequests to ensure we delete the
	// credentials in AWS before allowing the CredentialsRequest to be deleted in etcd.
	FinalizerDeprovision string = "cloudcredential.openshift.io/deprovision"

	// AnnotationCredentialsRequest is used on Secrets created as a target of CredentailsRequests.
	// The annotation value will map back to the namespace/name of the CredentialsRequest that created
	// or adopted the secret.
	AnnotationCredentialsRequest string = "cloudcredential.openshift.io/credentials-request"

	// AnnotationAWSPolicyLastApplied is added to target Secrets indicating the last AWS policy
	// we successfully applied. It is used to compare if changes are necessary, without requiring
	// AWS credentials to view the actual state.
	AnnotationAWSPolicyLastApplied string = "cloudcredential.openshift.io/aws-policy-last-applied"

	// CloudCredOperatorNamespace is the namespace where the credentials operator runs.
	CloudCredOperatorNamespace = "openshift-cloud-credential-operator"
)

// NOTE: Run "make" to regenerate code after modifying this file

// CredentialsRequestSpec defines the desired state of CredentialsRequest
type CredentialsRequestSpec struct {
	// SecretRef points to the secret where the credentials should be stored once generated.
	// +kubebuilder:validation:Required
	// +required
	SecretRef corev1.ObjectReference `json:"secretRef"`

	// ProviderSpec contains the cloud provider specific credentials specification.
	// +kubebuilder:pruning:PreserveUnknownFields
	ProviderSpec *runtime.RawExtension `json:"providerSpec,omitempty"`

	// ServiceAccountNames contains a list of ServiceAccounts that will use permissions associated with this
	// CredentialsRequest. This is not used by CCO, but the information is needed for being able to properly
	// set up access control in the cloud provider when the ServiceAccounts are used as part of the cloud
	// credentials flow.
	// +optional
	ServiceAccountNames []string `json:"serviceAccountNames,omitempty"`

	// cloudTokenPath is the path where the Kubernetes ServiceAccount token (JSON Web Token) is mounted
	// on the deployment for the workload requesting a credentials secret.
	// The presence of this field in combination with fields such as spec.providerSpec.stsIAMRoleARN
	// indicate that CCO should broker creation of a credentials secret containing fields necessary for
	// token based authentication methods such as with the AWS Secure Token Service (STS).
	// +optional
	CloudTokenPath string `json:"cloudTokenPath,omitempty"`
}

// CredentialsRequestStatus defines the observed state of CredentialsRequest
type CredentialsRequestStatus struct {
	// Provisioned is true once the credentials have been initially provisioned.
	// +kubebuilder:validation:Required
	// +required
	Provisioned bool `json:"provisioned"`

	// LastSyncTimestamp is the time that the credentials were last synced.
	LastSyncTimestamp *metav1.Time `json:"lastSyncTimestamp,omitempty"`

	// LastSyncGeneration is the generation of the credentials request resource
	// that was last synced. Used to determine if the object has changed and
	// requires a sync.
	// +kubebuilder:validation:Required
	// +required
	LastSyncGeneration int64 `json:"lastSyncGeneration"`

	// LastSyncCloudCredsSecretResourceVersion is the resource version of the
	// cloud credentials secret resource when the credentials request resource
	// was last synced. Used to determine if the the cloud credentials have
	// been updated since the last sync.
	// +optional
	LastSyncCloudCredsSecretResourceVersion string `json:"lastSyncCloudCredsSecretResourceVersion,omitempty"`

	// ProviderStatus contains cloud provider specific status.
	// +kubebuilder:pruning:PreserveUnknownFields
	ProviderStatus *runtime.RawExtension `json:"providerStatus,omitempty"`

	// Conditions includes detailed status for the CredentialsRequest
	// +optional
	Conditions []CredentialsRequestCondition `json:"conditions,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CredentialsRequest is the Schema for the credentialsrequests API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type CredentialsRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	// +required
	Spec   CredentialsRequestSpec   `json:"spec"`
	Status CredentialsRequestStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CredentialsRequestList contains a list of CredentialsRequest
type CredentialsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CredentialsRequest `json:"items"`
}

// CredentialsRequestCondition contains details for any of the conditions on a CredentialsRequest object
type CredentialsRequestCondition struct {
	// Type is the specific type of the condition
	// +kubebuilder:validation:Required
	// +required
	Type CredentialsRequestConditionType `json:"type"`
	// Status is the status of the condition
	// +kubebuilder:validation:Required
	// +required
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about the last transition
	Message string `json:"message,omitempty"`
}

// CredentialsRequestConditionType are the valid condition types for a CredentialsRequest
type CredentialsRequestConditionType string

// These are valid conditions for a CredentialsRequest
const (
	// InsufficientCloudCredentials is true when the cloud credentials are deemed to be insufficient
	// to either mint custom creds to satisfy the CredentialsRequest or insufficient to
	// be able to be passed along as-is to satisfy the CredentialsRequest
	InsufficientCloudCredentials CredentialsRequestConditionType = "InsufficientCloudCreds"
	// MissingTargetNamespace is true when the namespace specified to hold the resulting
	// credentials is not present
	MissingTargetNamespace CredentialsRequestConditionType = "MissingTargetNamespace"
	// CredentialsProvisionFailure is true whenver there has been an issue while trying
	// to provision the credentials (either passthrough or minting). Error message will
	// be stored directly in the condition message.
	CredentialsProvisionFailure CredentialsRequestConditionType = "CredentialsProvisionFailure"
	// CredentialsDeprovisionFailure is true whenever there is an error when trying
	// to clean up any previously-created cloud resources
	CredentialsDeprovisionFailure CredentialsRequestConditionType = "CredentialsDeprovisionFailure"
	// Ignored is true when the CredentialsRequest's ProviderSpec is for
	// a different infrastructure platform than what the cluster has been
	// deployed to. This is normal as the release image contains CredentialsRequests for all
	// possible clouds/infrastructure, and cloud-credential-operator will only act on the
	// CredentialsRequests where the cloud/infra matches.
	Ignored CredentialsRequestConditionType = "Ignored"
	// StaleCredentials is true when CredentialsRequest is no longer required and has to be cleaned ip
	StaleCredentials CredentialsRequestConditionType = "StaleCredentials"
	// OrphanedCloudResource is true when CCO was unable to delete a previously created
	// App Registration / Service Principal while pivoting from Mint mode to Passthrough
	OrphanedCloudResource CredentialsRequestConditionType = "OrphanedCloudResource"
)

var (
	// FailureConditionTypes is a list of all conditions where the overall controller status would not
	// be healthy.
	FailureConditionTypes = []CredentialsRequestConditionType{
		InsufficientCloudCredentials,
		MissingTargetNamespace,
		CredentialsProvisionFailure,
		CredentialsDeprovisionFailure,
	}
)
