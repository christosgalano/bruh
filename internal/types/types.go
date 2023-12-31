/*
Package types provides shared types used by multiple packages in the "bruh" application.
*/
package types

import (
	"fmt"
	"path/filepath"
)

// Resource contains information about a resource:
//   - ID: the resource ID (e.g. Microsoft.Network/virtualNetworks)
//   - Name: the resource name (e.g. virtualNetworks)
//   - Namespace: the resource namespace (e.g. Microsoft.Network)
//   - CurrentAPIVersion: the used API version (e.g. 2021-02-01)
//   - AvailableAPIVersions: the available API versions (e.g. [2021-02-01 2020-11-01])
type Resource struct {
	ID                   string
	Name                 string
	Namespace            string
	CurrentAPIVersion    string
	AvailableAPIVersions []string
}

// String returns a string representation of a types.Resource object.
func (r Resource) String() string {
	return fmt.Sprintf("%s:\n  - Name: %s\n  - Namespace: %s\n  - Current API Version: %s\n  - Available API Versions: %v\n",
		r.ID, r.Name, r.Namespace, r.CurrentAPIVersion, r.AvailableAPIVersions)
}

// BicepFile contains information about a bicep file:
//   - Path: the path to the bicep file (e.g. ./bicep/modules/virtualNetworks.bicep)
//   - Resources: the bicep resources defined in the bicep file
type BicepFile struct {
	Path      string
	Resources []Resource
}

// String returns a string representation of a types.BicepFile object.
func (file BicepFile) String() string {
	str := fmt.Sprintf("%s:\n  - Resources:\n", file.Path)
	for _, r := range file.Resources {
		str += fmt.Sprintf("     - %s:\n\t+ Name: %s\n\t+ Namespace: %s\n\t+ Current API Version: %s\n\t+ Available API Versions: %v\n",
			r.ID, r.Name, r.Namespace, r.CurrentAPIVersion, r.AvailableAPIVersions)
	}
	return str
}

// BicepDirectory contains information about a bicep directory:
//   - Path: the path to the bicep directory (e.g. ./bicep/modules)
//   - Files: the bicep files in the bicep directory
type BicepDirectory struct {
	Path  string
	Files []BicepFile
}

// String returns a string representation of a types.BicepDirectory object.
func (dir BicepDirectory) String() string {
	str := fmt.Sprintf("%s:\n", dir.Path)
	for _, file := range dir.Files {
		relName, err := filepath.Rel(dir.Path, file.Path)
		if err != nil {
			panic(err)
		}
		str += "- " + relName + ":\n  - Resources:\n"
		for _, r := range file.Resources {
			str += fmt.Sprintf("     - %s:\n\t+ Name: %s\n\t+ Namespace: %s\n\t+ Current API Version: %s\n\t+ Available API Versions: %v\n",
				r.ID, r.Name, r.Namespace, r.CurrentAPIVersion, r.AvailableAPIVersions)
		}
	}
	return str
}

// Mode represents the mode of the cli (scan or update).
type Mode int8

const (
	ModeScan   Mode = iota // ModeScan corresponds to the `bruh scan` command
	ModeUpdate             // ModeUpdate corresponds to the `bruh update` command
)

// String returns a string representation of a types.Mode object.
func (s Mode) String() string {
	switch s {
	case ModeScan:
		return "scan"
	case ModeUpdate:
		return "update"
	}
	return "unknown"
}
