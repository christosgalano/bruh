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
//   - CurrentAPIVersion: the used api version (e.g. 2021-02-01)
//   - AvailableAPIVersions: the available api versions (e.g. [2021-02-01 2020-11-01])
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
//   - Name: the name of the bicep file (e.g. virtualNetworks.bicep)
//   - Resources: the bicep resources defined in the bicep file
type BicepFile struct {
	Name      string
	Resources []Resource
}

// String returns a string representation of a types.BicepFile object.
func (file BicepFile) String() string {
	str := fmt.Sprintf("%s:\n  - Resources:\n", file.Name)
	for _, r := range file.Resources {
		str += fmt.Sprintf("     - %s:\n\t+ Name: %s\n\t+ Namespace: %s\n\t+ Current API Version: %s\n\t+ Available API Versions: %v\n",
			r.ID, r.Name, r.Namespace, r.CurrentAPIVersion, r.AvailableAPIVersions)
	}
	return str
}

// BicepDirectory contains information about a bicep directory:
//   - Name: the name of the bicep directory (e.g. virtualNetworks)
//   - Files: the bicep files in the bicep directory
type BicepDirectory struct {
	Name  string
	Files []BicepFile
}

// String returns a string representation of a types.BicepDirectory object.
func (dir BicepDirectory) String() string {
	str := fmt.Sprintf("%s:\n", dir.Name)
	for _, file := range dir.Files {
		relName, err := filepath.Rel(dir.Name, file.Name)
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

// Mode represents the mode of the cli (scan or update)
type Mode int8

const (
	ModeScan   Mode = iota // ModeScan corresponds maps to the `bruh scan` command
	ModeUpdate             // ModeUpdate corresponds maps to the `bruh update` command
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
