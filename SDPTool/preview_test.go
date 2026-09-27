package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sample = "language design-core version 0.5.\nunit Child.\ncontainer Main.\nMain contains Child.\n"

func fixture(t *testing.T) (string, string) {
	t.Helper()
	d := t.TempDir()
	p := filepath.Join(d, "source with spaces.design")
	if e := os.WriteFile(p, []byte(sample), 0600); e != nil {
		t.Fatal(e)
	}
	return p, filepath.Join(d, "output bundle")
}
func TestPreviewFreshSourceAndPreservation(t *testing.T) {
	p, d := fixture(t)
	o := PreviewOptions{Source: p, Output: d}
	r, e := Preview(context.Background(), o)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(r.Entry)
	if e != nil || !bytes.Contains(b, []byte("Main")) {
		t.Fatalf("entry %s %v", b, e)
	}
	m := map[string]any{}
	j, _ := os.ReadFile(filepath.Join(d, "manifest.json"))
	if e = json.Unmarshal(j, &m); e != nil {
		t.Fatal(e)
	}
	if len(m["diagrams"].([]any)) != 1 {
		t.Fatal("default must be bounded")
	}
	old, _ := os.ReadFile(p)
	if string(old) != sample {
		t.Fatal("source changed")
	}
	os.WriteFile(p, []byte(strings.ReplaceAll(sample, "Child", "Edited")), 0600)
	o.Revision = r.Revision
	if _, e = Preview(context.Background(), o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatalf("stale %v", e)
	}
	unchanged, _ := os.ReadFile(r.Entry)
	if !bytes.Equal(b, unchanged) {
		t.Fatal("failed generation replaced last valid preview")
	}
	o.Revision = ""
	r2, e := Preview(context.Background(), o)
	if e != nil || r.Revision == r2.Revision {
		t.Fatalf("regeneration %v", e)
	}
}
func TestPreviewFailuresPreserveSource(t *testing.T) {
	for _, kind := range []string{"invalid", "unsupported", "renderer", "selection", "containment", "symlink", "conflict", "extra"} {
		t.Run(kind, func(t *testing.T) {
			p, d := fixture(t)
			o := PreviewOptions{Source: p, Output: d}
			switch kind {
			case "invalid":
				os.WriteFile(p, []byte(sample+"unknown thing."), 0600)
			case "unsupported":
				os.WriteFile(p, []byte("language runtime-core version 0.1."), 0600)
			case "renderer":
				o.Renderer = "/missing/renderer"
			case "selection":
				o.Viewpoint = "VP99"
			case "containment":
				o.Output = filepath.Dir(p)
			case "symlink":
				alias := filepath.Join(t.TempDir(), "alias")
				os.Symlink(filepath.Dir(p), alias)
				o.Output = alias
			case "conflict":
				os.Mkdir(d, 0700)
				os.WriteFile(filepath.Join(d, "entry.md"), []byte("authored note"), 0600)
			case "extra":
				var out, err bytes.Buffer
				if Run(context.Background(), []string{"preview", p, "--output", d, "surplus"}, &out, &err) == 0 {
					t.Fatal("extra argument accepted")
				}
				return
			}
			before, _ := os.ReadFile(p)
			if _, e := Preview(context.Background(), o); e == nil {
				t.Fatal("expected failure")
			}
			after, _ := os.ReadFile(p)
			if !bytes.Equal(before, after) {
				t.Fatal("source changed")
			}
		})
	}
}
func TestCanceledAndStructuredDiagnostics(t *testing.T) {
	p, d := fixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := Preview(ctx, PreviewOptions{Source: p, Output: d}); e == nil {
		t.Fatal("canceled request accepted")
	}
	os.WriteFile(p, []byte(sample+"missing syntax"), 0600)
	var out, errs bytes.Buffer
	if Run(context.Background(), []string{"preview", p, "--output", d}, &out, &errs) == 0 {
		t.Fatal("invalid model accepted")
	}
	if !strings.Contains(errs.String(), "diagnostic") {
		t.Fatal(errs.String())
	}
}

func TestFailedRefreshKeepsPreviousBundle(t *testing.T) {
	p, d := fixture(t)
	o := PreviewOptions{Source: p, Output: d}
	r, e := Preview(context.Background(), o)
	if e != nil {
		t.Fatal(e)
	}
	before, _ := os.ReadFile(r.Entry)
	for _, kind := range []string{"parse", "render"} {
		os.WriteFile(p, []byte(sample), 0600)
		o.Renderer = ""
		if kind == "parse" {
			os.WriteFile(p, []byte("invalid source"), 0600)
		} else {
			o.Renderer = executable(t, "exit 9\n")
		}
		if _, e = Preview(context.Background(), o); e == nil {
			t.Fatal("expected failure", kind)
		}
		after, _ := os.ReadFile(r.Entry)
		if !bytes.Equal(before, after) {
			t.Fatal("failed refresh damaged last valid bundle", kind)
		}
	}
}
