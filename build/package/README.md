# Building and Packaging

All components of the Paloma CDP can be built and packaged using Docker. The following sections describe how to configure a deployement using Docker Compose.

Using the supplied `compose.yaml` file, the following profiles are available:

- `migration`: Creates the database schema and populates the database with initial data.
- `pipeline`: Spawns data ingestion, transformation and the web API service.


There is no graceful migration technique implemented. The `migration` profile will apply the configured migration version on the live instance, it's therefore advised to shutdown the pipeline services during this operation.

## Customizing deployments

To specify which version should be deployed, create a new file called `compose.override.yaml` in this directory with the following content:

```yaml
services:
  migrate:
    image: palomachain/cdp-migrate:VERSION
  rest:
    image: palomachain/cdp-rest:VERSION
    ports:
      - YOURPORT:8011
  ingest:
    image: palomachain/cdp-ingest:VERSION
  transform:
    image: palomachain/cdp-transform:VERSION
```

The contents will automatically be merged with the `compose.yaml` file, any further customization can be done in the same way.

## Controlling the pipeline

```sh
docker compose --profile migration up -d
docker compose --profile pipeline up -d
```

## Service configuration 

> [!WARNING]
> Local configuration and environment variables are not encrypted and should never be committed to version control!

All configuration is done through environment variables. The project ships with all available configuration keys mapped in individual `.env` files. The following variables are available:

 - `ingest.env`: All configuration keys for the ingest service.
 - `transform.env`: All configuration keys for the transform service.
 - `rest.env`: All configuration keys for the rest service.
 - `persistence.env`: Shared configuration among all services with dependencies on the database.
 - `postgres.env`: Configuration for the database.

 Similarly to the deployment customization, you can create a `SERVICENAME.local.env` file which overrides any or all of the configuration keys, like this:

```env
# persistence.env
CDP_PSQL_ADDRESS=postgres:5432
CDP_PSQL_USER=cdp
CDP_PSQL_PASSWORD=
CDP_PSQL_DATABASE=cdp
CDP_PSQL_TIMEOUT=5s

# persistence.local.env
CDP_PSQL_PASSWORD=mysecretpassword
```
