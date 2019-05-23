package mobilesecurityserviceapp

import (
	errs "errors"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/aerogear/mobile-security-service-operator/pkg/models"

	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityServiceApp_create(t *testing.T) {
	type fields struct {
		client client.Client
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		kind       string
		serviceURL string
		request reconcile.Request
		err     error
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reconcile.Result
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "should return a configmap",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				kind:     ConfigMap,
				err:      errors.NewInternalError(errs.New("Internal Server Error")),
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
		{
			name: "should return an error when type other than CONFIGMAP specified",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				kind:     "WRONG_KIND",
				err:      errors.NewInternalError(errs.New("Internal Server Error")),
			},
			want:      reconcile.Result{},
			wantErr:   true,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{
				&instance,
				&mssInstance,
				&route,
			}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", &instance.Namespace, "Request.Name", &instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityService.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(&instance, tt.args.kind, tt.args.serviceURL, reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceApp.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileMobileSecurityServiceApp_buildFactory(t *testing.T) {
	type fields struct {
		client client.Client
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		kind       string
		serviceURL string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reflect.Type
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "should create a ConfigMap",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&corev1.ConfigMap{}),
			args: args{
				kind:       ConfigMap,
				serviceURL: "service-url",
			},
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				kind:       "UNDEFINED",
				serviceURL: "service-url",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{
				&instance,
			}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", &instance.Namespace, "Request.Name", &instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceApp.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			configMap := r.buildFactory(reqLogger, &instance, tt.args.kind, tt.args.serviceURL)
			if configMap == nil {
				t.Errorf("ReconcileMobileSecurityServiceApp.buildFactory() - received no config type")
				return
			}

			gotType := reflect.TypeOf(configMap)
			expectedType := reflect.TypeOf(&corev1.ConfigMap{})
			if gotType != expectedType {
				t.Errorf("ReconcileMobileSecurityServiceApp.buildFactory() - received wrong config type")
				return
			}
		})
	}
}

func TestReconcileMobileSecurityServiceApp_Reconcile_FetchInstance(t *testing.T) {

	objs := []runtime.Object{
		&instanceInvalidName,
		&mssInstance,
		&route,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	// should return without error or requeueing request if App is not found
	if res.Requeue == true || err != nil {
		t.Fatalf("returned unexpectedly after attempting to fetch instance")
	}
}

func TestReconcileMobileSecurityServiceApp_Reconcile_Deletion(t *testing.T) {
	// required to test validity of namespace as namespace check won't be hit after
	// fetch instance call beforehand in reconcile loop
	objs := []runtime.Object{
		&instanceForDeletion,
		&mssInstance,
		&route,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	// mock fetchBindAppRestServiceByAppID http call
	fetchBindAppRestServiceByAppID = func(serviceURL string, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
		app := models.App{ID: "1234", AppName: "test", AppID: "test-app-id"}
		return &app, nil
	}

	service.DeleteAppFromServiceByRestAPI = func(serviceAPI string, id string, reqLogger logr.Logger) error {
		return nil
	}

	res, err := r.Reconcile(req)
	// should return without error or requeueing request if App is deleted
	if res.Requeue == true || err != nil {
		t.Fatalf("returned unexpectedly after attempting to delete app")
	}


}

func TestReconcileMobileSecurityServiceApp_Reconcile(t *testing.T) {
	// objects to track in the fake client
	objs := []runtime.Object{
		&instance,
		&mssInstance,
		&route,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	// mock fetchBindAppRestServiceByAppID http call
	fetchBindAppRestServiceByAppID = func(serviceURL string, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
		app := models.App{AppName: "test-app-22222", AppID: "test-app-id"}
		return &app, nil
	}

	// mock CreateAppByRestAPI http call
	service.CreateAppByRestAPI = func(serviceAPI string, app models.App, reqLogger logr.Logger) error {
		return nil
	}

	// mock UpdateAppNameByRestAPI http call
	service.UpdateAppNameByRestAPI = func(serviceAPI string, app *models.App, reqLogger logr.Logger) error {
		return nil
	}

	// Initial call to reconciler
	_, err := r.Reconcile(req)
	if err == nil {
		t.Error("expected an error when trying to update bind status")
	}

	// check configmap create 
	reqLogger := log.WithValues("Request.Namespace", &req.NamespacedName.Namespace, "Request.Name", &req.NamespacedName.Name)
	_, err = r.fetchSDKConfigMap(reqLogger, &instance)
	// configmap was not created properly
	if err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}

	SDKConfigMapStatus, err := r.updateSDKConfigMapStatus(reqLogger, req)
	// configmap was not updated properly
	if err != nil {
		t.Fatalf("update configmap: (%v)", err)
	}

	// mock UID so bind status can be updated successfully
	SDKConfigMapStatus.UID = "1234"

	err = r.updateBindStatus("http://mobile-security-service-application:1234/api", reqLogger, SDKConfigMapStatus, req)
	if err != nil {
		// bindStatus was not updated properly
		t.Fatalf("update bind status: (%v)", err)
	}

}

