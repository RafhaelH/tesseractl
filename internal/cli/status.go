package cli

import (
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/RafhaelH/cli_go/internal/config"
	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/spf13/cobra"
)

func newStatusCmd(flags *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "List running containers grouped by project",
		RunE: func(cmd *cobra.Command, args []string) error {
			containers, err := docker.ListContainers()
			if err != nil {
				return err
			}

			cfg, cfgErr := loadConfigOrNil(flags.configPath)
			if cfgErr != nil {
				return cfgErr
			}

			out := cmd.OutOrStdout()

			if cfg == nil {
				return printFlat(out, containers)
			}

			byName := indexByName(containers)
			return printPerProject(out, cfg,
				[]string{"SERVICE", "CONTAINER", "STATUS"},
				func(_, svc, container string) []string {
					status := "not running"
					if c, ok := byName[container]; ok {
						status = c.Status
					}
					return []string{svc, container, status}
				})
		},
	}
}

func loadConfigOrNil(explicit string) (*config.Config, error) {
	path, err := config.Resolve(explicit)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) && explicit == "" {
			return nil, nil
		}
		return nil, err
	}
	return config.Load(path)
}

func printFlat(out io.Writer, containers []docker.Container) error {
	if len(containers) == 0 {
		fmt.Fprintln(out, "No running containers.")
		return nil
	}
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tIMAGE\tSTATUS\tPORTS")
	for _, c := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.Names, c.Image, c.Status, c.Ports)
	}
	return w.Flush()
}
