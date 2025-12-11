#!/bin/bash

rm ./dev/kubernetes/local/*

kind delete cluster -n perses-kind

kind create cluster \
    --config="./dev/kubernetes/kind.yaml" \
    --name="perses-kind"

sh ./scripts/k8s/create-resources.sh
