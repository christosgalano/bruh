package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/christosgalano/bruh/internal/apiversions"
	"github.com/christosgalano/bruh/internal/bicep"
	"github.com/christosgalano/bruh/internal/types"
	"github.com/spf13/cobra"
)

var (
	scanPath           string
	output             string
	outdated           bool
	scanIncludePreview bool
)

// scanCmd represents the scan command.
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a bicep file or a directory containing bicep files",
	Long: `Scan a bicep file or a directory containing bicep files and
print out information regarding the API versions of Azure resources.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Invalid output format
		if output != "normal" && output != "table" && output != "markdown" {
			fmt.Fprintf(os.Stderr, "Error: invalid output format %s\n", output)
			cmd.Usage()
			os.Exit(1)
		}

		// Invalid path
		fs, err := os.Stat(scanPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "Error: no such file or directory %q\n", scanPath)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}

		// Scan file or directory
		if fs.IsDir() {
			err = scanDirectory()
		} else {
			err = scanFile()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

// init initializes the scan command.
func init() {
	// Local flags

	// path - required
	scanCmd.Flags().StringVarP(&scanPath, "path", "p", "", "path to bicep file or directory containing bicep files")
	scanCmd.MarkFlagRequired("path")

	// output - optional
	scanCmd.Flags().StringVarP(&output, "output", "o", "normal", "output format (normal, table, markdown)")

	// outdated - optional
	scanCmd.Flags().BoolVarP(&outdated, "outdated", "u", false, "show only outdated resources")

	// include-preview - optional
	scanCmd.Flags().BoolVarP(&scanIncludePreview, "include-preview", "r", false, "include preview API versions (if not set: only non-preview versions will be considered for the latest version)")

	// Examples
	scanCmd.Example = `
Scan a bicep file:
  bruh scan --path ./main.bicep

Scan a directory:
  bruh scan --path ./bicep/modules

Show only outdated resources in markdown format:
  bruh scan --path ./main.bicep --outdated --output markdown

Print output in table format including preview API versions:
  bruh scan --path ./bicep/modules --output table --include-preview`
}

// scanFile parses a file, fetches the latest API versions of Azure resources and then prints out information regarding the status of those resources.
// If outdated is true, only outdated resources are printed.
// If includePreview is true, preview API versions are also considered.
func scanFile() error {
	bicepFile, err := bicep.ParseFile(scanPath)
	if err != nil {
		return err
	}

	err = apiversions.UpdateBicepFile(bicepFile, scanIncludePreview)
	if err != nil {
		return err
	}

	switch output {
	case "normal":
		printFileNormal(bicepFile, bicepFile.Name, outdated, types.ModeScan)
	case "table":
		printFileTable(bicepFile, outdated)
	case "markdown":
		printFileMarkdown(bicepFile, outdated)
	}

	return nil
}

// scanDirectory parses a directory, fetches the latest API versions of Azure resources and then prints out information regarding the status of those resources.
// If outdated is true, only outdated resources are printed.
// If includePreview is true, preview API versions are also considered.
func scanDirectory() error {
	bicepDirectory, err := bicep.ParseDirectory(scanPath)
	if err != nil {
		return err
	}

	err = apiversions.UpdateBicepDirectory(bicepDirectory, scanIncludePreview)
	if err != nil {
		return err
	}

	switch output {
	case "normal":
		printDirectoryNormal(bicepDirectory, outdated, types.ModeScan)
	case "table":
		printDirectoryTable(bicepDirectory, outdated)
	case "markdown":
		printDirectoryMarkdown(bicepDirectory, outdated)
	}

	return nil
}
