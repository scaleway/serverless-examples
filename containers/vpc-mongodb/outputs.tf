output "mongodb_instance_public_endpoint" {
  description = "The public endpoint of the MongoDB instance"
  value       = scaleway_mongodb_instance.main.public_network[0].dns_record
}

output "server_container_endpoint" {
  description = "The endpoint of the server container"
  value       = "https://${scaleway_container.main.domain_name}"
}
