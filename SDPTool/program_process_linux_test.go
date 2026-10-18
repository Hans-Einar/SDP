//go:build linux

package sdptool

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// Exercise main's real signal wiring and real controlling-terminal semantics.
func TestProgramExecutableLifecycle(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "sdptool")
	if out, err := exec.Command("go", "build", "-o", binary, "./cmd/sdptool").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	t.Run("terminal-input", func(t *testing.T) {
		root, d := programFixture(t)
		if err := os.WriteFile(filepath.Join(root, "runner"), []byte("#!/bin/sh\necho $$ > runner-pid\necho ready > ready\nread value\nprintf '%s' \"$value\" > received\nexit 7\n"), 0700); err != nil {
			t.Fatal(err)
		}
		master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
		if err != nil {
			t.Fatal(err)
		}
		defer master.Close()
		var unlock, number uint32
		for _, op := range []struct {
			request uintptr
			value   *uint32
		}{{syscall.TIOCSPTLCK, &unlock}, {syscall.TIOCGPTN, &number}} {
			if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, master.Fd(), op.request, uintptr(unsafe.Pointer(op.value))); errno != 0 {
				t.Fatal(errno)
			}
		}
		slave, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|syscall.O_NOCTTY, 0)
		if err != nil {
			t.Fatal(err)
		}
		defer slave.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary, root, "run", "--program", d.ID)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = slave, slave, slave
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true, Ctty: 0}
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		// Kill a stopped runner as well if a regression causes the timeout.
		defer func() {
			if b, err := os.ReadFile(filepath.Join(root, "runner-pid")); err == nil {
				var pid int
				fmt.Sscanf(string(b), "%d", &pid)
				_ = syscall.Kill(-pid, syscall.SIGKILL)
			}
		}()
		if _, err := master.Write([]byte("terminal value\n")); err != nil {
			t.Fatal(err)
		}
		err = cmd.Wait()
		if ctx.Err() != nil {
			t.Fatal("terminal runner hung", ctx.Err())
		}
		if e, ok := err.(*exec.ExitError); !ok || e.ExitCode() != 7 {
			t.Fatalf("exit: %v", err)
		}
		b, err := os.ReadFile(filepath.Join(root, "received"))
		if err != nil || string(b) != "terminal value" {
			t.Fatalf("input=%q error=%v", b, err)
		}
	})
	t.Run("sigterm-descendants", func(t *testing.T) {
		root, d := programFixture(t)
		if err := os.WriteFile(filepath.Join(root, "runner"), []byte("#!/bin/sh\n(sleep 1; echo leaked > descendant) &\necho started > started\nwait\n"), 0700); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary, root, "run", "--program", d.ID)
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		deadline := time.Now().Add(5 * time.Second)
		for {
			if _, err := os.Stat(filepath.Join(root, "started")); err == nil {
				break
			}
			if time.Now().After(deadline) {
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
				t.Fatalf("did not start: %s", out.String())
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			t.Fatal(err)
		}
		if err := cmd.Wait(); err == nil || ctx.Err() != nil {
			t.Fatalf("termination: %v context=%v output=%s", err, ctx.Err(), out.String())
		}
		time.Sleep(1100 * time.Millisecond)
		if _, err := os.Stat(filepath.Join(root, "descendant")); !os.IsNotExist(err) {
			t.Fatal("descendant survived SIGTERM")
		}
	})
}
