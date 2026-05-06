package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RafhaelH/tesseractl/internal/docker"
)

// fakeDockerClient é um stub que satisfaz a interface dockerClient.
// Os campos permitem ao teste decidir o que ListContainers vai devolver.
type fakeDockerClient struct {
	containers []docker.Container
	err        error
}

func (f fakeDockerClient) ListContainers() ([]docker.Container, error) {
	return f.containers, f.err
}

// writeConfig escreve um YAML no diretório temporário do teste e devolve o path.
// t.TempDir() cria um diretório que é removido automaticamente ao fim do teste.
func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "tesserato.yaml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("escrevendo config: %v", err)
	}
	return path
}

// runStatus monta o comando raiz com o fake injetado, redireciona stdout/stderr
// para um buffer, executa "status [args...]" e devolve (output, erro).
func runStatus(t *testing.T, fake fakeDockerClient, args ...string) (string, error) {
	t.Helper()
	root := NewRootCmd(fake)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"status"}, args...))
	err := root.Execute()
	return out.String(), err
}

func TestStatus(t *testing.T) {
	defaultCfg := `projects:
  api:
    path: /any/path
    services:
      backend: api_backend
      db: api_db
`

	tests := []struct {
		name      string
		cfgYAML   string // se vazio, usa defaultCfg
		fake      fakeDockerClient
		wantInOut []string
		wantErr   bool
	}{
		{
			name: "container declarado existe e aparece com status",
			fake: fakeDockerClient{containers: []docker.Container{
				{Names: "api_backend", Image: "myapp:latest", Status: "Up 2 hours"},
			}},
			wantInOut: []string{"api", "api_backend", "Up 2 hours", "api_db", "parado"},
		},
		{
			name:      "nenhum container rodando: ambos aparecem como parados",
			fake:      fakeDockerClient{},
			wantInOut: []string{"api_backend", "parado", "api_db", "parado"},
		},
		{
			name:    "erro do docker propaga como erro do comando",
			fake:    fakeDockerClient{err: errors.New("docker daemon offline")},
			wantErr: true,
		},
		{
			name: "múltiplos projetos: output agrupado e ordenado",
			cfgYAML: `projects:
  api:
    path: /any/path
    services:
      backend: api_backend
  worker:
    path: /any/path
    services:
      celery: worker_celery
`,
			fake: fakeDockerClient{containers: []docker.Container{
				{Names: "api_backend", Status: "Up 1 hour"},
				{Names: "worker_celery", Status: "Up 30 minutes"},
			}},
			wantInOut: []string{"── api ──", "api_backend", "Up 1 hour", "── worker ──", "worker_celery", "Up 30 minutes"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			yaml := tc.cfgYAML
			if yaml == "" {
				yaml = defaultCfg
			}
			cfgPath := writeConfig(t, yaml)
			out, err := runStatus(t, tc.fake, "--config", cfgPath)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, recebi output=%q", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v\noutput:\n%s", err, out)
			}
			for _, want := range tc.wantInOut {
				if !strings.Contains(out, want) {
					t.Errorf("output não contém %q\noutput completo:\n%s", want, out)
				}
			}
		})
	}
}
