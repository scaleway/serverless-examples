terraform {
  # CONNECTS TO SCALEWAY CLOUD PROVIDER 
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

provider "scaleway" {
}

# CREATE PROJECT SNS-TUTORIAL
resource "scaleway_account_project" "sns_tutorial" {
  name = "sns-tutorial"
}

# TAKES YOUR PUBLIC SSH-KEY TO ACCESS INSTANCES
resource "scaleway_iam_ssh_key" "main" {
  name        = "sns-tutorial-public-ssh-key"
  project_id = scaleway_account_project.sns_tutorial.id
  public_key = var.public_ssh_key
}

# =====================   SNS   =====================

# ACTIVATE SNS
resource "scaleway_mnq_sns" "main" {
  project_id = scaleway_account_project.sns_tutorial.id
  region = "fr-par"
}

# CREATE CREDENTIALS
resource "scaleway_mnq_sns_credentials" "main" {
  project_id = scaleway_account_project.sns_tutorial.id
  permissions {
    can_manage = true // to set up the topic subject, the subscription to the topic
    can_receive = true // to subscribe a topic
    can_publish = true // to publish messages to the topic
  }
  # Wait for activation completion before creating the credentials.
  # If the credentials are created before the activation is completed, the project cannot be destroyed with terraform
  depends_on = [
    scaleway_mnq_sns.main
  ]
}

# CREATE TOPIC
resource "scaleway_mnq_sns_topic" "topic" {
  project_id = scaleway_account_project.sns_tutorial.id
  name = "sns-tutorial-topic"
  access_key = scaleway_mnq_sns_credentials.main.access_key
  secret_key = scaleway_mnq_sns_credentials.main.secret_key
}

# ==========   INSTANCES SECURITY GROUP   ===========

resource "scaleway_instance_security_group" "sns_www" {
  project_id = scaleway_account_project.sns_tutorial.id
  inbound_default_policy  = "drop"
  outbound_default_policy = "accept"

  # FOR SSH CONNEXIONS
  inbound_rule {
    action = "accept"
    port = "22"
  }

  # FOR HTTP CONNEXIONS
  inbound_rule {
    action = "accept"
    port = "8081"
  }
}

# =============   SUBSCRIBER SERVER   ===============

resource "scaleway_instance_ip" "subscriber_public_ip" {
  project_id = scaleway_account_project.sns_tutorial.id
  zone = "fr-par-1"
}

resource "scaleway_instance_server" "subscriber_sns_tuto_instance" {
  project_id = scaleway_account_project.sns_tutorial.id
  name = "suscriber-server"
  type  = "PLAY2-PICO"
  image = "debian_bookworm"
  ip_id = scaleway_instance_ip.subscriber_public_ip.id
  security_group_id= scaleway_instance_security_group.sns_www.id

  # USER DATA TO RUN THE SERVER AT START-UP
  user_data = {
    cloud-init = <<-EOF
    #cloud-config
    runcmd:
      - |
        apt-get update && apt-get install -y docker.io
        systemctl start docker
        systemctl enable docker
        docker pull rg.fr-par.scw.cloud/sns-x-instance-tutorial/subscriber-server:1.0
        docker run -d --restart=always \
          --name subscriber-server \
          -p 8081:8081 \
          rg.fr-par.scw.cloud/sns-x-instance-tutorial/subscriber-server:1.0
  EOF
  }
}

# TO WAIT FOR THE SUBSCRIBER PORT 8081 TO BE OPENED TO CONTINUE
resource "terraform_data" "bootstrap" {
  triggers_replace = [
    scaleway_instance_server.subscriber_sns_tuto_instance.id,
  ]

  provisioner "local-exec" {
    command = <<EOF
    #!/bin/bash
    TIMEOUT=180
    START_TIME=$(date +%s)
    while [ $(($(date +%s) - $START_TIME)) -lt $TIMEOUT ]; do
      nc -z ${scaleway_instance_server.subscriber_sns_tuto_instance.public_ip} 8081 && exit 0
      sleep 1
    done
    echo "Timeout reached"
    exit 1
    EOF
  }
}

# CREATE SNS TOPIC SUBSCRIPTION
resource scaleway_mnq_sns_topic_subscription main {
  project_id = scaleway_account_project.sns_tutorial.id
  access_key = scaleway_mnq_sns_credentials.main.access_key
  secret_key = scaleway_mnq_sns_credentials.main.secret_key
  topic_id = scaleway_mnq_sns_topic.topic.id
  protocol = "http"
  endpoint = "http://${scaleway_instance_ip.subscriber_public_ip.address}:8081/notifications"
}

# =================   PUBLISHER SERVER   =================

resource "scaleway_instance_ip" "publisher_public_ip" {
  project_id = scaleway_account_project.sns_tutorial.id
  zone = "fr-par-1"
}

resource "scaleway_instance_server" "publisher_sns_tuto_instance" {
  project_id = scaleway_account_project.sns_tutorial.id
  name       = "publisher-server"
  type       = "PLAY2-PICO"
  image      = "debian_bookworm"
  ip_id      = scaleway_instance_ip.publisher_public_ip.id
  security_group_id= scaleway_instance_security_group.sns_www.id

  # USER DATA TO RUN THE SERVER AT START-UP
  user_data = {
    cloud-init = <<-EOF
    #cloud-config
    write_files:
      - content: |
          #!/bin/bash
          export TOPIC_ARN="${scaleway_mnq_sns_topic.topic.arn}"
          export AWS_ACCESS_KEY="${scaleway_mnq_sns_credentials.main.access_key}"
          export AWS_SECRET_KEY="${scaleway_mnq_sns_credentials.main.secret_key}"
        owner: root:root
        path: /etc/profile.d/publisher-server_env.sh
        permissions: '0755'
    runcmd:
      - apt-get update && apt-get install -y docker.io
      - systemctl start docker
      - systemctl enable docker
      - chmod +x /etc/profile.d/publisher-server_env.sh
      # Using 'source' command doesn't work in cloud-init runcmd because each command runs in a separate shell.
      # this is why environment variables will be passed directly in the docker run command
      - docker pull rg.fr-par.scw.cloud/sns-x-instance-tutorial/publisher-server:1.0
      - |
        . /etc/profile.d/publisher-server_env.sh
        docker run -d --restart=always \
          --name publisher-server \
          -p 8081:8081 \
          -e TOPIC_ARN=$TOPIC_ARN \
          -e AWS_ACCESS_KEY=$AWS_ACCESS_KEY \
          -e AWS_SECRET_KEY=$AWS_SECRET_KEY \
          rg.fr-par.scw.cloud/sns-x-instance-tutorial/publisher-server:1.0
  EOF
  }
}
