# Image URL to use all building/pushing image targets
IMG_REPOSITORY ?= quay.io/fmatouschek/kubevirt-addon-manager
IMG_TAG ?= latest
IMG ?= ${IMG_REPOSITORY}:${IMG_TAG}

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

LOCALBIN ?= $(PROJECT_DIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

ifeq (,$(shell which oc))
KUBECTL = kubectl
else
KUBECTL = oc
endif

.PHONY:vendor
vendor:
	go mod vendor
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: vet golangci-lint ## Lint source code.
	$(GOLANGCILINT) run --timeout 4m0s ./...

# Build manager binary
.PHONY: manager
manager:
	go build -mod vendor -o bin/manager main.go

.PHONY: run
run: manager
	$(LOCALBIN)/manager controller

# Build the container image
container-build:
	podman build -t ${IMG} .

# Push the container image
container-push:
	podman push ${IMG}

.PHONY: deploy
deploy: kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cp config/manager/deployment.yaml config/manager/deployment.yaml.tmp
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -
	mv config/manager/deployment.yaml.tmp config/manager/deployment.yaml

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete -f -

.PHONY: golangci-lint
GOLANGCILINT := $(LOCALBIN)/golangci-lint
GOLANGCI_URL := https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
golangci-lint: $(GOLANGCILINT) ## Download golangci-lint
$(GOLANGCILINT): $(LOCALBIN)
	curl -sSfL $(GOLANGCI_URL) | sh -s -- -b $(LOCALBIN)

.PHONY: kustomize
KUSTOMIZE := $(LOCALBIN)/kustomize
kustomize: $(KUSTOMIZE) ## Download ginkgo
$(KUSTOMIZE): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/kustomize/kustomize/v4@latest
