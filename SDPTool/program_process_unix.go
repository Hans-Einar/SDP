//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly

package sdptool

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

func configureProgramProcess(cmd *exec.Cmd) (func() error, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// Cancellation owns the runner and its descendants, including a GUI started by make.
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = 2 * time.Second
	restore := func() error { return nil }
	var foreground int32
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCGPGRP, uintptr(unsafe.Pointer(&foreground)))
	if errno == syscall.ENOTTY || errno == syscall.EBADF {
		return restore, nil
	}
	if errno != 0 {
		return nil, fmt.Errorf("read terminal foreground group: %w", errno)
	}
	if int(foreground) != syscall.Getpgrp() {
		return nil, fmt.Errorf("run must own the terminal foreground group; bring it to the foreground first")
	}
	// A separate background group would stop on terminal input (SIGTTIN).
	cmd.SysProcAttr.Foreground = true
	cmd.SysProcAttr.Ctty = int(os.Stdin.Fd())
	restore = func() error {
		// Restoring the parent's foreground ownership is itself a background terminal operation.
		ignored := signal.Ignored(syscall.SIGTTOU)
		signal.Ignore(syscall.SIGTTOU)
		if !ignored {
			defer signal.Reset(syscall.SIGTTOU)
		}
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCSPGRP, uintptr(unsafe.Pointer(&foreground)))
		if errno != 0 {
			return fmt.Errorf("restore terminal foreground group: %w", errno)
		}
		return nil
	}
	return restore, nil
}
