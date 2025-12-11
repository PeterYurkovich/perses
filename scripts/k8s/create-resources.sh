#!/bin/bash

set -e

kind get kubeconfig --name perses-kind >./dev/kubernetes/local/kind-admin

kubectl apply -f ./dev/kubernetes/manifests/namespaces.yaml

# Create Perses backend service account and token, then save to a new kubeconfig file for
# ease of use
kubectl create serviceaccount perses-backend --namespace perses
PERSES_BACKEND_TOKEN="$(kubectl create token perses-backend --namespace perses --duration 8760h)"
kubectl --kubeconfig=./dev/kubernetes/local/kind-admin config set-credentials perses-backend --token=$PERSES_BACKEND_TOKEN
cp ./dev/kubernetes/local/kind-admin ./dev/kubernetes/local/perses-backend
kubectl --kubeconfig=./dev/kubernetes/local/perses-backend config set-context --current --user=perses-backend

# Create a service account "user" which has all permissions related to perses resources
kubectl create serviceaccount user --namespace perses
USER_TOKEN="$(kubectl create token user --namespace perses --duration 8760h)"
kubectl --kubeconfig=./dev/kubernetes/local/kind-admin config set-credentials user --token=$USER_TOKEN
cp ./dev/kubernetes/local/kind-admin ./dev/kubernetes/local/user
kubectl --kubeconfig=./dev/kubernetes/local/user config set-context --current --user=user

# Install the perses CRD's into the cluster
kubectl apply -f https://raw.githubusercontent.com/perses/perses-operator/main/config/crd/bases/perses.dev_perses.yaml
kubectl apply -f https://raw.githubusercontent.com/perses/perses-operator/main/config/crd/bases/perses.dev_persesdashboards.yaml
kubectl apply -f https://raw.githubusercontent.com/perses/perses-operator/main/config/crd/bases/perses.dev_persesdatasources.yaml
kubectl apply -f https://raw.githubusercontent.com/perses/perses-operator/main/config/crd/bases/perses.dev_persesglobaldatasources.yaml

# Give the user and perses-backend appropriate permission
kubectl apply -f ./dev/kubernetes/manifests/perses-backend-permissions.yaml
kubectl apply -f ./dev/kubernetes/manifests/user-permissions.yaml
