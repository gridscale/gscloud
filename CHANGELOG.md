# Changelog

## v0.7.0 (UNRELEASED)

FEATURES:

* gscloud-server-create learned `--profile` flag to specify a HW profile.
* gscloud-server-events subcommand has been added. You can now fetch event logs for a server ([#102](https://github.com/gridscale/gscloud/issues/102)).

FIXED:

* If no HW profile is specified when creating a server, `"q35"` is used now ([#89](https://github.com/gridscale/gscloud/issues/89)).

## v0.6.0 (November 23, 2020)

FEATURES:

* gscloud learned `gscloud server assign SERVER-ID IP-ADDR`
* gscloud learned `gscloud ip assign ID|ADDR` and  `gscloud ip release ID|ADDR` ([#85](https://github.com/gridscale/gscloud/issues/85)).
* Releases are now signed with our `gridscale GmbH <oss@gridscale.io>` GPG key (key ID: `4841EC2F6BC7BD4515F60C10047EC899C2DC3656`, [#72](https://github.com/gridscale/gscloud/issues/72)). Thanks @nvthongswansea!

FIXED:

* Lots of fixes in help texts. Better man pages.

## v0.5.0 (October 14, 2020)

FEATURES:

* Add a `gscloud template rm` command ([#80](https://github.com/gridscale/gscloud/issues/80)).
* Add basic support for IP addresses ([#78](https://github.com/gridscale/gscloud/issues/78)).
* Add support for storage resize ([#77](https://github.com/gridscale/gscloud/issues/77)).

FIXED:

* Weird nesting of JSON output ([#79](https://github.com/gridscale/gscloud/issues/79)).

## v0.4.0 (September 8, 2020)

Many bug fixes and additions in this release. We also dropped the "beta" from the version string. More to come.

FEATURES:

* Add `manpage` command to generate man-pages for gscloud.
* Add `gscloud template ls` to list available templates ([#59](https://github.com/gridscale/gscloud/issues/59)).
* Add `--noheading` flag to print tables without header ([#53](https://github.com/gridscale/gscloud/issues/53)).
* Add `gscloud server set` to allow changing server properties and hot-plugging (see [e48c149](https://github.com/gridscale/gscloud/commit/e48c149af4ff19fb846c7fb8288d0a6029880066)).

FIXED:

* Fixed working with multiple accounts ([#58](https://github.com/gridscale/gscloud/issues/58)).
* Fixed printing CHANGETIME column ([#60](https://github.com/gridscale/gscloud/issues/60)).
* `kubernetes` command's error handling has been improved ([#18](https://github.com/gridscale/gscloud/issues/18)).

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
