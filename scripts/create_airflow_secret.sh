#!/usr/bin/env bash

# Get project root path
ROOT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )/..

# Load libs
source ${ROOT_PATH}/scripts/lib/load_env.sh

# Apply secret to cluster
kubectl create secret generic \
  airflow-ssh-git-secret \
  --from-file=id_ed25519=$AIRFLOW_PRIVATE_KEY_PATH