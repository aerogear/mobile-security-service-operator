package mobilesecurityserviceapp

import (
	"context"
	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service/pkg/models"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityServiceApp_create(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		kind       string
		serviceURL string
		request    reconcile.Request
		err        error
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
				kind: ConfigMap,
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
				kind: "WRONG_KIND",
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
	if err == nil {
		t.Error("Should fail since the instance has not a valid name")
	}

	if res.Requeue {
		t.Error("reconcile requeue request which is not expected")
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
	fetchBindAppRestServiceByAppID = func(serviceURL string, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
		app := models.App{ID: "1234", AppName: mssApp.Spec.AppName, AppID: mssApp.Spec.AppId}
		return &app, nil
	}

	service.DeleteAppFromServiceByRestAPI = func(serviceAPI string, id string, reqLogger logr.Logger) error {
		return nil
	}

	res, err := r.Reconcile(req)

	// should return without error and not requeue if App is deleted
	if err != nil {
		t.Fatalf("returned unexpectedly after attempting to delete app")
	}

	if res.Requeue {
		t.Error("reconcile requeue request which is not expected")
	}

}

func TestReconcileMobileSecurityServiceApp_Reconcile(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstance,
		&instance,
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
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("reconcile requeue request which is not expected")
	}

	// mock fetchBindAppRestServiceByAppID http call
	fetchBindAppRestServiceByAppID = func(serviceURL string, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
		app := models.App{AppName: mssApp.Spec.AppName, AppID: mssApp.Spec.AppId}
		return &app, nil
	}

	// mock CreateAppByRestAPI http call
	service.CreateAppByRestAPI = func(serviceAPI string, app *models.App, reqLogger logr.Logger) error {
		return nil
	}

	// mock UpdateAppNameByRestAPI http call
	service.UpdateAppNameByRestAPI = func(serviceAPI string, app *models.App, reqLogger logr.Logger) error {
		return nil
	}

	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: getSDKConfigMapName(&instance), Namespace: instance.Namespace}, configMap)
	if err != nil {
		t.Fatalf("get configMap: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err == nil {
		t.Error("expected an error when trying to update bind status")
	}

	// Check the result of reconciliation to make sure it has the desired state
	if res.Requeue {
		t.Error("reconcile Requeue unexpected")
	}

	// Check if has more than one configMap with the same Id
	configMapList := &corev1.ConfigMapList{}
	listOps := &client.ListOptions{}
	listOps.InNamespace(instance.Namespace)
	listOps.MatchingLabels(getLabelsToFetch(&instance))
	err = r.client.List(context.TODO(), listOps, configMapList)
	if err != nil {
		t.Fatalf("error to get a list of configMaps: (%v)", err)
	}

	if len(configMapList.Items) > 1 {
		t.Fatalf("more than one configmap was found which is not expected: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err == nil {
		t.Error("expected an error when trying to update bind status")
	}

	// Check the result of reconciliation to make sure it has the desired state
	if res.Requeue {
		t.Error("reconcile Requeue unexpected")
	}

}

func TestReconcileMobileSecurityServiceApp_Reconcile_AppDeleteCR(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstance,
		&instanceForDeletion,
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
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	// Check the result of reconciliation to make sure it has the desired state
	if res.Requeue {
		t.Error("reconcile Requeue unexpected")
	}
}

func TestReconcileMobileSecurityServiceApp_Reconcile_MSSDeleteCR(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstanceForDeletion,
		&instanceWithFinalizer,
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
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state
	if res.Requeue {
		t.Error("reconcile Requeue unexpected")
	}
}

func TestReconcileMobileSecurityServiceApp_Reconcile_UpdateName(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstance,
		&instance,
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

	if res.Requeue {
		t.Error("reconcile requeue request which is not expected")
	}

	// mock fetchBindAppRestServiceByAppID http call
	fetchBindAppRestServiceByAppID = func(serviceURL string, mssApp *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
		app := models.App{AppName: "oldName", AppID: mssApp.Spec.AppId}
		return &app, nil
	}

	// mock CreateAppByRestAPI http call
	service.CreateAppByRestAPI = func(serviceAPI string, app *models.App, reqLogger logr.Logger) error {
		return nil
	}

	// mock UpdateAppNameByRestAPI http call
	service.UpdateAppNameByRestAPI = func(serviceAPI string, app *models.App, reqLogger logr.Logger) error {
		return nil
	}

	res, err = r.Reconcile(req)

	if res.Requeue {
		t.Error("reconcile requeue request which is not expected")
	}

	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: getSDKConfigMapName(&instance), Namespace: instance.Namespace}, configMap)
	if err != nil {
		t.Fatalf("get configMap: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err == nil {
		t.Error("expected an error when trying to update bind status")
	}

	// Check the result of reconciliation to make sure it has the desired state
	if res.Requeue {
		t.Error("reconcile Requeue unexpected")
	}

	// Check if has more than one configMap with the same Id
	configMapList := &corev1.ConfigMapList{}
	listOps := &client.ListOptions{}
	listOps.InNamespace(instance.Namespace)
	listOps.MatchingLabels(getLabelsToFetch(&instance))
	err = r.client.List(context.TODO(), listOps, configMapList)
	if err != nil {
		t.Fatalf("error to get a list of configMaps: (%v)", err)
	}

	if len(configMapList.Items) > 1 {
		t.Fatalf("more than one configmap was found which is not expected: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err == nil {
		t.Error("expected an error when trying to update bind status")
	}

	// Check the result of reconciliation to make sure it has the desired state
	if res.Requeue {
		t.Error("reconcile Requeue unexpected")
	}
}
