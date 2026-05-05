package cli

import (
	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/spf13/cobra"
)

func newDownCmd(flags *globalFlags) *cobra.Command {
	var volumes bool
	cmd := &cobra.Command{
		Use:   "down <project>",
		Short: "Stop and remove a project's containers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := lookupProject(flags.configPath, args[0])
			if err != nil {
				return err
			}

			composeArgs := []string{"down"}
			if volumes {
				composeArgs = append(composeArgs, "-v")
			}
			return docker.RunCompose(project.Path, project.ComposeFile, composeArgs...)
		},
	}
	cmd.Flags().BoolVar(&volumes, "volumes", false, "also remove named volumes")
	return cmd
}
