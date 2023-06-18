package parse

import (
	"reflect"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

const (
	azureDeployBicepFile      = "testdata/bicep/azure.deploy.bicep"
	azureDeployParametersFile = "testdata/bicep/azure.deploy.parameters.json"
)

func Test_validateBicepFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		testname string
		args     args
		wantErr  bool
	}{
		{
			"valid-bicep-file",
			args{azureDeployBicepFile},
			false,
		},
		{
			"invalid-bicep-file",
			args{azureDeployParametersFile},
			true,
		},
		{
			"non-existent-file",
			args{"testdata/bicep/non-existent-file"},
			true,
		},
		{
			"directory",
			args{"testdata/bicep"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			if err := validateBicepFile(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBicepFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile(t *testing.T) {
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
			"azure.deploy.bicep",
			args{azureDeployBicepFile},
			[]types.ResourceInfo{
				{
					ID:                "Microsoft.Resources/resourceGroups",
					Namespace:         "Microsoft.Resources",
					Resource:          "resourceGroups",
					CurrentAPIVersion: "2022-09-01",
				},
			},
			false,
		},
		{
			"compute.bicep",
			args{"testdata/bicep/modules/compute.bicep"},
			[]types.ResourceInfo{
				{
					ID:                "Microsoft.Web/serverfarms",
					Namespace:         "Microsoft.Web",
					Resource:          "serverfarms",
					CurrentAPIVersion: "2022-03-01",
				},
				{
					ID:                "Microsoft.Web/sites",
					Namespace:         "Microsoft.Web",
					Resource:          "sites",
					CurrentAPIVersion: "2022-03-01",
				},
			},
			false,
		},
		{
			"azure.deploy.parameters.json",
			args{azureDeployParametersFile},
			nil,
			true,
		},
		{
			"non-existent-file",
			args{"testdata/bicep/non-existent-file"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := File(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("File() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDir(t *testing.T) {
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
			"bicep",
			args{"testdata/bicep"},
			map[string][]types.ResourceInfo{
				azureDeployBicepFile: {
					{
						ID:                "Microsoft.Resources/resourceGroups",
						Namespace:         "Microsoft.Resources",
						Resource:          "resourceGroups",
						CurrentAPIVersion: "2022-09-01",
					},
				},
				"testdata/bicep/modules/compute.bicep": {
					{
						ID:                "Microsoft.Web/serverfarms",
						Namespace:         "Microsoft.Web",
						Resource:          "serverfarms",
						CurrentAPIVersion: "2022-03-01",
					},
					{
						ID:                "Microsoft.Web/sites",
						Namespace:         "Microsoft.Web",
						Resource:          "sites",
						CurrentAPIVersion: "2022-03-01",
					},
				},
			},
			false,
		},
		{
			"non-existent-dir",
			args{"testdata/non-existent-dir"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Directory(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dir() = %v, want %v", got, tt.want)
			}
		})
	}
}
