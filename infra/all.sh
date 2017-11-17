#!/usr/bin/env bash

# Create namespace
kubectl create -f namespaces/streaming-namespace.yml

# Add docker registry to secrets
docker-auth.sh mateuszdyminski Samsung!7 dyminski@gmail.com

# Create deployments
kubectl create -f deployments/streaming-api.yml
