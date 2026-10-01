package sdptool

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Host struct{ Viewer, SDLTool, Renderer string }

func program(value, fallback string) (string, error) {
	if value == "" {
		value = fallback
	}
	p, e := exec.LookPath(value)
	if e != nil {
		return "", failure("tool", e)
	}
	return filepath.Abs(p)
}

// boundedOutput avoids unbounded viewer diagnostics and capability probing.
type boundedOutput struct{ bytes.Buffer }

func (b *boundedOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 1<<20 {
		return 0, fmt.Errorf("child output limit")
	}
	return b.Buffer.Write(p)
}
func ViewPlan(ctx context.Context, p Project, modelID string, h Host) error {
	if p.Inventory.ImplementationPlan == "" {
		return failure("unavailable", fmt.Errorf("choose --plan: no unique ImplementationPlan discovered"))
	}
	plan, e := resolvePath(p.Root, p.Inventory.ImplementationPlan)
	if e != nil {
		return e
	}
	if _, e = boundedFile(plan); e != nil {
		return failure("source", e)
	}
	viewer, e := program(h.Viewer, "xfmd")
	if e != nil {
		return e
	}
	args := []string{plan}
	var session string
	if len(p.Inventory.Models) > 0 {
		m, source, err := p.model(modelID, false)
		if err != nil {
			return err
		}
		sdl, err := program(h.SDLTool, "sdl")
		if err != nil {
			return err
		}
		// Probe the selected host, not an unrelated installed viewer.
		probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var help boundedOutput
		c := exec.CommandContext(probeCtx, viewer, "--help")
		c.Stdout = &help
		c.Stderr = &help
		if err = c.Run(); err != nil {
			return failure("tool", fmt.Errorf("viewer capability probe: %w", err))
		}
		if !strings.Contains(help.String(), "--navigator") || !strings.Contains(help.String(), "--sdl-tool") {
			return failure("unsupported", fmt.Errorf("viewer lacks SDL navigation arguments"))
		}
		v, input, err := loadModel(source)
		if err != nil {
			return err
		}
		if v.Profile != m.Profile || (v.System != "" && v.System != m.System) {
			return failure("inventory", fmt.Errorf("discovered profile/System does not match validated source"))
		}
		parent := os.Getenv("XDG_RUNTIME_DIR")
		session, err = os.MkdirTemp(parent, "sdptool-")
		if err != nil {
			return failure("output", err)
		}
		defer os.RemoveAll(session)
		b, err := documents.Build(ctx, v, documents.Options{Navigator: true, Project: p.Inventory.ProjectID})
		if err != nil {
			return failure("model", err)
		}
		directory := filepath.Join(session, "navigation")
		if err = input.Outside(directory); err != nil {
			return failure("output", err)
		}
		if err = input.Fresh(); err != nil {
			return failure("stale", err)
		}
		if err = b.Publish(directory); err != nil {
			return failure("output", err)
		}
		args = append(args, "--navigator", filepath.Join(directory, "navigator.md"), "--sdl-tool", sdl, "--sdl-source", source, "--project", p.Inventory.ProjectID, "--window-id", filepath.Base(session))
		if h.Renderer != "" {
			r, err := program(h.Renderer, "")
			if err != nil {
				return err
			}
			args = append(args, "--renderer", r)
		}
	}
	if e = ctx.Err(); e != nil {
		return failure("canceled", e)
	}
	cmd := exec.CommandContext(ctx, viewer, args...)
	cmd.WaitDelay = time.Second
	var output boundedOutput
	cmd.Stdout = &output
	cmd.Stderr = &output
	if e = cmd.Run(); e != nil {
		if ctx.Err() != nil {
			return failure("canceled", ctx.Err())
		}
		return failure("viewer", fmt.Errorf("%w: %s", e, output.String()))
	}
	return nil
}
