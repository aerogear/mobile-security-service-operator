package mobilesecurityservice

import (
	"context"
	"reflect"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	instance = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service-app",
			Namespace: "mobile-security-service-operator",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:                    1,
			MemoryLimit:             "512Mi",
			MemoryRequest:           "512Mi",
			ClusterProtocol:         "http",
			ConfigMapName:           "mss-config",
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
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    getAppLabels(instance.Name),
		},
	}
)

func TestReconcileMobileSecurityService_update(t *testing.T) {
	type fields struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
		scheme   *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		want    reconcile.Result
		wantErr bool
	}{
		{
			name: "should requeue",
			fields: fields{
				instance: &instance,
				scheme:   scheme.Scheme,
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.instance}

			r, _ := getReconciler(objs)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.fields.instance.Name,
					Namespace: tt.fields.instance.Namespace,
				},
			}

			reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

			got, err := r.update(objs[0], reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileMobileSecurityService.update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityService_create(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
		kind     string
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
			name: "should create and return a new deployment",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &instance,
				kind:     DEEPLOYMENT,
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
		{
			name: "should fail to create an unknown kind",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &instance,
				kind:     "OBJECT",
			},
			want:      reconcile.Result{},
			wantErr:   true,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r, _ := getReconciler(objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityService.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, err := r.create(tt.args.instance, reqLogger, tt.args.kind)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReconcileMobileSecurityService.create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityService_buildFactory(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance *mobilesecurityservicev1alpha1.MobileSecurityService
		kind     string
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
			name: "should create a Deployment",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&v1beta1.Deployment{}),
			args: args{
				instance: &instance,
				kind:     DEEPLOYMENT,
			},
		},
		{
			name: "should create a ConfigMap",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&corev1.ConfigMap{}),
			args: args{
				instance: &instance,
				kind:     CONFIGMAP,
			},
		},
		{
			name: "should create a Service",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&corev1.Service{}),
			args: args{
				instance: &instance,
				kind:     SERVICE,
			},
		},
		{
			name: "should create a Route",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&routev1.Route{}),
			args: args{
				instance: &instance,
				kind:     ROUTE,
			},
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &instance,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r, _ := getReconciler(objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityService.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got, err := r.buildFactory(reqLogger, tt.args.instance, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.buildFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("ReconcileMobileSecurityService.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityService_Reconcile(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&instance,
		&route,
	}

	r, cl := getReconciler(objs)

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

	configMap := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), req.NamespacedName, configMap)

	// Check the result of reconciliation to make sure it has the desired state
	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state
	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	// check if the deployment has been created
	dep := &v1beta1.Deployment{}
	err = cl.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	// check if the service has been created
	service := &corev1.Service{}
	err = cl.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", service)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	// FIXME: There appears to be an issue with the Operator-SDK blocking this.
	// ERROR: get route: (no kind "Route" is registered for version "v1" in scheme "k8s.io/client-go/kubernetes/scheme/register.go:61")
	// See: https://github.com/operator-framework/operator-sdk/issues/1421
	// route := &routev1.Route{}
	err = cl.Get(context.TODO(), types.NamespacedName{Namespace: instance.Namespace, Name: utils.GetRouteName(&instance)}, &route)
	if err != nil {
		t.Fatalf("get route: (%v)", err)
	}
}

func TestReconcileMobileSecurityService_Reconcile_InvalidInstance(t *testing.T) {
	invalidInstance := &instance
	invalidInstance.Spec.ClusterProtocol = "ws"

	// objects to track in the fake client
	objs := []runtime.Object{
		&instance,
	}

	r, _ := getReconciler(objs)

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

	if !res.Requeue {
		t.Fatal("reconcile did not requeue request as expected")
	}
}

func getReconciler(objs []runtime.Object) (*ReconcileMobileSecurityService, client.Client) {
	s := scheme.Scheme

	routev1.Install(s)

	s.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, &mobilesecurityservicev1alpha1.MobileSecurityService{})
	s.AddKnownTypes(routev1.SchemeGroupVersion, &routev1.Route{})
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ConfigMap{}, &corev1.Service{})

	// create a fake client to mock API calls
	cl := fake.NewFakeClientWithScheme(s, objs...)
	// create a ReconcileMobileSecurityService object with the scheme and fake client
	return &ReconcileMobileSecurityService{client: cl, scheme: s}, cl
}
