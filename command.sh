#!/bin/bash
    
action="${1:-help}"

case $action in
    help)
        echo "Usage: command.sh [action]"
        echo "Actions:"
        echo "  help      Display this help message"
        echo "  start     Start the cluster - send WoL packets to nodes"
        echo "  stop      Stop the cluster - send shutdown commands to each node"
        echo "  bootstrap Run the bootstrap script against the current Kubeconfig context"
        ;;
    start)
        echo "Starting the cluster..."
        BOOTSTRAP_MAC_ADDRESSES=$(bw get notes homelab-mac-addresses)
        wol $BOOTSTRAP_MAC_ADDRESSES
        ;;
    stop)
        echo "Stopping the cluster..."
        NODES=$(bw get notes homelab-nodes)
        ENDPOINT=$(bw get notes homelab-endpoint)
        bw get notes homelab-talosconfig | base64 -d > .talosconfig
        talosctl shutdown -n $NODES -e $ENDPOINT --talosconfig .talosconfig --wait=false --force
        rm .talosconfig
        ;;
    bootstrap)
        echo "Running the bootstrap script against the current Kubeconfig context..."
        BOOTSTRAP_SECRET_BASE64=$(bw get notes homelab-bootstrap) ARGOCD_SSH_PRIVATE_KEY_FILE=$(bw get notes homelab-argocd-ssh-private-key) ./bootstrap/setup.sh
        ;;
    *)
        echo "Invalid action: $action"
        echo "Use 'command.sh help' for usage information"
        exit 1
        ;;
esac
