package mobilesecurityserviceapp

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"

)

func TestReconcileMobileSecurityServiceApp_updateBindStatusWithInvalidNamespace(t *testing.T) {
	type args struct {
		instance        *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return without an error when updating status",
			args: args{
				instance: 	&instance,
			},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{
				tt.args.instance,
			}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// mock request to simulate Reconcile() being called on an event for a watched resource
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      instance.Name,
					Namespace: instance.Namespace,
				},
			}

			if err := r.updateBindStatusWithInvalidNamespace(reqLogger, req); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceApp.updateBindStatusWithInvalidNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
