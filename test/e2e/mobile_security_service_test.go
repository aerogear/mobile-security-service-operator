package e2e

import (
	goctx "context"
	apis "github.com/aerogear/mobile-security-service-operator/pkg/apis"
	mssv1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobilesecurityservice/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	e2eutil "github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"testing"
	"time"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 200
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestMss(t *testing.T) {
	mssList := &mssv1alpha1.MobileSecurityServiceList{}
	if err := framework.AddToFrameworkScheme(apis.AddToScheme, mssList); err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	t.Run("mss-e2e", MssTest)
}
func MssTest(t *testing.T) {
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	// get namespace
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatalf("failed to get namespace %v", err)
	}
	// get global framework variables
	f := framework.Global
	if err := initializeMssResources(t, f, ctx, namespace); err != nil {
		t.Fatal(err)
	}
	//MSS database CR for testing
	mssdbName := "mobile-security-service-db"
	mssDBTestCR := &mssv1alpha1.MobileSecurityServiceDB{
		TypeMeta: metav1.TypeMeta{
			APIVersion: " mobile-security-service.aerogear.org/v1alpha1",
			Kind:       "MobileSecurityServiceDB",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mssdbName,
			Namespace: namespace,
		},
		Spec: mssv1alpha1.MobileSecurityServiceDBSpec{

			DatabaseName:           "mobile_security_service",
			DatabasePassword:       "postgres",
			DatabaseUser:           "postgresql",
			DatabaseNameParam:      "POSTGRESQL_DATABASE",
			DatabasePasswordParam:  "POSTGRESQL_PASSWORD",
			DatabaseUserParam:      "POSTGRESQL_USER",
			DatabasePort:           5432,
			Image:                  "centos/postgresql-96-centos7",
			ContainerName:          "database",
			DatabaseMemoryLimit:    "512Mi",
			DatabaseMemoryRequest:  "512Mi",
			DatabaseStorageRequest: "1Gi",
		},
		Status: mssv1alpha1.MobileSecurityServiceDBStatus{
			DatabaseStatus: "OK",
		},
	}
	//MSS CR struct for testing
	mssName := "mobile-security-service"
	mssTestCR := &mssv1alpha1.MobileSecurityService{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "mobile-security-service.aerogear.org/v1alpha1",
			Kind:       "MobileSecurityService",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mssName,
			Namespace: namespace,
		},
		Spec: mssv1alpha1.MobileSecurityServiceSpec{
			DatabaseName:                  "mobile_security_service",
			DatabasePassword:              "postgres",
			DatabaseUser:                  "postgresql",
			DatabaseHost:                  "mobile-security-service-db",
			LogLevel:                      "info",
			LogFormat:                     "json",
			AccessControlAllowOrigin:      "*",
			AccessControlAllowCredentials: "false",
			Port:                          3000,
			Image:                         "quay.io/aerogear/mobile-security-service:0.2.2",
			ContainerName:                 "application",
			ClusterProtocol:               "https",
		},
		Status: mssv1alpha1.MobileSecurityServiceStatus{
			AppStatus:      "OK",
			ConfigMapName:  "mobile-security-service-config",
			DeploymentName: "mobile-security-service",
		},
	}

	// Create MSS database Custom Resource
	if err := createMSSdbCustomeResource(t, f, ctx, mssDBTestCR, namespace, mssdbName); err != nil {
		t.Fatal(err)
	}
	//Create MSS Custom Resource
	if err := createMSSCustomResource(t, f, ctx, mssTestCR, namespace, mssName); err != nil {
		t.Fatal(err)
	}

	//Delete MSS database CR
	if err := deleteMSSdbCustomResource(t, f, ctx, mssDBTestCR, namespace, mssdbName); err != nil {
		t.Fatal(err)
	}

	//Delete MSS CR
	if err := deleteMSSCustomResource(t, f, ctx, mssTestCR, namespace, mssName); err != nil {
		t.Fatal(err)
	}

}
func initializeMssResources(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, namespace string) error {
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Successfully initialized cluster resources")

	// wait for mobile-security-service-operator to be ready
	if err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "mobile-security-service-operator", 1, retryInterval, timeout); err != nil {
		t.Fatal(err)
	}
	return err
}

func createMSSdbCustomeResource(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, testCR *mssv1alpha1.MobileSecurityServiceDB, namespace string, mssdbName string) error {
	err := f.Client.Create(goctx.TODO(), testCR, &framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	})
	if err != nil {
		return err
	}
	t.Log("Successfully created MSS database Custom Resource")
	// Ensure MSS database was deployed successfully
	if err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, mssdbName, 1, retryInterval, timeout); err != nil {
		t.Fatal(err)
	}
	t.Log("MSS database deployment was successful")
	return err
}
func createMSSCustomResource(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, testCR *mssv1alpha1.MobileSecurityService, namespace string, mssName string) error {
	err := f.Client.Create(goctx.TODO(), testCR, &framework.CleanupOptions{
		TestContext:   ctx,
		Timeout:       cleanupTimeout,
		RetryInterval: cleanupRetryInterval,
	})
	if err != nil {
		return err
	}
	t.Log("Successfully created MSS Custom Resource")
	//Ensure MSS was deployed successfully
	if err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, mssName, 1, retryInterval, timeout); err != nil {
		t.Fatal(err)
	}
	t.Log("MSS deployment was successful")

	return err
}
func deleteMSSdbCustomResource(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, mssDBTestCR *mssv1alpha1.MobileSecurityServiceDB, namespace string, mssdbName string) error {
	err := f.Client.Delete(goctx.TODO(), mssDBTestCR)
	if err != nil {
		return err
	}
	t.Log("Successfully deleted MSS database Custom Resource")

	return nil
}
func deleteMSSCustomResource(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, mssTestCR *mssv1alpha1.MobileSecurityService, namespace string, mssName string) error {
	err := f.Client.Delete(goctx.TODO(), mssTestCR)
	if err != nil {
		return err
	}
	t.Log("Successfully deleted MSS Custom Resource")
	return nil
}
