/*
Package bicep provides a comprehensive set of functions to manipulate Bicep files and directories.

It offers methods for parsing directories and files to extract valuable information regarding resource metadata, such as name and API version.

The package also includes functions to update the API versions of existing Bicep files in place or create new ones.
*/
package bicep

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/christosgalano/bruh/internal/types"
)

// UpdateFile receives a pointer to a BicepFile object and updates the file with the new API versions for each resource.
// inPlace determines whether the file should be updated in place or a new one should be created with suffix "_updated.bicep".
// includePreview determines whether preview API versions should be considered.
func UpdateFile(bicepFile *types.BicepFile, inPlace bool, includePreview bool) error {
	file, err := os.Stat(bicepFile.Name)
	if err != nil {
		return fmt.Errorf("failed to get file info %s", err)
	}

	data, err := os.ReadFile(bicepFile.Name)
	if err != nil {
		return fmt.Errorf("failed to read file %q", err)
	}
	content := string(data)

	for i := range bicepFile.Resources {
		latestAPIVersion := bicepFile.Resources[i].AvailableAPIVersions[0]

		// If we don't want to include preview versions, find the latest non-preview version
		if !includePreview && strings.HasSuffix(latestAPIVersion, "-preview") {
			for _, version := range bicepFile.Resources[i].AvailableAPIVersions {
				if !strings.HasSuffix(version, "-preview") {
					latestAPIVersion = version
					break
				}
			}
		}

		// Update the API version if needed
		if bicepFile.Resources[i].CurrentAPIVersion != latestAPIVersion {
			re := regexp.MustCompile(bicepFile.Resources[i].ID + "@" + bicepFile.Resources[i].CurrentAPIVersion)
			content = re.ReplaceAllString(content, bicepFile.Resources[i].ID+"@"+latestAPIVersion)

			// TODO: find better way to print update messages

			// Update object only if we want to update the file in place
			if inPlace {
				bicepFile.Resources[i].CurrentAPIVersion = latestAPIVersion
			}
		}

	}

	// If we don't want to update the file in place, create a new one
	modifiedFile := bicepFile.Name
	if !inPlace {
		modifiedFile = strings.Replace(modifiedFile, ".bicep", "_updated.bicep", 1)
	}

	err = os.WriteFile(modifiedFile, []byte(content), file.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to update file %s", err)
	}

	return nil
}

// UpdateDirectory receives a pointer to a BicepDirectory object and updates its files with the new API versions for each resource.
// inPlace determines whether the files should be updated in place or new ones should be created with suffix "_updated.bicep".
// includePreview determines whether preview API versions should be considered.
func UpdateDirectory(bicepDirectory *types.BicepDirectory, inPlace bool, includePreview bool) error {
	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	results := make(chan error)

	// Launch a goroutine for each file
	for i := range bicepDirectory.Files {
		wg.Add(1)
		go func(file *types.BicepFile, inPlace bool, includePreview bool) {
			defer wg.Done()
			err := UpdateFile(file, inPlace, includePreview)
			if err != nil {
				results <- err
			}
		}(&bicepDirectory.Files[i], inPlace, includePreview)
	}

	// Start a goroutine to close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Receive the results from the goroutines
	for err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}
