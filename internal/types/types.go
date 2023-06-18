// Package types provides shared types used by multiple packages in the "bruh" application.
package types

import "fmt"

// ResourceInfo contains information about a resource:
//   - ID: the resource ID (e.g. Microsoft.Network/virtualNetworks)
//   - Namespace: the resource namespace (e.g. Microsoft.Network)
//   - Resource: the resource name (e.g. virtualNetworks)
//   - CurrentAPIVersion: the used api version (e.g. 2021-02-01)
//   - AvailableAPIVersions: the available api versions (e.g. [2021-02-01 2020-11-01])
type ResourceInfo struct {
	ID                   string
	Namespace            string
	Resource             string
	CurrentAPIVersion    string
	AvailableAPIVersions []string
}

// String returns a string representation of a types.ResourceInfo.
func (r ResourceInfo) String() string {
	return fmt.Sprintf("%s:\n  - Namespace: %s\n  - Resource: %s\n  - Current API Version: %s\n  - Available API Versions: %v\n",
		r.ID, r.Namespace, r.Resource, r.CurrentAPIVersion, r.AvailableAPIVersions)
}
