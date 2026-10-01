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

func sessionTab(t *testing.T, p Project) Node {
	t.Helper()
	for _, n := range p.Navigation.Nodes {
		if n.ID == "sessions" {
			return n
		}
	}
	t.Fatal("missing Sessions tab")
	return Node{}
}

func TestSessionDiscoveryBrowseAndRefresh(t *testing.T) {
	root, _ := projectFixture(t)
	p, e := Discover(root)
	if e != nil || p.Capabilities["sessions"] != "absent" || sessionTab(t, p).State != "absent" {
		t.Fatal(p, e)
	}
	dir := filepath.Join(root, "SDP", "Sessions")
	if e = os.MkdirAll(filepath.Join(dir, "Archive"), 0700); e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(dir, "Archive", "session-#0001--A session.md")
	original := []byte("# A session\n\n| Status | active |\n")
	if e = os.WriteFile(path, original, 0600); e != nil {
		t.Fatal(e)
	}
	p, e = Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if p.Capabilities["sessions"] != "discovered" || p.Inventory.Sessions != "SDP/Sessions" {
		t.Fatal(p)
	}
	tab := sessionTab(t, p)
	if tab.State != "available" || len(tab.Children) != 1 || tab.Children[0] != "files/Sessions/Archive" {
		t.Fatal(tab)
	}
	var target *Target
	seen := map[string]bool{}
	for _, n := range p.Navigation.Nodes {
		if seen[n.ID] {
			t.Fatal("duplicate node", n.ID)
		}
		seen[n.ID] = true
		if n.Label == filepath.Base(path) {
			target = n.Target
		}
	}
	for _, n := range p.Navigation.Nodes {
		for _, id := range n.Children {
			if !seen[id] {
				t.Fatal("dangling child", id)
			}
		}
	}
	if target == nil || target.Operation != "open" || target.Path != path || len(target.Revision) != 64 {
		t.Fatal(target)
	}
	q, e := Discover(filepath.Join(root, "SDP"))
	if e != nil || q.Navigation.InventoryRevision != p.Navigation.InventoryRevision {
		t.Fatal("root/area mismatch", e)
	}
	for _, args := range [][]string{{root, "discover"}, {root, "tree"}, {root, "discover", "--json"}} {
		var out, errs bytes.Buffer
		if code := Run(context.Background(), args, &out, &errs); code != 0 {
			t.Fatal(code, errs.String())
		}
		if args[len(args)-1] == "--json" {
			var result Project
			if e = json.Unmarshal(out.Bytes(), &result); e != nil || result.Inventory.Sessions != "SDP/Sessions" {
				t.Fatal(e, out.String())
			}
		} else if !strings.Contains(strings.ToLower(out.String()), "sessions") {
			t.Fatal(out.String())
		}
		if args[1] == "tree" && !strings.Contains(out.String(), filepath.Base(path)) {
			t.Fatal(out.String())
		}
	}
	// Reads must not mutate the authored document or create an index.
	if b, e := os.ReadFile(path); e != nil || !bytes.Equal(b, original) {
		t.Fatal("document changed", e)
	}
	if _, e := os.Stat(filepath.Join(root, "SDP/navigation.json")); !os.IsNotExist(e) {
		t.Fatal("index created")
	}
	// Content changes invalidate the snapshot even with preserved size/mtime.
	st, _ := os.Stat(path)
	if e = os.WriteFile(path, bytes.ReplaceAll(original, []byte("active"), []byte("paused")), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.Chtimes(path, st.ModTime(), st.ModTime()); e != nil {
		t.Fatal(e)
	}
	q, e = Discover(root)
	if e != nil || q.Navigation.InventoryRevision == p.Navigation.InventoryRevision {
		t.Fatal("edit not refreshed", e)
	}
	renamed := filepath.Join(dir, "session-#0002--Renamed.md")
	if e = os.Rename(path, renamed); e != nil {
		t.Fatal(e)
	}
	q, e = Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	for _, n := range q.Navigation.Nodes {
		if n.Target != nil && n.Target.Path == path {
			t.Fatal("stale old target")
		}
	}
	if e = os.RemoveAll(dir); e != nil {
		t.Fatal(e)
	}
	q, e = Discover(root)
	if e != nil || sessionTab(t, q).State != "absent" {
		t.Fatal(q, e)
	}
}

func TestSessionDirectoryStatesAndReadLimits(t *testing.T) {
	root, _ := projectFixture(t)
	dir := filepath.Join(root, "SDP", "Sessions")
	if e := os.Mkdir(dir, 0700); e != nil {
		t.Fatal(e)
	}
	p, e := Discover(root)
	if e != nil || sessionTab(t, p).State != "empty" {
		t.Fatal(p, e)
	}
	if e = os.WriteFile(filepath.Join(dir, "large.md"), bytes.Repeat([]byte("x"), (1<<20)+1), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Guide\n"), 0600); e != nil {
		t.Fatal(e)
	}
	external := t.TempDir()
	if e = os.WriteFile(filepath.Join(external, "secret.md"), []byte("outside"), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(external, filepath.Join(dir, "linked")); e != nil {
		t.Fatal(e)
	}
	p, e = Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	for _, n := range p.Navigation.Nodes {
		switch n.Label {
		case "large.md":
			if n.State != "unavailable" || n.Target != nil || n.Diagnostic == "" {
				t.Fatal(n)
			}
		case "README.md":
			if n.Target == nil {
				t.Fatal(n)
			}
		case "linked":
			if n.Target != nil || n.State != "not-followed" {
				t.Fatal(n)
			}
		case "secret.md":
			t.Fatal("followed external symlink")
		}
	}
	if e = os.RemoveAll(dir); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(external, dir); e != nil {
		t.Fatal(e)
	}
	p, e = Discover(root)
	if e != nil || sessionTab(t, p).State != "unavailable" || p.Inventory.Sessions != "" {
		t.Fatal(p, e)
	}
	if e = os.Remove(dir); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(dir, []byte("not a directory"), 0600); e != nil {
		t.Fatal(e)
	}
	p, e = Discover(root)
	if e != nil || sessionTab(t, p).State != "unavailable" {
		t.Fatal(p, e)
	}
}
