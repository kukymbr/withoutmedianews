# WithoutMedia news 

Hi, we are nonexisting WithoutMediaNews team, 
and we are providing news without media files and with no any sense nor facts.

## Running the app

1. Start the Postgres with a schema from [withoutmedianews.sql](db/schema/withoutmedianews.sql) applied.
   As an option, a Docker configuration is available in the [.deps](.deps) directory:
   1. Set environment variables for the Postgres container using the `.env` file:
      ```shell
      cp .deps/example.env .deps/.env
      ```
   2. Start the container:
      ```shell
      make run_dev_deps
      ```
2. Build an application:
   ```shell
   make build
   ```
3. Initialize environment variables for an application; 
   all supported variables names could be found in the [example.env](example.env) file.
   ```shell
   cp example.env .env
   source .env
   ```
4. Run the binary:
   ```shell
   ./bin/withoutmedianews
   ```