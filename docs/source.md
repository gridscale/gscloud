# Installation from source

## Prerequisites

Verify that you have Go 1.14+ installed

    $ go version
    go version go1.15.7 linux/amd64

If Go is not installed, follow instructions on the [Go website](https://golang.org/doc/install) or, on Linux, download with your package manager.

## Get the source

Clone this repository:

    $ git clone https://github.com/gridscale/gscloud.git
    $ cd gsloud/

## Build and install

### Unix-like systems

```sh
# Installs to $GOPATH/bin by default
$ make
$ go install
```

### Windows

TODO

## Test

Run `gscloud version` to check if everything worked.

    $ gscloud version
    Version:        v0.7.1-14-g28b968e
    Git commit:     28b968e01ba1a89102e3637412846d2e0dbe7a9c
