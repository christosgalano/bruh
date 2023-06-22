# bruh

## Install

```bash
go install github.com/christosgalano/bruh/cmd/bruh
```

## Overview

bruh (Bicep Resource Update Helper) is a command-line tool for scanning and updating the API version of Azure resources in Bicep files.

```bash
./bruh
Usage:
  bruh [flags]
  bruh [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  scan        Scan a bicep file or a directory containing bicep files
  update      Update a bicep file or a directory containing bicep files

Flags:
  -h, --help      help for bruh
  -v, --version   version for bruh

Use "bruh [command] --help" for more information about a command.
```

It offers two main commands: [**scan**](#scan) and [**update**](#update).

## Scan

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout.

```bash
./bruh help scan
Scan a bicep file or a directory containing bicep files and
print out information regarding the API versions of Azure resources

Usage:
  bruh scan [flags]

Examples:

Scan a bicep file:
  bruh scan --path ./main.bicep

Scan a directory:
  bruh scan --path ./bicep/modules

Show only outdated resources:
  bruh scan --path ./main.bicep --outdated

Print output in table format:
  bruh scan --path ./bicep/modules --output table

Flags:
  -h, --help            help for scan
  -u, --outdated        show only outdated resources
  -o, --output string   output format (normal, table) (default "normal")
  -p, --path string     path to bicep file or directory containing bicep files
```

## Update

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.

```bash
Update a bicep file or a directory containing bicep files so that each Azure resource uses the latest API version available.
It is possible to update the files in place or create new files with "_updated.bicep" extension.

Usage:
  bruh update [flags]

Examples:

Update a bicep file in place:
  bruh update --path ./main.bicep --in-place

Update a directory including preview API versions:
  bruh update --path ./bicep/modules --include-preview

Use silent mode:
  bruh update --path ./main.bicep --silent

Flags:
  -h, --help              help for update
  -i, --in-place          update the bicep files in place (if not set: create new files with "_updated.bicep" extension)
  -r, --include-preview   include preview API versions (if not set: only non-preview versions will be considered)
  -p, --path string       path to bicep file or directory containing bicep files
  -s, --silent            silent mode (no output)
```

**NOTE**: all the API versions are fetched from the official [Microsoft Learn website](https://learn.microsoft.com/en-us/azure/templates/).
