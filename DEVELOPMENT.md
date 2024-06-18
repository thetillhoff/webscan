# Development

## How to start from scratch

```sh
go mod init github.com/thetillhoff/webscan
# go install github.com/spf13/cobra-cli@latest
cobra-cli init --viper
# cobra-cli add <command>
```

## How to test

`go test ./...`

## Building

- This application is available for Linux, Mac and Windows.
- A version variable is injected at build/release time, so the version displayed with `webscan version` automatically equals the git tag it was built with.

## Decisions / best practices

- Don't have global variables in a package -> they would be obstructed for the consumer and are not threadsafe
- Don't use functional options -> they require a lot of code / maintenance. Also, having functions to set a context object every time a function is called is tedious
- Use Context (called engine in this project). Not necessarily the go-context package, but implement "instance of package" as context and use that.
- For packages that have "global" variables / arguments, use Context (called "engine" in this project) as well.
- While making this tool available from commandline with the frameworks cobra and viper, I also tried WASM for running it in browsers.
- WASM doesn't work in this case, because the http package isn't available there.
