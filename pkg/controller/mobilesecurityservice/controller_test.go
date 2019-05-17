package mobilesecurityservice

import (
	"context"
	"reflect"
	"testing"

	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	instanceOne = mobilesecurityservicev1alpha1.MobileSecurityService{
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
			RouteName:               "route",
			SkipNamespaceValidation: true,
		},
	}

	instanceTwo = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service-app",
			Namespace: "mobile-security-service-operator-2",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:                    1,
			MemoryLimit:             "512Mi",
			MemoryRequest:           "512Mi",
			ClusterProtocol:         "http",
			ConfigMapName:           "mss-config",
			RouteName:               "route",
			SkipNamespaceValidation: true,
		},
	}

	route = routev1.Route{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "route.openshift.io/v1",
			Kind:       "Route",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetRouteName(&instanceOne),
			Namespace: instanceOne.Namespace,
			Labels:    getAppLabels(instanceOne.Name),
		},
	}
)

func TestReconcileMobileSecurityService_update(t *testing.T) {
	type fields struct {
		createdInstance  *mobilesecurityservicev1alpha1.MobileSecurityService
		instanceToUpdate *mobilesecurityservicev1alpha1.MobileSecurityService
		scheme           *runtime.Scheme
		namespace        string
	}
	tests := []struct {
		name    string
		fields  fields
		want    reconcile.Result
		wantErr bool
	}{
		{
			name: "should successfully update an instance",
			fields: fields{
				createdInstance:  &instanceOne,
				instanceToUpdate: &instanceOne,
				scheme:           scheme.Scheme,
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
		{
			name: "should give error when namespace not found",
			fields: fields{
				createdInstance:  &instanceOne,
				instanceToUpdate: &instanceTwo,
				scheme:           scheme.Scheme,
			},
			want:    reconcile.Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := getReconciler(objs)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.fields.createdInstance.Name,
					Namespace: tt.fields.createdInstance.Namespace,
				},
			}

			reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

			got, err := r.update(tt.fields.instanceToUpdate, reqLogger)
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
				instance: &instanceOne,
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
				instance: &instanceOne,
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

			r := getReconciler(objs)

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
				instance: &instanceOne,
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
				instance: &instanceOne,
				kind:     CONFIGMAP,
			},
		},
		{
			name: "should create the proxy Service",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&corev1.Service{}),
			args: args{
				instance: &instanceOne,
				kind:     PROXY_SERVICE,
			},
		},
		{
			name: "should create the application Service",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&corev1.Service{}),
			args: args{
				instance: &instanceOne,
				kind:     APPLICATION_SERVICE,
			},
		},
		{
			name: "should create a Route",
			fields: fields{
				scheme: scheme.Scheme,
			},
			want: reflect.TypeOf(&routev1.Route{}),
			args: args{
				instance: &instanceOne,
				kind:     ROUTE,
			},
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &instanceOne,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r := getReconciler(objs)

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
		&instanceOne,
	}

	r := getReconciler(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instanceOne.Name,
			Namespace: instanceOne.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), req.NamespacedName, configMap)

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
	err = r.client.Get(context.TODO(), req.NamespacedName, dep)
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

	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      utils.PROXY_SERVICE_INSTANCE_NAME,
		Namespace: instanceOne.Namespace,
	}, service)
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("get service: (%v)", service)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      utils.APPLICATION_SERVICE_INSTANCE_NAME,
		Namespace: instanceOne.Namespace,
	}, service)
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("get service: (%v)", service)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if !res.Requeue {
		t.Error("reconcile unexpectedly requeued request")
	}

	route := &routev1.Route{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetRouteName(&instanceOne), Namespace: instanceOne.Namespace}, route)
	if err != nil {
		t.Fatalf("get route: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
}

func TestReconcileMobileSecurityService_Reconcile_InvalidInstance(t *testing.T) {
	invalidInstance := &instanceOne
	invalidInstance.Spec.ClusterProtocol = "ws"

	// objects to track in the fake client
	objs := []runtime.Object{
		&instanceOne,
	}

	r := getReconciler(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instanceOne.Name,
			Namespace: instanceOne.Namespace,
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

func TestReconcileMobileSecurityService_Reconcile_UnknownNamespace(t *testing.T) {
	// objects to track in the fake client
	objs := []runtime.Object{
		&instanceOne,
	}

	r := getReconciler(objs)

	namespace := "unknown-namespace"

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instanceOne.Name,
			Namespace: namespace,
		},
	}

	_, err := r.Reconcile(req)
	if err == nil {
		t.Fatalf("expected not to find namespace '%v'", namespace)
	}
}

func getReconciler(objs []runtime.Object) *ReconcileMobileSecurityService {
	s := scheme.Scheme

	routev1.AddToScheme(s)

	s.AddKnownTypes(mobilesecurityservicev1alpha1.SchemeGroupVersion, &mobilesecurityservicev1alpha1.MobileSecurityService{})
	s.AddKnownTypes(routev1.SchemeGroupVersion, &routev1.Route{})
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ConfigMap{}, &corev1.Service{})

	// create a fake client to mock API calls
	cl := fake.NewFakeClient(objs...)
	// create a ReconcileMobileSecurityService object with the scheme and fake client
	return &ReconcileMobileSecurityService{client: cl, scheme: s}
}
