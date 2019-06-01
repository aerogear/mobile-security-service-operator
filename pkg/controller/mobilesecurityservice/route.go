package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//buildRoute returns the route resource
func (r *ReconcileMobileSecurityService) buildRoute(service *mobilesecurityservicev1alpha1.MobileSecurityService) *routev1.Route {

	ls := getAppLabels(service.Name)
	route := &routev1.Route{
		ObjectMeta: v1.ObjectMeta{
			Name:      service.Spec.RouteName,
			Namespace: service.Namespace,
			Labels:    ls,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: utils.ProxyServiceInstanceName,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("web"),
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationEdge,
			},
		},
	}

	// Set MobileSecurityService service as the owner and controller
	controllerutil.SetControllerReference(service, route, r.scheme)
	return route
}
