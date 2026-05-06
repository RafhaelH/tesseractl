package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RafhaelH/tesseractl/internal/docker"
	"github.com/spf13/cobra"
)

func newCleanCmd() *cobra.Command {
	var all, yes bool
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove containers parados e imagens dangling",
		Long: `Remove containers parados e imagens dangling (sem tag).

Com --all, roda o mais agressivo ` + "`docker system prune -a --volumes`" +
			` que também remove imagens não usadas e volumes anônimos.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			if all {
				fmt.Fprintln(out, "Isso irá remover:")
				fmt.Fprintln(out, "  - todos os containers parados")
				fmt.Fprintln(out, "  - todas as redes não usadas por pelo menos um container")
				fmt.Fprintln(out, "  - todas as imagens sem pelo menos um container associado")
				fmt.Fprintln(out, "  - todos os volumes anônimos não usados por pelo menos um container")
			} else {
				fmt.Fprintln(out, "Isso irá remover:")
				fmt.Fprintln(out, "  - todos os containers parados")
				fmt.Fprintln(out, "  - todas as imagens dangling (sem tag)")
			}

			if !yes {
				ok, err := confirm(out, os.Stdin, "Continuar?")
				if err != nil {
					return err
				}
				if !ok {
					fmt.Fprintln(out, "Cancelado.")
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
	cmd.Flags().BoolVar(&all, "all", false, "também remove imagens não usadas, redes e volumes anônimos")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "pula o prompt de confirmação")
	return cmd
}

func confirm(out io.Writer, in io.Reader, prompt string) (bool, error) {
	fmt.Fprintf(out, "%s [s/N] ", prompt)
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("erro lendo confirmação: %w", err)
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "s" || answer == "sim" || answer == "y" || answer == "yes", nil
}
