# Applications

Applications managed via ArgoCD. An app is created by the boostrap script that points to this directory, and has recursion disabled. This means that only ArgoCD applications should reside at the top level, and any manifests these applications reference should be in sub-directories.