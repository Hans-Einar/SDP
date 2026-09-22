package documents

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var markdownLink = regexp.MustCompile(`\]\(([^)]+)\)`)

func (b *Bundle) CheckLinks() error {
	for n, data := range b.Files {
		if !strings.HasSuffix(n, ".md") {
			continue
		}
		for _, m := range markdownLink.FindAllSubmatch(data, -1) {
			p, e := LocalPath(n, string(m[1]))
			if e != nil {
				return e
			}
			if p != "" {
				if _, ok := b.Files[p]; !ok {
					return fmt.Errorf("dead link %s -> %s", n, p)
				}
			}
		}
	}
	return nil
}
func safe(n string) bool {
	return n != "" && n != "." && !filepath.IsAbs(n) && filepath.Clean(n) == n && n != ".." && !strings.HasPrefix(n, "../") && !strings.Contains(n, "\\")
}

// Publish validates everything before replacing the directory. Unmanaged notes
// are copied unchanged; conflicting edits to generated files are rejected.
func (b *Bundle) Publish(output string) error {
	abs, e := filepath.Abs(output)
	if e != nil {
		return e
	}
	if abs == "/" {
		return fmt.Errorf("invalid output")
	}
	parent := filepath.Dir(abs)
	if e = os.MkdirAll(parent, 0700); e != nil {
		return e
	}
	lock := abs + ".lock"
	if e = os.Mkdir(lock, 0700); e != nil {
		return fmt.Errorf("output busy: %w", e)
	}
	defer os.Remove(lock)
	stage, e := os.MkdirTemp(parent, ".sdl-publish-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(stage)
	old := Manifest{}
	if j, err := os.ReadFile(filepath.Join(abs, "manifest.json")); err == nil {
		if e = json.Unmarshal(j, &old); e != nil {
			return fmt.Errorf("invalid previous manifest: %w", e)
		}
	}
	e = filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
		if os.IsNotExist(err) && p == abs {
			return nil
		}
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in output: %s", p)
		}
		if d.IsDir() {
			return nil
		}
		n, _ := filepath.Rel(abs, p)
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if n == "manifest.json" {
			return nil
		}
		if hash, ok := old.Outputs[n]; ok {
			if Hash(data) != hash {
				return fmt.Errorf("locally changed generated file: %s", n)
			}
			return nil
		}
		if _, ok := b.Files[n]; ok {
			return fmt.Errorf("unmanaged file conflicts: %s", n)
		}
		return write(stage, n, data)
	})
	if e != nil {
		return e
	}
	for n, data := range b.Files {
		if e = write(stage, n, data); e != nil {
			return e
		}
	}
	backup := stage + "-previous"
	exists := false
	if _, e = os.Lstat(abs); e == nil {
		exists = true
		if e = os.Rename(abs, backup); e != nil {
			return e
		}
	} else if !os.IsNotExist(e) {
		return e
	}
	if e = os.Rename(stage, abs); e != nil {
		if exists {
			os.Rename(backup, abs)
		}
		return e
	}
	if exists {
		return os.RemoveAll(backup)
	}
	return nil
}
func write(root, n string, b []byte) error {
	if !safe(n) {
		return fmt.Errorf("unsafe output path %q", n)
	}
	p := filepath.Join(root, n)
	if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		return e
	}
	return os.WriteFile(p, b, 0600)
}
