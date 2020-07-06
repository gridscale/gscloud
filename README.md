# gscloud: CLI for the gridscale cloud

[![Build Status](https://travis-ci.com/gridscale/gscloud.svg?branch=develop)](https://travis-ci.com/gridscale/gscloud)

## Overview

gscloud lets you manage objects on [gridscale.io](https://my.gridscale.io) via shell.

Note: this tool is still in the making and beta quality. Feel free to try it out. Feedback very welcome.

```txt
$ gscloud --help
gscloud is the CLI for the gridscale cloud.

Usage:
  gscloud [command]

Available Commands:
  help        Help about any command
  kubernetes  Operate managed Kubernetes clusters
  make-config Create a new configuration file
  network     Operations on networks
  server      Operations on servers
  ssh-key     Operations on SSH keys
  storage     Operations on storages
  version     Print the version
  completion  Generate shell Tab-Completion

Flags:
      --account string   Specify the account used; 'default' if none given
      --config string    Specify a configuration file; default ~/.config/gscloud/config.yaml
  -h, --help             Print usage
  -j, --json             Print JSON to stdout instead of a table
  -q, --quiet            Print only IDs of objects

Use "gscloud [command] --help" for more information about a command.
```

## Example Configuration

```yml
accounts:
- account:
  name: default
  userId: a13e84c9-852d-484-xxx-xxxx
  token: "9b6590592c65f7daa707d88"
- account:
  name: liveaccount
  userId: a13e84c9-852d-2222-xxx-xxxx
  token: 2222290592c6522222f7daa707d88
  url: https://api.gridscale.io

```

## Example configuration for ~/.kubeconfig/config

To use gscloud for user authentication in kubectl, here is a sample of kubeconfig

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
## Example add Tab-Completion

Uncomment compdef, otherwise it won't work.

`#compdef _gscloud gscloud`

zsh
```
$ ./gscloud completion zsh >> ~/.zshrc
```
bash
```
$ ./gscloud completion bash >> ~/.bash_profile
```
