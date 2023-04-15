# BoreScore API

A web API written in Golang for recording the results of board games.

## Local Development

Use `docker-compose up` to run the API locally - this will create the following containers:
- the API itself
- the Azurite Storage Emulator

Table data is stored in the local `.azurite` directory, which the Azurite container mounts as a Docker volume.
