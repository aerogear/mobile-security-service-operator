package v1alpha1

import (
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MobileSecurityServiceBackupSpec defines the desired state of MobileSecurityServiceBackup
// +k8s:openapi-gen=true
type MobileSecurityServiceBackupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Schedule        string `json:"schedule,omitempty"`
	Image           string `json:"image,omitempty"`
	DatabaseVersion string `json:"databaseVersion,omitempty"`
	ProductName     string `json:"productName,omitempty"`

	AwsS3BucketName               string `json:"awsS3BucketName,omitempty"`
	AwsAccessKeyId                string `json:"awsAccessKeyId,omitempty"`
	AwsSecretAccessKey            string `json:"awsSecretAccessKey,omitempty"`
	AwsCredentialsSecretName      string `json:"awsCredentialsSecretName,omitempty"`
	AwsCredentialsSecretNamespace string `json:"awsCredentialsSecretNamespace,omitempty"`

	EncryptionKeySecretName      string `json:"encryptionKeySecretName,omitempty"`
	EncryptionKeySecretNamespace string `json:"encryptionKeySecretNamespace,omitempty"`
	GpgPublicKey                 string `json:"cpgPublicKey,omitempty"`
	GpgEmail                     string `json:"gpgEmail,omitempty"`
	GpgTrustModel                string `json:"gpgTrustModel,omitempty"`
}

// MobileSecurityServiceBackupStatus defines the observed state of MobileSecurityServiceBackup
// +k8s:openapi-gen=true
type MobileSecurityServiceBackupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	BackupStatus                  string                `json:"backupStatus"`
	CronJobName                   string                `json:"cronJobName"`
	DBSecretName                  string                `json:"dbSecretName"`
	DBSecretData                  map[string]string     `json:"dbSecretData"`
	AWSSecretName                 string                `json:"awsSecretName"`
	AWSSecretData                 map[string]string     `json:"awsSecretData"`
	AwsCredentialsSecretNamespace string                `json:"awsCredentialsSecretNamespace"`
	EncryptionKeySecretName       string                `json:"encryptionKeySecretName"`
	EncryptionKeySecretNamespace  string                `json:"encryptionKeySecretNamespace"`
	EncryptionKeySecretData       map[string]string     `json:"encryptionKeySecretData"`
	HasEncryptionKey              bool                  `json:"hasEncryptionKey"`
	DatabasePodFound              bool                  `json:"databasePodFound"`
	DatabaseServiceFound          bool                  `json:"databaseServiceFound"`
	CronJobStatus                 v1beta1.CronJobStatus `json:"cronJobStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceBackup is the Schema for the mobilesecurityservicedbbackups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type MobileSecurityServiceBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MobileSecurityServiceBackupSpec   `json:"spec,omitempty"`
	Status MobileSecurityServiceBackupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MobileSecurityServiceBackupList contains a list of MobileSecurityServiceBackup
type MobileSecurityServiceBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MobileSecurityServiceBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MobileSecurityServiceBackup{}, &MobileSecurityServiceBackupList{})
}
