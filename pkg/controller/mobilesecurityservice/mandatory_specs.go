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
	oAuthPort                     = 4180
	oAuthImage                    = "docker.io/openshift/oauth-proxy:v1.1.0"
	oAuthContainerName            = "oauth-proxy"
	configMapName                 = "mobile-security-service-config"
	routeName                     = "route"
)

// addMandatorySpecsDefinitions will add the specs which are mandatory for Mobile Security Service CR in the case them
// not be applied
func addMandatorySpecsDefinitions(serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) {

	/*
		Environment Variables
		---------------------
		The following values are used to create the ConfigMap and the Environment Variables which will use these values
		These values are used for both the Mobile Security Service and its Database
	*/

	if serviceInstance.Spec.DatabaseName == "" {
		serviceInstance.Spec.DatabaseName = databaseName
	}

	if serviceInstance.Spec.DatabasePassword == "" {
		serviceInstance.Spec.DatabasePassword = databasePassword
	}

	if serviceInstance.Spec.DatabaseUser == "" {
		serviceInstance.Spec.DatabaseUser = databaseUser
	}

	if serviceInstance.Spec.DatabaseHost == "" {
		serviceInstance.Spec.DatabaseHost = databaseHost
	}

	if serviceInstance.Spec.Port == 0 {
		serviceInstance.Spec.Port = port
	}

	if serviceInstance.Spec.LogLevel == "" {
		serviceInstance.Spec.LogLevel = logLevel
	}

	if serviceInstance.Spec.LogFormat == "" {
		serviceInstance.Spec.LogFormat = logFormat
	}

	if serviceInstance.Spec.LogLevel == "" {
		serviceInstance.Spec.LogLevel = accessControlAllowOrigin
	}

	if serviceInstance.Spec.AccessControlAllowOrigin == "" {
		serviceInstance.Spec.AccessControlAllowCredentials = accessControlAllowCredentials
	}

	/*
		CR Service Resource
		---------------------
	*/

	if serviceInstance.Spec.Size == 0 {
		serviceInstance.Spec.Size = size
	}

	// The clusterProtocol is required and used to generated the Public Host URL
	// Options [http or https]
	if serviceInstance.Spec.ClusterProtocol == "" {
		serviceInstance.Spec.ClusterProtocol = clusterProtocol
	}

	if serviceInstance.Spec.MemoryLimit == "" {
		serviceInstance.Spec.MemoryLimit = memoryLimit
	}

	if serviceInstance.Spec.MemoryRequest == "" {
		serviceInstance.Spec.MemoryRequest = memoryRequest
	}

	if serviceInstance.Spec.RouteName == "" {
		serviceInstance.Spec.RouteName = routeName
	}

	if serviceInstance.Spec.ConfigMapName == "" {
		serviceInstance.Spec.ConfigMapName = configMapName
	}

	/*
		Service Container
		---------------------
	*/

	if serviceInstance.Spec.Image == "" {
		serviceInstance.Spec.Image = image
	}

	if serviceInstance.Spec.ContainerName == "" {
		serviceInstance.Spec.ContainerName = containerName
	}

	/*
		OAuth Container
		---------------------
	*/

	if serviceInstance.Spec.OAuthPort == 0 {
		serviceInstance.Spec.OAuthPort = oAuthPort
	}

	if serviceInstance.Spec.OAuthImage == "" {
		serviceInstance.Spec.OAuthImage = oAuthImage
	}

	if serviceInstance.Spec.OAuthContainerName == "" {
		serviceInstance.Spec.OAuthContainerName = oAuthContainerName
	}
}
