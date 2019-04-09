package mobilesecurityservice

import (
	"bytes"
	"encoding/json"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/models"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"unsafe"
)


// Returns an string map with the labels which wil be associated to the kubernetes/openshift objects
// which will be created and managed by this operator
func getAppLabels(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservice_cr": name}
}

// Returns an string name with the name of the configMap
func getConfigMapName(m *mobilesecurityservicev1alpha1.MobileSecurityService) string{
	if len(m.Spec.ConfigMapName) > 0 {
		return m.Spec.ConfigMapName
	}
	return m.Name
}

func getAppLabelsForSDKConfigMap(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservice_cr": name, "name": name+"-sdk-config"}
}

//To transform the object into a string with its json
func getSdkConfigStringJsonFormat(sdk *models.SDKConfig) string{
	jsonSdk, _ := json.Marshal(sdk)
	res:= strings.NewReader(string(jsonSdk))
	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	b := buf.Bytes()
	return *(*string)(unsafe.Pointer(&b))
}

// Helper to build the env vars which will be configured in the deployment of the Mobile Security Service Project
func buildAppEnvVars(m *mobilesecurityservicev1alpha1.MobileSecurityService) *[]corev1.EnvVar {
	res := []corev1.EnvVar{}
	for key := range getAppEnvVarsMap(m) {
		env := corev1.EnvVar{
			Name: key,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: getConfigMapName(m),
					},
					Key: key,
				},
			},
		}
		res = append(res, env)
	}
	return &res
}

// Helper to get a map[string]string with the key and values required/used to setup the Mobile Security Service Project
func getAppEnvVarsMap(m *mobilesecurityservicev1alpha1.MobileSecurityService) map[string]string {
	return map[string]string{
		"PGHOST":                           m.Spec.DatabaseHost,
		"LOG_LEVEL":                        m.Spec.LogLevel,
		"LOG_FORMAT":                       m.Spec.LogFormat,
		"ACCESS_CONTROL_ALLOW_ORIGIN":      m.Spec.AccessControlAllowOrigin,
		"ACCESS_CONTROL_ALLOW_CREDENTIALS": m.Spec.AccessControlAllowCredentials,
		"PGDATABASE":                       m.Spec.DatabaseName,
		"PGPASSWORD":                       m.Spec.DatabasePassword,
		"PGUSER":                           m.Spec.DatabaseUser,
	}
}
