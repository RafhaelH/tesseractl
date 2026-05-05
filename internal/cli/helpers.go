package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/RafhaelH/cli_go/internal/config"
	"github.com/RafhaelH/cli_go/internal/docker"
)

func lookupProject(configPath, name string) (*config.Project, error) {
	path, err := config.Resolve(configPath)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return nil, fmt.Errorf(
				"no tesserato config found.\n" +
					"  Looked for: ./tesserato.yaml, ./tesserato.yml, ~/.config/tesserato/config.yaml\n" +
					"  Create one or pass --config <path>.")
		}
		return nil, err
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, err
	}
	project, ok := cfg.Projects[name]
	if !ok {
		return nil, fmt.Errorf("project %q not found in %s.\n  Available: %s",
			name, path, availableProjects(cfg))
	}
	return &project, nil
}

func availableProjects(cfg *config.Config) string {
	if len(cfg.Projects) == 0 {
		return "(none defined)"
	}
	names := make([]string, 0, len(cfg.Projects))
	for n := range cfg.Projects {
		names = append(names, n)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func indexByName(containers []docker.Container) map[string]docker.Container {
	m := make(map[string]docker.Container, len(containers))
	for _, c := range containers {
		m[c.Names] = c
	}
	return m
}
