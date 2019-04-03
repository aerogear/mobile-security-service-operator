package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MobileSecurityServiceBindSpec defines the desired state of MobileSecurityServiceBind
// +k8s:openapi-gen=true
type MobileSecurityServiceBindSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Size int32  `json:"size"`
}

// MobileSecurityServiceBindStatus defines the observed state of MobileSecurityServiceBind
// +k8s:openapi-gen=true
type MobileSecurityServiceBindStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceBind is the Schema for the mobilesecurityservicebinds API
// +k8s:openapi-gen=true
type MobileSecurityServiceBind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MobileSecurityServiceBindSpec   `json:"spec,omitempty"`
	Status MobileSecurityServiceBindStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceBindList contains a list of MobileSecurityServiceBind
type MobileSecurityServiceBindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MobileSecurityServiceBind `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MobileSecurityServiceBind{}, &MobileSecurityServiceBindList{})
}
