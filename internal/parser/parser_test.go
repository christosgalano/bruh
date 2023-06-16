package parser

import (
	"reflect"
	"testing"
)

const (
	azureDeployBicepFile      = "testdata/bicep/azure.deploy.bicep"
	azureDeployParametersFile = "testdata/bicep/azure.deploy.parameters.json"
)

func TestValidateBicepFile(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			if err := ValidateBicepFile(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBicepFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsBicepFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		testname string
		args     args
		want     bool
	}{
		{
			"valid-bicep-file",
			args{azureDeployBicepFile},
			true,
		},
		{
			"invalid-bicep-file",
			args{azureDeployParametersFile},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			if got := IsBicepFile(tt.args.filename); got != tt.want {
				t.Errorf("IsBicepFile() = %v, want %v", got, tt.want)
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
		want    []ResourceInfo
		wantErr bool
	}{
		{
			"azure.deploy.bicep",
			args{azureDeployBicepFile},
			[]ResourceInfo{
				{
					ID:        "Microsoft.Resources/resourceGroups",
					Namespace: "Microsoft.Resources",
					Resource:  "resourceGroups",
					Version:   "2022-09-01",
				},
			},
			false,
		},
		{
			"compute.bicep",
			args{"testdata/bicep/modules/compute.bicep"},
			[]ResourceInfo{
				{
					ID:        "Microsoft.Web/serverfarms",
					Namespace: "Microsoft.Web",
					Resource:  "serverfarms",
					Version:   "2022-03-01",
				},
				{
					ID:        "Microsoft.Web/sites",
					Namespace: "Microsoft.Web",
					Resource:  "sites",
					Version:   "2022-03-01",
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

func TestParseDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]ResourceInfo
		wantErr bool
	}{
		{
			"bicep",
			args{"testdata/bicep"},
			map[string][]ResourceInfo{
				azureDeployBicepFile: {
					{
						ID:        "Microsoft.Resources/resourceGroups",
						Namespace: "Microsoft.Resources",
						Resource:  "resourceGroups",
						Version:   "2022-09-01",
					},
				},
				"testdata/bicep/modules/compute.bicep": {
					{
						ID:        "Microsoft.Web/serverfarms",
						Namespace: "Microsoft.Web",
						Resource:  "serverfarms",
						Version:   "2022-03-01",
					},
					{
						ID:        "Microsoft.Web/sites",
						Namespace: "Microsoft.Web",
						Resource:  "sites",
						Version:   "2022-03-01",
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
			got, err := ParseDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
