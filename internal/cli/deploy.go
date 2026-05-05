package cli

import (
	"fmt"

	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/RafhaelH/cli_go/internal/proc"
	"github.com/spf13/cobra"
)

func newDeployCmd(flags *globalFlags) *cobra.Command {
	var noPull bool
	var branch string
	cmd := &cobra.Command{
		Use:   "deploy <project>",
		Short: "Pull latest code and rebuild a project's containers",
		Long: `Runs the standard production deploy flow for a project:

  1. (optional) git checkout <branch>
  2. git pull --ff-only
  3. docker compose up -d --build

Use --no-pull to skip the git step.`,
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
	cmd.Flags().BoolVar(&noPull, "no-pull", false, "skip git pull, just rebuild")
	cmd.Flags().StringVar(&branch, "branch", "", "checkout this branch before pulling")
	return cmd
}
