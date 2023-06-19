/*
Package types provides shared types used by multiple packages in the "bruh" application.
*/
package types

import "fmt"

// ResourceInfo contains information about a resource:
//   - ID: the resource ID (e.g. Microsoft.Network/virtualNetworks)
//   - Name: the resource name (e.g. virtualNetworks)
//   - Namespace: the resource namespace (e.g. Microsoft.Network)
//   - CurrentAPIVersion: the used api version (e.g. 2021-02-01)
//   - AvailableAPIVersions: the available api versions (e.g. [2021-02-01 2020-11-01])
type ResourceInfo struct {
	ID                   string
	Name                 string
	Namespace            string
	CurrentAPIVersion    string
	AvailableAPIVersions []string
}

// String returns a string representation of a types.ResourceInfo.
func (r ResourceInfo) String() string {
	return fmt.Sprintf("%s:\n  - Name: %s\n  - Namespace: %s\n  - Current API Version: %s\n  - Available API Versions: %v\n",
		r.ID, r.Name, r.Namespace, r.CurrentAPIVersion, r.AvailableAPIVersions)
}

// UpdateResourceInfo contains information about a resource to be updated:
//   - Filename: the name of the file containing the resource (e.g. compute.bicep)
//   - Index: the index of the resource in the file (e.g. 0)
//   - APIVersions: the available api versions (e.g. [2021-02-01 2020-11-01])
type UpdateResourceInfo struct {
	Filename    string
	Index       int
	APIVersions []string
}

// String returns a string representation of a types.UpdateResourceInfo.
func (r UpdateResourceInfo) String() string {
	return fmt.Sprintf("%s:\n  - Index: %d\n  - Available API Versions: %v\n",
		r.Filename, r.Index, r.APIVersions)
}
