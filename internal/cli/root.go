package cli

import (
	"github.com/spf13/cobra"
)

var Version = "dev"

type globalFlags struct {
	configPath string
}

func NewRootCmd(dc dockerClient) *cobra.Command {
	flags := &globalFlags{}

	cmd := &cobra.Command{
		Use:     "tesserato",
		Version: Version,
		Short:   "Gerencie seus projetos Docker Compose em um só lugar",
		Long: `tesserato é uma CLI sobre Docker e Docker Compose
que centraliza operações do dia a dia em múltiplos projetos.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&flags.configPath, "config", "", "caminho do arquivo de configuração do tesserato")

	cmd.AddCommand(newStatusCmd(flags, dc))
	cmd.AddCommand(newUpCmd(flags))
	cmd.AddCommand(newDownCmd(flags))
	cmd.AddCommand(newLogsCmd(flags))
	cmd.AddCommand(newCleanCmd())
	cmd.AddCommand(newHealthCmd(flags))
	cmd.AddCommand(newDeployCmd(flags))

	return cmd
}
