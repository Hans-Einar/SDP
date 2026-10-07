package blueprints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRetainedRevisionsAndTampering(t *testing.T) {
	o, _ := setup(t)
	o.Output = ""
	catalogue := filepath.Join(t.TempDir(), "Blueprints")
	r, e := GenerateRetained(context.Background(), o, catalogue)
	if e != nil {
		t.Fatal(e)
	}
	again, e := GenerateRetained(context.Background(), o, catalogue)
	if e != nil || again.Path != r.Path {
		t.Fatal(again, e)
	}
	c := Scan(catalogue)
	if c.State != "available" || len(c.Entries) != 1 || c.Entries[0].State != "validated" {
		t.Fatalf("%+v", c)
	}
	old := read(t, filepath.Join(r.Path, "index.md"))
	data := read(t, o.Task)
	var task map[string]any
	json.Unmarshal(data, &task)
	task["Intent"] = "Changed intent"
	data, _ = json.Marshal(task)
	os.WriteFile(o.Task, data, 0600)
	newer, e := GenerateRetained(context.Background(), o, catalogue)
	if e != nil || newer.Path == r.Path {
		t.Fatal(e)
	}
	if string(old) != string(read(t, filepath.Join(r.Path, "index.md"))) {
		t.Fatal("old revision changed")
	}
	c = Scan(catalogue)
	if len(c.Entries) != 2 {
		t.Fatal(c)
	}
	os.WriteFile(filepath.Join(newer.Path, "index.md"), []byte("tampered"), 0600)
	c = Scan(catalogue)
	bad := 0
	for _, x := range c.Entries {
		if x.State == "invalid" {
			bad++
		}
	}
	if bad != 1 {
		t.Fatal(c)
	}
	if _, e = GenerateRetained(context.Background(), o, catalogue); e == nil {
		t.Fatal("overwrote retained edit")
	}
	if string(read(t, filepath.Join(newer.Path, "index.md"))) != "tampered" {
		t.Fatal("edit lost")
	}
}
func TestCatalogueInvalidCases(t *testing.T) {
	for _, kind := range []string{"missing", "extra", "symlink", "metadata", "duplicate", "unknown-schema"} {
		t.Run(kind, func(t *testing.T) {
			o, _ := setup(t)
			o.Output = ""
			dir := filepath.Join(t.TempDir(), "Blueprints")
			r, e := GenerateRetained(context.Background(), o, dir)
			if e != nil {
				t.Fatal(e)
			}
			switch kind {
			case "missing":
				os.Remove(filepath.Join(r.Path, "index.md"))
			case "extra":
				os.WriteFile(filepath.Join(r.Path, "notes"), []byte("unmanaged"), 0600)
			case "symlink":
				os.Remove(filepath.Join(r.Path, "index.md"))
				if e = os.Symlink(o.Task, filepath.Join(r.Path, "index.md")); e != nil {
					t.Skip(e)
				}
			case "metadata":
				os.WriteFile(filepath.Join(r.Path, "manifest.json"), []byte("{}"), 0600)
			case "unknown-schema":
				p := filepath.Join(r.Path, "manifest.json")
				b := read(t, p)
				b = []byte(strings.Replace(string(b), Schema, "future/99", 1))
				os.WriteFile(p, b, 0600)
			case "duplicate":
				dest := filepath.Join(filepath.Dir(r.Path), strings.Repeat("a", 64))
				if e = os.CopyFS(dest, os.DirFS(r.Path)); e != nil {
					t.Fatal(e)
				}
			}
			c := Scan(dir)
			found := false
			for _, x := range c.Entries {
				if x.State == "invalid" {
					found = true
				}
			}
			if !found {
				t.Fatal(c)
			}
		})
	}
}
func TestCatalogueAbsenceLimitsAndAliases(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "Blueprints")
	if Scan(dir).State != "absent" {
		t.Fatal("absence")
	}
	os.Mkdir(dir, 0700)
	if Scan(dir).State != "empty" {
		t.Fatal("empty")
	}
	for i := 0; i <= MaxCatalogueEntries; i++ {
		os.Mkdir(filepath.Join(dir, strings.Repeat("a", 60)+string(rune(1000+i))), 0700)
	}
	if Scan(dir).State != "unavailable" {
		t.Fatal("unbounded scan")
	}
	alias := filepath.Join(t.TempDir(), "link")
	if e := os.Symlink(dir, alias); e != nil {
		t.Skip(e)
	}
	if Scan(alias).State != "unavailable" {
		t.Fatal("alias followed")
	}
	b := &budget{}
	p := filepath.Join(t.TempDir(), "big")
	f, e := os.Create(p)
	if e != nil {
		t.Fatal(e)
	}
	f.Truncate(MaxCatalogueBytes + 1)
	f.Close()
	if _, e = b.read(p, MaxCatalogueBytes); e == nil {
		t.Fatal("byte budget not enforced")
	}
}

func TestCatalogueRejectsRehashedManifestEdits(t *testing.T) {
	for _, name := range []string{"sources/NOW/System.design", "index.md"} {
		t.Run(name, func(t *testing.T) {
			o, _ := setup(t)
			o.Output = ""
			dir := filepath.Join(t.TempDir(), "Blueprints")
			r, e := GenerateRetained(context.Background(), o, dir)
			if e != nil {
				t.Fatal(e)
			}
			file := filepath.Join(r.Path, filepath.FromSlash(name))
			data := append(read(t, file), []byte("changed")...)
			os.WriteFile(file, data, 0600)
			var manifest map[string]any
			mp := filepath.Join(r.Path, "manifest.json")
			json.Unmarshal(read(t, mp), &manifest)
			manifest["outputs"].(map[string]any)[name] = documents.Hash(data)
			raw, _ := json.Marshal(manifest)
			os.WriteFile(mp, raw, 0600)
			c := Scan(dir)
			if len(c.Entries) != 1 || c.Entries[0].State != "invalid" {
				t.Fatal("rehashed edit accepted", c)
			}
		})
	}
}
func TestCatalogueGlobalEntryBudget(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < MaxCatalogueEntries; i++ {
		task := filepath.Join(dir, fmt.Sprintf("%064x", i))
		os.Mkdir(task, 0700)
		if i == 0 {
			for j := 0; j < MaxCatalogueEntries; j++ {
				os.Mkdir(filepath.Join(task, fmt.Sprintf("%064x", j)), 0700)
			}
		}
	}
	c := Scan(dir)
	if len(c.Entries) > MaxCatalogueEntries {
		t.Fatal("entry limit exceeded", len(c.Entries))
	}
	if c.State != "unavailable" || len(c.Entries) != 1 {
		t.Fatal("scan must stop at first exhausted enumeration budget", c)
	}
}

func TestCatalogueDirectoryOverflowExhaustsSharedBudget(t *testing.T) {
	o, _ := setup(t)
	o.Output = ""
	r, err := GenerateRetained(context.Background(), o, filepath.Join(t.TempDir(), "Blueprints"))
	if err != nil {
		t.Fatal(err)
	}
	b := &budget{nodes: MaxCatalogueFiles - 2}
	if _, _, err = verify(r.Path, b); err == nil || b.nodes != MaxCatalogueFiles {
		t.Fatal("overflow must exhaust shared directory budget", err, b.nodes)
	}
	if _, _, err = verify(r.Path, b); err == nil || b.nodes != MaxCatalogueFiles {
		t.Fatal("next candidate must retain exhausted directory budget", err, b.nodes)
	}
}
