variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "project_id" {
  type = string
}

variable "image_version" {
  type = string
  default = "0.0.3"
}

variable "region" {
  type = string
  default = "fr-par"
}

variable "inference_cron_schedule" {
  type = string
  default = "*/15 * * * *"
}

variable "hf_model_file_name" {
  type = string
  default = "llama-2-7b.Q4_0.gguf"
}

variable "hf_model_download_source" {
  type = string
  default = "https://huggingface.co/TheBloke/Llama-2-7B-GGUF/resolve/main/llama-2-7b.Q4_0.gguf"
}
