/*
Package parser provides functions to parse Bicep files and directories,
and return information about the resources that are defined in them.
*/
package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// ResourceInfo contains information about a resource:
//   - ID: the resource ID (e.g. Microsoft.Network/virtualNetworks)
//   - Namespace: the resource namespace (e.g. Microsoft.Network)
//   - Resource: the resource name (e.g. virtualNetworks)
//   - Version: the used api version (e.g. 2021-02-01)
type ResourceInfo struct {
	ID        string
	Namespace string
	Resource  string
	Version   string
}

// ValidateBicepFile validates that a file exists has the .bicep extension.
// If the file does not exist, is a directory, or does not have the .bicep extension,
// an error is returned.
func ValidateBicepFile(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}

	if f.IsDir() {
		return fmt.Errorf("given path is a directory: %s", path)
	}

	if ext := filepath.Ext(path); ext != ".bicep" {
		return fmt.Errorf("invalid file extension: %s", ext)
	}
	return nil
}

// ParseFile parses a file and returns a slice of ResourceInfo.
func ParseFile(filename string) ([]ResourceInfo, error) {
	if err := ValidateBicepFile(filename); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(`(?P<namespace>Microsoft\.[a-zA-Z]+)/(?P<resource>[a-zA-Z]+)@(?P<version>[0-9]{4}-[0-9]{2}-[0-9]{2})`)
	if err != nil {
		return nil, err
	}

	results := []ResourceInfo{}

	matches := regex.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		results = append(results, ResourceInfo{
			ID:        match[1] + "/" + match[2],
			Namespace: match[1],
			Resource:  match[2],
			Version:   match[3],
		})
	}

	return results, nil
}

// ParseDir parses a directory and returns a map of filename to slice of ResourceInfo.
func ParseDir(dir string) (map[string][]ResourceInfo, error) {
	results := map[string][]ResourceInfo{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-bicep files
		if ValidateBicepFile(path) != nil {
			return nil
		}

		fileResults, err := ParseFile(path)
		if err != nil {
			return err
		}
		results[path] = fileResults

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
