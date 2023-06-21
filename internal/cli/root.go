/*
Package cli TODO: add description
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
}

// init initializes the root command
func init() {
	// Subcommands
	addSubCommands()

	// Version
	rootCmd.Version = "1.0.0"
}
