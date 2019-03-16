APP_NAME = mobile-security-service-operator
ORG_NAME = aerogear
PKG = github.com/$(ORG_NAME)/$(APP_NAME)
APP_FILE=./cmd/manager/main.go
BIN_DIR := $(GOPATH)/bin
BINARY ?= mobile-security-service-operator
TAG= 0.1.0
DOCKER-ORG=cmacedo
DOCKER-REPO=mobile-security-service-operator

# This follows the output format for goreleaser
BINARY_LINUX_64 = ./dist/linux_amd64/$(BINARY)

LDFLAGS=-ldflags "-w -s -X main.Version=${TAG}"

.PHONY: deploy
deploy:
	@echo Deploying Mobile Security Service Operator:
	oc create namespace mobile-security-service
	oc create -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	oc create -f deploy/cluster_role.yaml
	oc create -f deploy/cluster_role_binding.yaml
	oc create -f deploy/service_account.yaml
	oc create -f deploy/operator.yaml

.PHONY: undeploy
undeploy:
	@echo Undeploy Mobile Security Service Operator:
	oc delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_crd.yaml
	oc delete -f deploy/cluster_role.yaml
	oc delete -f deploy/cluster_role_binding.yaml
	oc delete -f deploy/service_account.yaml
	oc delete -f deploy/operator.yaml
	oc delete namespace mobile-security-service

.PHONY: deploy-app
deploy-app:
	@echo Deploying Mobile Security Service into project:
	oc apply -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml

.PHONY: undeploy-app
undeploy-app:
	@echo Undeploying Mobile Security Service from the project:
	oc delete -f deploy/crds/mobile-security-service_v1alpha1_mobilesecurityservice_cr.yaml

.PHONY: build
build:
	@echo Buinding operator with the tag $(TAG):
	operator-sdk build cmacedo/mobile-security-service-operator:$(TAG)

.PHONY: publish
publish:
	@echo Publishing operator in $(DOCKER-ORG)/$(DOCKER-REPO) with the tag $(TAG):
	docker push $(DOCKER-ORG)/$(DOCKER-REPO):$(TAG)