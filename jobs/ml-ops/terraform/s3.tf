resource "scaleway_object_bucket" "main" {
  name = "ml-ops-${random_string.random_suffix.result}"
}
