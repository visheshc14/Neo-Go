# `Neo`

`Neo` is a single file server. It responds to every `GET` request it receives with the content of a given file (specified by ENV or CLI argument), and for every other request (with any other HTTP method or path) it returns a 404.

`Neo` was invented to help with cases where a generally small file needs to be delivered at a certain path, for example [MTA STS's `/.well-known/mta-sts.txt`](https://en.wikipedia.org/wiki/MTA-STS). 

See also [`Neo-Rust`](https://github.com/visheshc14/Neo-rust)

# Quickstart

`Neo` only needs the path to a single file to run:

```console
$ Neo <file path>
```

By default, `Neo` will serve the file at host `127.0.0.1` on port `5000`. `Neo` can also take file content from STDIN like so:

```console
$ Neo <<EOF
> your file content
> goes here
> EOF
```

## Docker

To run `Neo` with docker:

```console
$ docker run --detach \
-p 5000:5000 \
-e HOST=0.0.0.0 \
-e FILE=/data/file \
-v /path/to/folder/with/your/file:/data \
--name Neo-go \
registry.gitlab.com/visheshc14/Neo-Go/cli:v2
```

Note that `v1` of `Neo-Go` used `net/http`, and `v2` uses [`fasthttp`](https://github.com/valyala/fasthttp).

# Usage

```console
$ ./Neo --help
Usage of ./Neo:
  -file string
        File to read
  -host string
        Host
  -port int
        Port (default -1)
  -stdin-read-timeout-seconds int
        Amount of seconds to wait for input on STDIN to serve (default -1)
```

# Environment Variables

| ENV variable                 | Default     | Example              | Description                                      |
|------------------------------|-------------|----------------------|--------------------------------------------------|
| `HOST`                       | `127.0.0.1` | `0.0.0.0`            | The host on which `Neo` will listen             |
| `PORT`                       | `5000`      | `3000`               | The port on which `Neo` will listen             |
| `FILE`                       | N/A         | `/path/to/your/file` | The path to the file that will be served         |
| `STDIN_READ_TIMEOUT_SECONDS` | `60`        | `10`                 | The amount of seconds to try and read from STDIN |

