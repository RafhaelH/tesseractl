package cli

import (
	"github.com/spf13/cobra"
)

var Version = "dev"

type globalFlags struct {
	configPath string
}

func NewRootCmd() *cobra.Command {
	flags := &globalFlags{}

	cmd := &cobra.Command{
		Use:     "tesserato",
		Version: Version,
		Short:   "Manage your Docker Compose projects from one place",
		Long: `tesserato is a CLI on top of Docker and Docker Compose
that centralizes day-to-day operations across multiple projects.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&flags.configPath, "config", "", "path to tesserato config file")

	cmd.AddCommand(newStatusCmd(flags))
	cmd.AddCommand(newUpCmd(flags))
	cmd.AddCommand(newDownCmd(flags))
	cmd.AddCommand(newLogsCmd(flags))
	cmd.AddCommand(newCleanCmd())
	cmd.AddCommand(newHealthCmd(flags))
	cmd.AddCommand(newDeployCmd(flags))

	return cmd
}
