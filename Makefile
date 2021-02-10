# The binary to build (just the basename).
BIN := $(shell basename $$PWD)

# Where to push the docker image.
REGISTRY ?= docker.io/absaoss

# Tag 
TAG ?= latest

# Image URL to use all building/pushing image targets
IMG ?= $(REGISTRY)/$(BIN)

# run all linters from .golangci.yaml; see: https://golangci-lint.run/usage/install/#local-installation
.PHONY: lint
lint:
	golangci-lint run

build:
	CGO_ENABLED=0 go build cmd/coredns.go

clean:
	go clean
	rm -f coredns

image: 
	docker build . -t ${IMG}:${TAG}

create-local-cluster:
	k3d cluster create coredns-crd -p "1053:30053/udp@server[0]" \
        --no-lb --k3s-server-arg "--no-deploy=traefik,servicelb,metrics-server"

import-image:
	k3d image import -c coredns-crd absaoss/k8s_crd:${TAG}

deploy-app: image import-image
	kubectl config use-context k3d-coredns-crd
	kubectl apply -f terratest/example/ns.yaml 
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/external-dns/master/docs/contributing/crd-source/crd-manifest.yaml
	helm repo add coredns https://coredns.github.io/helm
	helm repo update
	cd charts/coredns && helm dependency update
	helm upgrade -i coredns -n coredns charts/coredns \
		--set coredns.image.tag=${TAG}

.PHONY: terratest
terratest: deploy-app
	cd terratest/test/ && go test -v

