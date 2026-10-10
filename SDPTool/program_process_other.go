//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd && !dragonfly

package sdptool

import (
	"os/exec"
	"time"
)

func configureProgramProcess(cmd *exec.Cmd) (func() error, error) {
	cmd.WaitDelay = 2 * time.Second
	return func() error { return nil }, nil
}
