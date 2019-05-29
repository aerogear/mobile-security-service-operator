package mobilesecurityserviceapp

import (
	"github.com/go-logr/logr"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service/pkg/models"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityServiceApp_addFinalizer(t *testing.T) {
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
			name: "Should add the finalizer",
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

			if err := r.addFinalizer(reqLogger, tt.fields.instance, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconcileMobileSecurityServiceApp_removeFinalizer(t *testing.T) {
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
			name: "Should return success when the finalizer is called and the APP CR was not deleted",
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
		{
			name: "Should return success when the finalizer is called and the APP CR was deleted",
			fields: fields{
				instance: &instanceForDeletion,
				objs:     []runtime.Object{&instanceForDeletion},
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
			fetchBindAppRestServiceByAppID = func(serviceURL string, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
				app := models.App{ID: "1234", AppName: instance.Spec.AppName, AppID: instance.Spec.AppId}
				return &app, nil
			}

			if err := r.removeFinalizer("http://mobile-security-service-application:1234/api", reqLogger, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
