## v0.5.0

- Replicate xclip behaviour by @aosasona in https://github.com/aosasona/bore/pull/23
  - Invoking `bore` without any subcommand or argument will now:
    - paste the last item on the clipboard if there was no data piped in
    - copy if there was data piped in
- Load config in current directory (`bore.toml` or `.bore/config.toml` in the current directory) if present (as long as the user did not specify a config via the `-c` flag) by @aosasona in https://github.com/aosasona/bore/pull/24
- Use `source` as version by default when installed from source

## v0.4.0

- Use cleaner table output for config dump by @aosasona in https://github.com/aosasona/bore/pull/19
- Implemented "clear on paste" functionality by @aosasona in https://github.com/aosasona/bore/pull/20
- Refactored handler's public and private interfaces

## v0.3.0

- Implemented native clipboard passthrough support by @aosasona in https://github.com/aosasona/bore/pull/15

## v0.2.2

- Added support for copying `base64-encoded` content
- Refactored native clipboard stubs

## v0.2.1

- Fix: use `latest` version by default for `go install ...`
