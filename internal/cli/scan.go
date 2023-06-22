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
	scanPath string
	output   string
	outdated bool
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a bicep file or a directory containing bicep files",
	Long: `Scan a bicep file or a directory containing bicep files
and print out information regarding the API versions of Azure resources`,
	Run: func(cmd *cobra.Command, args []string) {
		// Invalid output format
		if output != "normal" && output != "table" {
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
			err = scanDirectory(scanPath, output, outdated)
		} else {
			err = scanFile(scanPath, output, outdated)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Local flags

	// path - required
	scanCmd.Flags().StringVarP(&scanPath, "path", "p", "", "path to bicep file or directory containing bicep files")
	scanCmd.MarkFlagRequired("path")

	// output - optional
	scanCmd.Flags().StringVarP(&output, "output", "o", "normal", "output format (normal, table)")

	// outdated - optional
	scanCmd.Flags().BoolVarP(&outdated, "outdated", "u", false, "show only outdated resources")
}

// scanFile parses a file, fetches the latest API versions of Azure resources and then prints out information regarding the status of those resources.
// If outdated is true, only outdated resources are printed.
func scanFile(path string, output string, outdated bool) error {
	bicepFile, err := bicep.ParseFile(path)
	if err != nil {
		return err
	}

	err = apiversions.UpdateBicepFile(bicepFile)
	if err != nil {
		return err
	}

	if output == "table" {
		printFileTable(bicepFile, outdated)
	} else {
		printFileNormal(bicepFile, bicepFile.Name, outdated, types.ModeScan)
	}

	return nil
}

// scanDirectory parses a directory, fetches the latest API versions of Azure resources and then prints out information regarding the status of those resources.
// If outdated is true, only outdated resources are printed.
func scanDirectory(path string, output string, outdated bool) error {
	bicepDirectory, err := bicep.ParseDirectory(path)
	if err != nil {
		return err
	}

	err = apiversions.UpdateBicepDirectory(bicepDirectory)
	if err != nil {
		return err
	}

	if output == "table" {
		printDirectoryTable(bicepDirectory, outdated)
	} else {
		printDirectoryNormal(bicepDirectory, outdated, types.ModeScan)
	}

	return nil
}
