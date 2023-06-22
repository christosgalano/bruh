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
	Short: "bruh is a tool for updating the API version of Azure resources in Bicep files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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
