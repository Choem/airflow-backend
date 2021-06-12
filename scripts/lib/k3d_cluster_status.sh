#!/usr/bin/env bash

# returns status 0=not existing, 1=stopped, 2=running
function k3d_cluster_status {
    cluster_entry=$(k3d cluster list | grep $1)
    rows=($cluster_entry)
    servers=($(echo ${rows[1]} | tr "/" "\n"))
    if [[ ! -z "$cluster_entry" ]];then
        if [[ ${servers[0]} == ${servers[1]} ]];then
            echo 2
            exit
        else
            echo 1
            exit
        fi
    else
        echo 0
        exit
    fi
}