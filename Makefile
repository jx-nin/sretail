LOCAL_ENV_FILE ?= .env.local
LOCAL_COMPOSE = docker compose --env-file $(LOCAL_ENV_FILE) -f compose.yml -f compose.local.yml
KIND_CLUSTER ?= sretail
KIND_CONTEXT ?= kind-$(KIND_CLUSTER)
KIND_DOCKER_CONTEXT ?= default
KUSTOMIZE ?= kustomize
K8S_LOCAL_CATALOG_TAG ?= $(shell git rev-parse --verify HEAD)
K8S_LOCAL_CATALOG_IMAGE ?= sretail-catalog:$(K8S_LOCAL_CATALOG_TAG)
K8S_LOCAL_CATALOG_OVERLAY = deploy/k8s/overlays/local/catalog

.PHONY: local-config local-build local-up local-test local-verify local-logs local-down \
	k8s-local-config k8s-local-render k8s-local-verify-clean-tree k8s-local-set-image k8s-local-build k8s-local-load k8s-local-deploy \
	k8s-local-status

local-config:
	@test -f "$(LOCAL_ENV_FILE)" || { \
		echo "Missing $(LOCAL_ENV_FILE). Copy .env.local.example first."; \
		exit 1; \
	}
	$(LOCAL_COMPOSE) config

local-build:
	$(LOCAL_COMPOSE) build

local-up:
	$(LOCAL_COMPOSE) up --build --detach --wait

local-test:
	@set -eu; \
		gateway_address="$$( $(LOCAL_COMPOSE) port gateway 8080 )"; \
		catalog_address="$$( $(LOCAL_COMPOSE) port catalog 8081 )"; \
		test "$$(curl --fail --silent --show-error "http://$$gateway_address/health")" = '{"service":"gateway","status":"ok"}'; \
		test "$$(curl --fail --silent --show-error --output /dev/null --write-out '%{http_code}' "http://$$catalog_address/health/ready")" = "204"; \
		test "$$(curl --fail --silent --show-error "http://$$catalog_address/products")" = '[{"id":"prod-1","name":"Mechanical Keyboard","price_cents":12999},{"id":"prod-2","name":"Ergonomic Mouse","price_cents":10999}]'; \
		echo "Local Compose smoke tests passed."

local-verify: local-up local-test

local-logs:
	$(LOCAL_COMPOSE) logs --follow

local-down:
	$(LOCAL_COMPOSE) down

k8s-local-config:
	@test "$$(kubectl config current-context)" = "$(KIND_CONTEXT)" || { \
		echo "Current kubectl context must be $(KIND_CONTEXT)."; \
		exit 1; \
	}
	@DOCKER_CONTEXT=$(KIND_DOCKER_CONTEXT) kind get clusters | grep -Fx "$(KIND_CLUSTER)" >/dev/null || { \
		echo "Kind cluster $(KIND_CLUSTER) was not found in Docker context $(KIND_DOCKER_CONTEXT)."; \
		exit 1; \
	}
	kubectl kustomize $(K8S_LOCAL_CATALOG_OVERLAY) >/dev/null

k8s-local-render: k8s-local-config
	kubectl kustomize $(K8S_LOCAL_CATALOG_OVERLAY)

k8s-local-verify-clean-tree:
	@test -z "$$(git status --porcelain)" || { \
		echo "Working tree has uncommitted changes. Commit them before building an image tagged with HEAD."; \
		exit 1; \
	}

k8s-local-set-image:
	@command -v "$(KUSTOMIZE)" >/dev/null || { \
		echo "Missing Kustomize. Install the standalone kustomize CLI to update image tags."; \
		exit 1; \
	}
	cd $(K8S_LOCAL_CATALOG_OVERLAY) && $(KUSTOMIZE) edit set image registry.invalid/sretail/catalog=$(K8S_LOCAL_CATALOG_IMAGE)

k8s-local-build:
	docker --context $(KIND_DOCKER_CONTEXT) build -f services/catalog/Dockerfile -t $(K8S_LOCAL_CATALOG_IMAGE) .

k8s-local-load:
	$(MAKE) k8s-local-verify-clean-tree
	$(MAKE) k8s-local-config
	$(MAKE) k8s-local-set-image
	$(MAKE) k8s-local-config
	$(MAKE) k8s-local-build
	DOCKER_CONTEXT=$(KIND_DOCKER_CONTEXT) kind load docker-image $(K8S_LOCAL_CATALOG_IMAGE) --name $(KIND_CLUSTER)

k8s-local-deploy: k8s-local-load
	kubectl apply -k $(K8S_LOCAL_CATALOG_OVERLAY)
	kubectl rollout status deployment/catalog --timeout=60s

k8s-local-status: k8s-local-config
	kubectl get deployment,pods,service,endpointslice -l app.kubernetes.io/name=catalog,app.kubernetes.io/instance=sretail
