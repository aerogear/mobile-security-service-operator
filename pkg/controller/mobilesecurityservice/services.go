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
func (r *ReconcileMobileSecurityService) buildProxyService(service *mobilesecurityservicev1alpha1.MobileSecurityService) *corev1.Service {
	ls := getAppLabels(service.Name)
	targetPort := intstr.FromInt(int(service.Spec.OAuthPort))
	ser := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.ProxyServiceInstanceName,
			Namespace: service.Namespace,
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
	// Set MobileSecurityService service as the owner and controller
	controllerutil.SetControllerReference(service, ser, r.scheme)
	return ser
}

func (r *ReconcileMobileSecurityService) buildApplicationService(service *mobilesecurityservicev1alpha1.MobileSecurityService) *corev1.Service {
	ls := getAppLabels(service.Name)
	ser := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.ApplicationServiceInstanceName,
			Namespace: service.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: service.Spec.Port,
					Name: "server",
				},
			},
		},
	}
	// Set MobileSecurityService service as the owner and controller
	controllerutil.SetControllerReference(service, ser, r.scheme)
	return ser
}
