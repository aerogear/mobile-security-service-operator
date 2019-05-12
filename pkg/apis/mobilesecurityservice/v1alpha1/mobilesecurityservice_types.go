package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	routev1 "github.com/openshift/api/route/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MobileSecurityServiceSpec defines the desired state of MobileSecurityService
// +k8s:openapi-gen=true
type MobileSecurityServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	//Enviroment Variables for the Service
	DatabaseName                  string `json:"databaseName"`
	DatabasePassword              string `json:"databasePassword"`
	DatabaseUser                  string `json:"databaseUser"`
	DatabaseHost                  string `json:"databaseHost"`
	Port                          int32  `json:"port,omitempty"`
	LogLevel                      string `json:"logLevel,omitempty"`
	LogFormat                     string `json:"logFormat,omitempty"`
	AccessControlAllowOrigin      string `json:"accessControlAllowOrigin"`
	AccessControlAllowCredentials string `json:"accessControlAllowCredentials"`

	//CR mandatory configuration values
	Size                          int32  `json:"size"`
	Image                         string `json:"image"`
	ContainerName                 string `json:"containerName"`
	ClusterProtocol               string `json:"clusterProtocol"`
	MemoryLimit                   string `json:"memoryLimit"`
	MemoryRequest                 string `json:"memoryRequest"`

	//CR optional configuration values
	ConfigMapName           string `json:"configMapName,omitempty"`
	RouteName               string `json:"routeName,omitempty"`
	SkipNamespaceValidation bool   `json:"skipNamespaceValidation,omitempty"`
}

// MobileSecurityServiceStatus defines the observed state of MobileSecurityService
// +k8s:openapi-gen=true
type MobileSecurityServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	ConfigMapName string `json:"configMapName"`
	DeploymentName string `json:"deploymentName"`
	DeploymentStatus v1beta1.DeploymentStatus `json:"deploymentStatus"`
	ServiceName string `json:"serviceName"`
	ServiceStatus v1.ServiceStatus `json:"serviceStatus"`
	RouteName string `json:"routeName"`
	RouteStatus routev1.RouteStatus `json:"routeStatus"`
	AppStatus string `json:"appStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityService is the Schema for the mobilesecurityservices API
// +k8s:openapi-gen=true
type MobileSecurityService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MobileSecurityServiceSpec   `json:"spec,omitempty"`
	Status MobileSecurityServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceList contains a list of MobileSecurityService
type MobileSecurityServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MobileSecurityService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MobileSecurityService{}, &MobileSecurityServiceList{})
}
