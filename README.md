# homelab
Infrastructure code for my home lab

For background documentation about the homelab project, decisions and design choices you can [read my blog article](https://blog.turong.dev/blog/planning-a-homelab/).

[homelab-architecture](https://blog.turong.dev/blog/homelab-architecture.png)

## Introduction

This repository contains the code needed to build my homelab physical cluster, bootstrap it with ArgoCD, then hand over to ArgoCD to manage the deployment of services and sub (virtual) clusters. As far as is practical every configuration is managed here, with sensitive information held in ad external secret provider (Bitwarden Secrets Manager).

## To-do

- #2 At present the code is very static and doesn't support different configurations of services for different environments, so I'll be improving this over time.
- With more flexible configurations in place, I'll be able to create local (#3) and virtual (#4) development environments, allowing me to adopt a branching and merging strategy - an important development practice.
- #5 Enable automated testing of changes in a deployment pipeline