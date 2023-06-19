package bicep

import (
	"reflect"
	"strings"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

func TestParseUpdateFile(t *testing.T) {
	type args struct {
		filename       string
		resources      []types.ResourceInfo
		inPlace        bool
		includePreview bool
	}
	tests := []struct {
		name    string
		args    args
		initial []types.ResourceInfo
		final   []types.ResourceInfo
		wantErr bool
	}{
		{
			name: "azure.deploy.bicep",
			args: args{
				filename: "testdata/parse_update/azure.deploy.bicep",
				resources: []types.ResourceInfo{
					{
						AvailableAPIVersions: []string{
							"2022-09-01",
							"2021-04-01",
							"2021-01-01",
							"2020-10-01",
						},
					},
				},
				inPlace:        false,
				includePreview: false,
			},
			initial: []types.ResourceInfo{
				{
					ID:                "Microsoft.Resources/resourceGroups",
					Name:              "resourceGroups",
					Namespace:         "Microsoft.Resources",
					CurrentAPIVersion: "2021-01-01",
				},
			},
			final: []types.ResourceInfo{
				{
					ID:                "Microsoft.Resources/resourceGroups",
					Name:              "resourceGroups",
					Namespace:         "Microsoft.Resources",
					CurrentAPIVersion: "2022-09-01",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First parse
			got, err := ParseFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.initial) {
				t.Errorf("First parse: ParseFile() = %v, want %v", got, tt.initial)
			}

			// Inject available API versions
			for i := range got {
				got[i].AvailableAPIVersions = tt.args.resources[i].AvailableAPIVersions
			}

			// Update file
			err = UpdateFile(tt.args.filename, got, tt.args.inPlace, tt.args.includePreview)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Second parse
			updatedFile := strings.Replace(tt.args.filename, ".bicep", "_updated.bicep", 1)
			got, err = ParseFile(updatedFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.final) {
				t.Errorf("Second parse: ParseFile() = %v, want %v", got, tt.final)
			}

			// Cleanup
			err = deleteFile(updatedFile)
			if err != nil {
				t.Errorf("deleteFile() error = %v", err)
			}
		})
	}
}
