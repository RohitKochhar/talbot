# talbot

Quick and easy scaffolding for go servers

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
