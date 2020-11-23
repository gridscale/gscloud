# Development Notes

- All errors should go to stderr, always.
- All other output should go to stdout.
- Print tables by default.
- Print JSON on stdout when -j|--json was given so that user can pipe to jq(1) if needed. Make sure you always print valid JSON or nothing at all.
- Set return code != 0 in case of error. Errors don't have to be JSON.
- Don't be too noisy. If there is nothing to tell, don't print anything.
- When using IP addresses in examples use [203.0.113.0/24](https://tools.ietf.org/html/rfc5737) prefix for IPv4 addresses and [2001:db8::/32](https://tools.ietf.org/html/rfc3849) prefix for IPv6 addresses.

## Building and testing

- Run `make` to build and `make test` to run tests for your platform
- Run `goreleaser build --snapshot --rm-dist` to test build for all supported platforms. Make sure to install [GoReleaser](https://goreleaser.com/) before.
