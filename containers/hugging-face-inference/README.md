# Hugging Face Models

### Deploy models in Serverless Containers

- Export these variables:

```bash
export SCW_ACCESS_KEY="access-key" SCW_SECRET_KEY="secret-key" SCW_PROJECT_ID="project-id" REGION="fr-par"
```

- Add/remove Hugging Face models (with `.gguf` extension) in `terraform/hf-models.json` file.

- Run script to deploy multiple hugging face models using terraform workspaces:

```bash
cd terraform && bash terraform.sh -a
```

### Benchmark models

Check your models were deployed on the console and copy your container endpoints to the `terraform/hf-models.json` file, then perform the following command:

```bash
python benchmark-models.py
```

This will generate a box plot to analyze response time per model family, and a `csv` file containing textual responses per each model.

### Destroy terraform resources for all models

```bash
bash terraform.sh -d
```