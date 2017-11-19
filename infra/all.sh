#!/usr/bin/env bash

# Create namespace
kubectl create -f namespaces/streaming-namespace.yml

# Create limits
kubectl create -f limits/streaming-limits.yml

# Add docker registry to secrets
sh ./docker-auth.sh mateuszdyminski "$1" dyminski@gmail.com

# Create deployments
kubectl create -f deployments/streaming-api.yml

# Create service
kubectl create -f services/api.yml
