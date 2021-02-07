# gscloud

`gscloud` is the command-line tool which let you manage your virtual infrastructure on the [gridscale](https://gridscale.io)'s cloud platform. It's also possible to make use of it to work with Kubernetes clusters running on top of the platform as well.

It's entirely written in Go, so everything you need is just the `gscloud` binary! Compile it all by yourself or get the tool at the [release](https://github.com/gridscale/gscloud/releases) page.

```txt
Usage:
  gscloud [command]

Available Commands:
  completion  Generate shell completion scripts
  help        Help about any command
  ip          Operations on IP addresses
  iso-image   Operations on ISO images
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
      --config string    Path to configuration file (default "$XDG_CONFIG_HOME/gscloud/config.yaml")
      --debug            Debug mode
  -h, --help             Print usage
  -j, --json             Print JSON to stdout instead of a table
      --noheading        Do not print column headings
  -q, --quiet            Print only object IDs

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

## Kubernetes

To use `gscloud` combined with `kubectl`, here is an example configuration (~/.kubeconfig/config):

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

## Exit Codes

`gscloud` returns zero exit code on success, non-zero on failure. Following exit codes map to these failure modes:

1. The requested command failed.
2. Reading the configuration file failed.
3. The configuration could not be parsed.
4. The account specified does not exist in the configuration file.

## Shell Completions

Generate shell completion scripts for zsh and bash.

* bash

```shell
$ gscloud completion bash >> ~/.bash_profile
```

* zsh

```shell
$ gscloud completion zsh >> ~/.zshrc
```

## Install man-pages

Generate man-pages and install them. Example:

```shell
$ sudo gscloud manpage /usr/local/share/man/man1/
```

## Development

Please create an [issue](https://github.com/gridscale/gscloud/issues) if you have questions, want to start a discussion, or want to start work on something.

[Pull requests](https://github.com/gridscale/gscloud/pulls) are always welcome. Make sure to create an issue first to signal others that you are working on something. Also make sure to take a look at the [Development Notes](development.md).

Have fun!
