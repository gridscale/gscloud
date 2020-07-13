# Changelog

## v0.3.0-beta (July 13, 2020)

FEATURES:

* We added `gscloud {server,storage,network,ssh-key}` commands. These commands allow you to list and manipulate the objects in various ways.
* You can now output all data as JSON by passing `--json` flag on the command line.
* There are now shell completions available for bash and zsh.

And much more.

## v0.2.0-beta (March 11, 2020)

FEATURES:

* Use standard user-level cache directory [#11](https://github.com/gridscale/gscloud/issues/11).

## v0.1.0-beta (January 8, 2020)

Initial release of gscloud.

FEATURES:

* Support make-config for creating a new configuration file
* Support Kubernetes cluster sub-commands: save-kubeconfig and exec-credential for managing a cluster's authentication
