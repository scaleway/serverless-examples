
.PHONY: tf-init
tf-init:
	cd terraform && terraform init

.PHONY: tf-plan
tf-plan:
	cd terraform && terraform plan -var-file=vars/main.tfvars

.PHONY: tf-apply
tf-apply:
	cd terraform && terraform apply -auto-approve -var-file=vars/main.tfvars

.PHONY: tf-destroy
tf-destroy:
	cd terraform && terraform destroy -var-file=vars/main.tfvars
