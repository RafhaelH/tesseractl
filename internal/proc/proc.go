package proc

import (
	"os"
	"os/exec"
)

// Run executa name+args fazendo streaming do stdio para o processo pai.
// Se workdir não estiver vazio, o comando roda nesse diretório.
func Run(workdir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if workdir != "" {
		cmd.Dir = workdir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
