# Mobile Security Service - Standard Operating Procedures

# Critical

## MobileSecurityServiceConsoleDown

Troubleshoot using the following steps via either the console or the cli:

Console:
1. Check the operator pod is present as it is responsible for managing the service pod as described in [MobileSecurityServiceOperatorDown](https://github.com/aerogear/mobile-security-service-operator/SOP/SOP-operator.md)
    1. If resolving the MobileSecurityServiceOperatorDown doesn't resolve the issue, please continue with the below steps

2. Check that the Service Pod is deployed in the same namespace as the operator
3. Check the status of the Mobile Security Service Custom Resources
    1. Go To -> Resources -> Other Resources -> Choose a resource to list -> Mobile Security Service -> mobile-security-service -> Actions -> Edit yaml
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
4. Check the status of the Mobile Security Service DB Custom Resource. If the Database Pod is not available this will cause the Service pod to error
    1. Go To -> Resources -> Other Resources -> Choose a resource to list -> Mobile Security Service DB-> mobile-security-service-db -> Actions -> Edit yaml
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
6. If the service pod is present check the logs of the OAuth Proxy Container
    1. Go To -> Applications -> Pods -> mobile-security-service-<xyz123> -> Logs -> Container -> oauth-proxy
7. If the service pod is present check the logs of the Application Container
    1. Go To -> Applications -> Pods -> mobile-security-service-<xyz123> -> Logs -> Container -> application


CLI:
1. Check the operator pod is present as it is responsible for managing the service pod as described in [MobileSecurityServiceOperatorDown](https://github.com/aerogear/mobile-security-service-operator/SOP/SOP-operator.md)
    1. If resolving the MobileSecurityServiceOperatorDown doesn't resolve the issue, please continue with the below steps

2. Check that the Service Pod is deployed in the same namespace as the operator
2. Check the status of the Mobile Security Service Custom Resources
    1. `oc get MobileSecurityService mobile-security-service -o yaml`
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
3. Check the status of the Mobile Security Service DB Custom Resource. If the Database Pod is not available this will cause the Service pod to error
    1. `oc get MobileSecurityServiceDB mobile-security-service-db -o yaml`
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
6. If the service pod is present check the logs of the OAuth Proxy Container
    1. Get the service pod name -> `oc get pods | grep mobile-security-service`
    2. `oc logs <service-podname> -c oauth-proxy`
    3. Save the logs by running `oc logs <service-podname> -c oauth-proxy > <filename>.log`
7. If the service pod is present check the logs of the Application Container
    1. `oc logs <service-podname> -c application`
    2. Save the logs by running `oc logs <service-podname> -c application > <filename>.log`


## MobileSecurityServiceDown

1. Check the steps for the [MobileSecurityServiceConsoleDown](https://github.com/aerogear/mobile-security-service-operator/SOP/SOP-operator.md) in order to troubleshoot this alert

## MobileSecurityServiceDatabaseDown

Troubleshoot using the following steps via either the console or the cli:

Console:
1. Check the operator pod is present as it is responsible for managing the service pod as described in [MobileSecurityServiceOperatorDown](https://github.com/aerogear/mobile-security-service-operator/SOP/SOP-operator.md)
    1. If resolving the MobileSecurityServiceOperatorDown doesn't resolve the issue, please continue with the below steps

2. Check that the Database Pod is deployed in the same namespace as the operator
3. Check the status of the Mobile Security Service DB Custom Resource. If the Database Pod is not available this will cause the Service pod to error
    1. Go To -> Resources -> Other Resources -> Choose a resource to list -> Mobile Security Service DB-> mobile-security-service-db -> Actions -> Edit yaml
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
4. Check that the DB pod is using values correct from a config map
    1. Navigate to Resources -> Config Maps -> mobile-security-service-config to ensure this config map exists and contains values for ACCESS_CONTROL_ALLOW_CREDENTIALS, ACCESS_CONTROL_ALLOW_ORIGIN, LOG_FORMAT, LOG_LEVEL, PGDATABASE, PGHOST, PGPASSWORD, PGUSER.
    2. Navigate to Applications -> Deployments -> mobile-security-service-db -> Environment
    3. Verify that the Container dtabase is sourcing POSTGRESQL_DATABASE, POSTGRESQL_USER, POSTGRESQL_PASSWORD from the config map mobile-security-service-config Config Map
3. If the DB pod is present check the logs of the pod for any errors
    1. Navigate to Applications -> Pods -> mobile-security-service-db-<xyz123> -> Logs


CLI:
1. Check the operator pod is present as it is responsible for managing the service pod as described in [MobileSecurityServiceOperatorDown](https://github.com/aerogear/mobile-security-service-operator/SOP/SOP-operator.md)
    1. If resolving the MobileSecurityServiceOperatorDown doesn't resolve the issue, please continue with the below steps

2. Check that the Database Pod is deployed in the same namespace as the operator
2. Check the status of the Mobile Security Service DB Custom Resource. If the Database Pod is not available this will cause the Service pod to error
    1. `oc get MobileSecurityServiceDB mobile-security-service-db -o yaml`
    2. Check that the status' are as expected as described in the [README](https://github.com/aerogear/mobile-security-service-operator#status-definition-per-types)
3. If the DB pod is present check the logs of the pod for any errors
    1. Get the service pod name -> `oc get pods | grep mobile-security-service-db`
    2. `oc logs <database-podname>`
    3. Save the logs by running `oc logs <database-podname> > <filename>.log`

# Warning

## MobileSecurityServicePodCPUHigh ## MobileSecurityServicePodMemoryHigh

1. Capture a snapshot of the 'Mobile Security Service Application' Grafana dashboard and track it over time. The metrics can be useful for identifying performance isssues over time.


## MobileSecurityServiceApiHighRequestDuration, MobileSecurityServiceApiHighRequestFailure , MobileSecurityServiceApiHighConcurrentRequests

1. Capture a snapshot of the 'Mobile Security Service Application' Grafana dashboard. The metrics can be useful for identifying errors and performance bottlenecks.

Troubleshoot using the following steps via either the console or the cli:

Console:
1. Capture a snapshot of the 'Mobile Security Service Application' Grafana dashboard. The metrics can be useful for identifying errors and performance bottlenecks.
2. If necessary, recreate the syndesis-server pod to restore service.
    1. Navigate to Application -> Pods -> mobile-security-service-<xyz123> -> Actions -> Delete -> Delete

CLI:
1. Capture application logs for analysis.
    1. Get the service pod name -> `oc get pods | grep mobile-security-service`
    2. `oc logs <service-podname> -c application`
    3. Save the logs by running `oc logs <service-podname> -c application > <filename>.log`
2. If necessary, recreate the syndesis-server pod to restore service.
    1. Get the service pod name -> `oc get pods | grep mobile-security-service`
    2. oc delete pod <service-podname>
