pipeline {
    agent {
        node {
            label "psi_rhel7_openshift311"
        }
    }

    libraries {
        lib('fh-pipeline-library')
        lib('qe-pipeline-library')
    }
    
    environment {
        GOPATH = "${env.WORKSPACE}/"
        PATH = "${env.PATH}:${env.WORKSPACE}/bin:/usr/local/go/bin"
        GOOS = "linux"
        GOARCH = "amd64"
        CGO_ENABLED = 0
        OPERATOR_NAME = "mobile-security-service-operator"
        OPERATOR_CONTAINER_IMAGE_CANDIDATE_NAME = "quay.io/aerogear/${env.OPERATOR_NAME}:candidate-${env.BRANCH_NAME}"
        OPERATOR_CONTAINER_IMAGE_NAME = "quay.io/aerogear/${env.OPERATOR_NAME}:${env.BRANCH_NAME}"
        OPERATOR_CONTAINER_IMAGE_NAME_LATEST = "quay.io/aerogear/${env.OPERATOR_NAME}:latest"
        OPENSHIFT_PROJECT_NAME = "mobile-security-service"
        CLONED_REPOSITORY_PATH = "src/github.com/aerogear/mobile-security-service-operator"
        CREDENTIALS_ID = "quay-aerogear-bot"
    }

    options {
        checkoutToSubdirectory("src/github.com/aerogear/mobile-security-service-operator")
    }

    stages {

        stage("Trust"){
            steps{
                enforceTrustedApproval('aerogear')
            }
            post{
                failure{
                    echo "====++++'Trust' execution failed++++===="
                    echo "You are not authorized to run this job"
                }
        
            }
        }
        
        stage("Run oc-cluster-up"){
            steps{
                // qe-pipeline-library step
                ocClusterUp()
            }
            post{
                failure{
                    echo "====++++Run oc-cluster-up execution failed++++===="
                    echo "Try to rerun the job"
                }
        
            }
        }

        stage("Install Operator SDK") {
            steps {
                // qe-pipeline-library step
                installOperatorSdk version: "v0.10.0"
            }
            post {
                failure {
                    echo "====++++'Install Operator SDK' execution failed++++===="
                    echo "Please check if the version of operator-sdk you provided exists"
                    echo "https://github.com/operator-framework/operator-sdk/releases"
                }
            }
        }


        stage("Create an OpenShift project") {
            steps {
                // qe-pipeline-library step
                newOpenshiftProject "${env.OPENSHIFT_PROJECT_NAME}"
            }
        }
         stage("Install dependencies") {
            steps {
                dir("${env.CLONED_REPOSITORY_PATH}") {
                    sh "mkdir -p ${env.WORKSPACE}/bin"
                    sh "curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh"
                    sh "make setup"
                }
            }
        }

        stage("Build code binary"){
            steps{
                dir("${env.CLONED_REPOSITORY_PATH}") {
                    sh "make code/compile"
                }
            }
            post{
                failure{
                    echo "====++++'Build code binary' execution failed++++===="
                    echo "Try to run 'make code/compile' locally and make sure it pass"
                }
            }
        }

        stage("Build & push container image") {
            steps{
                dir("${env.CLONED_REPOSITORY_PATH}") {
                    // qe-pipeline-library step
                    dockerBuildAndPush(
                        credentialsId: "${env.CREDENTIALS_ID}",
                        containerRegistryServerName: "quay.io",
                        containerImageName: "${env.OPERATOR_CONTAINER_IMAGE_CANDIDATE_NAME}",
                        pathToDockerfile: "build/Dockerfile"
                    )
                }
            }
            post{
                failure{
                    echo "====++++'Build & push container image' execution failed++++===="
                }
            }
        }
        
        stage("Build test binary"){
            steps{
                dir("${env.CLONED_REPOSITORY_PATH}") {
                    script {
                        sh "make test/compile"
                    }
                }
            }
            post{
                failure{
                    echo "====++++'Build test binary' execution failed++++===="
                    echo "Try to run 'make test/compile' locally and make sure it pass"
                }
            }
        }
        stage("Test operator") {
            steps{
                dir("${env.CLONED_REPOSITORY_PATH}") {
                    // qe-pipeline-library step
                    runOperatorTestWithImage (
                        containerImageName: "${env.OPERATOR_CONTAINER_IMAGE_CANDIDATE_NAME}",
                        namespace: "${env.OPENSHIFT_PROJECT_NAME}",
                        namespacedManifest: true,
                        globalManifest: true
                    )
                }
            }
            post{
                failure{
                    echo "====++++Test operator execution failed++++===="
                }
            }
        }
        stage("Retag the image if the test passed and delete an old tag") {
            steps{
                // qe-pipeline-library step
                tagRemoteContainerImage(
                    credentialsId: "${env.CREDENTIALS_ID}",
                    sourceImage: "${env.OPERATOR_CONTAINER_IMAGE_CANDIDATE_NAME}",
                    targetImage: "${env.OPERATOR_CONTAINER_IMAGE_NAME}",
                    deleteOriginalImage: true
                )
            }
        }
        stage("Create a 'latest' tag from 'master'") {
            when {
                branch 'master'
            }
            steps{
                // qe-pipeline-library step
                tagRemoteContainerImage(
                    credentialsId: "${env.CREDENTIALS_ID}",
                    sourceImage: "${env.OPERATOR_CONTAINER_IMAGE_NAME}",
                    targetImage: "${env.OPERATOR_CONTAINER_IMAGE_NAME_latest}",
                    deleteOriginalImage: false
                )
            }
        }
    }
    post {
        failure {
            mail(
                to: 'stethoma@redhat.com',
                subject: 'MSS Operator build failed',
                body: "See the pipeline here: ${env.RUN_DISPLAY_URL}"
            )
        }
    }
}
