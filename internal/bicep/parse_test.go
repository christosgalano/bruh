package bicep

import (
	"reflect"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

const (
	parseAzureDeployBicepFile      = "testdata/parse/azure.deploy.bicep"
	parseAzureDeployParametersFile = "testdata/parse/azure.deploy.parameters.json"
)

func Test_validateBicepFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid-bicep-file",
			args:    args{parseAzureDeployBicepFile},
			wantErr: false,
		},
		{
			name:    "invalid-bicep-file",
			args:    args{parseAzureDeployParametersFile},
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
			if err := validateBicepFile(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBicepFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []types.ResourceInfo
		wantErr bool
	}{
		{
			name: parseAzureDeployBicepFile,
			args: args{parseAzureDeployBicepFile},
			want: []types.ResourceInfo{
				{
					ID:                "Microsoft.Resources/resourceGroups",
					Name:              "resourceGroups",
					Namespace:         "Microsoft.Resources",
					CurrentAPIVersion: "2021-01-01",
				},
			},
			wantErr: false,
		},
		{
			name: "compute.bicep",
			args: args{"testdata/parse/modules/compute.bicep"},
			want: []types.ResourceInfo{
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
			wantErr: false,
		},
		{
			name:    parseAzureDeployParametersFile,
			args:    args{parseAzureDeployParametersFile},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "non-existent-file",
			args:    args{"testdata/parse/non-existent-file"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFile() = %v, want %v", got, tt.want)
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
		want    map[string][]types.ResourceInfo
		wantErr bool
	}{
		{
			name: "testdata-parse",
			args: args{"testdata/parse"},
			want: map[string][]types.ResourceInfo{
				parseAzureDeployBicepFile: {
					{
						ID:                "Microsoft.Resources/resourceGroups",
						Name:              "resourceGroups",
						Namespace:         "Microsoft.Resources",
						CurrentAPIVersion: "2021-01-01",
					},
				},
				"testdata/parse/modules/compute.bicep": {
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
				"testdata/parse/modules/identity.bicep": {
					{
						ID:                "Microsoft.ManagedIdentity/userAssignedIdentities",
						Name:              "userAssignedIdentities",
						Namespace:         "Microsoft.ManagedIdentity",
						CurrentAPIVersion: "2022-01-31-preview",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "non-existent-dir",
			args:    args{"testdata/parse/non-existent-dir"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDirectory(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
