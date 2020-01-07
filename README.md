# gscloud: The Command Line Interface for the gridcsale cloud

## Supported Services

    - kubernetes

## Overview

```txt
gscloud --help
gscloud is the command line interface for the gridscale cloud.

Usage:
  gscloud [command]

Available Commands:
  help        Help about any command
  kubernetes  Operate managed Kubernetes clusters
  make-config Create a new configuration file

Flags:
      --account string   the account used, 'default' if none given
      --config string    configuration file, default /home/bk/.config/gridscale/config.yaml
  -h, --help             help for gscloud

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
        - "$HOME/.gscloud/config.yaml"
        - "--account"
        - "test"
        - "kubernetes"
        - "cluster"
        - "exec-credential"
        - "--cluster"
        - "9489f3a7-c8f8-4b38-bc9b-aa472a1c0d2a"
```
