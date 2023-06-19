/*
TODO: add description
*/
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/christosgalano/bruh/internal/azapiversions"
	"github.com/christosgalano/bruh/internal/bicep"
	"github.com/christosgalano/bruh/internal/types"
)

func main() {
	dir := "/Users/galano/Developer/Christos/Development/Go/bruh/testdata"
	dirResults, err := bicep.ParseDirectory(dir)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	type info struct {
		filename string
		index    int
		resource types.ResourceInfo
	}

	// Create a channel to receive the updated resource information
	ch := make(chan types.UpdateResourceInfo)

	// Launch a goroutine for each file's resources
	for filename, resources := range dirResults {
		wg.Add(1)
		go func(filename string, resources []types.ResourceInfo) {
			defer wg.Done()
			for index, resource := range resources {
				versions, err := azapiversions.GetAPIVersions(resource)
				if err != nil {
					log.Fatalf("failed to update API versions: %s", err)
				}
				resource.AvailableAPIVersions = versions
				ch <- types.UpdateResourceInfo{Filename: filename, Index: index, APIVersions: versions}
			}
		}(filename, resources)
	}

	// Start a goroutine to close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Receive the updated resource information and update the map
	for info := range ch {
		dirResults[info.Filename][info.Index].AvailableAPIVersions = info.APIVersions
	}

	for filename, resources := range dirResults {
		fmt.Printf("%s:\n", filepath.Base(filename))
		for _, resource := range resources {
			fmt.Println(resource)
		}
	}

	err = bicep.UpdateDirectory(dirResults, false, false)
	if err != nil {
		log.Fatalf("failed to update API versions: %s", err)
	}
}
