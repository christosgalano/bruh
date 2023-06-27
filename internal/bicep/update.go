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
// inPlace determines whether the function will update the file in place or create a new one with the suffix "_updated.bicep".
func UpdateFile(bicepFile *types.BicepFile, inPlace bool) error {
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

		// Update the API version if needed
		if bicepFile.Resources[i].CurrentAPIVersion != latestAPIVersion {
			re := regexp.MustCompile(bicepFile.Resources[i].ID + "@" + bicepFile.Resources[i].CurrentAPIVersion)
			content = re.ReplaceAllString(content, bicepFile.Resources[i].ID+"@"+latestAPIVersion)
			bicepFile.Resources[i].CurrentAPIVersion = latestAPIVersion
		}

	}

	// If we don't want to update the file in place, create a new one
	modifiedFile := bicepFile.Name
	if !inPlace {
		modifiedFile = strings.Replace(modifiedFile, ".bicep", "_updated.bicep", 1)

		// Entry now points to the new file
		bicepFile.Name = modifiedFile
	}

	err = os.WriteFile(modifiedFile, []byte(content), file.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to update file %s", err)
	}

	return nil
}

// UpdateDirectory receives a pointer to a BicepDirectory object and updates its files with the new API versions for each resource.
// inPlace determines whether the function will update each file in place or create a new one with the suffix "_updated.bicep".
func UpdateDirectory(bicepDirectory *types.BicepDirectory, inPlace bool) error {
	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	results := make(chan error)

	// Launch a goroutine for each file
	for i := range bicepDirectory.Files {
		wg.Add(1)
		go func(file *types.BicepFile, inPlace bool) {
			defer wg.Done()
			err := UpdateFile(file, inPlace)
			if err != nil {
				results <- err
			}
		}(&bicepDirectory.Files[i], inPlace)
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
