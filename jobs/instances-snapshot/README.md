# Serverless Jobs for automatic Instance snapshot

This project shows how it's possible to automate tasks using Serverless Jobs.

This example is very simple, it generates snapshots of your desired Instance.

# Set-up

## Requirements

- Scaleway Account
- Docker daemon running to build the image
- Container registry namespace created, for this example we assume that your namespace name is `jobs-snapshot`: [doc here](https://www.scaleway.com/en/docs/containers/container-registry/how-to/create-namespace/)
- API keys generated, Access Key and Secret Key [doc here](https://www.scaleway.com/en/docs/iam/how-to/create-api-keys/)

## Step 1 : Build and push to Container registry

Serverless Jobs, like Serverless Containers (which are suited for HTTP applications), works
with containers. So first, use your terminal reach this folder and run the following commands:

```shell
# First command is to login to container registry, you can find it in Scaleway console
docker login rg.fr-par.scw.cloud/jobs-snapshot -u nologin --password-stdin <<< "$SCW_SECRET_KEY"

# Here we build the image to push
docker build -t rg.fr-par.scw.cloud/jobs-snapshot/jobs-snapshot:v1 .

## TIP: for Apple Silicon or other ARM processors, please use the following command as Serverless Jobs supports amd64 architecture
# docker buildx build --platform linux/amd64 -t rg.fr-par.scw.cloud/jobs-snapshot/jobs-snapshot:v1 .

# Push the image online to be used on Serverless Jobs
docker push rg.fr-par.scw.cloud/jobs-snapshot/jobs-snapshot:v1
```
> [!TIP]
> As we do not expose a web server and we do not require features such as auto-scaling, Serverless Jobs are perfect for this use case.

To check if everyting is ok, on the Scaleway Console you can verify if your tag is present in Container Registry.

## Step 2: Creating the Job Definition

On Scaleway Console on the following link you can create a new Job Definition: https://console.scaleway.com/serverless-jobs/jobs/create?region=fr-par

1. On Container image, select the image you created in the step before.
1. You can set the image name to something clear like `jobs-snapshot` too.
1. For the region you can select the one you prefer :)
1. Regarding the resources you can keep the default values, this job is fast and do not require specific compute power or memory.
1. To schedule your job for example every night at 2am, you can set the cron to `0 2 * * *`.
1. Important: advanced option, you need to set the following environment variables:

> [!TIP]
> For sensitive data like `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` we recommend to inject them via Secret Manager, [more info here](https://www.scaleway.com/en/docs/serverless/jobs/how-to/reference-secret-in-job/).

- `INSTANCE_ID`: grab the instance ID you want to create snapshots from
- `SCW_ZONE`: you need to give the ZONE of you instance, like `fr-par-2`
- `SCW_ACCESS_KEY`: your access key
- `SCW_SECRET_KEY`: your secret key
- `SCW_DEFAULT_ORGANIZATION_ID`: your organzation ID
- `SCW_DEFAULT_PROJECT_ID`: your project ID

* Then click "create job"

## Step 3: Run the job

On your created Job Definition, just click the button "Run Job" and within seconds it should be successful.

## Troubleshooting

If your Job Run state goes in error, you can use the "Logs" tab in Scaleway Console to get more informations about the error.

# Possible improvements

You can exercice by adding the following features:

- Instead of managing a single instance, make it account wide
- Add disk backups
- Add alerts if something goes wrong
- Use secret manager instead of job environment variables

# Additional content

- [Jobs Documentation](https://www.scaleway.com/en/docs/serverless/jobs/how-to/create-job-from-scaleway-registry/)
- [Other methods to deploy Jobs](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/deploy-job/)
- [Secret key / access key doc](https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/)
- [CRON schedule help](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/cron-schedules/)
