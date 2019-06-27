# Mobile Security Service Operator- Standard Operating Procedures

## MobileSecurityServiceOperatorDown

Troubleshoot using the following steps via either the console or the cli:

Console:
1. Log into the OpenShift console
2. Switch to the Mobile Security Service namespace
3. Check the status of the Mobile Security Service Custom Resources
    1. Go To -> Resources -> Other Resources -> Choose a resource to list -> Mobile Security Service -> mobile-security-service -> Actions -> Edit yaml
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
4. Check the status of the Mobile Security Service DB Custom Resource
    1. Go To -> Resources -> Other Resources -> Choose a resource to list -> Mobile Security Service DB-> mobile-security-service-db -> Actions -> Edit yaml
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
5. If the operator pod is present check the logs of the pod for any errors
    1. Navigate to Applications -> Pods -> mobile-security-service-operator-<xyz123> -> Logs
    2. Check for any errors in the logs
6. If the operator pod is not present scale up the pod
    1. Navigate to Overview -> Click Dropdown for Deployment mobile-security-service-operator -> Click on up arrow to scale the pod to 1
7. In the Monitoring Tab View any Error Events


CLI:
1. Login
    1. `oc login <openshift-url>:8443 -u <username> -p <password>`
2. Switch to the Mobile Security Service namespace
    1. `oc project mobile-security-service`
3. Check the status of the Mobile Security Service Custom Resources
    1. `oc get MobileSecurityService mobile-security-service -o yaml`
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
4. Check the status of the Mobile Security Service DB Custom Resource
    1. `oc get MobileSecurityServiceDB mobile-security-service-db -o yaml`
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
5. If the operator pod is present check the logs of the pod for any errors
    1. Get the operator pod name -> `oc get pods | grep operator`
    2. `oc logs <operator-podname>`
    3. Save the logs by running `oc logs <operator-podname> > <filename>.log`
6. If the operator pod is not preset scale up the pod
    1. `oc scale deployments mobile-security-service-operator --replicas=1`


## MobileSecurityServicePodCount

Troubleshoot using the following steps via either the console or the cli:

Console:
1. Log into the OpenShift cluster
2. Switch to the Mobile Security Service namespace
3. Navigate to the Overview using the side navigation
4. Check the following
    1. There is one pod running under the application mobilesecurityservice
    2. There is one pod running under the application mobilesecurityservicedb
    3. There is one pod running under the Other Resources
5. The operator should  maintain the number of pods running for the Mobile Security Service, please resolve the operator pod numbers first by scaling it up or down as necessary.
    1. Navigate to Overview -> Click Dropdown for Deployment mobile-security-service-operator -> Click on the correct arrow to scale the pod to 1

CLI:
1. Login
    1. `oc login <openshift-url>:8443 -u <username> -p <password>`
2. Switch to the Mobile Security Service namespace
    1. `oc project mobile-security-service`
4. Check the following
    1. Run `oc get pods`
    1. There is one pod running under READY for pod mobile-security-service-<xyz123>
    2. There is one pod running under READY for pod mobile-security-service-db-<xyz123>
    3. There is one pod running under READY for pod mobile-security-service-operator<xyz123>
5. The operator should  maintain the number of pods running for the mobilesecurity service, please resolve the operator pod numbers first by scaling it up or down as necessary.
    1. `oc scale deployments mobile-security-service-operator --replicas=1`