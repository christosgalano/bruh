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
	"strings"
	"sync"

	"github.com/christosgalano/bruh/internal/types"
)

const (
	// pattern is the regex pattern used to match resource IDs in Bicep files
	pattern = `(?P<namespace>Microsoft\.[a-zA-Z]+)/(?P<resource>[a-zA-Z]+)@(?P<version>[0-9]{4}-[0-9]{2}-[0-9]{2}-preview|[0-9]{4}-[0-9]{2}-[0-9]{2})`
)

var (
	// cache is a synchronized map used to store the contents of Bicep files
	cache sync.Map
)

// readBicepFile reads a Bicep file and returns its contents as a byte slice.
// If the file does not exist, is a directory, or does not have the .bicep extension, the function returns an error.
func readBicepFile(filePath string) ([]byte, error) {
	f, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file does not exist %q", filePath)
		}
		return nil, err
	}

	if f.IsDir() {
		return nil, fmt.Errorf("given path is a directory %q", filePath)
	}

	if ext := filepath.Ext(filePath); ext != ".bicep" {
		return nil, fmt.Errorf("invalid file extension %q", ext)
	}

	// Check if the file is already cached
	if data, ok := cache.Load(filePath); ok {
		return data.([]byte), nil
	}

	// File is not cached, read it
	filePath = filepath.Clean(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Cache the file content
	cache.Store(filePath, data)

	return data, nil
}

// ParseFile parses a file and returns a pointer to a BicepFile object.
func ParseFile(filePath string) (*types.BicepFile, error) {

	data, err := readBicepFile(filePath)
	if err != nil {
		return nil, err
	}
	content := string(data)

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	results := []types.Resource{}

	matches := regex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		results = append(results, types.Resource{
			ID:                match[1] + "/" + match[2],
			Name:              match[2],
			Namespace:         match[1],
			CurrentAPIVersion: match[3],
		})
	}

	bicepFile := types.BicepFile{
		Path:      filePath,
		Resources: results,
	}

	return &bicepFile, nil
}

// ParseDirectory parses a directory and returns a pointer to a BicepDirectory object.
func ParseDirectory(dirPath string) (*types.BicepDirectory, error) {
	bicepDir := types.BicepDirectory{
		Path: dirPath,
	}

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		file, err := ParseFile(path)
		if err != nil {
			// Ignore directories and files with invalid extensions
			if strings.Contains(err.Error(), "given path is a directory") || strings.Contains(err.Error(), "invalid file extension") {
				return nil
			}
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
