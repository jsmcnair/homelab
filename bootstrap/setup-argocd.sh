#!/usr/bin/env bash

ARGOCD_SSH_KEYFILE="${ARGOCD_SSH_KEYFILE:?Must set ARGOCD_SSH_KEYFILE to path of and SSH private key}"

# Check helm is available
if ! command -v helm &> /dev/null; then
    echo "❌ Helm is not installed. Please install Helm before proceeding."
    exit 1
fi
echo "✅ Helm is installed."

# Ensure helm has ArgoCD repository added
if ! helm repo list | grep -q argo; then
    if ! helm repo add argo https://argoproj.github.io/argo-helm; then
        echo "❌ Failed to add ArgoCD Helm repository."
        exit 1
    else
        echo "✅ ArgoCD Helm repository added successfully."
    fi
else
    echo "✅ ArgoCD Helm respository already added."
fi
if ! helm repo update &>/dev/null; then
    echo "❌ Failed to update Helm repositories."
    exit 1
else
    echo "✅ Helm repositories updated."
fi

# Ensure helm has ArgoCD installed
if ! helm status argocd -n argocd &> /dev/null; then
    if helm install argocd argo/argo-cd --namespace argocd --create-namespace > /dev/null; then
        echo "✅ ArgoCD installed successfully."
    else
        echo "❌ Failed to install ArgoCD."
        exit 1
    fi
else
    echo "✅ ArgoCD is already installed."
fi

# External secrets isn't available yet so we declaratively create the ArgoCD repository secret if it doesn't exist
if ! kubectl get secret argocd-repo-credentials -n argocd &> /dev/null; then
    REPOSITORY=$(git remote get-url origin)
    if ! kubectl create secret generic argocd-repo-credentials \
        --namespace argocd \
        --from-literal=url="$REPOSITORY" \
        --from-literal=sshPrivateKey="$(cat $ARGOCD_SSH_KEYFILE)" &> /dev/null; then
        echo "❌ Failed to create ArgoCD repository secret."
        exit 1
    else
        echo "✅ ArgoCD repository secret created successfully."
    fi
    if ! kubectl label secret argocd-repo-credentials -n argocd argocd.argoproj.io/secret-type=repository &> /dev/null; then
        echo "❌ Failed to label ArgoCD repository secret."
        exit 1
    else
        echo "✅ ArgoCD repository secret labeled successfully."
    fi
else
    echo "✅ ArgoCD repository secret already exists."
fi
