/*
Copyright © 2023 Christos Galanopoulos

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
bruh (Bicep Resource Update Helper) is a command-line tool for scanning and updating the API version of Azure resources in bicep files.

It offers two main commands: scan and update.

# Scan

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout.

Example usage:

Scan a bicep file:

	bruh scan --path ./main.bicep

Scan a directory:

	bruh scan --path ./bicep/modules

Show only outdated resources:

	bruh scan --path ./main.bicep --outdated

Print output in table format:

	bruh scan --path ./bicep/modules --output table

For full usage details, run "bruh scan --help" or "bruh help scan".

# Update

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.

Example usage:

Update a bicep file in place:

	bruh update --path ./main.bicep --in-place

Create a new bicep file with the "_updated.bicep" extension:

	bruh update --path ./main.bicep

Update a directory in-place including preview API versions:

	bruh update --path ./bicep/modules --in-place --include-preview

Use silent mode:

	bruh update --path ./main.bicep --silent

For full usage details, run "bruh update --help" or "bruh help update"

Note: all the API versions are fetched from the official Microsoft learn website (https://learn.microsoft.com/en-us/azure/templates/).
*/
package main

import (
	"github.com/christosgalano/bruh/internal/cli"
)

func main() {
	cli.Execute()
}
