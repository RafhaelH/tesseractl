package cli

import "github.com/RafhaelH/tesseractl/internal/docker"

type dockerClient interface {
	ListContainers() ([]docker.Container, error)
}
