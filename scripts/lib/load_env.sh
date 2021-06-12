#!/usr/bin/env bash

# Get project root path
ROOT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )/../..

# Env file path
ENV_FILE=$ROOT_PATH/.env

# Export the vars in .env into your shell:
export $(egrep -v '^#' $ENV_FILE | xargs)

# Export registry name
export REGISTRY_NAME=$CLUSTER_NAME-registry