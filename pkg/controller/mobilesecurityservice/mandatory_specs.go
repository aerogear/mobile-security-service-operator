package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

const (
	databaseName                  = "mobile_security_service"
	databasePassword              = "postgres"
	databaseUser                  = "postgresql"
	databaseHost                  = "mobile-security-service-db"
	port                          = 3000
	logLevel                      = "info"
	logFormat                     = "json"
	accessControlAllowOrigin      = "*"
	accessControlAllowCredentials = "false"
	size                          = 1
	clusterProtocol               = "http"
	memoryLimit                   = "512Mi"
	memoryRequest                 = "512Mi"
	image                         = "aerogear/mobile-security-service:master"
	containerName                 = "application"
	oAuthImage                    = "docker.io/openshift/oauth-proxy:v1.1.0"
	oAuthContainerName            = "oauth-proxy"
	configMapName                 = "mobile-security-service-config"
	routeName                     = "route"
)

// addMandatorySpecsDefinitions will add the specs which are mandatory for Mobile Security Service CR in the case them
// not be applied
func addMandatorySpecsDefinitions(mss *mobilesecurityservicev1alpha1.MobileSecurityService) {

	/*
		Environment Variables
		---------------------
		The following values are used to create the ConfigMap and the Environment Variables which will use these values
		These values are used for both the Mobile Security Service and its Database
	*/

	if mss.Spec.DatabaseName == "" {
		mss.Spec.DatabaseName = databaseName
	}

	if mss.Spec.DatabasePassword == "" {
		mss.Spec.DatabasePassword = databasePassword
	}

	if mss.Spec.DatabaseUser == "" {
		mss.Spec.DatabaseUser = databaseUser
	}

	if mss.Spec.DatabaseHost == "" {
		mss.Spec.DatabaseHost = databaseHost
	}

	if mss.Spec.Port == 0 {
		mss.Spec.Port = port
	}

	if mss.Spec.LogLevel == "" {
		mss.Spec.LogLevel = logLevel
	}

	if mss.Spec.LogFormat == "" {
		mss.Spec.LogFormat = logFormat
	}

	if mss.Spec.LogLevel == "" {
		mss.Spec.LogLevel = accessControlAllowOrigin
	}

	if mss.Spec.AccessControlAllowOrigin == "" {
		mss.Spec.AccessControlAllowCredentials = accessControlAllowCredentials
	}

	/*
		CR Service Resource
		---------------------
	*/

	if mss.Spec.Size == 0 {
		mss.Spec.Size = size
	}

	// The clusterProtocol is required and used to generated the Public Host URL
	// Options [http or https]
	if mss.Spec.ClusterProtocol == "" {
		mss.Spec.ClusterProtocol = clusterProtocol
	}

	if mss.Spec.MemoryLimit == "" {
		mss.Spec.MemoryLimit = memoryLimit
	}

	if mss.Spec.MemoryRequest == "" {
		mss.Spec.MemoryRequest = memoryRequest
	}

	if mss.Spec.RouteName == "" {
		mss.Spec.RouteName = routeName
	}

	if mss.Spec.ConfigMapName == "" {
		mss.Spec.ConfigMapName = configMapName
	}

	/*
		Service Container
		---------------------
	*/

	if mss.Spec.Image == "" {
		mss.Spec.Image = image
	}

	if mss.Spec.ContainerName == "" {
		mss.Spec.ContainerName = containerName
	}

	/*
		OAuth Container
		---------------------
	*/

	if mss.Spec.OAuthImage == "" {
		mss.Spec.OAuthImage = oAuthImage
	}

	if mss.Spec.OAuthContainerName == "" {
		mss.Spec.OAuthContainerName = oAuthContainerName
	}
}
