package main

import (
	"path/filepath"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	args := append([]string{"run", ".."}, os.Args[1:]...)
	cmd := exec.Command("go", args...)
	_, file, _, ok := runtime.Caller(0)
	if ok {
		cmd.Dir = filepath.Dir(filepath.Dir(file))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}

		os.Exit(1)
	}
}
