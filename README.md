# homelab
Infrastructure code and documentation for my home lab

For background documentation about the homelab project, decisions and design choices you can [read my blog article](https://blog.turong.dev/blog/planning-a-homelab/).

https://blog.turong.dev/blog/homelab-architecture.png

## Introduction

This repository contains the code needed to build my homelab physical cluster, bootstrap it with ArgoCD, then hand over to ArgoCD to manage the deployment of services and sub (virtual) clusters. As far as is practical every configuration is managed here, with sensitive information held in ad external secret provider (Bitwarden Secrets Manager).

## To-do

- At present the code is very static and doesn't support different configurations of services for different environments, so I'll be improving this over time.
- With more flexible configurations in place, I'll be able to create local and virtual development environments, allowing me to adopt a branching and merging strategy - an important development practice.