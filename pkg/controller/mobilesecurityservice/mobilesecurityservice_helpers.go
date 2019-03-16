package mobilesecurityservice

import (
	corev1 "k8s.io/api/core/v1"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
)

// labelsForMobileSecurityService returns the labels for selecting the resources
// belonging to the given MobileSecurityService CR name.
func labelsForMobileSecurityService(name string) map[string]string {
	return map[string]string{"app": "mobilesecurityservice", "mobilesecurityservice_cr": name}
}

func getAnnotationsForMobileSecurityIngress(name string) map[string]string {
	return map[string]string{"kubernetes.io/ingress.class": "nginx", "mobilesecurityservice_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// return the list of environment variables of the Mobile Security Service Project
func getConfigMapForMobileSecurityService() map[string]string {
	return map[string]string{
		"PGHOST": "postgresql",
		"PGUSER": "postgresql",
		"PGPASSWORD": "postgres",
		"PGDATABASE": "mobile_security_service",
		"PORT": "3000",
		"LOG_LEVEL":                        "info",
		"LOG_FORMAT":                       "json",
		"ACCESS_CONTROL_ALLOW_ORIGIN":      "*",
		"ACCESS_CONTROL_ALLOW_CREDENTIALS": "false",
		"STATIC_FILES_DIR":                 "static",
		"PGPORT":                           "5432",
		"PGSSLMODE":                        "disable",
		"PGCONNECT_TIMEOUT":                "5",
		"PGAPPNAME":                        "",
		"PGSSLCERT":                        "",
		"PGSSLKEY":                         "",
		"PGSSLROOTCERT":                    "",
		"DB_MAX_CONNECTIONS":               "100",
	}
}

// return the Env Var for the project
func getAllEnvVarsToSetupMobileSecurityService(m *mobilesecurityservicev1alpha1.MobileSecurityService) *[]corev1.EnvVar {
	res := []corev1.EnvVar{}
	for key := range getConfigMapForMobileSecurityService() {
		env := corev1.EnvVar{
			Name: key,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.Name,
					},
					Key: key,
				},
			},
		}
		res = append(res, env)
	}
	return &res
}

