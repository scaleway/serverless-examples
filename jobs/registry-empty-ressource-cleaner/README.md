# Scaleway Container Registry Cleaner

This project helps you clean up your Container Registry by deleting namespaces that do not contain any images.

## Requirements

- Scaleway Account
- Docker daemon running to build the image
- Container registry namespace created, for this example we assume that your namespace name is `registry-cleaner`: [doc here](https://www.scaleway.com/en/docs/containers/container-registry/how-to/create-namespace/)
- API keys generated, Access Key and Secret Key [doc here](https://www.scaleway.com/en/docs/iam/how-to/create-api-keys/)

## Step 1: Build and Push to Container Registry

Serverless Jobs, like Serverless Containers (which are suited for HTTP applications), works
with containers. So first, use your terminal reach this folder and run the following commands:

```shell
# The first command logs in to the container registry; you can find it in the Scaleway console
docker login rg.fr-par.scw.cloud/registry-cleaner -u nologin --password-stdin <<< "$SCW_SECRET_KEY"

# The next command builds the image to push
docker build -t rg.fr-par.scw.cloud/registry-cleaner/empty-namespaces:v1 .

## TIP: For Apple Silicon or other ARM processors, please use the following command as Serverless Jobs supports amd64 architecture
# docker buildx build --platform linux/amd64 -t rg.fr-par.scw.cloud/registry-cleaner/empty-namespaces:v1 .

# This command pushes the image online to be used on Serverless Jobs
docker push rg.fr-par.scw.cloud/registry-cleaner/empty-namespaces:v1
```

> [!TIP]
> As we do not expose a web server and we do not require features such as auto-scaling, Serverless Jobs are perfect for this use case.

To check if everyting is ok, on the Scaleway Console you can verify if your tag is present in Container Registry.

## Step 2: Creating the Job Definition

On Scaleway Console on the following link you can create a new Job Definition: https://console.scaleway.com/serverless-jobs/jobs/create?region=fr-par

1. On Container image, select the image you created in the step before.
2. You can set the image name to something clear like `registry-namespace-cleaner` too.
3. For the region you can select the one you prefer :)
4. Regarding the resources you can keep the default values, this job is fast and do not require specific compute power or memory.
5. To schedule your job for example every night at 2am, you can set the cron to `0 2 * * *`.
6. Important: advanced option, you need to set the following environment variables:

> [!TIP]
> For sensitive data like `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` we recommend to inject them via Secret Manager, [more info here](https://www.scaleway.com/en/docs/serverless/jobs/how-to/reference-secret-in-job/).

- **Environment Variables**: Set the required environment variables:
  - `SCW_DEFAULT_ORGANIZATION_ID`: Your Scaleway organization ID.
  - `SCW_ACCESS_KEY`: Your Scaleway API access key.
  - `SCW_SECRET_KEY`: Your Scaleway API secret key.
  - `SCW_PROJECT_ID`: Your Scaleway project ID.
  - `SCW_NO_DRY_RUN`: Set to `true` to delete namespaces; otherwise, it will perform a dry run.

* Then click "Create Job"

## Step 3: Run the job

On your created Job Definition, just click the button "Run Job" and within seconds it should be successful.

## Troubleshooting

If your Job Run state goes in error, you can use the "Logs" tab in Scaleway Console to get more informations about the error.

# Additional content

- [Jobs Documentation](https://www.scaleway.com/en/docs/serverless/jobs/how-to/create-job-from-scaleway-registry/)
- [Other methods to deploy Jobs](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/deploy-job/)
- [Secret key / access key doc](https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/)
- [CRON schedule help](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/cron-schedules/)
