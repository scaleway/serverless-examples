# SNS Instances Notification System

This repository contains the source code for the this tutorial: [Creating a notification system with Scaleway SNS and Instances.](https://www.scaleway.com/en/docs/tutorials/sns-instances-notification-system)

## Requirements

To complete the actions presented below, you must have:

- A Scaleway account logged into the [console](https://console.scaleway.com)
- [Set up authentication for the Terraform provider](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs#authentication)
- An [SSH key](/identity-and-access-management/organizations-and-projects/how-to/create-ssh-key/)
- [Set up the Scaleway CLI](/developer-tools/scaleway-cli/quickstart/)


## Context

This example shows you how to set up a Notification System with Terraform across Instances using Scaleway's products, through a simulated CPU monitoring example.

## Setup

Once you have cloned this repository, export your public SSH key to a `TF_VAR_public_ssh_key` environement variable:

```
export TF_VAR_public_ssh_key=$(cat ~/.ssh/id_ed25519.pub)
```

Then you can run these commands:

```console
cd terraform
terraform init
terraform plan
terraform apply
```

## Testing

- Connect to the subscriber Instance with this URL ```http://<your_subscriber_ip_address>:8081``` (you can find the address on the Console Instances page)
- Click on `Confirm subscription`. If you get an error because the URL hasn't been received, you can reload the page, it should take less than 30s to appear. You will see a AWS xml page.
- Go back to the home page, and click on `Notifications`. The notifications received will be displayed here.
- Then connect to the publisher Instance with that URL ```http://<your_subscriber_ip_address>:8081```
- Go to the home page of your publisher server: ```http://<your_publisher_ip_address>:8081```
- Click on a CPU behavior, and check the notification page of your subscriber server. A notification about the behavior should have appeared.
- Once you're done testing, you can apply `terraform destroy` to clean and remove the project.
