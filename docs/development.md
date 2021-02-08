# Development Notes

gscloud tries to be a well-behaved command-line tool. We try to follow the [Command Line Interface Guidelines](https://clig.dev/) as much as possible. Feel free to open an issue if you think we could do better.

Most notably:

- All errors should go to stderr, always.
- All other results should go to stdout. Informational messages should go to stderr.
- Print tables by default.
- Print JSON on stdout when -j or --json was given so that user can pipe to jq(1) if needed. Make sure you always print valid JSON or nothing at all.
- Set return code != 0 in case of error. Errors are simple messages meant for the user to read and understand quickly.
- Don't be too noisy. If there is nothing to tell, don't print anything.
- When using IP addresses in examples, use [203.0.113.0/24](https://tools.ietf.org/html/rfc5737) prefix for IPv4 addresses and [2001:db8::/32](https://tools.ietf.org/html/rfc3849) prefix for IPv6 addresses.
- When using email addresses or domain names in examples, use `example.com`.

## Building and testing

- Run `make` to build and `make test` to run tests for your platform.
- Run `goreleaser build --snapshot --rm-dist` to test build for all supported platforms. Make sure to install [GoReleaser](https://goreleaser.com/) before.

## Don't hesitate to get in touch

There are no stupid questions. Feel free to use GitHub issues if you have a problem or question. If you want to work on something, make sure to create an issue first to let others know.
