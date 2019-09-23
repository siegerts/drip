# Overview

`drip` is a utility that will monitor your [Plumber](https://www.rplumber.io) applications for any changes in your source and automatically restart your server.

> This project is under development and subject to change. All feedback and issues are welcome. 🍻

The key features of drip are:

- Automatic restarting of Plumber applications on file changes 🚀
- Distributed as a single binary. Install drip by unzipping it and moving it to a directory included in your system's PATH
- Ignore specific directories
- Generate route maps

## Plumber Application Structure

drip requires that the Plumber application structure make use of an `entrypoint.R` that references a `plumber.R` app.

```r
plumber::plumb("plumber.R")$run("0.0.0.0", port=8000)
```

```r
# Packages ----
library(plumber)

# Plumb API ----
pr <- plumb("plumber.R")

pr$run("0.0.0.0",port=8000)
```

## Use

### Command: drip

drip is a utility that will monitor your Plumber applications for any changes in
your source and automatically restart your server.

#### Usage

- `drip [flags]`
- `drip [command]`

Available Commands:

- `help` Help about any command
- `routes` Display all routes in your Plumber application
- `version` Print the version number of drip
- `watch` Watch the current directory for any changes

Flags:

- `-h`, `--help` help for drip

### Command: watch

Watch and rebuild the source if any changes are made across subdirectories

#### Usage

Usage: `drip watch [flags]`

The list of available flags are:

- `-d`, `--dir` (_string_) Source directory to watch
- `-e`, `--entry` (_string_) Plumber application entrypoint file (default "`entrypoint.r`")
- `-f`, `--filter` (_string_) Filter endpoints by prefix match
- `-h`, `--help` help for watch
- `--host` (_string_) Display route endpoints with a specific host (default "127.0.0.1")
- `--port` (_int_) Display route endpoints with a specific port (default 8000)
- `--routes` Display route map alongside file watcher
- `--showHost` Display absolute route endpoint in output
- `-s`, `--skip` (_strings_) A comma-separated list of directories to not watch. (default [node_modules,.Rproj.user])

#### Examples

```sh
drip watch  --routes --showHost --host http://test.com/ --port 5464 -f sum
```

### Command: routes

A quick way to visualize your application's routing structure

#### Usage

Usage: `drip routes [flags]`

- `-e`, `--entry` (_string_) Plumber application entrypoint file (default "entrypoint.r")
- `-h`, `--help` help for routes
- All available flags for `drip watch`

#### Examples

```sh
drip watch --routes --showHost --host http://test.com/ --port 5464
```

## Developing

If you want to work on drip, you'll first need [Go](https://golang.org/) installed on your machine.

For local development, first make sure Go is properly installed and that a GOPATH has been set. You will also need to add $GOPATH/bin to your $PATH.

Next, using Git, clone this repository into \$GOPATH/src/github.com/siegerts/drip.

```sh
$ git clone github.com/siegerts/drip
```

```sh
$ go build -o build/drip  github.com/siegerts/drip
$ go install github.com/siegerts/drip
```
