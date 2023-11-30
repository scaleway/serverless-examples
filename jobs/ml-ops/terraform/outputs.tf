output "inference_api_endpoint" {
  value = scaleway_container.inference_api_container.domain_name
}

output "inference_api_token" {
  value     = scaleway_container_token.inference_api_token.token
  sensitive = true
}
