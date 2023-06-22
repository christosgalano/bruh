/*
Package cli provides a command-line interface (CLI) for the bruh tool, utilizing cobra-cli. It offers two main commands: scan and update.

The scan command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and prints the results to stdout. For full usage details, run "bruh scan --help" or "bruh help scan".

The update command parses the given bicep file or directory, fetches the latest API versions for each Azure resource referenced in the file(s),
and updates the file(s) in place or creates new ones with the "_updated.bicep" extension.
For full usage details, run "bruh update --help" or "bruh help update".
*/
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bruh",
	Short: "bruh is a command-line tool for updating the API version of Azure resources in Bicep files",
	Long: `bruh (Bicep Resource Update Helper) is a command-line tool for updating the API version of Azure resources in Bicep files.

It can be used to scan a Bicep file or directory and print out information regarding the API versions of used Azure resources.
bruh can also be used to update all the resources to the latest API version available either in place or by creating new files with the "_updated.bicep" extension.
All the API versions are fetched from the official Microsoft learn website (https://learn.microsoft.com/en-us/azure/templates/).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute executes the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// addSubCommands adds subcommands to the root command
func addSubCommands() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(updateCmd)
}

// init initializes the root command
func init() {
	// Subcommands
	addSubCommands()

	// Version
	rootCmd.Version = "1.0.0"
}
