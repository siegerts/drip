# drip

`drip` is a utility that will monitor your Plumber applications for any changes in
your source and automatically restart your server. Perfect for development.

## Developing

```sh
$ git clone github.com/siegerts/drip
```

```sh
$ go build -o build/drip  github.com/siegerts/drip
$ go install github.com/siegerts/drip
```

## Use

```
drip is a utility that will monitor your Plumber applications for any changes in
your source and automatically restart your server. Perfect for development.
Complete documentation is available at x

Usage:
  drip [flags]
  drip [command]

Available Commands:
  help        Help about any command
  routes      Display all routes in your Plumber application
  version     Print the version number of drip
  watch       Watch the current directory for any changes

Flags:
  -h, --help   help for drip

Use "drip [command] --help" for more information about a command.
```

### watch

```
drip watch -h
Watch and rebuild the source if any changes are made across subdirectories

Usage:
  drip watch [flags]

Flags:
  -d, --dir string      Source directory to watch
  -e, --entry string    Plumber application entrypoint file (default "entrypoint.r")
  -f, --filter string   Filter endpoints by prefix match
  -h, --help            help for watch
      --host string     Display route endpoints with a specific host (default "127.0.0.1")
      --port int        Display route endpoints with a specific port (default 8000)
      --routes          Display route map alongside file watcher
      --showHost        Display absolute route endpoint in output
  -s, --skip strings    A comma-separated list of directories to not watch. (default [node_modules,.Rproj.user])
```

### routes

```
drip routes -h
A quick way to visualize your application's routing structure

Usage:
drip routes [flags]

Flags:
-e, --entry string   Plumber application entrypoint file (default "entrypoint.r")
-h, --help           help for routes
```

### building
