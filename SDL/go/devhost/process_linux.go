package devhost

import (
	"os"
	"os/exec"
	"syscall"
)

func group(command *exec.Cmd)             { command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} }
func interrupt(process *os.Process) error { return syscall.Kill(-process.Pid, syscall.SIGINT) }
func kill(process *os.Process) error      { return syscall.Kill(-process.Pid, syscall.SIGKILL) }
