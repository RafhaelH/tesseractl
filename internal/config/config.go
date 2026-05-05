package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ErrNotFound = errors.New("arquivo de configuração não encontrado")

type Project struct {
	Path        string            `yaml:"path"`
	ComposeFile string            `yaml:"compose_file"`
	Services    map[string]string `yaml:"services"`
}

type Config struct {
	Projects map[string]Project `yaml:"projects"`
}

func Resolve(explicit string) (string, error) {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			return "", fmt.Errorf("configuração %q: %w", explicit, err)
		}
		return explicit, nil
	}

	for _, name := range []string{"tesserato.yaml", "tesserato.yml"} {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		candidate := filepath.Join(home, ".config", "tesserato", "config.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", ErrNotFound
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro lendo configuração %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("erro interpretando configuração %q: %w", path, err)
	}
	return &cfg, nil
}
