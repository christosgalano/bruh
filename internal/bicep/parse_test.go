package bicep

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

func Test_readBicepFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid-bicep-file",
			args:    args{"testdata/parse/azure.deploy.bicep"},
			wantErr: false,
		},
		{
			name:    "invalid-bicep-file",
			args:    args{"testdata/parse/azure.deploy.parameters.json"},
			wantErr: true,
		},
		{
			name:    "non-existent-file",
			args:    args{"testdata/non-existent-file"},
			wantErr: true,
		},
		{
			name:    "directory",
			args:    args{"testdata"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := readBicepFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Fatalf("ValidateBicepFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    types.BicepFile
		wantErr bool
	}{
		{
			name: "testdata/parse/azure.deploy.bicep",
			args: args{"testdata/parse/azure.deploy.bicep"},
			want: types.BicepFile{
				Path: filepath.FromSlash("testdata/parse/azure.deploy.bicep"),
				Resources: []types.Resource{
					{
						ID:                "Microsoft.Resources/resourceGroups",
						Name:              "resourceGroups",
						Namespace:         "Microsoft.Resources",
						CurrentAPIVersion: "2021-01-01",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "compute.bicep",
			args: args{"testdata/parse/modules/compute.bicep"},
			want: types.BicepFile{
				Path: filepath.FromSlash("testdata/parse/modules/compute.bicep"),
				Resources: []types.Resource{
					{
						ID:                "Microsoft.Web/serverfarms",
						Name:              "serverfarms",
						Namespace:         "Microsoft.Web",
						CurrentAPIVersion: "2021-01-15",
					},
					{
						ID:                "Microsoft.Web/sites",
						Name:              "sites",
						Namespace:         "Microsoft.Web",
						CurrentAPIVersion: "2019-08-01",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "testdata/parse/azure.deploy.parameters.json",
			args:    args{"testdata/parse/azure.deploy.parameters.json"},
			want:    types.BicepFile{},
			wantErr: true,
		},
		{
			name:    "non-existent-file",
			args:    args{"testdata/parse/non-existent-file"},
			want:    types.BicepFile{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("ParseFile() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestParseDirectory(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    types.BicepDirectory
		wantErr bool
	}{
		{
			name: "testdata-parse",
			args: args{"testdata/parse"},
			want: types.BicepDirectory{
				Path: "testdata/parse",
				Files: []types.BicepFile{
					{
						Path: filepath.FromSlash("testdata/parse/azure.deploy.bicep"),
						Resources: []types.Resource{
							{
								ID:                "Microsoft.Resources/resourceGroups",
								Name:              "resourceGroups",
								Namespace:         "Microsoft.Resources",
								CurrentAPIVersion: "2021-01-01",
							},
						},
					},
					{
						Path: filepath.FromSlash("testdata/parse/modules/compute.bicep"),
						Resources: []types.Resource{
							{
								ID:                "Microsoft.Web/serverfarms",
								Name:              "serverfarms",
								Namespace:         "Microsoft.Web",
								CurrentAPIVersion: "2021-01-15",
							},
							{
								ID:                "Microsoft.Web/sites",
								Name:              "sites",
								Namespace:         "Microsoft.Web",
								CurrentAPIVersion: "2019-08-01",
							},
						},
					},
					{
						Path: filepath.FromSlash("testdata/parse/modules/identity.bicep"),
						Resources: []types.Resource{
							{
								ID:                "Microsoft.ManagedIdentity/userAssignedIdentities",
								Name:              "userAssignedIdentities",
								Namespace:         "Microsoft.ManagedIdentity",
								CurrentAPIVersion: "2022-01-31-preview",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "non-existent-dir",
			args:    args{"testdata/parse/non-existent-dir"},
			want:    types.BicepDirectory{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDirectory(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("ParseDirectory() = %v, want %v", *got, tt.want)
			}
		})
	}
}
