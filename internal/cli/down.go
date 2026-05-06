package cli

import (
	"github.com/RafhaelH/tesseractl/internal/docker"
	"github.com/spf13/cobra"
)

func newDownCmd(flags *globalFlags) *cobra.Command {
	var volumes bool
	cmd := &cobra.Command{
		Use:   "down <projeto>",
		Short: "Para e remove os containers de um projeto",
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
	cmd.Flags().BoolVar(&volumes, "volumes", false, "também remove os volumes nomeados")
	return cmd
}
