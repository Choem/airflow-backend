#!/usr/bin/env bash

# Get project root path
ROOT_PATH=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )/..

# Load libs
source ${ROOT_PATH}/scripts/lib/load_env.sh
source ${ROOT_PATH}/scripts/lib/k3d_cluster_status.sh

CLUSTER_STATUS=$(k3d_cluster_status $CLUSTER_NAME)

# Based on the cluster status create, start or print message
if [ $CLUSTER_STATUS -eq 0 ]; then
  echo "cluster: $CLUSTER_NAME does not exist!"
elif [ $CLUSTER_STATUS -ge 1 ]; then
  k3d cluster delete $CLUSTER_NAME
fi

