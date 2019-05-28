package mobilesecurityserviceapp

import (
	"time"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	routev1 "github.com/openshift/api/route/v1"
)

var (
	instance = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "test-app",
			AppId:   "test-app-id",
		},
	}

	instanceNoAppName = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "",
			AppId:   "test-app-id",
		},
	}

	instanceNoAppId = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "test-app",
			AppId:   "",
		},
	}

	instanceInvalidName = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "invalid",
			Namespace: "mobile-security-service",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "test-app",
			AppId:   "test-app-id",
		},
	}

	instanceInvalidNameSpace = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "invalid",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "test-app",
			AppId:   "test-app-id",
		},
	}

	instanceForDeletion = mobilesecurityservicev1alpha1.MobileSecurityServiceApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service",
			DeletionTimestamp: &metav1.Time{time.Now()},
			Finalizers: []string{
				"test-finalizer",
			},				
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceAppSpec{
			AppName: "test-app",
			AppId:   "test-app-id",
		},
	}



	mssInstance = mobilesecurityservicev1alpha1.MobileSecurityService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mobile-security-service",
			Namespace: "mobile-security-service-proxy",
		},
		Spec: mobilesecurityservicev1alpha1.MobileSecurityServiceSpec{
			Size:            1,
			MemoryLimit:     "512Mi",
			MemoryRequest:   "512Mi",
			ClusterProtocol: "http",
			ConfigMapName:   "mss-config",
			Port:            1234,
			RouteName:       "mss-route",
		},
	}

	route = routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GetRouteName(&mssInstance),
			Namespace: mssInstance.Namespace,
			Labels:    getAppLabels(mssInstance.Name),
		},
		Status: routev1.RouteStatus{
			Ingress: []routev1.RouteIngress{
				{
					Host:           "testhost",
					RouterName:     "",
					WildcardPolicy: "",
				},
			},
		},
	}
)
