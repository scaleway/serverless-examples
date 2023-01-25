# Terraform Python example

In this tutorial you will discover an example of Instance automation using Python and Terraform:

* Automatically shut down / start instances

## Requirements

* You have an account and are logged into the [Scaleway console](https://console.scaleway.com)
* You have [generated an API key](/console/my-project/how-to/generate-api-key/)
* You have [Python](https://www.python.org/) installed on your machine
* You have [Terraform](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs) installed on your machine

## Context

In this tutorial, we will simulate a project with a production environment that will be running all the time and a development environment that will be turn off on week-ends to save costs.

## Initialize your Terraform project

1. Create a 'Terraform' folder to store your configuration as explained in the terraform documentation.
2. Create 5 files to configure your infrastructure:
  a. 'main.tf': will contain the main set of configurations for your project. Here, it will be our instance
  b. 'provider.tf': Terraform relies on plugins called “providers” to interact with remote systems
  c. 'backend.tf': each Terraform configuration can specify a backend, which defines where the state file of the current infrastructure will be stored. Thanks to this file, Terraform keeps track of the managed resources. This state can be stored locally or remotely. Configuring a remote backend allows multiple people to work on the same infrastructure
  d. 'variables.tf': will contain the variable definitions for your project. Since all Terraform values must be defined, any variables that are not given a default value will become required arguments
  e. 'terraform.tfvars': allows you to set the actual value of the variables
3. Create the following folder:
  a. 'function': to store your function code
  b. 'files': to temporarily store your zip function code

   Your folder should now look like this:

```bash
Terraform
| -- files
| -- function
-- main.tf
-- backend.tf
-- provider.tf
-- terraform.tfvars
-- variables.tf
```

4. Edit the 'backend.tf' file to enable distant configuration backup

```hcl
terraform {
  backend "s3" {
    bucket                      = "XXXXXXXXX"
    key                         = "terraform.tfstate"
    region                      = "fr-par"
    endpoint                    = "https://s3.fr-par.scw.cloud"
    skip_credentials_validation = true
    skip_region_validation      = true
  }
}

/*
For the credentials part:
==> Create a ~/.aws/credentials:
[default]
aws_access_key_id=<SCW_ACCESS_KEY>
aws_secret_access_key=<SCW_SECRET_KEY>
region=fr-par
*/
```

5. Edit the 'provider.tf' file and add Scaleway as a provider

```hcl
terraform {
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "2.9.1"
    }
  }
  required_version = ">= 0.13"
}
```

6. Specify the following variable in the 'variables.tf' file

```hcl
variable "zone" {
  type = string
}

variable "region" {
  type = string
}

variable "env" {
  type = string
}

variable "project_id" {
  type        = string
  description = "Your project ID"
}
```

7. Add the variables value to 'terraform.tfvars'

```ini
zone                           = "fr-par-1"
region                         = "fr-par"
env                            = "dev"
project_id                     = "Your project ID"
```

## Writing the code

For this example, we will use the native python library `urllib`, which will enable us to keep the package slim.
All information about the Instance API can be found in the [developers documentation](https://developers.scaleway.com/en/products/instance/api/#get-2c1c6f).

1. In the 'function folder", create a `handler.py` file as follows:

```py
import os
from urllib import request,parse,error
import json

auth_token=os.environ['X-AUTH-TOKEN']

def handle(event, context):
    ## get information from cron
    event_body=eval(event["body"])
    zone=event_body["zone"]
    server_id=event_body["server_id"]
    action=event_body["action"] #action should be "poweron" or "poweroff"
    
    #create request
    url=f"https://api.scaleway.com/instance/v1/zones/{zone}/servers/{server_id}/action"
    data=json.dumps({"action":action}).encode('ascii')
    req = request.Request(url, data=data,  method="POST")
    req.add_header('Content-Type', 'application/json')
    req.add_header('X-Auth-Token',auth_token)
    
    #Sending request to Instance API
    try: 
        res=request.urlopen(req).read().decode()
    except error.HTTPError as e:
        res=e.read().decode()
        
    return {
        "body": json.loads(res),
        "statusCode": 200,
    }
```

## Configure your infrastructure

Edit 'main.tf' to add:

* A production instance using a GP1-S named "Prod"
  
```hcl
## Configuring Producion environment 
resource "scaleway_instance_ip" "public_ip-prod" {
    project_id = var.project_id
}

resource "scaleway_instance_server" "scw-instance-prod" {
  name="prod"
  project_id = var.project_id
  type  = "GP1-S"
  image = "ubuntu_focal"

  tags = ["terraform instance", "scw-instance", "production"]

  ip_id = scaleway_instance_ip.public_ip-prod.id

  root_volume {
    # The local storage of a DEV1-L instance is 80 GB, subtract 30 GB from the additional l_ssd volume, then the root volume needs to be 50 GB.
    size_in_gb = 200
  }
}
```

* A development instance using a DEV1-L named "Dev"
  
```hcl
## Configuring Development environment that will be automatically turn off on week-ends and turn on monday mornings
resource "scaleway_instance_ip" "public_ip-dev" {
    project_id = var.project_id
}

resource "scaleway_instance_server" "scw-instance-dev" {
  name="dev"
  project_id = var.project_id
  type  = "DEV1-L"
  image = "ubuntu_focal"

  tags = ["terraform instance", "scw-instance", "dev"]

  ip_id = scaleway_instance_ip.public_ip-dev.id


  root_volume {
    size_in_gb = 80
  }
}
```

* An IAM application and API key so that the function can manage instances

```hcl
# Create an IAM application to provide the a scoped token for the function
resource "scaleway_iam_application" "instance_power_toggler" {
  name        = "Instance power toggler"
  description = "IAM application for the function to toggle instances"
}

resource "scaleway_iam_policy" "can_manage_instances" {
  name        = "can-manage-instances"
  description = "policy to manage instances in a specific project"
  application_id = scaleway_iam_application.instance_power_toggler.id
  rule {
    project_ids = [var.project_id]
    # Currently the least possible privilege to manage instances
    permission_set_names = ["InstancesFullAccess"]
  }
}

resource "scaleway_iam_api_key" "instance_toggler_key" {
  application_id = scaleway_iam_application.instance_power_toggler.id
  description    = "API key for the function to toggle instances"
}
```

* A function that will run code you've just written
  
```hcl
# Creating function code archive that will then be updated
data "archive_file" "source_zip" {
  type             = "zip"
  source_dir       = "${path.module}/function"
  output_path      = "${path.module}/files/function.zip"
}

# Creating the function namespace
resource "scaleway_function_namespace" "main"  {
  name        = "instance-management"
  description = "namespace to gather all functions dedicated to instance management"
  project_id = var.project_id
}

# Creating the function
resource "scaleway_function" "main" {
  namespace_id = scaleway_function_namespace.main.id
  name         = "instancewake"
  runtime      = "python310"
  handler      = "handler.handle"
  privacy      = "public"
  zip_file = data.archive_file.source_zip.output_path # this enable to automatically zip your code and each time you edit it
  zip_hash = filesha256(data.archive_file.source_zip.output_path)
  deploy = true
  max_scale    = "5"
  environment_variables = {
      "X-AUTH-TOKEN" = scaleway_iam_api_key.instance_toggler_key.secret_key
  }
}
```

* A cronjob attached to the function to turn your function off every Friday evening

```hcl
# Adding a first cron to turn off the instance every friday evening (11:30 pm)
resource "scaleway_function_cron" "turn-off" {
    function_id = scaleway_function.main.id
    schedule = "30 23 * * 5"
    args = jsonencode({
        "zone": scaleway_instance_server.scw-instance-dev.zone,
        "server_id": regex("/([^/]+$)", scaleway_instance_server.scw-instance-dev.id)[0], # We use the dev instance id and strip it from the region
        "action":"poweroff"
        }
    )
}
```

* A cronjob attached to the function to turn your function on every Monday mornings

```hcl
# Adding a second cron to turn on the instance every monday morning (7:00 am)
resource "scaleway_function_cron" "turn-on" {
    function_id = scaleway_function.main.id
    schedule = "0 7 * * 1"
    args = jsonencode({
        "zone": scaleway_instance_server.scw-instance-dev.zone,
        "server_id": regex("/([^/]+$)", scaleway_instance_server.scw-instance-dev.id)[0],
        "action":"poweron"
        }
    )
}
```

## Deploy your infrastructure

Now that everything is set up, deploy everything using Terraform

1. Add your Scaleway credentials to your environment variables

```bash
export SCW_ACCESS_KEY="my-access-key"
export SCW_SECRET_KEY="my-secret-key"
```

2. Initialize Terraform

```bash
terraform init 
```

3. Let terraform verify your configuration

```bash
terraform plan
```

4. Deploy your infrastructure
  
```bash
terraform apply
````
