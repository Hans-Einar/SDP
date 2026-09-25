// Package sdptool coordinates language-owned services and document consumers.
package sdptool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

const Version = "sdptool/0.1"

type Failure struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Diagnostic any    `json:"diagnostic,omitempty"`
}

func (f *Failure) Error() string         { return f.Code + ": " + f.Message }
func failure(code string, e error) error { return &Failure{Code: code, Message: e.Error()} }

type Result struct {
	Schema    string `json:"schema"`
	Operation string `json:"operation"`
	Source    string `json:"source"`
	Profile   string `json:"profile"`
	Revision  string `json:"revision"`
	Entry     string `json:"entry"`
	Directory string `json:"directory"`
	Ownership string `json:"ownership"`
}
type PreviewOptions struct {
	Source, Output, Renderer, Viewpoint, URI, Revision string
	operation                                          string
}

func readSource(name string) ([]byte, error) {
	f, e := os.Open(name)
	if e != nil {
		return nil, failure("source", e)
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, failure("source", fmt.Errorf("expected regular file"))
	}
	b, e := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
	if e != nil {
		return nil, e
	}
	if len(b) > parser.MaxBytes {
		return nil, failure("limit", fmt.Errorf("source exceeds %d bytes", parser.MaxBytes))
	}
	return b, nil
}
func loadModel(name string) (*viewpoint.Views, []byte, error) {
	b, e := readSource(name)
	if e != nil {
		return nil, nil, e
	}
	v, e := viewpoint.New(string(b))
	if e != nil {
		f := &Failure{Code: "model", Message: e.Error()}
		if d, ok := e.(parser.Diagnostic); ok {
			f.Diagnostic = parser.Data(d)
		}
		return nil, nil, f
	}
	return v, b, nil
}

// physical resolves existing ancestors too, so an absent output under a symlink
// cannot hide source containment from the publication guard.
func physical(p string) (string, error) {
	a, e := filepath.Abs(p)
	if e != nil {
		return "", e
	}
	if _, e = os.Lstat(a); e == nil {
		return filepath.EvalSymlinks(a)
	} else if !os.IsNotExist(e) {
		return "", e
	}
	parent, e := physical(filepath.Dir(a))
	if e != nil {
		return "", e
	}
	return filepath.Join(parent, filepath.Base(a)), nil
}
func outside(output, source string) error {
	a, e := physical(output)
	if e != nil {
		return e
	}
	b, e := physical(source)
	if e != nil {
		return e
	}
	rel, e := filepath.Rel(a, b)
	if e != nil {
		return e
	}
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("source must be outside output directory")
	}
	return nil
}
func renderer(name string) (documents.Renderer, error) {
	if name == "" {
		return nil, nil
	}
	p, e := exec.LookPath(name)
	if e != nil {
		return nil, failure("tool", e)
	}
	p, e = filepath.Abs(p)
	if e != nil {
		return nil, e
	}
	r, e := documents.NewMmdr(p)
	if e != nil {
		return nil, failure("tool", e)
	}
	return r, nil
}
func Preview(ctx context.Context, o PreviewOptions) (Result, error) {
	var result Result
	if o.Source == "" || o.Output == "" {
		return result, failure("arguments", fmt.Errorf("saved source and --output are required"))
	}
	if e := ctx.Err(); e != nil {
		return result, failure("canceled", e)
	}
	if e := outside(o.Output, o.Source); e != nil {
		return result, failure("output", e)
	}
	source, e := physical(o.Source)
	if e != nil {
		return result, failure("source", e)
	}
	v, _, e := loadModel(source)
	if e != nil {
		return result, e
	}
	if o.Revision != "" && o.Revision != v.Revision {
		return result, failure("stale", fmt.Errorf("source revision changed; refresh inventory"))
	}
	q := viewpoint.Query{Viewpoint: o.Viewpoint, Depth: 2, Direction: "both"}
	target := "main"
	if o.URI != "" {
		if o.Viewpoint != "" {
			return result, failure("arguments", fmt.Errorf("choose --uri or --viewpoint"))
		}
		s, err := viewpoint.ParseURI(o.URI)
		if err != nil {
			return result, failure("selection", err)
		}
		q = s.Query
		target = s.Target
	}
	if q.Viewpoint == "" {
		// A compact, content-derived first diagram; users select other views explicitly.
		for _, vp := range []string{"VP02", "VP01", "VP03", "VP04", "VP06", "VP09", "VP08", "VP10", "VP11"} {
			for _, d := range v.Diagrams {
				if strings.HasPrefix(d.ID, vp+"-") && len(d.Nodes)+len(d.Elements) > 0 {
					q.Viewpoint = vp
					q.Diagram = d.ID
					break
				}
			}
			if q.Viewpoint != "" {
				break
			}
		}
		if q.Viewpoint == "" {
			q.Viewpoint = "VP11"
		}
	}
	selected, e := v.Select(q)
	if e != nil {
		return result, failure("selection", e)
	}
	if len(selected.Diagrams) > 24 {
		return result, failure("limit", fmt.Errorf("selection exceeds 24 diagrams; select a diagram or focus"))
	}
	r, e := renderer(o.Renderer)
	if e != nil {
		return result, e
	}
	b, e := documents.Selected(ctx, v, q, r)
	if e != nil {
		return result, failure("render", e)
	}
	b.Put("delivery.txt", target+"\n")
	output, e := filepath.Abs(o.Output)
	if e != nil {
		return result, e
	}
	operation := o.operation
	if operation == "" {
		operation = "preview"
	}
	result = Result{Version, operation, source, "design-core/0.5", v.Revision, filepath.Join(output, "entry.md"), output, "caller-owned; remove directory after consumer release"}
	j, _ := json.MarshalIndent(result, "", "  ")
	b.Files["sdptool.json"] = append(j, '\n')
	b.Seal()
	size := 0
	for _, data := range b.Files {
		size += len(data)
	}
	if len(b.Files) > 128 || size > 32<<20 {
		return Result{}, failure("limit", fmt.Errorf("bundle exceeds 128 files or 32 MiB"))
	}
	if e = ctx.Err(); e != nil {
		return Result{}, failure("canceled", e)
	}
	current, e := readSource(source)
	if e != nil {
		return Result{}, e
	}
	if documents.Hash(current) != v.Revision {
		return Result{}, failure("stale", fmt.Errorf("source changed during generation"))
	}
	if e = b.Publish(output); e != nil {
		return Result{}, failure("output", e)
	}
	return result, nil
}
