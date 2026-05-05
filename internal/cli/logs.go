package cli

import (
	"fmt"

	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/spf13/cobra"
)

func newLogsCmd(flags *globalFlags) *cobra.Command {
	var follow bool
	var tail string
	cmd := &cobra.Command{
		Use:   "logs <project> <service>",
		Short: "Show logs for a project's service",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName, serviceName := args[0], args[1]

			project, err := lookupProject(flags.configPath, projectName)
			if err != nil {
				return err
			}

			containerName, ok := project.Services[serviceName]
			if !ok {
				return fmt.Errorf("service %q not found in project %q", serviceName, projectName)
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
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow log output")
	cmd.Flags().StringVar(&tail, "tail", "", "number of lines to show from end (e.g. 100)")
	return cmd
}
