package azapiversions

import (
	"reflect"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

func Test_fetchResourcePage(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
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
				t.Errorf("fetchResourcePage() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				body:    "<a href=\"2021-02-01/virtualnetworks\">2021-02-01</a>",
				pattern: `href="(\d{4}-\d{2}-\d{2})/virtualnetworks"`,
			},
			want:    []string{"2021-02-01"},
			wantErr: false,
		},
		{
			name: "two-versions",
			args: args{
				body:    "<a href=\"2021-02-01/virtualnetworks\">2021-02-01</a><a href=\"2022-02-01/virtualnetworks\">2022-02-01</a>",
				pattern: `href="(\d{4}-\d{2}-\d{2})/virtualnetworks"`,
			},
			want:    []string{"2022-02-01", "2021-02-01"}, // Sorted in descending order
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
				t.Errorf("extractAPIVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractAPIVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAPIVersions(t *testing.T) {
	type args struct {
		resource types.ResourceInfo
	}
	tests := []struct {
		name    string
		args    args
		subset  []string
		wantErr bool
	}{
		{
			name: "server-farms",
			args: args{
				resource: types.ResourceInfo{
					ID:        "Microsoft.Web/serverFarms",
					Namespace: "Microsoft.Web",
					Resource:  "serverFarms",
				},
			},
			subset:  []string{"2022-03-01", "2021-03-01", "2021-02-01", "2021-01-15", "2021-01-01", "2020-12-01", "2020-10-01"},
			wantErr: false,
		},
		{
			name: "invalid-resource",
			args: args{
				resource: types.ResourceInfo{
					ID:        "Microsoft.Web/invalid",
					Namespace: "Microsoft.Web",
					Resource:  "invalid",
				},
			},
			subset:  nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAPIVersions(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isSubset(tt.subset, got) {
				t.Errorf("GetAPIVersions() = %v is not superset of %v", got, tt.subset)
			}
		})
	}
}

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
