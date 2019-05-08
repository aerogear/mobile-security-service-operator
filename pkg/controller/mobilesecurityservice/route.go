package mobilesecurityservice

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	routev1 "github.com/openshift/api/route/v1"
)

//buildRoute returns the route resource
func (r *ReconcileMobileSecurityService) buildRoute(m *mobilesecurityservicev1alpha1.MobileSecurityService) *routev1.Route {
	ls := getAppLabels(m.Name)


	route := &routev1.Route{
		TypeMeta: v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Route",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      utils.GetRouteName(m),
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: m.Name ,
			},
		},
	}

	if len(m.Spec.RoutePath) > 0 {
		route.Spec.Path = m.Spec.RoutePath
	}

	// Set MobileSecurityService instance as the owner and controller
	controllerutil.SetControllerReference(m, route, r.scheme)
	return route
}