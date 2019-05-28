package utils

import (
	"fmt"
	"os"
	"strings"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	// AppNamespaceEnvVar is the constant for env variable AppNamespaceEnvVar
	// which is the namespace where the APP CR can applied.
	// The namespaces should be informed split by ";".
	AppNamespaceEnvVar = "APP_NAMESPACES"
	// OperatorNamespaceForLocalEnv is valid and used just in the local env and for the tests.
	OperatorNamespaceForLocalEnv   = "mobile-security-service-proxy"
	ProxyServiceInstanceName       = "mobile-security-service-proxy"
	ApplicationServiceInstanceName = "mobile-security-service-application"
	InitEndpoint                   = "/init"
)

// The MobileSecurityServiceCRName has the name of the CR which should not be changed.
const MobileSecurityServiceCRName = "mobile-security-service"

var log = logf.Log.WithName("mobile-security-service-operator.utils")

//GetPublicServiceAPIURL returns the public service URL API
func GetPublicServiceAPIURL(route *routev1.Route, serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	return fmt.Sprintf("%v://%v%v", serviceInstance.Spec.ClusterProtocol, route.Status.Ingress[0].Host, InitEndpoint)
}

//GetRouteName returns an string name with the name of the router
func GetRouteName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	if len(m.Spec.RouteName) > 0 {
		return m.Spec.RouteName
	}
	return m.Name
}

// GetConfigMapName returns an string name with the name of the configMap
func GetConfigMapName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string {
	if len(m.Spec.ConfigMapName) > 0 {
		return m.Spec.ConfigMapName
	}
	return m.Name
}

// GetAppNamespaces returns the namespace the operator should be watching for changes
func GetAppNamespaces() (string, error) {
	ns, found := os.LookupEnv(AppNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", AppNamespaceEnvVar)
	}
	return ns, nil
}

// IsValidAppNamespace return true when the namespace informed is declared in the ENV VAR APP_NAMESPACES
// DEPRECATED
func IsValidAppNamespace(namespace string) (bool, error) {
	appNamespacesEnvVar, err := GetAppNamespaces()
	if err != nil {
		log.Error(err, "Unable to check if is app namespace %s is valid", namespace)

		// To skip when it is local env or a unit test
		_, err := k8sutil.GetOperatorNamespace()
		if err != nil {
			//Return true for the local env and for the unit tests
			if err == k8sutil.ErrNoNamespace { //
				log.Info("Allow to continue because it is outside of the cluster")
				return true, nil
			}
			return false, err
		}

		return false, err
	}
	for _, ns := range strings.Split(appNamespacesEnvVar, ";") {
		if ns == namespace {
			return true, nil
		}
	}
	err = fmt.Errorf("Invalid Namespace")
	return false, err
}

// IsValidOperatorNamespace return true when the namespace informed is declared in the ENV VAR APP_NAMESPACES
// DEPRECATED
func IsValidOperatorNamespace(namespace string) (bool, error) {
	//FIXME: this check is used to bypass validation of namespace.
	ns, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		//Return true for the local env and for the unit tests
		if err == k8sutil.ErrNoNamespace {
			log.Info("Allow to continue because it is outside of the cluster")
			return true, nil
		}
		return false, err
	}
	if ns == namespace {
		return true, nil
	}
	err = fmt.Errorf("Invalid Namespace")
	return false, err
}
