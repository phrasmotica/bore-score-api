# Deployment to Azure App Services

## Building and pushing the Docker container

Run `build.ps1 -Push` locally. This will push the Docker container to Docker Hub, and the *bore-score-api* App Service deployment will consume it automatically.

## Configuring app settings

Set the following Application settings:

- `GIN_MODE`=`release`
- `BORESCORE_ENV`=`production`

Set **at least one** of the following Application settings:

- `AZURE_TABLES_CONNECTION_STRING`=`<connection string for Azure Table Storage>`
- `MONGODB_URI`=`<URI for Mongo DB>`

Set the following General settings:

- enable HTTPS Only
