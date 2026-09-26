//go:build !linux

package devhost

import (
	"os"
	"os/exec"
)

func group(command *exec.Cmd)             {}
func interrupt(process *os.Process) error { return process.Signal(os.Interrupt) }
func kill(process *os.Process) error      { return process.Kill() }
