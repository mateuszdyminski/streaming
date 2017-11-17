#!/bin/bash

NAMESPACES=`kubectl get namespace -o name | awk -F '/' '{print $2}'`
DOCKER_REGISTRY_SERVER='https://index.docker.io/v1/'
DOCKER_USER="$1"
DOCKER_PASSWORD="$2"
DOCKER_EMAIL="$3"

function create-auth-for-namespace {
    kubectl create secret docker-registry streaming-registry \
        --docker-server=$DOCKER_REGISTRY_SERVER  \
        --docker-username=$DOCKER_USER \
        --docker-password=$DOCKER_PASSWORD \
        --docker-email=$DOCKER_EMAIL \
        --namespace=$1
}

for n in $NAMESPACES; do
    create-auth-for-namespace $n
done