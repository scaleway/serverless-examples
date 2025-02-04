# Serverless Jobs for cleaning old snapshots

This project shows how it's possible to automate tasks using Serverless Jobs.

This simple example shows how to clean up snapshots after X days, it's useful to avoid a growing list of snapshots.

# Set-up

## Requirements

- Scaleway Account
- Docker daemon running to build the image
- Container registry namespace created, for this example we assume that your namespace name is `jobs-snapshot-cleaner`: [doc here](https://www.scaleway.com/en/docs/containers/container-registry/how-to/create-namespace/)
- API keys generated, Access Key and Secret Key [doc here](https://www.scaleway.com/en/docs/iam/how-to/create-api-keys/)

## Step 1 : Build and push to Container registry

Serverless Jobs, like Serverless Containers (which are suited for HTTP applications), works
with containers. So first, use your terminal reach this folder and run the following commands:

```shell
# First command is to login to container registry, you can find it in Scaleway console
docker login rg.fr-par.scw.cloud/jobs-snapshot-cleaner -u nologin --password-stdin <<< "$SCW_SECRET_KEY"

# Here we build the image to push
docker build -t rg.fr-par.scw.cloud/jobs-snapshot-cleaner/jobs-snapshot-cleaner:v1 .

## TIP: for Apple Silicon or other ARM processors, please use the following command as Serverless Jobs supports amd64 architecture
# docker buildx build --platform linux/amd64 -t rg.fr-par.scw.cloud/jobs-snapshot-cleaner/jobs-snapshot-cleaner:v1 .

# Push the image online to be used on Serverless Jobs
docker push rg.fr-par.scw.cloud/jobs-snapshot-cleaner/jobs-snapshot-cleaner:v1
```
> [!TIP]
> As we do not expose a web server and we do not require features such as auto-scaling, Serverless Jobs are perfect for this use case.

To check if everyting is ok, on the Scaleway Console you can verify if your tag is present in Container Registry.

## Step 2: Creating the Job Definition

On Scaleway Console on the following link you can create a new Job Definition: https://console.scaleway.com/serverless-jobs/jobs/create?region=fr-par

1. On Container image, select the image you created in the step before.
1. You can set the image name to something clear like `jobs-snapshot-cleaner` too.
1. For the region you can select the one you prefer :)
1. Regarding the resources you can keep the default values, this job is fast and do not require specific compute power or memory.
1. To schedule your job for example every two days at 2am, you can set the cron to `0 2 */2 * *`.
1. Important: advanced option, you need to set the following environment variables:

> [!TIP]
> For sensitive data like `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` we recommend to inject them via Secret Manager, [more info here](https://www.scaleway.com/en/docs/serverless/jobs/how-to/reference-secret-in-job/).

- `SCW_DELETE_AFTER_DAYS`: number of days after the snapshots will be deleted
- `SCW_PROJECT_ID`: project you want to clean up
- `SCW_ZONE`: you need to give the ZONE of your snapshot you want to clean, like `fr-par-2`
- `SCW_ACCESS_KEY`: your access key
- `SCW_SECRET_KEY`: your secret key
- `SCW_DEFAULT_ORGANIZATION_ID`: your organzation ID

* Then click "create job"

## Step 3: Run the job

On your created Job Definition, just click the button "Run Job" and within seconds it should be successful.

## Troubleshooting

If your Job Run state goes in error, you can use the "Logs" tab in Scaleway Console to get more informations about the error.

# Possible improvements

You can exercice by adding the following features:

- Add tags to exclude
- Add alerts if a Job goes in error
- Use Secret Manager instead of job environment variables
- Support multiple zones dans projects

# Additional content

- [Jobs Documentation](https://www.scaleway.com/en/docs/serverless/jobs/how-to/create-job-from-scaleway-registry/)
- [Other methods to deploy Jobs](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/deploy-job/)
- [Secret key / access key doc](https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/)
- [CRON schedule help](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/cron-schedules/)
