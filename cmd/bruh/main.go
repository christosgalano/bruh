/*
The bruh command line tool is a simple utility that parses Bicep files and directories
and prints out the resources that are defined in them.
*/
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/christosgalano/bruh/internal/azapiversions"
	"github.com/christosgalano/bruh/internal/parse"
	"github.com/christosgalano/bruh/internal/types"
)

func main() {
	dir := "/Users/galano/Developer/Christos/Development/Go/bruh/internal/parse/testdata/bicep"
	dirResults, err := parse.Directory(dir)
	if err != nil {
		log.Fatalf("parse error: %s", err)
	}

	fmt.Printf("\nDirectory: %s\n", dir)
	for filename, results := range dirResults {
		fmt.Printf("\nFile: %s\n", filename)
		for _, result := range results {
			fmt.Println(result)
		}
	}
	fmt.Printf("\n\n")

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Create a channel to receive the API versions
	ch := make(chan types.ResourceInfo)

	// Launch a goroutine for each file's resources
	for filename, resources := range dirResults {
		wg.Add(1)
		go func(filename string, resources []types.ResourceInfo) {
			defer wg.Done()
			for _, resource := range resources {
				versions, err := azapiversions.GetAPIVersions(resource)
				if err != nil {
					log.Fatalf("failed to update API versions: %s", err)
				}
				resource.AvailableAPIVersions = versions
				ch <- resource
			}
		}(filename, resources)
	}

	// Start a goroutine to close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	fmt.Printf("Available API Versions:\n\n")
	for resource := range ch {
		fmt.Println(resource)
	}
}
