#!/usr/bin/env bash

# Get project root path
ROOT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )/..

# Load libs
source ${ROOT_PATH}/scripts/lib/load_env.sh

# Place json in daemon file
bash -c 'cat <<EOT >>/etc/docker/daemon.json
{
  "insecure-registries" : ["localhost:$REGISTRY_PORT"]
}
EOT
'