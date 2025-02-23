# Bootstrapping

You have to start somewhere...

This folder needs to be as minimal as possible. It contains only enough configuration to boostrap the automation tools that manage the rest of the configuration.

## Setup ArgoCD

An SSH key with access to this repository needs to be created and the path to the private key should be exported to the `ARGOCD_SSH_KEYFILE` environment variable. It automatically discovers the GitHub repository URL based on the current git repository.

```shell
export ARGOCD_SSH_KEYFILE=~/.ssh/id_rsa
./setup-argocd.sh
```

The script is idempotent so you can run it multiple times to check everything is set up.