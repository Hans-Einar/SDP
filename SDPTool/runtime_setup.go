package sdptool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type SetupResult struct {
	Schema     string   `json:"schema"`
	Operation  string   `json:"operation"`
	Entrypoint string   `json:"entrypoint"`
	Executable string   `json:"executable"`
	Version    string   `json:"version"`
	Revision   string   `json:"revision"`
	Status     string   `json:"status"`
	Warnings   []string `json:"warnings,omitempty"`
}

// RegisterRuntime registers this engine without replacing another installation.
func RegisterRuntime(ctx context.Context) (SetupResult, error) {
	return registerRuntime(ctx, os.Executable)
}
func registerRuntime(ctx context.Context, executable func() (string, error)) (SetupResult, error) {
	r := SetupResult{Schema: Version, Operation: "setup", Version: BuildVersion, Revision: BuildRevision}
	if runtime.GOOS != "linux" {
		return r, fmt.Errorf("setup currently supports Linux; this platform has no qualified registration contract")
	}
	if err := ctx.Err(); err != nil {
		return r, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return r, err
	}
	if !filepath.IsAbs(home) {
		return r, fmt.Errorf("HOME must be an absolute directory")
	}
	source, err := executable()
	if err != nil {
		return r, err
	}
	source, err = filepath.EvalSymlinks(source)
	if err != nil {
		return r, fmt.Errorf("running executable is no longer available: %w", err)
	}
	info, err := os.Stat(source)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		return r, fmt.Errorf("running executable must remain a regular executable file")
	}
	// A temporary source outside HOME cannot be a durable update target. Isolated
	// HOME fixtures remain valid; production gh installs live beneath HOME.
	within := func(parent, path string) bool {
		rel, e := filepath.Rel(parent, path)
		return e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
	}
	if within(os.TempDir(), source) && !within(home, source) {
		return r, fmt.Errorf("refusing transient executable %s; install the extension, then run gh sdp setup", source)
	}
	r.Executable = source
	r.Entrypoint = filepath.Join(home, ".local", "bin", "sdptool")
	if source == r.Entrypoint {
		return r, fmt.Errorf("sdptool is an unmanaged file at the public entrypoint; move it aside explicitly, then run gh sdp setup")
	}
	for _, dir := range []string{filepath.Join(home, ".local"), filepath.Join(home, ".local", "bin")} {
		if err = os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return r, err
		}
		st, e := os.Lstat(dir)
		if e != nil || !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
			return r, fmt.Errorf("registration directory must be a real directory: %s", dir)
		}
	}
	r.Status = "registered"
	if err = os.Symlink(source, r.Entrypoint); err != nil {
		if !os.IsExist(err) {
			return r, err
		}
		target, e := os.Readlink(r.Entrypoint)
		if e != nil || target != source {
			return r, fmt.Errorf("refusing to replace unmanaged entrypoint %s; inspect and move it aside explicitly before retrying setup", r.Entrypoint)
		}
		r.Status = "already-registered"
	}
	shell, err := exec.LookPath("sdptool")
	if err != nil {
		r.Warnings = append(r.Warnings, "Add $HOME/.local/bin to PATH for shell use; XFMD uses the registered absolute path directly.")
	} else if resolved, e := filepath.EvalSymlinks(shell); e != nil || resolved != source {
		r.Warnings = append(r.Warnings, fmt.Sprintf("PATH selects another sdptool (%s); remove that obsolete installation or put $HOME/.local/bin first. XFMD uses %s.", shell, r.Entrypoint))
	}
	return r, nil
}
