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
				createdInstance:  &dbInstance,
				instanceToUpdate: &dbInstance,
			},
			wantErr: false,
		},
		{
			name: "should give an error when the namespace is not found",
			fields: fields{
				createdInstance:  &dbInstance,
				instanceToUpdate: &dbInstanceNonDefaultNamespace,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.createdInstance}

			r := buildReconcileWithFakeClientWithMocks(objs, t)

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
				instance:        &dbInstance,
				serviceInstance: &serviceInstance,
				kind:            Deployment,
			},
			wantErr: false,
		},
		{
			name: "should fail to create an unknown type",
			fields: fields{
				scheme.Scheme,
			},
			args: args{
				instance:        &dbInstance,
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

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceDB.create() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			err := r.create(tt.args.instance, tt.args.kind, reqLogger)
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
				instance:        &dbInstance,
				serviceInstance: &serviceInstance,
				kind:            Deployment,
			},
			want: reflect.TypeOf(&v1beta1.Deployment{}),
		},
		{
			name: "should create a Persistent Volume Claim",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:        &dbInstance,
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
				instance:        &dbInstance,
				serviceInstance: &serviceInstance,
				kind:            Service,
			},
			want: reflect.TypeOf(&corev1.Service{}),
		},
		{
			name: "Should panic when trying to create unrecognized object type",
			fields: fields{
				scheme: scheme.Scheme,
			},
			args: args{
				instance:        &dbInstance,
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

			r := buildReconcileWithFakeClientWithMocks(objs, t)

			reqLogger := log.WithValues("Request.Namespace", tt.args.instance.Namespace, "Request.Name", tt.args.instance.Name)

			// testing if the panic() function is executed
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("ReconcileMobileSecurityServiceDB.buildFactory() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			got := r.buildFactory(tt.args.instance, tt.args.kind, reqLogger)

			if gotType := reflect.TypeOf(got); !reflect.DeepEqual(gotType, tt.want) {
				t.Errorf("ReconcileMobileSecurityServiceDB.buildFactory() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestReconcileMobileSecurityServiceDB_Reconcile(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
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
		t.Error("did not expect request to requeue")
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
		t.Error("did not expect request to requeue")
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
		t.Error("did not expect request to requeue")
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
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: "unknown",
		},
	}

	res, err := r.Reconcile(req)
	if err == nil {
		t.Error("should fail since the instance do not exist in the <unknown> nammespace")
	}

	if res.Requeue {
		t.Fatalf("did not expected reconcile to requeue.")
	}
}

func TestReconcileMobileSecurityServiceDB_Reconcile_UsingMSSConfigMapToCreateEnvVars(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
		&serviceInstance,
		&configMap,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
		},
	}

	_ = r.client.Create(context.TODO(), &serviceInstance)
	_ = r.client.Create(context.TODO(), &configMap)

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment := &v1beta1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom == nil {
		t.Error("deployment envvar did not came from service instance config map")
	}

	if deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name != configMap.Name {
		t.Fatalf("deployment envvar did not came from service instance config map: (%v,%v)", deployment.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name, configMap.Name)
	}

	if res.Requeue {
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
		t.Error("did not expect request to requeue")
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
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}

func TestReconcileMobileSecurityServiceDB_Reconcile_ReplicasSizes(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstance,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstance.Name,
			Namespace: dbInstance.Namespace,
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
		t.Error("did not expect request to requeue")
	}

	//Mock Replicas wrong size
	size := int32(3)
	deployment.Spec.Replicas = &size

	// Update
	err = r.client.Update(context.TODO(), deployment)
	if err != nil {
		t.Fatalf("fails when ttry to update deployment replicas: (%v)", err)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	deployment = &v1beta1.Deployment{}
	err = r.client.Get(context.TODO(), req.NamespacedName, deployment)
	if err != nil {
		t.Fatalf("get deployment: (%v)", err)
	}

	if *deployment.Spec.Replicas != dbInstance.Spec.Size {
		t.Error("Replicas size was not respected")
	}
}

func TestReconcileMobileSecurityServiceDB_Reconcile_InstanceWithoutSpec(t *testing.T) {

	// objects to track in the fake client
	objs := []runtime.Object{
		&dbInstanceWithoutSpec,
	}

	r := buildReconcileWithFakeClientWithMocks(objs, t)

	// mock request to simulate Reconcile() being called on an event for a watched resource
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      dbInstanceWithoutSpec.Name,
			Namespace: dbInstanceWithoutSpec.Namespace,
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
		t.Error("did not expect request to requeue")
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
		t.Error("did not expect request to requeue")
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
		t.Error("did not expect request to requeue")
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	if res.Requeue {
		t.Error("did not expect request to requeue")
	}
}
