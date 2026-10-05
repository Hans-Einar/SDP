package model

import (
	"os"
	"path/filepath"
	"strings"
)

// SourceView is an owned read snapshot. Bytes are never aliases into mutable storage.
// Reading WORK does not commit or freeze it. Consumers must label it preliminary.
type SourceView struct {
	Result
	Files Files `json:"files"` // encoding/json emits byte values as base64
}

func Snapshot(area, ref string) (SourceView, error) {
	p, e := resolve(area, ref)
	if e != nil {
		return SourceView{}, e
	}
	a, f, e := capture(p)
	if e != nil {
		return SourceView{}, e
	}
	return SourceView{result("snapshot", p, a, f), f}, nil
}

// Inspect recognizes only artifact-named directories with their own metadata.
// A malformed advertised artifact returns an error, never ungoverned sources.
func Inspect(dir string) (*Artifact, error) {
	name := filepath.Base(dir)
	kind := false
	for _, p := range []string{"WORK--", "PROPOSAL--", "CANDIDATE--", "RELEASE--"} {
		kind = kind || strings.HasPrefix(name, p)
	}
	if !kind {
		return nil, nil
	}
	if _, e := os.Lstat(filepath.Join(dir, metadataName(dir))); os.IsNotExist(e) {
		return nil, nil
	} else if e != nil {
		return nil, e
	}
	a, _, e := capture(dir)
	if e != nil {
		return nil, e
	}
	return &a, nil
}

// OperationArea distinguishes reserved model transaction storage from an
// arbitrary folder named .model-operations. Writers always create its lock file.
func OperationArea(dir string) bool {
	st, e := os.Lstat(filepath.Join(dir, "lock"))
	return filepath.Base(dir) == operations && e == nil && st.Mode().IsRegular()
}
