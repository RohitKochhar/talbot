# talbot

Quick and easy scaffolding for go servers and microservices

## Installation

### Building from source

To build `talbot` from source, select the proper `GOOS` and `GOARCH` for your platform and use the following instructions:
 
```bash
$ git clone git@github.com:RohitKochhar/talbot.git
$ cd talbot
$ GOOS={darwin,linux,windows} GOARCH={amd64,arm64} go build
$ mv talbot /usr/local/bin
```

## Usage

```bash
Usage:
  talbot make [flags]

Aliases:
  make, m

Flags:
  -n, --app-name string   Name of application (Required)
  -d, --dir string        Path to target app directory (default "./")
  -h, --help              help for make
  -m, --mod-name string   Name of top-level application go module (default $app-name)
```

## Example

Running:

```bash
$ talbot make -n example-output -d . -m github.com/rohitkochhar/talbot-output
```

Generates the application seen in this repository's `./example-output` directory.

## Contributing

For information on contributing to this project, please see [CONTRIBUTING.md](./CONTRIBUTING.md)
