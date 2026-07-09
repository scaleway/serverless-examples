output "container_public_endpoint" {
  description = "Public HTTPS endpoint of the deployed container"
  value       = "https://${scaleway_container.main.domain_name}"
}
