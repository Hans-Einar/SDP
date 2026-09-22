// Package devhost rebuilds and restarts an explicitly selected local Go app.
// Failed compilation leaves the running process intact. This is not hot swapping.
package devhost

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Config struct {
	Root, Package, Go, Tags string
	Args                    []string
	Log                     io.Writer
}
type Manager struct {
	mu       sync.Mutex
	config   Config
	dir      string
	revision uint64
	process  *exec.Cmd
	done     chan error
	closed   bool
}

func New(config Config) (*Manager, error) {
	if config.Root == "" {
		config.Root = "."
	}
	root, err := filepath.Abs(config.Root)
	if err != nil {
		return nil, err
	}
	config.Root = root
	if config.Package == "" {
		return nil, fmt.Errorf("Select an explicit Go package")
	}
	if config.Go == "" {
		config.Go = filepath.Join(runtime.GOROOT(), "bin", "go")
	}
	if config.Log == nil {
		config.Log = io.Discard
	}
	dir, err := os.MkdirTemp("", "sdl-dev-")
	if err != nil {
		return nil, err
	}
	return &Manager{config: config, dir: dir}, nil
}
func (m *Manager) Rebuild(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return fmt.Errorf("closed: Go development host")
	}
	m.revision++
	binary := filepath.Join(m.dir, fmt.Sprintf("app-%d", m.revision))
	args := []string{"build", "-o", binary}
	if m.config.Tags != "" {
		args = append(args, "-tags", m.config.Tags)
	}
	args = append(args, m.config.Package)
	command := exec.CommandContext(ctx, m.config.Go, args...)
	command.Dir = m.config.Root
	group(command)
	command.Cancel = func() error {
		if command.Process == nil {
			return os.ErrProcessDone
		}
		return kill(command.Process)
	}
	command.WaitDelay = 2 * time.Second
	output, err := command.CombinedOutput()
	if err != nil {
		os.Remove(binary)
		return fmt.Errorf("go-build: %w\n%.32000s", err, output)
	}
	if err = ctx.Err(); err != nil {
		os.Remove(binary)
		return err
	}
	m.stop()
	child := exec.Command(binary, m.config.Args...)
	child.Dir = m.config.Root
	group(child)
	child.Stdout = m.config.Log
	child.Stderr = m.config.Log
	if err = child.Start(); err != nil {
		os.Remove(binary)
		return fmt.Errorf("go-start: %w", err)
	}
	m.process = child
	m.done = make(chan error, 1)
	done := m.done
	go func() { done <- child.Wait(); os.Remove(binary) }()
	return nil
}
func (m *Manager) stop() {
	if m.process == nil {
		return
	}
	_ = interrupt(m.process.Process)
	select {
	case <-m.done:
	case <-time.After(2 * time.Second):
		_ = kill(m.process.Process)
		<-m.done
	}
	m.process = nil
	m.done = nil
}
func (m *Manager) PID() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.process == nil {
		return 0
	}
	return m.process.Process.Pid
}
func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return
	}
	m.closed = true
	m.stop()
	os.RemoveAll(m.dir)
}
