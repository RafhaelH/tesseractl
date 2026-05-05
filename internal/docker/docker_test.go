package docker

import (
	"strings"
	"testing"
)

func TestParseContainers(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Container
		wantErr bool
	}{
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name: "single container",
			input: `{"ID":"abc","Names":"web","Image":"nginx:alpine","Status":"Up 5 seconds","Ports":"0.0.0.0:8080->80/tcp","State":"running"}
`,
			want: []Container{{
				ID: "abc", Names: "web", Image: "nginx:alpine",
				Status: "Up 5 seconds", Ports: "0.0.0.0:8080->80/tcp", State: "running",
			}},
		},
		{
			name: "multiple containers with blank line",
			input: `{"ID":"a","Names":"web","Image":"nginx","Status":"Up","Ports":"","State":"running"}

{"ID":"b","Names":"db","Image":"postgres","Status":"Up","Ports":"","State":"running"}
`,
			want: []Container{
				{ID: "a", Names: "web", Image: "nginx", Status: "Up", State: "running"},
				{ID: "b", Names: "db", Image: "postgres", Status: "Up", State: "running"},
			},
		},
		{
			name:    "invalid JSON",
			input:   `{not valid json}`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseContainers(strings.NewReader(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("len mismatch: want %d, got %d", len(tc.want), len(got))
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("container[%d]: want %+v, got %+v", i, tc.want[i], got[i])
				}
			}
		})
	}
}
