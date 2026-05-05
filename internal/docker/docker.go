package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/RafhaelH/cli_go/internal/proc"
)

type Container struct {
	ID     string `json:"ID"`
	Names  string `json:"Names"`
	Image  string `json:"Image"`
	Status string `json:"Status"`
	Ports  string `json:"Ports"`
	State  string `json:"State"`
}

func ListContainers() ([]Container, error) {
	return runDockerPs(false)
}

func ListAll() ([]Container, error) {
	return runDockerPs(true)
}

func runDockerPs(includeStopped bool) ([]Container, error) {
	args := []string{"ps", "--no-trunc", "--format", "{{json .}}"}
	if includeStopped {
		args = append(args, "-a")
	}
	cmd := exec.Command("docker", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := bytes.TrimSpace(stderr.Bytes())
		if len(msg) == 0 {
			return nil, fmt.Errorf("docker ps failed: %w", err)
		}
		return nil, fmt.Errorf("docker ps failed: %s", msg)
	}
	return parseContainers(&stdout)
}

func parseContainers(r io.Reader) ([]Container, error) {
	var containers []Container
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var c Container
		if err := json.Unmarshal(line, &c); err != nil {
			return nil, fmt.Errorf("parse docker ps output: %w", err)
		}
		containers = append(containers, c)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read docker ps output: %w", err)
	}
	return containers, nil
}

func RunCompose(workdir, composeFile string, args ...string) error {
	full := []string{"compose"}
	if composeFile != "" {
		full = append(full, "-f", composeFile)
	}
	full = append(full, args...)
	return proc.Run(workdir, "docker", full...)
}

func RunDocker(args ...string) error {
	return proc.Run("", "docker", args...)
}
