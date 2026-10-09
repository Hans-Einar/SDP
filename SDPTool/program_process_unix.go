//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly

package sdptool

import (
	"os/exec"
	"syscall"
	"time"
)

func configureProgramProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// A build runner may start the actual GUI as a descendant. Cancellation must
	// stop the owned process group, not leave that GUI behind when make exits.
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = 2 * time.Second
}
