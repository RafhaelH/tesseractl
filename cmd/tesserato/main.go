package main

import (
	"fmt"
	"os"

	"github.com/RafhaelH/tesseractl/internal/cli"
	"github.com/RafhaelH/tesseractl/internal/docker"
)

func main() {
	if err := cli.NewRootCmd(docker.Client{}).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
