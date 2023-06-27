/*
Package apiversions provides functions to fetch and update API versions for Azure resources in a bicep file or directory.
The API versions are fetched from the official Microsoft Learn website (https://learn.microsoft.com/en-us/azure/templates/).
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
	"time"

	"github.com/christosgalano/bruh/internal/types"
)

const (
	// dateFormat is the format used for parsing the dates in the API versions.
	dateFormat = "2006-01-02"
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

// extractAPIVersions extracts all the API versions using a regex pattern and returns them sorted in descending order.
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

	// Sort the versions in descending order (newest first).
	// Non-preview versions should be sorted after preview versions.
	// For example, "2021-02-02-preview" should be after "2021-02-02".
	sort.Slice(versions, func(i, j int) bool {
		// Extract the dates from the strings
		dateI := strings.Split(versions[i], "-preview")[0]
		dateJ := strings.Split(versions[j], "-preview")[0]

		// Parse the dates
		parsedDateI, errorI := time.Parse(dateFormat, dateI)
		parsedDateJ, errorJ := time.Parse(dateFormat, dateJ)

		// If the dates cannot be parsed, return the string comparison
		if errorI != nil || errorJ != nil {
			return versions[i] > versions[j]
		}

		// Compare the dates
		if parsedDateI.Equal(parsedDateJ) {
			return !strings.HasSuffix(versions[i], "-preview")
		}

		return parsedDateI.After(parsedDateJ)
	})

	return versions, nil
}

// UpdateResource updates the available API versions for a given resource.
// If includePreview is true, preview API versions will be included.
func UpdateResource(resource *types.Resource, includePreview bool) error {
	url := "https://learn.microsoft.com/en-us/azure/templates/" + strings.ToLower(resource.Namespace) + "/" + strings.ToLower(resource.Name)

	var pattern string
	if includePreview {
		pattern = `href="(\d{4}-\d{2}-\d{2}-preview|\d{4}-\d{2}-\d{2})/` + strings.ToLower(resource.Name) + `"`
	} else {
		pattern = `href="(\d{4}-\d{2}-\d{2})/` + strings.ToLower(resource.Name) + `"`
	}

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
// If includePreview is true, preview API versions will be included.
func UpdateBicepFile(bicepFile *types.BicepFile, includePreview bool) error {
	for i := range bicepFile.Resources {
		err := UpdateResource(&bicepFile.Resources[i], includePreview)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateBicepDirectory updates the available API versions for all resources in all bicep files of a given bicep directory.
// If includePreview is true, preview API versions will be included.
func UpdateBicepDirectory(bicepDirectory *types.BicepDirectory, includePreview bool) error {
	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	results := make(chan error)

	// Launch a goroutine for each file
	for i := range bicepDirectory.Files {
		wg.Add(1)
		go func(file *types.BicepFile) {
			defer wg.Done()
			err := UpdateBicepFile(file, includePreview)
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
