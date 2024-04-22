## Deploy Hugging Face Models in Serverless Containers

- Export these variables:

```bash
export SCW_ACCESS_KEY="access-key" SCW_SECRET_KEY="secret-key" SCW_PROJECT_ID="project-id"
```

- Run script to deploy multiple hugging face models using terraform workspaces:

```bash
bash ./deploy-models.sh
```