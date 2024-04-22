#!/bin/bash

set -e

export SCW_ACCESS_KEY=${SCW_ACCESS_KEY} \
       SCW_SECRET_KEY=${SCW_SECRET_KEY} \
       SCW_PROJECT_ID=${SCW_PROJECT_ID}

declare -A hf_models 

hf_models["llama-2-7b.Q2_K.gguf"]="https://huggingface.co/TheBloke/Mistral-7B-Instruct-v0.2-GGUF/resolve/main/mistral-7b-instruct-v0.2.Q8_0.gguf"
hf_models["mistral-7b-instruct-v0.2.Q2_K.gguf"]="https://huggingface.co/TheBloke/Mistral-7B-Instruct-v0.2-GGUF/resolve/main/mistral-7b-instruct-v0.2.Q8_0.gguf"

terraform init

for model_file_name in "${!hf_models[@]}";
do
  terraform workspace new $model_file_name
  export TF_VAR_hf_model_file_name=$model_file_name \
         TF_VAR_hf_model_download_source=${hf_models[$model_file_name]} \
         TF_VAR_access_key=$SCW_ACCESS_KEY \
         TF_VAR_secret_key=$SCW_SECRET_KEY \
         TF_VAR_project_id=$SCW_PROJECT_ID
  terraform plan -var-file=testing.tfvars
  terraform apply -var-file=testing.tfvars -auto-approve
done
