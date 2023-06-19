/*
Package bicep provides a comprehensive set of functions to manipulate Bicep files and directories.

It offers methods for parsing directories and files to extract valuable information regarding resource metadata, such as name and API version.

The package also includes functions to update the API versions of existing Bicep files in place or create new ones.
*/
package bicep

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/christosgalano/bruh/internal/types"
)

const (
	// pattern is the regex pattern used to match resource IDs in Bicep files
	pattern = `(?P<namespace>Microsoft\.[a-zA-Z]+)/(?P<resource>[a-zA-Z]+)@(?P<version>[0-9]{4}-[0-9]{2}-[0-9]{2}-preview|[0-9]{4}-[0-9]{2}-[0-9]{2})`
)

// validateBicepFile validates that a file exists has the .bicep extension.
// If the file does not exist, is a directory, or does not have the .bicep extension,
// an error is returned.
func validateBicepFile(path string) error {
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

// ParseFile parses a file and returns a slice of types.ResourceInfo.
func ParseFile(filename string) ([]types.ResourceInfo, error) {
	if err := validateBicepFile(filename); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	results := []types.ResourceInfo{}

	matches := regex.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		results = append(results, types.ResourceInfo{
			ID:                match[1] + "/" + match[2],
			Name:              match[2],
			Namespace:         match[1],
			CurrentAPIVersion: match[3],
		})
	}

	return results, nil
}

// ParseDirectory parses a directory and returns a map of filename to slice of types.ResourceInfo.
func ParseDirectory(dir string) (map[string][]types.ResourceInfo, error) {
	results := map[string][]types.ResourceInfo{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-bicep files
		if validateBicepFile(path) != nil {
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
