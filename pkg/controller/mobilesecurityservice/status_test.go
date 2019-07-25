package mobilesecurityservice

import (
	"reflect"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityService_updateStatus(t *testing.T) {
	type fields struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
		objs     []runtime.Object
		scheme   *runtime.Scheme
	}
	type args struct {
		request            reconcile.Request
		configMap          *corev1.ConfigMap
		deployment         *appsv1.Deployment
		proxyService       *corev1.Service
		applicationService *corev1.Service
		route              *routev1.Route
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should return error as the resources have not been created",
			fields: fields{
				instance: &mssInstance,
				objs:     []runtime.Object{&mssInstance},
				scheme:   scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      mssInstance.Name,
						Namespace: mssInstance.Namespace,
					},
				},
				configMap:          &configMap,
				deployment:         &appsv1.Deployment{},
				proxyService:       &corev1.Service{},
				applicationService: &corev1.Service{},
				route:              &routev1.Route{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := buildReconcileWithFakeClientWithMocks(tt.fields.objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.request.Namespace, "Request.Name", tt.args.request.Name)

			if err := r.updateStatus(reqLogger, tt.args.configMap, tt.args.deployment, tt.args.proxyService, tt.args.applicationService, tt.args.route, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.updateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconcileMobileSecurityService_updateConfigMapStatus(t *testing.T) {

	type fields struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
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
				instance: &mssInstance,
				scheme:   scheme.Scheme,
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should update the ConfigMap status",
			fields: fields{
				instance: &mssInstance,
				scheme:   scheme.Scheme,
			},
			args: args{
				request: reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      mssInstance.Name,
						Namespace: mssInstance.Namespace,
					},
				},
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      mssInstance.Spec.ConfigMapName,
					Namespace: mssInstance.Namespace,
					Labels:    getAppLabels(mssInstance.Name),
				},
				Data: getAppEnvVarsMap(&mssInstance),
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
				r.create(&mssInstance, reqLogger, ConfigMap)
			}

			got, err := r.updateConfigMapStatus(reqLogger, tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.updateConfigMapStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil && tt.want != nil) && !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("ReconcileMobileSecurityService.updateConfigMapStatus() = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}

func TestReconcileMobileSecurityService_updateBindStatusWithInvalidNamespace(t *testing.T) {
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should return without an error when updating status",
			args: args{
				instance: &mssInstance,
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
					Name:      mssInstance.Name,
					Namespace: mssInstance.Namespace,
				},
			}

			if err := r.updateStatusWithInvalidNamespace(reqLogger, req); (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceApp.updateBindStatusWithInvalidNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
