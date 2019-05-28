package mobilesecurityservicedb

import (
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

//buildReconcileWithFakeClientWithMocks return reconcile with fake client, schemes and mock objects
func buildReconcileWithFakeClientWithMocks(objs []runtime.Object, t *testing.T) *ReconcileMobileSecurityServiceDB {
	s := scheme.Scheme

	s.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, &mobilesecurityservicev1alpha1.MobileSecurityServiceDB{})
	s.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, &mobilesecurityservicev1alpha1.MobileSecurityService{})

	// create a fake client to mock API calls with the mock objects
	cl := fake.NewFakeClientWithScheme(s, objs...)

	// create a ReconcileMobileSecurityService object with the scheme and fake client
	return &ReconcileMobileSecurityServiceDB{client: cl, scheme: s}
}
