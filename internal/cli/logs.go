package cli

import (
	"fmt"

	"github.com/RafhaelH/tesseractl/internal/docker"
	"github.com/spf13/cobra"
)

func newLogsCmd(flags *globalFlags) *cobra.Command {
	var follow bool
	var tail string
	cmd := &cobra.Command{
		Use:   "logs <projeto> <serviço>",
		Short: "Mostra os logs de um serviço de um projeto",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName, serviceName := args[0], args[1]

			project, err := lookupProject(flags.configPath, projectName)
			if err != nil {
				return err
			}

			containerName, ok := project.Services[serviceName]
			if !ok {
				return fmt.Errorf("serviço %q não encontrado no projeto %q", serviceName, projectName)
			}

			dockerArgs := []string{"logs"}
			if follow {
				dockerArgs = append(dockerArgs, "-f")
			}
			if tail != "" {
				dockerArgs = append(dockerArgs, "--tail", tail)
			}
			dockerArgs = append(dockerArgs, containerName)
			return docker.RunDocker(dockerArgs...)
		},
	}
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "acompanha a saída dos logs em tempo real")
	cmd.Flags().StringVar(&tail, "tail", "", "número de linhas a partir do fim (ex: 100)")
	return cmd
}
