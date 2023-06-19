/*
Package bicep provides a comprehensive set of functions to manipulate Bicep files and directories.

It offers methods for parsing directories and files to extract valuable information regarding resource metadata, such as name and API version.

The package also includes functions to update the API versions of existing Bicep files in place or create new ones.
*/
package bicep

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/christosgalano/bruh/internal/types"
)

// UpdateFile receives the name of a file and its resources, and updates the file with the new API versions.
// inPlace determines whether the file should be updated in place or a new one should be created with suffix "_updated.bicep".
// includePreview determines whether preview API versions should be considered.
func UpdateFile(filename string, resources []types.ResourceInfo, inPlace bool, includePreview bool) error {
	err := validateBicepFile(filename)
	if err != nil {
		return fmt.Errorf("failed to validate file: %s", err)
	}

	file, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("failed to get file info: %s", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	content := string(data)

	// TODO: find better way to print update messages
	fmt.Printf("%s:\n", filepath.Base(filename))

	for _, resource := range resources {
		latestAPIVersion := resource.AvailableAPIVersions[0]

		// If we don't want to include preview versions, find the latest non-preview version
		if !includePreview && strings.HasSuffix(latestAPIVersion, "-preview") {
			for _, version := range resource.AvailableAPIVersions {
				if !strings.HasSuffix(version, "-preview") {
					latestAPIVersion = version
					break
				}
			}
		}

		// Update the API version if needed
		if resource.CurrentAPIVersion != latestAPIVersion {
			re := regexp.MustCompile(resource.ID + "@" + resource.CurrentAPIVersion)
			content = re.ReplaceAllString(content, resource.ID+"@"+latestAPIVersion)

			// TODO: find better way to print update messages
			fmt.Printf(" - updated %s from %s to %s\n", resource.ID, resource.CurrentAPIVersion, latestAPIVersion)
		}

	}

	// If we don't want to update the file in place, create a new one
	if !inPlace {
		filename = strings.Replace(filename, ".bicep", "_updated.bicep", 1)
	}

	err = os.WriteFile(filename, []byte(content), file.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to update file: %s", err)
	}

	return nil
}

// UpdateDirectory receives a map of filenames and their resources and updates the files with the new API versions.
// inPlace determines whether the files should be updated in place or new ones should be created with suffix "_updated.bicep".
// includePreview determines whether preview API versions should be considered.
func UpdateDirectory(files map[string][]types.ResourceInfo, inPlace bool, includePreview bool) error {
	for filename, resources := range files {
		err := UpdateFile(filename, resources, inPlace, includePreview)
		if err != nil {
			return fmt.Errorf("failed to update file: %s", err)
		}
		fmt.Println()
	}
	return nil
}
