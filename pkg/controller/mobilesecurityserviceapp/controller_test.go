package mobilesecurityserviceapp

import (
	errs "errors"
	"reflect"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	routev1 "github.com/openshift/api/route/v1"

	corev1 "k8s.io/api/core/v1"
	// "k8s.io/api/extensions/v1beta1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/aerogear/mobile-security-service-operator/pkg/models"

	"github.com/aerogear/mobile-security-service-operator/pkg/service"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	instance = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName:                 "test-app",
			AppId:                   "test-app-id",
			SkipNamespaceValidation: true,
		},
	}

	instanceTwo = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service-2",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "test-app",
			AppId:   "test-app-id",
		},
	}

	mssInstance = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service-operator",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:                    1,
			MemoryLimit:             "512Mi",
			MemoryRequest:           "512Mi",
			ClusterProtocol:         "http",
			ConfigMapName:           "mss-config",
			Port:                    1234,
			RouteName:               "mss-route",
			SkipNamespaceValidation: true,
		},
	}

	route = routev1.Route{
		TypeMeta: v1.TypeMeta{
			APIVersion: "route.openshift.io/v1",
			Kind:       "Route",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      utils.GetRouteName(&mssInstance),
			Namespace: mssInstance.Namespace,
			Labels:    getAppLabels(mssInstance.Name),
		},
	}
)

const envVar = "envVar"

func TestReconcileMobileSecurityServiceApp_create(t *testing.T) {
	type fields struct {
		client client.Client
		scheme *runtime.Scheme
	}
	type args struct {
		instance   *mobilesecurityservicev1alpha1.MobileSecurityServiceApp
		kind       string
		serviceURL string
		// reqLogger  logr.Logger
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
				instance: &instance,
				kind:     CONFIGMAP,
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
				instance: &instance,
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

			objs := []runtime.Object{tt.args.instance}

			tt.fields.scheme.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, tt.args.instance)

			cl := fake.NewFakeClient(objs...)

			r := &ReconcileMobileSecurityServiceApp{
				client: cl,
				scheme: tt.fields.scheme,
			}

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityService.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, err := r.create(tt.args.instance, tt.args.kind, tt.args.serviceURL, reqLogger, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceApp.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileMobileSecurityServiceApp.create() = %v, want %v", got, tt.want)
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
				instance:   &instance,
				kind:       CONFIGMAP,
				serviceURL: "service-url",
			},
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:   &instance,
				kind:       "UNDEFINED",
				serviceURL: "service-url",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			tt.fields.scheme.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, tt.args.instance)

			cl := fake.NewFakeClient(objs...)

			r := &ReconcileMobileSecurityServiceApp{
				client: cl,
				scheme: tt.fields.scheme,
			}

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceApp.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, err := r.buildFactory(reqLogger, tt.args.instance, tt.args.kind, tt.args.serviceURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceApp.buildFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("ReconcileMobileSecurityServiceApp.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityServiceApp_Reconcile(t *testing.T) {
	// objects to track in the fake client
	objs := []runtime.Object{
		&instance,
		&instanceTwo,
		&mssInstance,
		&route,
	}

	r := getReconciler(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	// Initial call to reconciler - should fail on config map, create it, and requeue request
	res, err := r.Reconcile(req)
	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	// check configmap created
	reqLogger := log.WithValues("Request.Namespace", &req.NamespacedName.Namespace, "Request.Name", &req.NamespacedName.Name)
	_, err = r.fetchSDKConfigMap(reqLogger, &instance)
	// configmap was not created properly
	if err != nil {
		// configmap was not created properly
		t.Fatalf("get configmap: (%v)", err)
	}

	// mock `fetchBindAppRestServiceByAppID` restful call
	fetchBindAppRestServiceByAppID = func(serviceURL string, instance *mobilesecurityservicev1alpha1.MobileSecurityServiceApp, reqLogger logr.Logger) (*models.App, error) {
		app := models.App{AppName: "test-app", AppID: "test-app-id"}
		return &app, nil
	}

	// mock CreateAppByRestAPI
	service.CreateAppByRestAPI = func(serviceAPI string, app models.App, reqLogger logr.Logger) error {
		return nil
	}

	// call reconcile again - configmap should be created
	res, err = r.Reconcile(req)
	if res.Requeue {
		t.Error("reconcile did not requeue request as expected after creating app")
	}

	// call reconcile again after creating app
	res, err = r.Reconcile(req)

	SDKConfigMapStatus, err := r.updateSDKConfigMapStatus(reqLogger, req)
	// configmap was not updated properly
	if err != nil {
		// configmap was not updated properly
		t.Fatalf("update configmap: (%v)", err)
	}

	// mock updateBindStatus
	// updateBindStatus = func() error {
	// 	return nil
	// }

	// 
	SDKConfigMapStatus.UID = "1234"
	err = r.updateBindStatus("http://mobile-security-service-application:1234/api", reqLogger, SDKConfigMapStatus, req)
	if err != nil {
		// bindStatus was not updated properly
		t.Fatalf("update bind status: (%v)", err)
	}

	if res.Requeue {
		t.Error("reconcile requeued instead of returning successfully")
	} 
}

func getReconciler(objs []runtime.Object) *ReconcileMobileSecurityServiceApp {
	s := scheme.Scheme

	routev1.AddToScheme(s)

	s.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, &mobilesecurityservicev1alpha1.MobileSecurityServiceApp{})

	s.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, &mobilesecurityservicev1alpha1.MobileSecurityService{})
	s.AddKnownTypes(routev1.SchemeGroupVersion, &routev1.Route{})
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ConfigMap{}, &corev1.Service{})

	// create a fake client to mock API calls
	cl := fake.NewFakeClient(objs...)
	// create a ReconcileMobileSecurityServiceApp object with the scheme and fake client
	return &ReconcileMobileSecurityServiceApp{client: cl, scheme: s}
}
