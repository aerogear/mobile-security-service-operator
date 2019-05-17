package mobilesecurityservicedb

import (
	"context"
	"reflect"
	"testing"

	"github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileMobileSecurityServiceDB_update(t *testing.T) {
	type fields struct {
		createdInstance  *v1alpha1.MobileSecurityServiceDB
		instanceToUpdate *v1alpha1.MobileSecurityServiceDB
		scheme           *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should successfully update the instance",
			fields: fields{
				createdInstance:  &instanceOne,
				instanceToUpdate: &instanceOne,
			},
			wantErr: false,
		},
		{
			name: "should give an error when the namespace is not found",
			fields: fields{
				createdInstance:  &instanceOne,
				instanceToUpdate: &instanceTwo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := newReconcilerWithFakeClient(objs)

			reqLogger := log.WithValues("Request.Namespace", tt.fields.instanceToUpdate.Namespace, "Request.Name", tt.fields.createdInstance.Name)

			err := r.update(tt.fields.instanceToUpdate, reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceDB.update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileMobileSecurityServiceDB_create(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance        *mobilesecurityservicev1alpha1.MobileSecurityServiceDB
		serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService
		kind            string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "should create and return a new deployment",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:        &instanceOne,
				serviceInstance: &serviceInstance,
				kind:            DEEPLOYMENT,
			},
			wantErr: false,
		},
		{
			name: "should fail to create an unknown type",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:        &instanceOne,
				serviceInstance: &serviceInstance,
				kind:            "UNKNOWN",
			},
			wantErr:   false,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r := newReconcilerWithFakeClient(objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceDB.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, tt.args.serviceInstance, tt.args.kind, reqLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileMobileSecurityServiceDB.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReconcileMobileSecurityServiceDB_buildFactory(t *testing.T) {
	type fields struct {
		scheme *runtime.Scheme
	}
	type args struct {
		instance        *mobilesecurityservicev1alpha1.MobileSecurityServiceDB
		serviceInstance *mobilesecurityservicev1alpha1.MobileSecurityService
		kind            string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      reflect.Type
		wantPanic bool
	}{
		{
			name: "should create a Deployment",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:        &instanceOne,
				serviceInstance: &serviceInstance,
				kind:            DEEPLOYMENT,
			},
			want: reflect.TypeOf(&v1beta1.Deployment{}),
		},
		{
			name: "should create a Persistent Volume Claim",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:        &instanceOne,
				serviceInstance: &serviceInstance,
				kind:            PVC,
			},
			want: reflect.TypeOf(&corev1.PersistentVolumeClaim{}),
		},
		{
			name: "should create a Service",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:        &instanceOne,
				serviceInstance: &serviceInstance,
				kind:            SERVICE,
			},
			want: reflect.TypeOf(&corev1.Service{}),
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:        &instanceOne,
				serviceInstance: &serviceInstance,
				kind:            "UNDEFINED",
			},
			wantPanic: true,
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.args.instance}

			r := newReconcilerWithFakeClient(objs)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceDB.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got := r.buildFactory(tt.args.instance, tt.args.serviceInstance, tt.args.kind, reqLogger)

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("ReconcileMobileSecurityServiceDB.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityServiceDB_Reconcile(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&instanceOne,
	}

	r := newReconcilerWithFakeClient(objs)

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

	deployment := &v1beta1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	service := &corev1.Service{}
	err = r.client.Get(context.TODO(), req.NamespacedName, service)
	if err != nil {
		t.Fatalf("get service: (%v)", err)
	}

	if res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), req.NamespacedName, pvc)
	if err != nil {
		t.Fatalf("get pvc: (%v)", err)
	}

	if res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileMobileSecurityServiceDB_Reconcile_NotFound(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&instanceOne,
	}

	r := newReconcilerWithFakeClient(objs)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      instanceOne.Name,
			Namespace: "unknown",
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Fatalf("did not expected reconcile to requeue.")
	}
}

// Creates a new reconciler with a fake client
func newReconcilerWithFakeClient(objs []runtime.Object) *ReconcileMobileSecurityServiceDB {
	s := scheme.Scheme

	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.MobileSecurityServiceDB{})

	// create a fake client to mock API calls
	cl := fake.NewFakeClient(objs...)
	// create a ReconcileMobileSecurityService object with the scheme and fake client
	return &ReconcileMobileSecurityServiceDB{client: cl, scheme: s}
}
