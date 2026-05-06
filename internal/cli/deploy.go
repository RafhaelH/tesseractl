package cli

import (
	"fmt"

	"github.com/RafhaelH/tesseractl/internal/docker"
	"github.com/RafhaelH/tesseractl/internal/proc"
	"github.com/spf13/cobra"
)

func newDeployCmd(flags *globalFlags) *cobra.Command {
	var noPull bool
	var branch string
	cmd := &cobra.Command{
		Use:   "deploy <projeto>",
		Short: "Atualiza o código e reconstrói os containers de um projeto",
		Long: `Executa o fluxo padrão de deploy em produção para um projeto:

  1. (opcional) git checkout <branch>
  2. git pull --ff-only
  3. docker compose up -d --build

Use --no-pull para pular o passo do git.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project, err := lookupProject(flags.configPath, args[0])
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()

			if !noPull {
				if branch != "" {
					fmt.Fprintf(out, "→ git checkout %s\n", branch)
					if err := proc.Run(project.Path, "git", "checkout", branch); err != nil {
						return fmt.Errorf("git checkout: %w", err)
					}
				}
				fmt.Fprintln(out, "→ git pull --ff-only")
				if err := proc.Run(project.Path, "git", "pull", "--ff-only"); err != nil {
					return fmt.Errorf("git pull: %w", err)
				}
			}

			fmt.Fprintln(out, "→ docker compose up -d --build")
			return docker.RunCompose(project.Path, project.ComposeFile, "up", "-d", "--build")
		},
	}
	cmd.Flags().BoolVar(&noPull, "no-pull", false, "pula o git pull, apenas reconstrói")
	cmd.Flags().StringVar(&branch, "branch", "", "faz checkout desta branch antes do pull")
	return cmd
}
