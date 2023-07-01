# bruh

[![ci](https://github.com/christosgalano/bruh/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/christosgalano/bruh/actions/workflows/ci.yaml)

## Table of contents

- [Installation](#installation)
- [Usage](#usage)
- [GitHub Action](#github-action)

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

**bruh (Bicep Resource Update Helper)** is a command-line tool for scanning and updating the API version of Azure resources in Bicep files.

It offers two main commands: [**scan**](#scan) and [**update**](#update).

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

**NOTE**: all the API versions are fetched from the official [Microsoft Learn website](https://learn.microsoft.com/en-us/azure/templates/).

## GitHub Action

bruh can also be used as a GitHub Action to scan and update bicep files in a repository.

### Syntax

```yaml
  uses: christosgalano/bruh@v1.0.0
  with:
    command: scan | update              # command to execute (required)
    path: ./...                         # path to the bicep file or directory (required)
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

- name: Push changes
  run: |
    git config --global user.name 'Your Name'
    git config --global user.email 'your-username@users.noreply.github.com'
    git commit -am "Updated API versions of Azure resources"
    git push
```

A complete example of scanning a bicep directory, updating the outdated files, and pushing the changes:

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
      uses: actions/checkout@v3

    - name: Update bicep directory with bruh
      uses: christosgalano/bruh@v1.0.0
      with:
        command: update
        path: ${{ github.workspace }}/bicep # need path relative to workspace
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
    
    # Everything works correctly, so we can commit and push the changes.
    - name: Push changes
      run: |
        git config --global user.name 'Your Name'
        git config --global user.email 'your-username@users.noreply.github.com'
        git commit -am "Updated API versions of Azure resources"
        git push
```
