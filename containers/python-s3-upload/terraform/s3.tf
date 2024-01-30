resource "scaleway_object_bucket" "main" {
  name = "python-s3-example-${random_string.suffix.result}"
}
