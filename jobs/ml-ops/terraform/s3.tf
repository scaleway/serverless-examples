resource "scaleway_object_bucket" "data_store" {
  name = "data-store-${random_string.random_suffix.result}"
}

resource "scaleway_object_bucket" "model_registry" {
  name = "model-registry-${random_string.random_suffix.result}"
}

resource "scaleway_object_bucket" "performance_monitoring_record" {
  name = "performance-monitoring-${random_string.random_suffix.result}"
}
