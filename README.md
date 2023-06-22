# bruh

## Install

```bash
go install github.com/christosgalano/bruh/cmd/bruh
```

## Overview

bruh (Bicep Resource Update Helper) is a command-line tool for scanning and updating the API version of Azure resources in bicep files.

It offers two main commands: scan and update.

## Scan

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout.

Example usage:

Scan a bicep file:

```bash
bruh scan --path ./main.bicep
```

Scan a directory:

```bash
bruh scan --path ./bicep/modules
```

Show only outdated resources:

```bash
bruh scan --path ./main.bicep --outdated
```

Print output in table format:

```bash
bruh scan --path ./bicep/modules --output table
```

For full usage details, run `bruh scan --help` or `bruh help scan`.

## Update

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.

Example usage:

Update a bicep file in place:

```bash
bruh update --path ./main.bicep --in-place
```

Create a new bicep file with the "_updated.bicep" extension:

```bash
bruh update --path ./main.bicep
```

Update a directory in place including preview API versions:

```bash
bruh update --path ./bicep/modules --in-place --include-preview
```

Use silent mode:

```bash
bruh update --path ./main.bicep --silent
```

For full usage details, run `bruh update --help` or `bruh help update`.

**NOTE**: all the API versions are fetched from the official [Microsoft Learn website](https://learn.microsoft.com/en-us/azure/templates/).
