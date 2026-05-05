package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantErr  bool
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name: "valid config with one project",
			yaml: `
projects:
  demo:
    path: /tmp/demo
    compose_file: docker-compose.yml
    services:
      web: demo_web
      db: demo_db
`,
			validate: func(t *testing.T, cfg *Config) {
				if len(cfg.Projects) != 1 {
					t.Fatalf("want 1 project, got %d", len(cfg.Projects))
				}
				p, ok := cfg.Projects["demo"]
				if !ok {
					t.Fatal(`want project "demo"`)
				}
				if p.Path != "/tmp/demo" {
					t.Errorf("path: want /tmp/demo, got %q", p.Path)
				}
				if p.Services["web"] != "demo_web" {
					t.Errorf(`services["web"]: want demo_web, got %q`, p.Services["web"])
				}
			},
		},
		{
			name:    "invalid YAML",
			yaml:    "projects: : :",
			wantErr: true,
		},
		{
			name: "empty config is valid (no projects)",
			yaml: ``,
			validate: func(t *testing.T, cfg *Config) {
				if len(cfg.Projects) != 0 {
					t.Errorf("want 0 projects, got %d", len(cfg.Projects))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "tesserato.yaml")
			if err := os.WriteFile(path, []byte(tc.yaml), 0o644); err != nil {
				t.Fatalf("write temp config: %v", err)
			}

			cfg, err := Load(path)
			if tc.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.validate != nil {
				tc.validate(t, cfg)
			}
		})
	}
}

func TestResolve_NotFound(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	_, err := Resolve("")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestResolve_ExplicitMissing(t *testing.T) {
	_, err := Resolve(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if errors.Is(err, ErrNotFound) {
		t.Fatal("explicit-missing should NOT be ErrNotFound")
	}
}
