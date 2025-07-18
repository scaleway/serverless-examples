output "metabase_container_url" {
  description = "The URL to connect to Metabase"
  value       = "https://${scaleway_container.main.domain_name}"
}
