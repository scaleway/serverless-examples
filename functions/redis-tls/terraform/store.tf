resource "scaleway_redis_cluster" "weather_store" {
  name         = "serverless-weather-redis-example"
  version      = "7.0.5"
  node_type    = "RED1-MICRO"
  user_name    = var.redis_user
  password     = var.redis_password
  tags         = ["serverless-examples", "weather", "redis"]
  cluster_size = 1
  tls_enabled  = "true"

  // Due to the nature of serverless functions, 
  // the IPs we use are unpredictable. 
  // Here we will restrict to IPs from Scaleway's AS (http://as12876.net/)
  dynamic "acl" {
    for_each = [
      "62.210.0.0/16",
      "195.154.0.0/16",
      "212.129.0.0/18",
      "62.4.0.0/19",
      "212.83.128.0/19",
      "212.83.160.0/19",
      "212.47.224.0/19",
      "163.172.0.0/16",
      "51.15.0.0/16",
      "151.115.0.0/16",
      "51.158.0.0/15",
    ]
    content {
      ip          = acl.value
      description = "Allow Scaleway IPs"
    }
  }
}
