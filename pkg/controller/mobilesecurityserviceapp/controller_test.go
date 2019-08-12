package mobilesecurityserviceapp

import (
	"testing"

	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"github.com/go-logr/logr"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service/pkg/models"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

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
	if err == nil {
		t.Error("expected an error when trying to update bind status")
	}

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
}
