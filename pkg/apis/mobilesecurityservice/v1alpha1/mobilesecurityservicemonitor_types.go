package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MobileSecurityServiceMonitorSpec defines the desired state of MobileSecurityServiceMonitor
// +k8s:openapi-gen=true
type MobileSecurityServiceMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	NamespaceSelector string   `json:namespaceSelector,omitempty"`
	LabelSelector          *metav1.LabelSelector   `json:"labelSelector,omitempty"`
	ClusterHost                   string `json:"clusterHost"`
	HostSufix                     string `json:"hostSufix"`
	Protocol                      string `json:"protocol"`
}

// MobileSecurityServiceMonitorStatus defines the observed state of MobileSecurityServiceMonitor
// +k8s:openapi-gen=true
type MobileSecurityServiceMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceMonitor is the Schema for the mobilesecurityservicemonitors API
// +k8s:openapi-gen=true
type MobileSecurityServiceMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MobileSecurityServiceMonitorSpec   `json:"spec,omitempty"`
	Status MobileSecurityServiceMonitorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceMonitorList contains a list of MobileSecurityServiceMonitor
type MobileSecurityServiceMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MobileSecurityServiceMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MobileSecurityServiceMonitor{}, &MobileSecurityServiceMonitorList{})
}
