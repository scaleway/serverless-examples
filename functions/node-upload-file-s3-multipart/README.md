# Node function to read files and upload to S3

This function does the following steps :

* Read a file from an HTTP request
* Send the file to an S3 bucket

## Requirements

If you want to enable S3 upload, ensure to create a bucket and have the following secrets variables available in your environment:

```env
S3_REGION = # Default: fr-par
SCW_ACCESS_KEY =
SCW_SECRET_KEY =
BUCKET_NAME =
```

To call the function (replace `README.md` with the file you want to upload):

```sh
curl -X POST -F "data=@README.md" -F "another=@package.json" http://localhost:8080
```
