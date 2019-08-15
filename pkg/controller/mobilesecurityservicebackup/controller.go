package mobilesecurityservicebackup

import (
	"context"
	mobilesecurityservicev1alpha1 "github.com/aerogear/mobile-security-service-operator/pkg/apis/mobile-security-service/v1alpha1"
	"github.com/aerogear/mobile-security-service-operator/pkg/utils"
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

const (
	CronJob   = "CronJob"
	DBSecret  = "DBSecret"
	AwsSecret = "AwsSecret"
	EncSecret = "EncSecret"
)

var log = logf.Log.WithName("controller_mobilesecurityservicebackup")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MobileSecurityServiceBackup Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMobileSecurityServiceBackup{client: mgr.GetClient(), scheme: mgr.GetScheme(), config: mgr.GetConfig()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mobilesecurityservicebackup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MobileSecurityServiceBackup
	err = c.Watch(&source.Kind{Type: &mobilesecurityservicev1alpha1.MobileSecurityServiceBackup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch CronJob
	if err := watchCronJob(c); err != nil {
		return err
	}

	// Watch watchSecret
	if err := watchSecret(c); err != nil {
		return err
	}

	// Watch Pod
	if err := watchPod(c); err != nil {
		return err
	}

	// Watch Service
	if err := watchService(c); err != nil {
		return err
	}

	return nil
}

// Update the object and reconcile it
func (r *ReconcileMobileSecurityServiceBackup) update(obj runtime.Object, reqLogger logr.Logger) error {
	err := r.client.Update(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to update Object", "obj:", obj)
		return err
	}
	reqLogger.Info("Object updated", "obj:", obj)
	return nil
}

// blank assignment to verify that ReconcileMobileSecurityServiceBackup implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMobileSecurityServiceBackup{}

// ReconcileMobileSecurityServiceBackup reconciles a MobileSecurityServiceBackup object
type ReconcileMobileSecurityServiceBackup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	config    *rest.Config
	scheme    *runtime.Scheme
	dbPod     *corev1.Pod
	dbService *corev1.Service
}

// Create the object and reconcile it
func (r *ReconcileMobileSecurityServiceBackup) create(bkp *mobilesecurityservicev1alpha1.MobileSecurityServiceBackup, kind string, reqLogger logr.Logger) error {
	obj, err := r.buildFactory(bkp, kind, reqLogger)
	if err != nil {
		reqLogger.Error(err, "Failed to build object ", "kind", kind, "Namespace", bkp.Namespace)
		return err
	}
	reqLogger.Info("Creating a new ", "kind", kind, "Namespace", bkp.Namespace)
	err = r.client.Create(context.TODO(), obj)
	if err != nil {
		reqLogger.Error(err, "Failed to create new ", "kind", kind, "Namespace", bkp.Namespace)
		return err
	}
	reqLogger.Info("Created successfully", "kind", kind, "Namespace", bkp.Namespace)
	return nil
}

// buildFactory will return the resource according to the kind defined
func (r *ReconcileMobileSecurityServiceBackup) buildFactory(bkp *mobilesecurityservicev1alpha1.MobileSecurityServiceBackup, kind string, reqLogger logr.Logger) (runtime.Object, error) {
	reqLogger.Info("Check "+kind, "into the namespace", bkp.Namespace)
	switch kind {
	case CronJob:
		return r.buildCronJob(bkp), nil
	case DBSecret:
		// build Database secret data
		secretData, err := r.buildDBSecretData(bkp)
		if err != nil {
			reqLogger.Error(err, "Unable to create DB Data secret")
			return nil, err
		}
		return r.buildSecret(bkp, dbSecretPrefix, secretData, nil), nil
	case AwsSecret:
		secretData := buildAwsSecretData(bkp)
		return r.buildSecret(bkp, awsSecretPrefix, secretData, nil), nil
	case EncSecret:
		secretData, secretStringData := buildEncSecretData(bkp)
		return r.buildSecret(bkp, encryptionKeySecret, secretData, secretStringData), nil
	default:
		msg := "Failed to recognize type of object" + kind + " into the Namespace " + bkp.Namespace
		panic(msg)
	}
}

// Reconcile reads that state of the cluster for a MobileSecurityServiceBackup object and makes changes based on the state read
// and what is in the MobileSecurityServiceBackup.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMobileSecurityServiceBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MobileSecurityServiceBackup")

	// Fetch the MobileSecurityService DB Backup
	bkp := &mobilesecurityservicev1alpha1.MobileSecurityServiceBackup{}
	bkp, err := r.fetchBkpInstance(reqLogger, request)
	if err != nil {
		reqLogger.Error(err, "Failed to get Mobile Security Service Backup")
		return reconcile.Result{}, err
	}

	// Check if the DB BKP CR was applied in the same namespace of the operator
	if isValidNamespace, err := utils.IsValidOperatorNamespace(bkp.Namespace); err != nil || isValidNamespace == false {
		operatorNamespace, _ := k8sutil.GetOperatorNamespace()
		reqLogger.Error(err, "Unable to reconcile Mobile Security Service Backup", "bkp.Namespace", bkp.Namespace, "isValidNamespace", isValidNamespace, "Operator.Namespace", operatorNamespace)

		//Update status with Invalid Namespace
		if err := r.updateStatusWithInvalidNamespace(reqLogger, request); err != nil {
			return reconcile.Result{}, err
		}

		// Stop reconcile
		return reconcile.Result{}, nil
	}
	reqLogger.Info("Valid namespace for Mobile Security Service Backup", "request.namespace", request.Namespace)

	// Add const values for mandatory specs
	addMandatorySpecsDefinitions(bkp)

	// Get database pod
	dbPod, err := r.fetchBDPod(reqLogger, request)
	if err != nil || dbPod == nil {
		reqLogger.Error(err, "Unable to find the database pod", "request.namespace", request.Namespace)
		return reconcile.Result{RequeueAfter: time.Second * 10}, err
	}

	// set in the reconcile
	r.dbPod = dbPod

	// Get database service
	dbService, err := r.fetchServiceDB(reqLogger, request)
	if err != nil || dbService == nil {
		reqLogger.Error(err, "Unable to find the database service", "request.namespace", request.Namespace)
		return reconcile.Result{RequeueAfter: time.Second * 10}, err
	}

	// set in the reconcile
	r.dbService = dbService

	// Check if the secret for the database is created, if not create one
	if _, err := r.fetchSecret(reqLogger, bkp.Namespace, dbSecretPrefix+bkp.Name); err != nil {
		if err := r.create(bkp, DBSecret, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if the secret for the s3 is created, if not create one
	if _, err := r.fetchSecret(reqLogger, getAwsSecretNamespace(bkp), getAWSSecretName(bkp)); err != nil {
		if bkp.Spec.AwsCredentialsSecretName != "" {
			reqLogger.Error(err, "Unable to find AWS secret informed and will not be created by the operator", "SecretName", bkp.Spec.AwsCredentialsSecretName)
			return reconcile.Result{}, err
		}
		if err := r.create(bkp, AwsSecret, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if the secret for the encryptionKey is created, if not create one just when the data is informed
	if hasEncryptionKeySecret(bkp) {
		if _, err := r.fetchSecret(reqLogger, getEncSecretNamespace(bkp), getEncSecretName(bkp)); err != nil {
			if bkp.Spec.EncryptionKeySecretName != "" {
				reqLogger.Error(err, "Unable to find EncryptionKey secret informed and will not be created by the operator", "SecretName", bkp.Spec.EncryptionKeySecretName)
				return reconcile.Result{}, err
			}
			if err := r.create(bkp, EncSecret, reqLogger); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	// Check if the CronJob is created, if not create one
	if _, err := r.fetchCronJob(reqLogger, bkp); err != nil {
		if err := r.create(bkp, CronJob, reqLogger); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update status for pod database found
	if err := r.updatePodDatabaseFoundStatus(reqLogger, request, dbPod); err != nil {
		return reconcile.Result{}, err
	}

	// Update status for service database found
	if err := r.updateServiceDatabaseFoundStatus(reqLogger, request, dbService); err != nil {
		return reconcile.Result{}, err
	}

	// Update status for CronJobStatus
	cronJobStatus, err := r.updateCronJobStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for database secret
	dbSecretStatus, err := r.updateDBSecretStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for aws secret
	awsSecretStatus, err := r.updateAWSSecretStatus(reqLogger, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update status for pod database found
	if err := r.updateEncSecretStatus(reqLogger, request); err != nil {
		return reconcile.Result{}, err
	}

	// Update status for Backup
	if err := r.updateBackupStatus(reqLogger, cronJobStatus, dbSecretStatus, awsSecretStatus, dbPod, dbService, request); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
