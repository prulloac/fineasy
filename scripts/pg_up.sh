#!/bin/sh

# spin up postgres container with docker
# USER: postgres
# PASSWORD: 1234
# DATABASE: postgres
# PORT: 5432
# sslmode=disable
docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=postgres -p 5432:5432 -d postgres:alpine