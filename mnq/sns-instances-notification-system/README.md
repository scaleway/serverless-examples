# SNS Instances Notification System

This repository contains the source code for the this tutorial: [Creating a notification system with Scaleway SNS and Instances.](https://github.com/scaleway/docs-content/blob/int-add-mnq-tuto/tutorials/sns-instances-notification-system/index.mdx)

## Requirements

This example assumes that you are familiar with:

- how System Notification Service work. You can find out more in the [SNS Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/messaging/quickstart/)
- how Instances work. If needed, you can visit the [Instances Quickstart documentation](https://www.scaleway.com/en/docs/compute/instances/quickstart/).
- how to create a SSH key. If not done already, you can follow the first 6 steps of [this walkthrough](https://www.scaleway.com/en/docs/identity-and-access-management/organizations-and-projects/how-to/create-ssh-key/#how-to-upload-the-public-ssh-key-to-the-scaleway-interface).

## Context

This example shows you how to set up a Notification System with Terraform across Instances using Scaleway's products, through a simulated CPU monitoring example.

## Setup

Once you have cloned this repository you can run:

```console
terraform init
terraform plan
terraform apply
```
You will be prompted to enter your public ssh. You will find them in the `~/.ssh/` folder:

```console
cat ~/.ssh/*.pub
```

## Testing

Connect to the subscriber Instance with this URL ```http://<your_subscriber_ip_address>:8081``` (you can find the address on the Console Instances page)
Click on `Confirm subscription`. If you have an error because the URL hasn't been received, you can reload the page. You will see a AWS xml page.
Go back to the home page, and click on `Notifications`. The notifications received will be displayed here.
Then connect to the publisher Instance with that URL ```http://<your_subscriber_ip_address>:8081```
Go to the home page of your publisher server: ```http://<your_publisher_ip_address>:8081```
Click on a CPU behavior, and check the notification page of your subscriber server. A notification about the behavior should have appeared.
Once you're done testing, you can apply `terraform destroy` to clean and remove the project.