package mobilesecurityserviceapp

import (
	"testing"

	"github.com/go-logr/logr"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobile-security-service/v1alpha1"
	"github.com/aerogear/mobile-security-service/pkg/models"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityServiceApp_updateBindStatusWithInvalidNamespace(t *testing.T) {
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return without an error when updating status",
			args: args{
				instance: &instance,
			},
			wantErr: false,
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

func TestReconcileMobileSecurityServiceApp_updateStatus(t *testing.T) {
	type fields struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		objs     []runtime.Object
		scheme   *runtime.Scheme
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should work when the resources are created",
			fields: fields{
				instance: &instance,
				objs:     []runtime.Object{&instance},
				scheme:   scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      instance.Name,
						Namespace: instance.Namespace,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			// mock fetchBindAppRestServiceByAppID http call
			fetchBindAppRestServiceByAppID = func(serviceURL string, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
				app := models.App{ID: "1234", AppName: mssApp.Spec.AppName, AppID: mssApp.Spec.AppId}
				return &app, nil
			}

			if err := r.updateBindStatus("http://hostest", reqLogger, tt.fields.instance, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
