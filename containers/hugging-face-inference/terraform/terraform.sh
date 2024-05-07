#!/bin/bash

set -e

# Common environment variables
export TF_VAR_access_key=${SCW_ACCESS_KEY} \
       TF_VAR_secret_key=${SCW_SECRET_KEY} \
       TF_VAR_project_id=${SCW_PROJECT_ID}

# Associative list of models to deploy using json data
declare -A hf_models
eval "$(jq -r '.[]|.[]|"hf_models[\(.file)]=\(.source)"' hf-models.json)"

# Login to docker Scaleway's registry
docker login "rg.$REGION.scw.cloud" -u nologin --password-stdin <<< "$SCW_SECRET_KEY"

# Initialize, plan, and deploy each model in a Terraform workspace
apply() {
       terraform init
       for model_file_name in "${!hf_models[@]}";
       do
         terraform workspace select -or-create $model_file_name
         export TF_VAR_hf_model_file_name=$model_file_name \
                TF_VAR_hf_model_download_source=${hf_models[$model_file_name]}
         terraform plan
         terraform apply -auto-approve
       done
}

# Destroy resources of each Terraform workspace
destroy(){
       for model_file_name in "${!hf_models[@]}";
       do
         terraform workspace select $model_file_name
         export TF_VAR_hf_model_file_name=$model_file_name \
                TF_VAR_hf_model_download_source=${hf_models[$model_file_name]}
         terraform destroy -auto-approve
       done
}

# Script actions per flag
while getopts "ad" option; do
  case $option in
    a)
      echo "deploying models"
      apply
      ;;
    d)
      echo "destroying models"
      destroy
      ;;
    *)
      echo "flag is not provided"
      exit 1
  esac
done