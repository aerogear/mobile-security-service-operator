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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
				createdInstance:  &mssInstance,
				instanceToUpdate: &mssInstance,
				scheme:           scheme.Scheme,
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
		{
			name: "should give error when namespace not found",
			fields: fields{
				createdInstance:  &mssInstance,
				instanceToUpdate: &mssInstance2,
				scheme:           scheme.Scheme,
			},
			want:    reconcile.Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.fields.createdInstance.Name,
					Namespace: tt.fields.createdInstance.Namespace,
				},
			}

			res, err := r.Reconcile(req)
			if err != nil {
				t.Fatalf("reconcile: (%v)", err)
			}

			reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
			err = r.update(tt.fields.instanceToUpdate, reqLogger)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) == tt.wantErr && !res.Requeue {
				t.Errorf("reconcile did not requeue request as expected")
				return
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
				instance: &mssInstance,
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
				instance: &mssInstance,
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

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityService.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, reqLogger, tt.args.kind)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityService.create() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				instance: &mssInstance,
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
				instance: &mssInstance,
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
				instance: &mssInstance,
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
				instance: &mssInstance,
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
				instance: &mssInstance,
				kind:     ROUTE,
			},
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance: &mssInstance,
				kind:     "UNDEFINED",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityService.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got := r.buildFactory(reqLogger, tt.args.instance, tt.args.kind)
			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("ReconcileMobileSecurityService.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityService_Reconcile(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      mssInstance.Name,
			Namespace: mssInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetConfigMapName(&mssInstance), Namespace: mssInstance.Namespace}, configMap)
	if err != nil {
		t.Fatalf("get configMap: (%v)", err)
	}

	// Check the result of reconciliation to make sure it has the desired state
	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	// check if the deployment has been created
	dep := &v1beta1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	// Check if the quantity of Replicas for this deployment is equals the specification
	dsize := *dep.Spec.Replicas
	if dsize != mssInstance.Spec.Size {
		t.Errorf("dep size (%d) is not the expected size (%d)", dsize, mssInstance.Spec.Size)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      utils.APPLICATION_SERVICE_INSTANCE_NAME,
		Namespace: mssInstance.Namespace,
	}, service)
	if err != nil {
		t.Fatalf("get application service: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	// check if the service has been created
	proxyService := &corev1.Service{}

	err = r.client.Get(context.TODO(), types.NamespacedName{
		Name:      utils.PROXY_SERVICE_INSTANCE_NAME,
		Namespace: mssInstance.Namespace,
	}, proxyService)
	if err != nil {
		t.Fatalf("get proxy service: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	route := &routev1.Route{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: utils.GetRouteName(&mssInstance), Namespace: mssInstance.Namespace}, route)
	if err != nil {
		t.Fatalf("get route: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

}

func TestReconcileMobileSecurityService_Reconcile_InvalidSpec(t *testing.T) {
	mssInstance.Spec.ClusterProtocol = "invalid"

	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      mssInstance.Name,
			Namespace: mssInstance.Namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("reconcile requeue request which is not expected")
	}
}

func TestReconcileMobileSecurityService_Reconcile_UnknownNamespace(t *testing.T) {
	// objects to track in the fake client
	objs := []runtime.Object{
		&mssInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	namespace := "unknown-namespace"

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      mssInstance.Name,
			Namespace: namespace,
		},
	}

	_, err := r.Reconcile(req)
	if err == nil {
		t.Fatalf("expected not to find namespace '%v'", namespace)
	}
}
