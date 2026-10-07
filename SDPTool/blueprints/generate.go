// Package blueprints adapts captured model sources to the pure SDL analyzer.
package blueprints

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/sourcegraph"
)

const Schema = "sdp-blueprint/1"

type Options struct{ Area, From, To, Entry, Task, Output string }
type Capture struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Head           string `json:"head"`
	MetadataDigest string `json:"metadataDigest"`
	SourceDigest   string `json:"sourceDigest"`
	Preliminary    bool   `json:"preliminary"`
}
type CompilerIdentity struct{ Profile, SourceGraph, ExecutableSHA256, GoVersion string }

type Document struct {
	Schema      string `json:"schema"`
	Revision    string `json:"revision"`
	TaskID      string `json:"taskId"`
	TaskDigest  string `json:"taskDigest"`
	Producer    string `json:"producer"`
	Now, Target Capture
	Analysis    *blueprint.Analysis `json:"analysis"`
	Status      string              `json:"status"`
	Entry       string              `json:"entry"`
	Compiler    CompilerIdentity    `json:"compiler"`
}
type Result struct {
	RetainedRevision string `json:"retainedRevision,omitempty"`
	Schema           string `json:"schema"`
	Operation        string `json:"operation"`
	Status           string `json:"status"`
	Revision         string `json:"revision"`
	Path             string `json:"path"`
	Preliminary      bool   `json:"preliminary"`
}

func capture(v model.SourceView) Capture {
	return Capture{v.Artifact.ID, v.Artifact.Kind, v.Artifact.Head, v.Artifact.MetadataDigest, v.LiveDigest, v.Preliminary}
}
func readTask(p string) ([]byte, blueprint.Task, error) {
	var task blueprint.Task
	if err := noSymlinks(p); err != nil {
		return nil, task, err
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, task, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil {
		return nil, task, e
	}
	if !st.Mode().IsRegular() {
		return nil, task, fmt.Errorf("task must be a regular file")
	}
	b, e := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if e != nil {
		return nil, task, e
	}
	if len(b) > 1<<20 {
		return nil, task, fmt.Errorf("task exceeds 1 MiB")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if e = d.Decode(&task); e != nil {
		return nil, task, e
	}
	var extra any
	if e = d.Decode(&extra); e != io.EOF {
		return nil, task, fmt.Errorf("task requires one JSON object")
	}
	return b, task, nil
}
func compile(ctx context.Context, v model.SourceView, entry string) (*sourcegraph.Snapshot, error) {
	if !sourcegraph.ValidPath(entry) {
		return nil, fmt.Errorf("invalid entry path")
	}
	cache := parser.SyntaxCache{}
	files := []*parser.File{}
	seen := map[string]bool{}
	total := 0
	var visit func(string) error
	visit = func(name string) error {
		if e := ctx.Err(); e != nil {
			return e
		}
		if seen[name] {
			return nil
		}
		seen[name] = true
		if !sourcegraph.ValidPath(name) {
			return fmt.Errorf("invalid source path")
		}
		b, ok := v.Files[name]
		if !ok {
			return fmt.Errorf("missing captured source %s", name)
		}
		total += len(b)
		if total > parser.MaxBytes || len(seen) > sourcegraph.MaxFiles {
			return fmt.Errorf("source limits exceeded")
		}
		f, e := cache.Parse(name, string(b))
		if e != nil {
			return e
		}
		files = append(files, f)
		for _, dep := range sourcegraph.Dependencies(f.Model()) {
			if e = visit(dep.Path); e != nil {
				return e
			}
		}
		return nil
	}
	if e := visit(entry); e != nil {
		return nil, e
	}
	s, issues := sourcegraph.Compile(entry, files, false)
	if len(issues) > 0 {
		return nil, issues
	}
	return s, nil
}
func noSymlinks(p string) error {
	abs, e := filepath.Abs(p)
	if e != nil {
		return e
	}
	for p = abs; ; p = filepath.Dir(p) {
		st, e := os.Lstat(p)
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if e == nil && st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink path %s", p)
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return nil
}
func inside(a, b string) bool {
	rel, e := filepath.Rel(a, b)
	return e == nil && (rel == "." || rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
func outputSafe(output string, paths ...string) error {
	if e := noSymlinks(output); e != nil {
		return e
	}
	abs, e := filepath.Abs(output)
	if e != nil {
		return e
	}
	for _, p := range paths {
		if e := noSymlinks(p); e != nil {
			return e
		}
		p, e = filepath.Abs(p)
		if e != nil {
			return e
		}
		if inside(abs, p) || inside(p, abs) {
			return fmt.Errorf("output overlaps input %s", p)
		}
	}
	return nil
}
func Generate(ctx context.Context, o Options) (Result, error) { return generate(ctx, o, nil) }

// beforePublish is an internal failure/freshness test seam, never public task data.
func generate(ctx context.Context, o Options, beforePublish func()) (Result, error) {
	var result Result
	if o.Area == "" || o.From == "" || o.To == "" || o.Entry == "" || o.Task == "" || o.Output == "" {
		return result, fmt.Errorf("all blueprint options are required")
	}
	if e := ctx.Err(); e != nil {
		return result, e
	}
	taskBytes, task, e := readTask(o.Task)
	if e != nil {
		return result, e
	}
	now, e := model.Snapshot(o.Area, o.From)
	if e != nil {
		return result, e
	}
	target, e := model.Snapshot(o.Area, o.To)
	if e != nil {
		return result, e
	}
	if e = outputSafe(o.Output, now.Path, target.Path, o.Task, filepath.Join(o.Area, ".model-operations")); e != nil {
		return result, e
	}
	ns, e := compile(ctx, now, o.Entry)
	if e != nil {
		return result, e
	}
	ts, e := compile(ctx, target, o.Entry)
	if e != nil {
		return result, e
	}
	analysis, e := blueprint.Analyze(ctx, blueprint.Input{Snapshot: ns}, blueprint.Input{Snapshot: ts}, task, blueprint.DefaultLimits())
	if e != nil {
		return result, e
	}
	doc := Document{Schema: Schema, TaskID: task.ID, TaskDigest: documents.Hash(taskBytes), Producer: Schema + "/" + blueprint.Version + "/" + blueprint.PolicyVersion, Now: capture(now), Target: capture(target), Analysis: analysis, Status: "diagnostic-preview"}
	doc.Entry = o.Entry
	executable, e := os.Executable()
	if e != nil {
		return result, e
	}
	f, e := os.Open(executable)
	if e != nil {
		return result, e
	}
	h := sha256.New()
	_, e = io.Copy(h, f)
	closeErr := f.Close()
	if e != nil {
		return result, e
	}
	if closeErr != nil {
		return result, closeErr
	}
	doc.Compiler = CompilerIdentity{ns.Profile(), sourcegraph.Version, fmt.Sprintf("%x", h.Sum(nil)), runtime.Version()}
	identity, _ := json.Marshal(doc)
	doc.Revision = documents.Hash(identity)
	bundle := &documents.Bundle{Files: blueprint.Documents(analysis, task), Manifest: documents.Manifest{Version: Schema, Revision: doc.Revision}}
	bundle.Files["assignment.json"] = taskBytes
	data, _ := json.MarshalIndent(doc, "", "  ")
	bundle.Files["blueprint.json"] = append(data, '\n')
	for _, side := range []struct {
		name string
		v    model.SourceView
	}{{"NOW", now}, {"TARGET", target}} {
		for p, b := range side.v.Files {
			// Model capture already validates paths. Preserve source bytes without treating
			// authored Markdown as generator-owned links during link validation.
			bundle.Files["sources/"+side.name+"/"+p] = append([]byte{}, b...)
		}
	}
	check := *bundle
	check.Files = map[string][]byte{}
	for p, b := range bundle.Files {
		if strings.HasPrefix(p, "sources/") {
			check.Files[p] = nil
		} else {
			check.Files[p] = b
		}
	}
	if e = check.CheckLinks(); e != nil {
		return result, e
	}
	bundle.Seal()
	if beforePublish != nil {
		beforePublish()
	}
	if e = ctx.Err(); e != nil {
		return result, e
	}
	for _, input := range []struct {
		ref string
		old model.SourceView
	}{{o.From, now}, {o.To, target}} {
		current, err := model.Snapshot(o.Area, input.ref)
		if err != nil {
			return result, err
		}
		if capture(current) != capture(input.old) {
			return result, fmt.Errorf("stale model input")
		}
	}
	current, _, e := readTask(o.Task)
	if e != nil {
		return result, e
	}
	if !bytes.Equal(current, taskBytes) {
		return result, fmt.Errorf("stale task input")
	}
	if e = outputSafe(o.Output, now.Path, target.Path, o.Task, filepath.Join(o.Area, ".model-operations")); e != nil {
		return result, e
	}
	if e = bundle.Publish(o.Output); e != nil {
		return result, e
	}
	abs, _ := filepath.Abs(o.Output)
	return Result{Schema: Schema, Operation: "create-blueprint", Status: "diagnostic-preview", Revision: doc.Revision, Path: abs, Preliminary: now.Preliminary || target.Preliminary}, nil
}
