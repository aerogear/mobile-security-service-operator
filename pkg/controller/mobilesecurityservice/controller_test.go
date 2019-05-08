package mobilesecurityservice

import (
	"reflect"
	"testing"

	"github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// func TestReconcileMobileSecurityService_Reconcile(t *testing.T) {

// 	var (
// 		name                                = "mobile-security-service-operator"
// 		namespace                           = "mobile-security-service-app"
// 		size                          int32 = 1
// 		image                               = "aerogear/mobile-security-service:master"
// 		databaseName                        = "mobile_security_service"
// 		databaseHost                        = "mobile-security-service-db"
// 		port                          int32 = 3000
// 		logLevel                            = "info"
// 		logFormat                           = "text"
// 		accessControlAllowOrigin            = "*"
// 		accessControlAllowCredentials       = "false"
// 		memoryLimit                         = "512Mi"
// 		memoryRequest                       = "512Mi"
// 		clusterHost                         = "192.168.0.1"
// 		protocol                            = "http"
// 		hostSuffix                          = ".nip.io"
// 		configMapName                       = "mss-config"
// 	)

// 	instance := &mssv1alpha1.MobileSecurityService{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      name,
// 			Namespace: namespace,
// 		},
// 		Spec: v1alpha1.MobileSecurityServiceSpec{
// 			Size:                          size,
// 			Image:                         image,
// 			DatabaseHost:                  databaseHost,
// 			DatabaseName:                  databaseName,
// 			Port:                          port,
// 			LogLevel:                      logLevel,
// 			LogFormat:                     logFormat,
// 			AccessControlAllowOrigin:      accessControlAllowOrigin,
// 			AccessControlAllowCredentials: accessControlAllowCredentials,
// 			MemoryLimit:                   memoryLimit,
// 			MemoryRequest:                 memoryRequest,
// 			ClusterHost:                   clusterHost,
// 			Protocol:                      protocol,
// 			HostSufix:                     hostSuffix,
// 			ConfigMapName:                 configMapName,
// 		},
// 	}

// 	objs := []runtime.Object{
// 		instance,
// 	}

// 	s := scheme.Scheme
// 	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, instance)
// 	cl := fake.NewFakeClient(objs...)

// 	r := &ReconcileMobileSecurityService{cl, s}

// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      name,
// 			Namespace: namespace,
// 		},
// 	}

// 	res, err := r.Reconcile(req)
// 	if err != nil {
// 		t.Fatalf("reconcile: (%v)", err)
// 	}
// 	if !res.Requeue {
// 		t.Error("reconcile did not requeue")
// 	}

// 	dep := &v1beta1.Deployment{}
// 	err = cl.Get(context.TODO(), req.NamespacedName, dep)

// 	if err != nil {
// 		// FAILS HERE
// 		t.Fatalf("get deployment: (%v)", err)
// 	}

// 	dSize := *dep.Spec.Replicas
// 	if dSize != size {
// 		t.Errorf("dep size (%d) is not the expected size (%d)", dSize, size)
// 	}
// }

func TestReconcileMobileSecurityService_update(t *testing.T) {
	type fields struct {
		instance *v1alpha1.MobileSecurityService
		scheme   *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		want    reconcile.Result
		wantErr bool
	}{
		{
			fields: fields{
				instance: &v1alpha1.MobileSecurityService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mobile-security-service-app",
						Namespace: "mobile-security-service-operator",
					},
					Spec: v1alpha1.MobileSecurityServiceSpec{
						Size: 1,
					},
				},
				scheme: scheme.Scheme,
			},
			want:    reconcile.Result{Requeue: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			objs := []runtime.Object{tt.fields.instance}

			tt.fields.scheme.AddKnownTypes(v1alpha1.SchemeGroupVersion, tt.fields.instance)

			cl := fake.NewFakeClient(objs...)

			r := &ReconcileMobileSecurityService{cl, tt.fields.scheme}

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
