/*
Package cli TODO: add description
*/
package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/christosgalano/bruh/internal/apiversions"
	"github.com/christosgalano/bruh/internal/bicep"
	"github.com/spf13/cobra"
)

var (
	path     string
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
			fmt.Printf("Error: invalid output format %s\n", output)
			cmd.Usage()
			os.Exit(1)
		}

		// Invalid path
		fs, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("Error: no such file or directory %q\n", path)
			} else {
				fmt.Println(err)
			}
			os.Exit(1)
		}

		// Scan file or directory
		if fs.IsDir() {
			err = scanDirectory(path, output, outdated)
		} else {
			err = scanFile(path, output, outdated)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// Local flags

	// path - required
	scanCmd.Flags().StringVarP(&path, "path", "p", "", "path to bicep file or directory containing bicep files")
	scanCmd.MarkFlagRequired("path")

	// output - optional
	scanCmd.Flags().StringVarP(&output, "output", "o", "normal", "output format (normal, table)")

	// outdated - optional
	scanCmd.Flags().BoolVarP(&outdated, "outdated", "u", false, "show only outdated resources")
}

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
		printFileNormal(bicepFile, bicepFile.Name, outdated)
	}

	return nil
}

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
		printDirectoryNormal(bicepDirectory, outdated)
	}

	return nil
}
