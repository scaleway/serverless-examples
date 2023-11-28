resource "scaleway_object_bucket" "data_store_bucket" {
  name = "data-store-1"
}

resource "scaleway_object_bucket" "model_registry_bucket" {
  name = "model-registry"
}

resource "scaleway_object_bucket" "performance_monitoring_bucket" {
  name = "performance-monitoring"
}