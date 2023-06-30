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
	data, ok := cache.Load(bicepFile.Path)

	// If the file is not cached, read it (this should never happen)
	if !ok {
		d, err := readBicepFile(bicepFile.Path)
		if err != nil {
			return err
		}
		data = d
	}
	content := string(data.([]byte))

	// Update the API versions for each resource - if needed
	for i := range bicepFile.Resources {
		latestAPIVersion := bicepFile.Resources[i].AvailableAPIVersions[0]
		if bicepFile.Resources[i].CurrentAPIVersion != latestAPIVersion {
			re := regexp.MustCompile(bicepFile.Resources[i].ID + "@" + bicepFile.Resources[i].CurrentAPIVersion)
			content = re.ReplaceAllString(content, bicepFile.Resources[i].ID+"@"+latestAPIVersion)
			bicepFile.Resources[i].CurrentAPIVersion = latestAPIVersion
		}

	}

	// Use the same permissions as the original file
	f, err := os.Stat(bicepFile.Path)
	if err != nil {
		return err
	}

	// If the file is not updated in place, create a new one with the suffix "_updated.bicep",
	// and remove the original file from the cache
	if !inPlace {
		cache.Delete(bicepFile.Path)
		bicepFile.Path = strings.Replace(bicepFile.Path, ".bicep", "_updated.bicep", 1)
	}

	// Write the updated content to the file
	err = os.WriteFile(bicepFile.Path, []byte(content), f.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to update file %s", err)
	}

	// Cache the new content appropriately
	cache.Store(bicepFile.Path, []byte(content))

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
			if err := UpdateFile(file, inPlace); err != nil {
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
