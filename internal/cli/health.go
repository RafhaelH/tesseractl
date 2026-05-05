package cli

import (
	"github.com/RafhaelH/cli_go/internal/config"
	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/spf13/cobra"
)

func newHealthCmd(flags *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check declared services across all projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.Resolve(flags.configPath)
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			containers, err := docker.ListAll()
			if err != nil {
				return err
			}
			byName := indexByName(containers)

			return printPerProject(cmd.OutOrStdout(), cfg,
				[]string{"SERVICE", "CONTAINER", "STATE", "IMAGE"},
				func(_, svc, container string) []string {
					state, image := "missing", "-"
					if c, ok := byName[container]; ok {
						state = c.Status
						image = c.Image
					}
					return []string{svc, container, state, image}
				})
		},
	}
}
