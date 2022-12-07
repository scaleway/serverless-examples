DOCKER_REGISTRY=rg.fr-par.scw.cloud
DOCKER_REGISTRY_NS=cors-demo

GATEWAY_VERSION=0.0.6
SERVER_VERSION=0.0.1

GATEWAY_TAG=$(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_NS)/gateway:$(GATEWAY_VERSION)
SERVER_TAG=$(DOCKER_REGISTRY)/$(DOCKER_REGISTRY_NS)/server:$(SERVER_VERSION)

.PHONY: create-namespace
create-namespace:
	scw registry namespace create name=$(DOCKER_REGISTRY_NS)

.PHONY: docker-login
docker-login:
	scw registry login

.PHONY: build-server
build-server:
	cd server && docker build -t $(SERVER_TAG) .

.PHONY: build-gateway
build-gateway:
	cd gateway && docker build -t $(GATEWAY_TAG) .

.PHONY: push-server
push-server:
	docker push $(SERVER_TAG)

.PHONY: push-gateway
push-gateway:
	docker push $(GATEWAY_TAG)

.PHONY: tf-init
tf-init:
	cd terraform && terraform init

.PHONY: tf-plan
tf-plan:
	cd terraform && terraform plan -var-file=vars/main.tfvars

.PHONY: tf-apply
tf-apply:
	cd terraform && terraform apply -auto-approve -var-file=vars/main.tfvars

