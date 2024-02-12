resource "scaleway_job_definition" "main" {
  name         = "hello_jobs"
  cpu_limit    = 1000
  memory_limit = 1024

  image_uri    = docker_image.main.name

  command      = "sh hello.sh"

  env = {
    "MESSAGE" : "Hello from your Job!",
  }

  # TODO - enable once CRONs implemented
  # cron {
  #   schedule = "/5 * * * *"
  #   timezone = "Europe/Paris"
  # }
}
