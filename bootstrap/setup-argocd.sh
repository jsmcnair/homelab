#!/usr/bin/env bash

ARGOCD_SSH_KEYFILE="${ARGOCD_SSH_KEYFILE:?Must set ARGOCD_SSH_KEYFILE to path of and SSH private key}"
REPOSITORY=$(git remote get-url origin)

# Check helm is available
if ! command -v helm &> /dev/null; then
    echo "❌ Helm is not installed. Please install Helm before proceeding."
    exit 1
fi
echo "✅ Helm is installed."


# Ensure helm has ArgoCD installed
if ! helm status argocd -n argocd &> /dev/null; then
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
    if helm install argocd argo/argo-cd --namespace argocd --create-namespace --values=bootstrap/argocd-values.yaml &> /dev/null; then
        echo "✅ ArgoCD installed successfully."
    else
        echo "❌ Failed to install ArgoCD."
        exit 1
    fi
else
    if ! helm upgrade argocd argo/argo-cd --namespace argocd --values=bootstrap/argocd-values.yaml &> /dev/null; then
        echo "❌ Failed to upgrade ArgoCD."
        exit 1
    else
        echo "✅ ArgoCD upgraded successfully."
    fi
fi

# External secrets isn't available yet so we declaratively create the ArgoCD repository secret if it doesn't exist
if ! kubectl get secret argocd-repo-credentials -n argocd &> /dev/null; then
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


kubectl config set-context --current --namespace argocd > /dev/null
if ! argocd app get applications > /dev/null; then
    if argocd app create applications \
        --dest-name in-cluster \
        --dest-namespace default \
        --repo $REPOSITORY \
        --path applications \
        --directory-recurse \
        --sync-policy auto \
        --auto-prune \
        --self-heal \
        --directory-include "**application.yaml"; then
        echo "✅ ArgoCD bootstrap application created"
    else
        echo "❌ Failed to create ArgoCD bootstrap application."
        exit 1
    fi
else
    echo "✅ ArgoCD bootstrap application already exists"
fi