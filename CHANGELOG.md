# Changelog

## v0.11.1 (2022-XX-XX)

FEATURES:

* `gscloud info --json` now includes the sum of total cores, memory, and storage capacity in the output.

FIXED:
* GitHub does not support `%(describe)`, yet, so we have to upload release artifacts manually for now (see https://github.community/t/support-for-describe-in-export-subst/196618 and [#131](https://github.com/gridscale/gscloud/issues/131)).

## v0.11.0 (2021-09-22)

FEATURES:

* Release tarballs now include version strings. With this you can simply use make(1) in a build environment without git(1) installed and still have `gscloud version` produce correct output (Fixes [#131](https://github.com/gridscale/gscloud/issues/131)).
* gscloud-server-create `--with-template` learned to accept IDs in addition to template names (see [PR #133](https://github.com/gridscale/gscloud/pull/133)).

FIXED:

* Fixed the behavior of `iso-image ls --quiet` ([#134](https://github.com/gridscale/gscloud/issues/134)). Thanks [@ghostwheel42](https://github.com/ghostwheel42)!
* Fixed output of `gscloud kubernetes cluster -h` ([#137](https://github.com/gridscale/gscloud/issues/137)).

## v0.10.0 (2021-04-06)

FEATURES:

* We added `gscloud kubernetes releases` and `gscloud postgresql releases` subcommands that gives you a list of all currently available Managed Kubernetes releases ([#113](https://github.com/gridscale/gscloud/issues/113)) and PostgreSQL releases ([#122](https://github.com/gridscale/gscloud/issues/122)).
* gscloud gained a `gscloud info` command that shows you a quick account summary as well as the API tokens and user account in use. Example:

```raw
$ gscloud --account=dev@example.com info
SETTING    VALUE
Account    dev@example.com
User ID    7ff8003b-55c5-45c5-bf0c-3746735a4f99
API token  <redacted>
URL        https://api.gridscale.io
Getting information about used resources…
OBJECT             COUNT
Platform services  0
Servers            18
Storages           24
IP addresses       2
```

FIXED:

* This release also fixes the build on OpenBSD ([3be8074](https://github.com/gridscale/gscloud/commit/3be807415a17d1ea29ea7d1bdb0493d5f825ba48)).

## v0.9.0 (2021-02-27)

FEATURES:

* We removed the `--password` flag when. Passwords are now auto-generated when creating servers ([#103](https://github.com/gridscale/gscloud/issues/103)).
* We added builds for Apple M1 to our releases ([#112](https://github.com/gridscale/gscloud/issues/112)).

## v0.8.0 (2021-02-19)

FEATURES:

* You can now create networks with gscloud-network-create ([PR #107](https://github.com/gridscale/gscloud/pull/107)).
* Make gscloud-server-events a bit more useful by adding initiator column and removing other less useful ones ([#110](https://github.com/gridscale/gscloud/issues/110)).
* gscloud server rm learned a `--include-related` flag that includes storages and assigned IP addresses when removing servers ([#98](https://github.com/gridscale/gscloud/issues/98)).
* Added examples to the README to get started more quickly ([#93](https://github.com/gridscale/gscloud/issues/93)).

FIXED:

* Removing an object will print a message now to let you know what happens.
* gscloud-server-create does now leave a clean state if server creation fails ([#97](https://github.com/gridscale/gscloud/issues/97)).

## v0.7.1 (2021-01-13)

FIXED:

* Fix a bug in gscloud-server-create when --password was not given ([#104](https://github.com/gridscale/gscloud/issues/104)).

## v0.7.0 (2021-01-10)

FEATURES:

* gscloud-server-create will now auto generate passwords when `--with-template=` is given and `--password` is not explicitly given on the command line. `--password` flag itself has been marked as deprecated and will be removed in a future release. Passwords are generated on the client and should be sufficiently secure (we use [github.com/sethvargo/go-password](https://pkg.go.dev/github.com/sethvargo/go-password)). The password is printed after the storage is created (See [#90](https://github.com/gridscale/gscloud/issues/90) for more).
* gscloud-server-create learned `--profile` flag to specify a HW profile.
* gscloud-server-create learned `--availability-zone` flag to influence a server's physical distance ([#91](https://github.com/gridscale/gscloud/issues/91)).
* gscloud-server-create also learned `--auto-recovery` flag to specify auto-recovery behavior ([#92](https://github.com/gridscale/gscloud/issues/92)).
* gscloud-server-events subcommand has been added. You can now fetch event logs for a server ([#102](https://github.com/gridscale/gscloud/issues/102)).
* Added new `gscloud iso-image` subcommand to list, create, and eventually delete ISO image objects ([#101](https://github.com/gridscale/gscloud/issues/101)).

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
