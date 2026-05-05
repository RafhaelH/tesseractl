package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/RafhaelH/cli_go/internal/config"
)

type rowFunc func(projectName, serviceName, containerName string) []string

func printPerProject(out io.Writer, cfg *config.Config, columns []string, rowFor rowFunc) error {
	if len(cfg.Projects) == 0 {
		fmt.Fprintln(out, "Nenhum projeto definido na configuração.")
		return nil
	}

	header := "  " + strings.Join(columns, "\t")

	projectNames := make([]string, 0, len(cfg.Projects))
	for n := range cfg.Projects {
		projectNames = append(projectNames, n)
	}
	sort.Strings(projectNames)

	for _, projectName := range projectNames {
		project := cfg.Projects[projectName]
		fmt.Fprintf(out, "── %s ──\n", projectName)

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, header)

		serviceNames := make([]string, 0, len(project.Services))
		for s := range project.Services {
			serviceNames = append(serviceNames, s)
		}
		sort.Strings(serviceNames)

		for _, svc := range serviceNames {
			cells := rowFor(projectName, svc, project.Services[svc])
			fmt.Fprintln(w, "  "+strings.Join(cells, "\t"))
		}
		if err := w.Flush(); err != nil {
			return err
		}
		fmt.Fprintln(out)
	}
	return nil
}
