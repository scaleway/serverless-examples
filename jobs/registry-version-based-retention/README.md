# Scaleway Container Registry Tag Cleaner

This project aims to clean up Scaleway Container Registry tags to keep only the N latest tags for each image. It is useful for managing disk space and keeping the registry organized.

## Usage

1. **Build the Application:**
   ```bash
   go build -o reg-clean
   ```

2. **Run the Application:**
   ```bash
   ./reg-clean
   ```

## Environment Variables

The application requires the following environment variables to be set:

- `SCW_DEFAULT_ORGANIZATION_ID`: Your Scaleway organization ID.
- `SCW_ACCESS_KEY`: Your Scaleway access key.
- `SCW_SECRET_KEY`: Your Scaleway secret key.
- `SCW_PROJECT_ID`: Your Scaleway project ID.
- `SCW_NUMBER_VERSIONS_TO_KEEP`: The number of latest tags to keep for each image.
- `SCW_NO_DRY_RUN` (optional): Set to `true` to perform actual deletions. If not set, the application will run in dry-run mode, only logging the actions that would be taken.

## Example

To run the application in dry-run mode and keep the 5 latest tags for each image, set the following environment variables and run:

```bash
export SCW_DEFAULT_ORGANIZATION_ID=your-organization-id
export SCW_ACCESS_KEY=your-access-key
export SCW_SECRET_KEY=your-secret-key
export SCW_PROJECT_ID=your-project-id
export SCW_NUMBER_VERSIONS_TO_KEEP=5

./reg-clean
```

To run the application and actually delete the tags, set `SCW_NO_DRY_RUN` to `true`:

```bash
export SCW_NO_DRY_RUN=true

./reg-clean
```