package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//buildService returns the service resource
func (r *ReconcileMobileSecurityService) buildProxyService(m *mobilesecurityservicev1alpha1.MobileSecurityService) *corev1.Service {
	ls := getAppLabels(m.Name)
	targetPort := intstr.FromInt(int(m.Spec.OAuthPort))
	ser := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.PROXY_SERVICE_INSTANCE_NAME,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					TargetPort: targetPort,
					Port:       80,
					Name:       "web",
				},
			},
		},
	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, ser, r.scheme)
	return ser
}

func (r *ReconcileMobileSecurityService) buildApplicationService(m *mobilesecurityservicev1alpha1.MobileSecurityService) *corev1.Service {
	ls := getAppLabels(m.Name)
	ser := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.APPLICATION_SERVICE_INSTANCE_NAME,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: m.Spec.Port,
					Name: "server",
				},
			},
		},
	}
	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, ser, r.scheme)
	return ser
}
