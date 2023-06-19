/*
Package azapiversions provides functions to fetch and extract API versions for Azure resources.
*/
package azapiversions

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"

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

// GetAPIVersions returns a sorted list of API versions for a given resource.
func GetAPIVersions(resource types.ResourceInfo) ([]string, error) {
	url := "https://learn.microsoft.com/en-us/azure/templates/" + strings.ToLower(resource.Namespace) + "/" + strings.ToLower(resource.Name)
	pattern := `href="(\d{4}-\d{2}-\d{2}-preview|\d{4}-\d{2}-\d{2})/` + strings.ToLower(resource.Name) + `"`

	body, err := fetchResourcePage(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource page: %w", err)
	}

	versions, err := extractAPIVersions(body, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to extract API versions: %w", err)
	}

	return versions, nil
}
