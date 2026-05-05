package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RafhaelH/cli_go/internal/docker"
	"github.com/spf13/cobra"
)

func newCleanCmd() *cobra.Command {
	var all, yes bool
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Prune stopped containers and dangling images",
		Long: `Remove stopped containers and dangling (untagged) images.

With --all, runs the more aggressive ` + "`docker system prune -a --volumes`" +
			` which also removes unused images and anonymous volumes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			if all {
				fmt.Fprintln(out, "This will remove:")
				fmt.Fprintln(out, "  - all stopped containers")
				fmt.Fprintln(out, "  - all networks not used by at least one container")
				fmt.Fprintln(out, "  - all images without at least one container associated to them")
				fmt.Fprintln(out, "  - all anonymous volumes not used by at least one container")
			} else {
				fmt.Fprintln(out, "This will remove:")
				fmt.Fprintln(out, "  - all stopped containers")
				fmt.Fprintln(out, "  - all dangling (untagged) images")
			}

			if !yes {
				ok, err := confirm(out, os.Stdin, "Continue?")
				if err != nil {
					return err
				}
				if !ok {
					fmt.Fprintln(out, "Aborted.")
					return nil
				}
			}

			if all {
				return docker.RunDocker("system", "prune", "-a", "--volumes", "-f")
			}
			if err := docker.RunDocker("container", "prune", "-f"); err != nil {
				return err
			}
			return docker.RunDocker("image", "prune", "-f")
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "also remove unused images, networks, and anonymous volumes")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation prompt")
	return cmd
}

func confirm(out io.Writer, in io.Reader, prompt string) (bool, error) {
	fmt.Fprintf(out, "%s [y/N] ", prompt)
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("read confirmation: %w", err)
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}
