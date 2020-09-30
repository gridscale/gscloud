# gscloud: CLI for the gridscale cloud

## Overview

gscloud lets you manage objects on [gridscale.io](https://my.gridscale.io) via shell.

```txt
$ gscloud --help
gscloud lets you manage objects on gridscale.io via command line. It provides a Docker-CLI comparable command line that allows you to create, manipulate, and remove objects on gridscale.io.

Usage:
  gscloud [command]

Available Commands:
  completion  Generate shell completion scripts
  help        Help about any command
  kubernetes  Operate managed Kubernetes clusters
  make-config Create a new configuration file
  manpage     Create man-pages for gscloud
  network     Operations on networks
  server      Operations on servers
  ssh-key     Operations on SSH keys
  storage     Operations on storages
  template    Operations on templates
  version     Print the version

Flags:
      --account string   Specify the account used (default "default")
      --config string    Specify a configuration file (default "~/.config/gscloud/config.yaml")
  -h, --help             Print usage
  -j, --json             Print JSON to stdout instead of a table
      --noheading        Do not print column headings
  -q, --quiet            Print only IDs of objects

Use "gscloud [command] --help" for more information about a command.
```

## Configuration

You can use `gscloud make-config` to generate a new config file. Make sure to add your user ID and API token here.

Example config:

```yml
accounts:
- account:
  name: default
  userId: 2727b9ab-65ff-4d1e-af5e-d08d682bd1fa
  token: 6eb139b3b6515515a6f358d3a635e9b38f05935782602d4fd5c1b5716af54526
- account:
  name: liveaccount
  userId: 2727b9ab-65ff-4d1e-af5e-d08d682bd1fa
  token: 6eb139b3b6515515a6f358d3a635e9b38f05935782602d4fd5c1b5716af54526
  url: https://api.gridscale.io
```

## Example configuration for ~/.kubeconfig/config

To use gscloud for user authentication in kubectl, here is an example kubeconfig:

```yml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tL
    server: https://185.102.93.54:6443
  name: k8s-1-15-5-gs0-de-vte1-bbm49m
contexts:
- context:
    cluster: k8s-1-15-5-gs0-de-vte1-bbm49m
    user: k8s-1-15-5-gs0-de-vte1-bbm49m-admin
  name: k8s-1-15-5-gs0-de-vte1-bbm49m-admin@k8s-1-15-5-gs0-de-vte1-bbm49m
current-context: k8s-1-15-5-gs0-de-vte1-bbm49m-admin@k8s-1-15-5-gs0-de-vte1-bbm49m
kind: Config
preferences: {}
users:
- name: k8s-1-15-5-gs0-de-vte1-bbm49m-admin
  user:
    exec:
      apiVersion: "client.authentication.k8s.io/v1beta1"
      command: $HOME/gscloud
      args:
        - "--config"
        - "$HOME/.config/gscloud/config.yaml"
        - "--account"
        - "test"
        - "kubernetes"
        - "cluster"
        - "exec-credential"
        - "--cluster"
        - "9489f3a7-c8f8-4b38-bc9b-aa472a1c0d2a"
```

## Add completion scripts to your shell

zsh

```shell
$ gscloud completion zsh >> ~/.zshrc
```

bash

```shell
$ gscloud completion bash >> ~/.bash_profile
```
