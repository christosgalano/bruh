/*
Package bicep provides a set of functions to manipulate Bicep files and directories.

It offers methods for parsing directories and files to extract valuable information regarding resource metadata, such as name and API version.
The two main functions are ParseDirectory and ParseFile, which receive a directory or file path, and return a pointer to a BicepDirectory or BicepFile object.

The package also includes functions to update the API versions of existing Bicep files in place or create new ones.
This can be done by calling UpdateDirectory or UpdateFile, which receive a pointer to a BicepDirectory or BicepFile object.
*/
package bicep

import (
	"errors"
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

// validateBicepFile validates that a file exists and has the .bicep extension.
// If the file does not exist, is a directory, or does not have the .bicep extension, the function returns an error.
func validateBicepFile(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no such file or directory %q", path)
		}
		return err
	}

	if f.IsDir() {
		return fmt.Errorf("given path is a directory %q", path)
	}

	if ext := filepath.Ext(path); ext != ".bicep" {
		return fmt.Errorf("invalid file extension %q", ext)
	}

	return nil
}

// ParseFile parses a file and returns a pointer to a BicepFile object.
func ParseFile(filename string) (*types.BicepFile, error) {
	if err := validateBicepFile(filename); err != nil {
		return nil, err
	}

	filename = filepath.Clean(filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	results := []types.Resource{}

	matches := regex.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		results = append(results, types.Resource{
			ID:                match[1] + "/" + match[2],
			Name:              match[2],
			Namespace:         match[1],
			CurrentAPIVersion: match[3],
		})
	}

	bicepFile := types.BicepFile{
		Name:      filename,
		Resources: results,
	}

	return &bicepFile, nil
}

// ParseDirectory parses a directory and returns a pointer to a BicepDirectory object.
func ParseDirectory(dir string) (*types.BicepDirectory, error) {
	bicepDir := types.BicepDirectory{
		Name: dir,
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-bicep files
		if validateBicepFile(path) != nil {
			return nil
		}

		file, err := ParseFile(path)
		if err != nil {
			return err
		}
		bicepDir.Files = append(bicepDir.Files, *file)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &bicepDir, nil
}
