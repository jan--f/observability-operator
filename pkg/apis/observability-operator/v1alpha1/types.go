// +groupName=observability-operator.rhobs
// +kubebuilder:rbac:groups=observability-operator.rhobs,resources=configs,verbs=list;get;watch
// +kubebuilder:rbac:groups=observability-operator.rhobs,resources=configs/status;configs/finalizers,verbs=get;update

package v1alpha1

import (
	// monv1 "github.com/rhobs/obo-prometheus-operator/pkg/apis/monitoring/v1"
	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config is the Schema for the observability-operator config API
// +k8s:openapi-gen=true
// +kubebuilder:resource
// +kubebuilder:subresource:status
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSpec   `json:"spec,omitempty"`
	Status ConfigStatus `json:"status,omitempty"`
}

// ConfigList contains a list of Configs
// +kubebuilder:resource
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config `json:"items"`
}

// ConfigSpec is the specification for the observability-opertor
type ConfigSpec struct {
}

// ConfigStatus defines the observed state of the Config.
// It should always be reconstructable from the state of the operator and/or outside world.
type ConfigStatus struct {
	// Conditions provide status information about the Config
	// +listType=atomic
	// Conditions []Condition `json:"conditions"`
}
