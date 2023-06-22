/*
Package apiversions provides functions to fetch and update API versions for Azure resources in a bicep file or directory.
*/
package apiversions

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/christosgalano/bruh/internal/types"
)

// fetchResourcePage fetches the HTML content of a given URL.
func fetchResourcePage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extractAPIVersions extracts and returns all the API versions sorted using a regex pattern.
func extractAPIVersions(body string, pattern string) ([]string, error) {
	versions := []string{}

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(body, -1)

	for _, match := range matches {
		versions = append(versions, match[1])
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no API versions found")
	}

	// Sort the versions in descending order (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})

	return versions, nil
}

// UpdateResource updates the available API versions for a given resource.
func UpdateResource(resource *types.Resource) error {
	url := "https://learn.microsoft.com/en-us/azure/templates/" + strings.ToLower(resource.Namespace) + "/" + strings.ToLower(resource.Name)
	pattern := `href="(\d{4}-\d{2}-\d{2}-preview|\d{4}-\d{2}-\d{2})/` + strings.ToLower(resource.Name) + `"`

	body, err := fetchResourcePage(url)
	if err != nil {
		return err
	}

	versions, err := extractAPIVersions(body, pattern)
	if err != nil {
		return err
	}

	resource.AvailableAPIVersions = versions

	return nil
}

// UpdateBicepFile updates the available API versions for all resources in a given bicep file.
func UpdateBicepFile(bicepFile *types.BicepFile) error {
	for i := range bicepFile.Resources {
		err := UpdateResource(&bicepFile.Resources[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateBicepDirectory updates the available API versions for all resources in all bicep files in a given bicep directory.
func UpdateBicepDirectory(bicepDirectory *types.BicepDirectory) error {
	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	results := make(chan error)

	// Launch a goroutine for each file
	for i := range bicepDirectory.Files {
		wg.Add(1)
		go func(file *types.BicepFile) {
			defer wg.Done()
			err := UpdateBicepFile(file)
			if err != nil {
				results <- err
			}
		}(&bicepDirectory.Files[i])
	}

	// Start a goroutine to wait for all goroutines to finish
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
