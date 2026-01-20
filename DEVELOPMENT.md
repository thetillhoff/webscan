# Development

## Important commands

```sh
# List available make targets
make list

# Install dependencies
make install

# Run tests
make test

# Format code
make format

# Upgrade dependencies
make upgrade
```

## Building

- This application is available for Linux, Mac and Windows.
- A version variable is injected at build/release time, so the version automatically equals the git tag it was built with.

## Workflows and release process

**1. Verification Workflow** (`build-golang-executable-on-push.yaml`): Runs on all PRs and `renovate/**` branches. Builds for all platforms/architectures, injects commit hash as version, and verifies version injection works. On success, renovate may automerge.

**2. Tag-on-Main Workflow** (`tag-on-main.yaml`): Runs on pushes to `main` by `renovate[bot]`. Calculates next patch version, updates CHANGELOG.md (inserts new entry before current version), verifies build (linux/amd64), then creates and pushes git tag.

**3. Release Workflow** (`release-golang-executable-on-tag.yaml`): Runs on tag push (`v[0-9]+.[0-9]+.[0-9]+`). Builds binaries for all platforms/architectures, verifies version, creates GitHub release, builds/pushes Docker images to ghcr.io, and triggers external updates (goreportcard, homebrew tap). On failure, automatically deletes the tag.

All versions are stored as git tags. Pushing a tag triggers the release pipeline.

## Decisions / best practices

- Don't have global variables in a package -> they would be obstructed for the consumer and are not threadsafe
- Don't use functional options -> they require a lot of code / maintenance. Also, having functions to set a context object every time a function is called is tedious
- Use Context (called engine in this project). Not necessarily the go-context package, but implement "instance of package" as context and use that.
- For packages that have "global" variables / arguments, use Context (called "engine" in this project) as well.
- While making this tool available from commandline with the frameworks cobra and viper, I also tried WASM for running it in browsers.
  - WASM doesn't work in this case, because the http package isn't available there.
- Although cobra and viper were used in early versions, they seem unmaintained and were since replaced with `urfave/cli`.
