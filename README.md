# Overview

`drip` is an easy-to-use development utility that will monitor your [Plumber](https://www.rplumber.io) applications for any changes in your source and automatically restart your server.

> This project is under development and subject to change. All feedback and issues are welcome. 🍻

The key features of drip are:

- Automatic restarting of Plumber applications on file changes 🚀
- Distributed as a single binary. Install drip by unzipping it and moving it to a directory included in your system's PATH
- Ignore specific directories
- Generate and watch route maps

# Requirements

drip utilizes [Rscript](https://support.rstudio.com/hc/en-us/articles/218012917-How-to-run-R-scripts-from-the-command-line) to run the Plumber application process. For that reason, R is required for the CLI to correctly execute.

# Plumber Application Structure

drip requires that the Plumber application structure make use of an `entrypoint.R` that references a `plumber.R` app.

```r
# entrypoint.R

plumber::plumb("plumber.R")$run("0.0.0.0", port=8000)
```

```r
# entrypoint.R

library(plumber)

pr <- plumb("plumber.R")

pr$run("0.0.0.0", port=8000)
```

# Use

## Command: drip

Watch the current directory for changes using default option flag parameters.

### Usage

- `drip [flags]`
- `drip [command]`

Available Commands:

- `help` Help about any command
- `routes` Display all routes in your Plumber application
- `version` Print the version number of drip
- `watch` Watch the current directory for any changes

Flags:

- `-h`, `--help` help for drip

### Example

```sh
# cd into project
$ drip

[project-dir] skipping directory: .Rproj.user
[project-dir] skipping directory: node_modules
[project-dir] plumbing...
[project-dir] running: Rscript /project-dir/entrypoint.r
[project-dir] watching...
Starting server to listen on port 8000

[project-dir] modified file: /project-dir/plumber.R
[project-dir] plumbing...
[project-dir] running: Rscript /project-dir/entrypoint.r
[project-dir] watching...

Starting server to listen on port 8000

```

## Command: watch

Watch and rebuild the source if any changes are made across subdirectories

### Usage

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

### Examples

```sh
# cd into project
$ drip watch  --routes

[project-dir] skipping directory: .Rproj.user
[project-dir] skipping directory: node_modules
[project-dir] plumbing...
[project-dir] running: Rscript /project-dir/entrypoint.r
[project-dir] routing...

+--------------+----------------------------+---------------+
| PLUMBER VERB |          ENDPOINT          |    HANDLER    |
+--------------+----------------------------+---------------+
| @get         | /echo                      | function      |
| @get         | /dynamic/<param1>/<param2> | function      |
| @get         | /two                       | function      |
| @get         | /plot                      | function      |
| @post        | /sum                       | function      |
| @get         | /req                       | function      |
| @assets      | ./files/static             | static assets |
+--------------+----------------------------+---------------+

[project-dir] watching...
Starting server to listen on port 8000
```

Or, display routes with an absolute URI and port.

```sh
# cd into project
$ drip watch --routes --showHost --host http://localhost

[project-dir] skipping directory: .Rproj.user
[project-dir] skipping directory: node_modules
[project-dir] plumbing...
[project-dir] running: Rscript /project-dir/entrypoint.r
[project-dir] routing...

+--------------+-------------------------------------------------+---------------+
| PLUMBER VERB |                    ENDPOINT                     |    HANDLER    |
+--------------+-------------------------------------------------+---------------+
| @get         | http://localhost:8000/echo                      | function      |
| @get         | http://localhost:8000/dynamic/<param1>/<param2> | function      |
| @get         | http://localhost:8000/plot                      | function      |
| @post        | http://localhost:8000/sum                       | function      |
| @get         | http://localhost:8000/req                       | function      |
| @assets      | http://localhost:8000/files/static              | static assets |
+--------------+-------------------------------------------------+---------------+

[project-dir] watching...
Starting server to listen on port 8000
```

## Command: routes

A quick way to visualize your application's routing structure without starting the watcher

### Usage

Usage: `drip routes [flags]`

- `-e`, `--entry` (_string_) Plumber application entrypoint file (default "entrypoint.r")
- `-h`, `--help` help for routes

### Examples

```sh
$ drip routes
```

## Command: completion

Generate `bash` completion commands for drip

### Usage

Usage: `drip completion [flags]`

- `-h`, `--help` help for completion

# Developing

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
