# bruh

## Install

```bash
go install github.com/christosgalano/bruh/cmd/bruh@latest
```

## Overview

**bruh (Bicep Resource Update Helper)** is a command-line tool for scanning and updating the API version of Azure resources in Bicep files.

It offers two main commands: [**scan**](#scan) and [**update**](#update).

## Scan

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout.

It can be used to detect drift between the API versions used in the bicep files and the latest available ones.

### Example usage

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
{absolute-path}/bicep:

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

## Update

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.

### Example usage

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
{absolute-path}/bicep:

modules/compute_updated.bicep:
  + Updated Microsoft.Web/serverfarms from version 2022-03-01 to 2022-03-01
  + Updated Microsoft.Web/sites from version 2022-03-01 to 2022-03-01

modules/identity_updated.bicep:
  + Updated Microsoft.ManagedIdentity/userAssignedIdentities from version 2023-01-31 to 2023-01-31
```

**NOTE**: all the API versions are fetched from the official [Microsoft Learn website](https://learn.microsoft.com/en-us/azure/templates/).
