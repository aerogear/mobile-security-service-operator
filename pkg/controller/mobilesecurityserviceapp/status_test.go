package mobilesecurityserviceapp

import (
	"github.com/go-logr/logr"
	"reflect"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service/pkg/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		request   reconcile.Request
		configMap *corev1.ConfigMap
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
				configMap: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-app-security",
						Namespace: instance.Namespace,
						Labels:    getSDKAppLabels(&instance),
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

			if err := r.updateBindStatus("http://hostest", reqLogger, tt.args.configMap, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconcileMobileSecurityServiceApp_updateConfigMapStatus(t *testing.T) {

	type fields struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		scheme   *runtime.Scheme
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         *corev1.ConfigMap
		wantErr      bool
		shouldCreate bool
	}{
		{
			name: "should fail to find the ConfigMap",
			fields: fields{
				instance: &instance,
				scheme:   scheme.Scheme,
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should update the ConfigMap status",
			fields: fields{
				instance: &instance,
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
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app-security",
					Namespace: instance.Namespace,
					Labels:    getSDKAppLabels(&instance),
				},
			},
			shouldCreate: true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objs := []runtime.Object{tt.fields.instance}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			if tt.shouldCreate {
				r.create(&instance, ConfigMap, "http://testhost", reqLogger, tt.args.request)
			}

			got, err := r.updateConfigMapStatus(reqLogger, tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceApp.updateConfigMapStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil && tt.want != nil) && !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("ReconcileMobileSecurityServiceApp.updateConfigMapStatus() = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}
