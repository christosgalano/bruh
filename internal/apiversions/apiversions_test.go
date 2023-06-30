package apiversions

import (
	"reflect"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

/// Unit Tests ///

func Test_fetchResourcePage(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal-page",
			args: args{
				url: "https://learn.microsoft.com/en-us/azure/templates/microsoft.network/virtualnetworks",
			},
			wantErr: false,
		},
		{
			name: "invalid-page",
			args: args{
				url: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchResourcePage(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("fetchResourcePage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extractAPIVersions(t *testing.T) {
	type args struct {
		body    string
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "one-version",
			args: args{
				body:    "href=\"2021-02-01/virtualnetworks\"",
				pattern: `href="(\d{4}-\d{2}-\d{2})/virtualnetworks"`,
			},
			want:    []string{"2021-02-01"},
			wantErr: false,
		},
		{
			name: "two-versions",
			args: args{
				body:    "href=\"2021-02-01/virtualnetworks\", href=\"2022-02-01/virtualnetworks\", href=\"2022-02-01-preview/virtualnetworks\"",
				pattern: `href="(\d{4}-\d{2}-\d{2}|\d{4}-\d{2}-\d{2}-preview)/virtualnetworks"`,
			},
			want:    []string{"2022-02-01", "2022-02-01-preview", "2021-02-01"}, // Sorted in descending order (non-preview first)
			wantErr: false,
		},
		{
			name: "no-versions",
			args: args{
				body:    "invalid",
				pattern: `href="(\d{4}-\d{2}-\d{2})/virtualnetworks"`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractAPIVersions(tt.args.body, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Fatalf("extractAPIVersions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractAPIVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateResource(t *testing.T) {
	type args struct {
		resource       *types.Resource
		includePreview bool
	}
	tests := []struct {
		name    string
		args    args
		subset  []string
		wantErr bool
	}{
		{
			name: "valid-resource",
			args: args{
				resource: &types.Resource{
					ID:        "Microsoft.Web/serverFarms",
					Name:      "serverFarms",
					Namespace: "Microsoft.Web",
				},
				includePreview: true,
			},
			subset:  []string{"2022-03-01", "2021-03-01", "2021-02-01", "2021-01-15", "2021-01-01"},
			wantErr: false,
		},
		{
			name: "with-preview",
			args: args{
				resource: &types.Resource{
					ID:        "Microsoft.Insights/diagnosticSettings",
					Name:      "diagnosticSettings",
					Namespace: "Microsoft.Insights",
				},
				includePreview: true,
			},
			subset:  []string{"2021-05-01-preview", "2020-01-01-preview", "2017-05-01-preview", "2016-09-01", "2015-07-01"},
			wantErr: false,
		},
		{
			name: "without-preview",
			args: args{
				resource: &types.Resource{
					ID:        "Microsoft.Insights/diagnosticSettings",
					Name:      "diagnosticSettings",
					Namespace: "Microsoft.Insights",
				},
				includePreview: false,
			},
			subset: []string{"2016-09-01", "2015-07-01"},
		},
		{
			name: "invalid-resource",
			args: args{
				resource: &types.Resource{
					ID:        "Microsoft.Web/invalid",
					Name:      "invalid",
					Namespace: "Microsoft.Web",
				},
				includePreview: true,
			},
			subset:  nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateResource(tt.args.resource, true); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateResource() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !isSubset(tt.subset, tt.args.resource.AvailableAPIVersions) {
				t.Errorf("UpdateResource() = %v is not superset of %v", tt.args.resource.AvailableAPIVersions, tt.subset)
			}
		})
	}
}

func TestUpdateBicepFile(t *testing.T) {
	type args struct {
		bicepFile *types.BicepFile
	}
	tests := []struct {
		name    string
		args    args
		subset  [][]string
		wantErr bool
	}{
		{
			name: "valid-file",
			args: args{
				bicepFile: &types.BicepFile{
					Path: "compute.bicep",
					Resources: []types.Resource{
						{
							ID:        "Microsoft.Web/serverFarms",
							Name:      "serverFarms",
							Namespace: "Microsoft.Web",
						},
						{
							ID:        "Microsoft.Web/sites",
							Name:      "sites",
							Namespace: "Microsoft.Web",
						},
					},
				},
			},
			subset: [][]string{
				{"2022-03-01", "2021-03-01", "2021-02-01", "2021-01-15", "2021-01-01", "2020-12-01", "2020-10-01"},
				{"2022-03-01", "2021-03-01", "2021-02-01", "2021-01-15", "2021-01-01", "2020-12-01", "2020-10-01"},
			},
			wantErr: false,
		},
		{
			name: "invalid-file",
			args: args{
				bicepFile: &types.BicepFile{
					Path: "invalid.bicep",
					Resources: []types.Resource{
						{
							ID:        "Microsoft.Web/invalid",
							Name:      "invalid",
							Namespace: "Microsoft.Web",
						},
					},
				},
			},
			subset:  nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateBicepFile(tt.args.bicepFile, true); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateBicepFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			for i, resource := range tt.args.bicepFile.Resources {
				if !isSubset(tt.subset[i], resource.AvailableAPIVersions) {
					t.Errorf("UpdateBicepFile() = %v is not superset of %v", resource.AvailableAPIVersions, tt.subset)
				}
			}
		})
	}
}

func TestUpdateBicepDirectory(t *testing.T) {
	type args struct {
		bicepDirectory *types.BicepDirectory
	}
	tests := []struct {
		name    string
		args    args
		subset  [][][]string
		wantErr bool
	}{
		{
			name: "valid-directory",
			args: args{
				bicepDirectory: &types.BicepDirectory{
					Path: "compute",
					Files: []types.BicepFile{
						{
							Path: "compute.bicep",
							Resources: []types.Resource{
								{
									ID:        "Microsoft.Web/serverFarms",
									Name:      "serverFarms",
									Namespace: "Microsoft.Web",
								},
								{
									ID:        "Microsoft.Web/sites",
									Name:      "sites",
									Namespace: "Microsoft.Web",
								},
							},
						},
					},
				},
			},
			subset: [][][]string{
				{
					{"2022-03-01", "2021-03-01", "2021-02-01", "2021-01-15", "2021-01-01", "2020-12-01", "2020-10-01"},
					{"2022-03-01", "2021-03-01", "2021-02-01", "2021-01-15", "2021-01-01", "2020-12-01", "2020-10-01"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid-directory",
			args: args{
				bicepDirectory: &types.BicepDirectory{
					Path: "invalid",
					Files: []types.BicepFile{
						{
							Path: "invalid.bicep",
							Resources: []types.Resource{
								{
									ID:        "Microsoft.Web/invalid",
									Name:      "invalid",
									Namespace: "Microsoft.Web",
								},
							},
						},
					},
				},
			},
			subset:  nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateBicepDirectory(tt.args.bicepDirectory, true); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateBicepDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			for i, file := range tt.args.bicepDirectory.Files {
				for j, resource := range file.Resources {
					if !isSubset(tt.subset[i][j], resource.AvailableAPIVersions) {
						t.Errorf("UpdateBicepDirectory() = %v is not superset of %v", resource.AvailableAPIVersions, tt.subset)
					}
				}
			}
		})
	}
}

/// Benchmarks ///

//revive:disable:unhandled-error

func Benchmark_fetchResourcePage(b *testing.B) {
	url := "https://learn.microsoft.com/en-us/azure/templates/microsoft.network/virtualnetworks"
	for i := 0; i < b.N; i++ {
		fetchResourcePage(url)
	}
}

func Benchmark_extractAPIVersions(b *testing.B) {
	body := "href=\"2021-02-01/virtualnetworks\", href=\"2022-02-01/virtualnetworks\", href=\"2022-02-01/virtualnetworks\""
	pattern := `href="(\d{4}-\d{2}-\d{2})/virtualnetworks"`
	for i := 0; i < b.N; i++ {
		extractAPIVersions(body, pattern)
	}
}

func BenchmarkUpdateResource(b *testing.B) {
	resource := &types.Resource{
		ID:        "Microsoft.Web/serverFarms",
		Name:      "serverFarms",
		Namespace: "Microsoft.Web",
	}
	for i := 0; i < b.N; i++ {
		UpdateResource(resource, true)
	}
}

func BenchmarkUpdateBicepFile(b *testing.B) {
	bicepFile := &types.BicepFile{
		Path: "test.bicep",
		Resources: []types.Resource{
			{
				ID:        "Microsoft.Web/serverFarms",
				Name:      "serverFarms",
				Namespace: "Microsoft.Web",
			},
			{
				ID:        "Microsoft.Web/sites",
				Name:      "sites",
				Namespace: "Microsoft.Web",
			},
			{
				ID:        "Microsoft.Network/virtualNetworks",
				Name:      "virtualNetworks",
				Namespace: "Microsoft.Network",
			},
		},
	}
	for i := 0; i < b.N; i++ {
		UpdateBicepFile(bicepFile, true)
	}
}

func BenchmarkUpdateBicepDirectory(b *testing.B) {
	bicepDirectory := &types.BicepDirectory{
		Path: "test",
		Files: []types.BicepFile{
			{
				Path: "test.bicep",
				Resources: []types.Resource{
					{
						ID:        "Microsoft.Web/serverFarms",
						Name:      "serverFarms",
						Namespace: "Microsoft.Web",
					},
					{
						ID:        "Microsoft.Web/sites",
						Name:      "sites",
						Namespace: "Microsoft.Web",
					},
					{
						ID:        "Microsoft.Network/virtualNetworks",
						Name:      "virtualNetworks",
						Namespace: "Microsoft.Network",
					},
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		UpdateBicepDirectory(bicepDirectory, true)
	}
}

/// Helping Functions ///

// isSubset returns true if slice1 is a subset of slice2.
func isSubset(slice1, slice2 []string) bool {
	set := make(map[string]bool)
	for _, item := range slice2 {
		set[item] = true
	}
	for _, item := range slice1 {
		if !set[item] {
			return false
		}
	}
	return true
}
