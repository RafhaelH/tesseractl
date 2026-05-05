package proc

import (
	"os"
	"os/exec"
)

// Run executes name+args, streaming stdio to the parent process.
// If workdir is non-empty, the command runs there.
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
