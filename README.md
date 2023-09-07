# bruh

[![ci](https://github.com/christosgalano/bruh/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/christosgalano/bruh/actions/workflows/ci.yaml)
[![Code Coverage](https://img.shields.io/badge/coverage-90.9%25-31C754)](https://img.shields.io/badge/coverage-90.9%25-31C754)
[![Go Report Card](https://goreportcard.com/badge/github.com/christosgalano/bruh)](https://goreportcard.com/report/github.com/christosgalano/bruh)
[![Go Reference](https://pkg.go.dev/badge/github.com/christosgalano/bruh.svg)](https://pkg.go.dev/github.com/christosgalano/bruh)

## Table of contents

- [Description](#description)
- [Installation](#installation)
- [Usage](#usage)
- [Autocompletion](#autocompletion)
- [GitHub Action](#github-action)
- [License](#license)

## Description

**bruh (Bicep Resource Update Helper)** is a command-line tool for scanning and updating the API version of Azure resources in Bicep files.

## Installation

### Homebrew

```bash
brew tap christosgalano/christosgalano
brew install bruh
```

### Go

```bash
go install github.com/christosgalano/bruh/cmd/bruh@latest
```

### Binary

Download the latest binary from the [releases page](https://github.com/christosgalano/bruh/releases/latest).

## Usage

bruh offers two main commands: [**scan**](#scan) and [**update**](#update).

> **NOTE**: bruh does not validate if your current resource declaration matches with the new API schema.

### Scan

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout.

It can be used to detect drift between the API versions used in the bicep files and the latest available ones.

Example usage:

Scan a bicep file and print the results using the normal format:

```text
> bruh scan --path ./bicep/modules/compute.bicep
./bicep/modules/compute.bicep:
  - Microsoft.Web/serverfarms is using 2021-01-15 while the latest version is 2022-03-01
  - Microsoft.Web/sites is using 2019-08-01 while the latest version is 2022-03-01
```

Scan a directory and print only outdated resources using the table format:

```text
> bruh scan --path ./bicep --output table --outdated
./bicep:

+------------------------+--------------------------------------------------+---------------------+--------------------+
|          FILE          |                     RESOURCE                     | CURRENT API VERSION | LATEST API VERSION |
+------------------------+--------------------------------------------------+---------------------+--------------------+
| modules/compute.bicep  | Microsoft.Web/serverfarms                        |     2022-03-01      |     2022-03-01     |
+                        +--------------------------------------------------+---------------------+--------------------+
|                        | Microsoft.Web/sites                              |     2022-03-01      |     2022-03-01     |
+------------------------+--------------------------------------------------+---------------------+--------------------+
| modules/identity.bicep | Microsoft.ManagedIdentity/userAssignedIdentities | 2022-01-31-preview  |     2023-01-31     |
+------------------------+--------------------------------------------------+---------------------+--------------------+
```

### Update

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.

Example usage:

Update a bicep file in place:

```text
> bruh update --path ./bicep/modules/compute.bicep --in-place
./bicep/modules/compute.bicep:
  + Updated Microsoft.Web/serverfarms from version 2022-03-01 to 2022-03-01
  + Updated Microsoft.Web/sites from version 2022-03-01 to 2022-03-01
```

Update a directory and create new files with the "_updated.bicep" extension, including preview API versions:

```text
> bruh update --path ./bicep --include-preview
./bicep:

modules/compute_updated.bicep:
  + Updated Microsoft.Web/serverfarms from version 2022-03-01 to 2022-03-01
  + Updated Microsoft.Web/sites from version 2022-03-01 to 2022-03-01

modules/identity_updated.bicep:
  + Updated Microsoft.ManagedIdentity/userAssignedIdentities from version 2023-01-31 to 2023-01-31
```

> **NOTE**: all the API versions are fetched from the official [Microsoft Learn website](https://learn.microsoft.com/en-us/azure/templates/).

## Autocompletion

bruh provides autocompletion support. You can generate the autocompletion script for bruh specific to your shell by using the `bruh completion` command.

Supported shells are `bash`, `fish`, `zsh`, and `powershell`.

## GitHub Action

bruh can also be used as a GitHub Action to scan and update bicep files in a repository.

### Syntax

```yaml
  uses: christosgalano/bruh@v1.0.0
  with:
    command: scan | update              # command to execute (required)
    path: ./...                         # path to the bicep file or directory (required), relative to github.workspace
    include-preview: true | false       # whether to include preview API versions (optional, default: false)
    summary: true | false               # whether to print a step summary of the results (optional, default: false)
    
    # scan command only
    output: normal | table | markdown   # output format for scan command (optional, default: normal)
    outdated: true | false              # whether to print only outdated resources with scan command (optional, default: false)
    
    # update command only
    in-place: true | false              # whether to update the bicep file(s) in place or create new ones with the "_updated.bicep" extension (optional, default: true)
    silent: true | false                # whether to suppress all output (optional, default: false)
```

### Examples

Scan a bicep directory, print the results using the normal format, and generate a step summary:

```yaml
- name: Scan bicep directory with bruh
  uses: christosgalano/bruh@v1.0.0
  with:
    command: scan
    path: ./bicep
    output: normal
    summary: true
```

Update a bicep file in place and suppress all output:

```yaml
- name: Update bicep file with bruh
  uses: christosgalano/bruh@v1.0.0
  with:
    command: update
    path: ./bicep/modules/compute.bicep
    in-place: true
    silent: true
```

A complete example can be found below. It consists of the following steps:

1. Checkout the repository
2. Update the API versions of Azure resources in the bicep directory (in place)
3. Lint the main template
4. Validate the main template
5. Commit the changes - if any
6. Push the changes - if needed

```yaml
validate:
  runs-on: ubuntu-latest
  permissions:
    contents: write
  defaults:
    run:
      shell: bash
      working-directory: bicep
  steps:
    - name: Checkout
      uses: actions/checkout@v4

    - name: Update bicep directory with bruh
      uses: christosgalano/bruh@v1.0.0
      with:
        command: update
        path: ./bicep # path relative to workspace
        summary: true
        in-place: true
        include-preview: true

    # Here we catch errors that might occur if the new API versions are not
    # compatible with the used declarations.
    - name: Lint template
      run: az bicep build --file main.bicep
    
    - name: Validate template
      run: |
        az deployment sub validate \
        --name "${{ vars.DEPLOYMENT_NAME }}" \
        --location "${{ vars.LOCATION }}" \
        --template-file main.bicep \
        --parameters main.parameters.json
    
    # Everything works correctly, so we can commit and push the changes - if any.
    - name: Commit changes
      id: check-changes
      run: |
        git config --local user.name "github-actions[bot]"
        git config --local user.email "github-actions[bot]@users.noreply.github.com"
        git add .
        git commit -m "Updated API versions of Azure resources"
        git diff --quiet --exit-code -- .
        if [ $? -ne 0 ]; then
          echo "changed=true" >> $GITHUB_ENV
        fi
    
    - name: Push changes - if needed
      if: steps.check-changes.outputs.changed == 'true'
      uses: ad-m/github-push-action@master
      with:
        branch: ${{ github.ref }}
        github_token: ${{ secrets.GITHUB_TOKEN }}
```

## License

This project is licensed under the [Apache License 2.0 License](LICENSE).
