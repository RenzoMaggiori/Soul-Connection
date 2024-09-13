#!/bin/sh

# This SCRIPT is ONLY to be used in the DOCKERFILE

./scripts/wait-for-it.sh db:5432 5432 '-- echo "Postgres is ready"'
./scripts/wait-for-it.sh file-storage:27017 27017 '-- echo "Mongo is ready"'
./api/api -env-path .env -port 8000

