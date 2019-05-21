package mobilesecurityservicedb

import (
	"github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Centralized mock objects for use in tests
var (
	instanceOne = v1alpha1.MobileSecurityServiceDB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service-db",
			Namespace: "mobile-security-service",
		},
		Spec: v1alpha1.MobileSecurityServiceDBSpec{
			Image:                   "centos/postgresql-96-centos7",
			Size:                    1,
			ContainerName:           "database",
			DatabaseNameParam:       "POSTGRESQL_DATABASE",
			DatabasePasswordParam:   "POSTGRESQL_PASSWORD",
			DatabaseUserParam:       "POSTGRESQL_USER",
			DatabasePort:            5432,
			DatabaseMemoryLimit:     "512Mi",
			DatabaseMemoryRequest:   "512Mi",
			DatabaseStorageRequest:  "1Gi",
			DatabaseName:            "mobile_security_service",
			DatabasePassword:        "postgres",
			DatabaseUser:            "postgresql",
			SkipNamespaceValidation: true,
		},
	}

	instanceTwo = v1alpha1.MobileSecurityServiceDB{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service-db",
			Namespace: "mobile-security-service-namespace",
		},
		Spec: v1alpha1.MobileSecurityServiceDBSpec{
			Image:                   "centos/postgresql-96-centos7",
			Size:                    1,
			ContainerName:           "database",
			DatabaseNameParam:       "POSTGRESQL_DATABASE",
			DatabasePasswordParam:   "POSTGRESQL_PASSWORD",
			DatabaseUserParam:       "POSTGRESQL_USER",
			DatabasePort:            5432,
			DatabaseMemoryLimit:     "512Mi",
			DatabaseMemoryRequest:   "512Mi",
			DatabaseStorageRequest:  "1Gi",
			DatabaseName:            "mobile_security_service",
			DatabasePassword:        "postgres",
			DatabaseUser:            "postgresql",
			SkipNamespaceValidation: true,
		},
	}

	serviceInstance = v1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service-app",
			Namespace: "mobile-security-service",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:                    1,
			MemoryLimit:             "512Mi",
			MemoryRequest:           "512Mi",
			ClusterProtocol:         "http",
			ConfigMapName:           "mss-config",
			RouteName:               "mss-route",
			SkipNamespaceValidation: true,
		},
	}
)
