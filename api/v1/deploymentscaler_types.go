package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentTarget identifies a Deployment to control
type DeploymentTarget struct {
	// Name of the Deployment (required)
	Name string `json:"name"`
	// Namespace of the Deployment (optional; default: same as this CR)
	Namespace *string `json:"namespace,omitempty"`
}

type DeploymentScalerSpec struct {
	// The target Deployment to scale
	Target DeploymentTarget `json:"target"`
	// Desired replicas for the target Deployment (>=0)
	Replicas int32 `json:"replicas"`
}
// DeploymentScalerStatus defines the observed state of DeploymentScaler.
type DeploymentScalerStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the DeploymentScaler resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DeploymentScaler is the Schema for the deploymentscalers API
type DeploymentScaler struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of DeploymentScaler
	// +required
	Spec DeploymentScalerSpec `json:"spec"`

	// status defines the observed state of DeploymentScaler
	// +optional
	Status DeploymentScalerStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// DeploymentScalerList contains a list of DeploymentScaler
type DeploymentScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploymentScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploymentScaler{}, &DeploymentScalerList{})
}
