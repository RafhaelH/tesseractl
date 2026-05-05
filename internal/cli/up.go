package cli

import (
	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/spf13/cobra"
)

func newUpCmd(flags *globalFlags) *cobra.Command {
	var build bool
	cmd := &cobra.Command{
		Use:   "up <project>",
		Short: "Start a project's containers via docker compose",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := lookupProject(flags.configPath, args[0])
			if err != nil {
				return err
			}

			composeArgs := []string{"up", "-d"}
			if build {
				composeArgs = append(composeArgs, "--build")
			}
			return docker.RunCompose(project.Path, project.ComposeFile, composeArgs...)
		},
	}
	cmd.Flags().BoolVar(&build, "build", false, "rebuild images before starting")
	return cmd
}
