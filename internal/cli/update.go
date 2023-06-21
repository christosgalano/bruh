/*
Package cli TODO: add description
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a bicep file or a directory containing bicep files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		// inPlace, _ := cmd.Flags().GetBool("in-place")
		silent, _ := cmd.Flags().GetBool("silent")

		if silent {
			os.Stdout, _ = os.Open(os.DevNull)
		}

		f, err := os.Stat(path)
		if err != nil {
			fmt.Printf("no such file or directory: %s\n", path)
			os.Exit(1)
		}

		// TODO: add logic for file or directory
		if f.IsDir() {
			fmt.Println("updating directory...")
		} else {
			fmt.Println("updating file...")
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
	updateCmd.Flags().StringVarP(&path, "path", "p", "", "path to bicep file or directory containing bicep files")
	updateCmd.MarkFlagRequired("path")

	// in-place - optional
	updateCmd.Flags().BoolP("in-place", "i", false, "update the bicep files in place (if not set: create new files with \"_updated.bicep\" extension)")

	// silent - optional
	updateCmd.Flags().BoolP("silent", "s", false, "silent mode (no output)")

}
