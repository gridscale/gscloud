# Changelog

## 3.0.1 (May 19, 2020)

FEATURES:
* Add K8S support (get k8s config, renew k8s cluster's credentials).

BUG FIXES:
* Fixed wrong import path in exmaples.
* Fixed wrong import path in README.

## 3.0.0 (May 6, 2020)

FEATURES:
* Add storage cloning.

DEPRECATED FEATURES:
* `requestCheckTimeoutSecs` is removed from `NewConfiguration` function.

IMPROVEMENTS:
* Every function (mostly) can be controlled through context.
* `ShutdownServer` does not run powering off when the server cannot be shut down gracefully. To ungracefully power off a server, the `StopServer` function should be used when `ShutdownServer` fails.
* Reduce size of `vendor` directory by removing unnecessary packages.

## 2.2.2 (April 8, 2020)

DEPRECATED FEATURES:
* Deprecated and removed labels create/delete options.

BUG FIXES:
* Fixed "context is expired but still retrying".
* Fixed some typos.

## 2.2.1 (January 21, 2020)

BUG FIXES:
* (Hot fix) Fix incompatible major version when using gomod due to missing `/v2` suffix in module path of `go.mod` file

## 2.2.0 (January 21, 2020)

IMPROVEMENTS:
* Retry requests in case of network issues (timeouts, connection resets, connection refused, etc)
* Simple requests backoff in case of retrying errors
* Increase defaultDelayIntervalMilliSecs to 1000 to reduce stress on API
* Better variables/functions' names
* Remove `LocationUUID` as objects' location depends on Project's location
* Add gomod

BUG FIXES:
* Fix "cannot remove all labels of an object"
* Fix resource leak due to not closing response's body

## 2.1.0 (November 05, 2019)

IMPROVEMENTS:
* Errors that are from http requests now include request UUIDs
* No need to create structs when exporting snapshots to S3
* Waiting for asynchronous requests is now faster and more memory-friendly

BUG FIXES:
* Fix README
* Fix missing JSON properties

## 2.0.0 (October 07, 2019)

FEATURES:
* Add `sync` mode
* The library now can be fully controlled through `context`
* Auto retry when server returns 5xx and 424 http codes
* Add a default configuration for the client

IMPROVEMENTS:
* Server type is now limited to pre-defined values
* Storage type is now limited to pre-defined values
* IP address family is now limited to pre-defined values
* Loadbalancer algorithm is now limited to pre-defined values
* All time-related properties are now type of GSTime (a custom type of time.Time)
* Friendly godoc.

BUG FIXES:
* Fixed bugs when unmarshalling JSON to concrete type (mismatched type)

## 1.0.0 (September 05, 2019)

FEATURES:
* Add support for Locations
* Add support for Events
* Add support for Labels
* Add support for Deletes

IMPROVEMENTS:
* Heavily code refactoring to improve code quality
* Achieve 95% test coverage
* Achieve 100% compliant golang code conventions based on goreportcard
* Power-off server if graceful shutdown fails
* Backward compatibility for server creation API

## 0.2.0 (August 23, 2019)

FEATURES:

* Add support for LBaaS (GH-2)
* Add support for PaaS (GH-6)
* Add support for ISO Image Handling (GH-8)
* Add support for Object Storage (GH-11)
* Add support for Snapshots (GH-12) and Snapshot Scheduler (GH-13)
* Add support for Firewall Handling (GH-14)

IMPROVEMENTS:

* Unit Tests for all functionality
* Logging support
* Many examples have been added
* Consistency as well as language styles improved

## 0.1.0 (January 2, 2019)

- Initial release of gsclient-go.
