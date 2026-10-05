package model

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type transaction struct {
	ID        string `yaml:"id"`
	Target    string `yaml:"target"`
	Expected  string `yaml:"expected"`
	NewDigest string `yaml:"newDigest"`
	State     string `yaml:"state"`
}

var faultHook func(string)

func fault(s string) {
	if faultHook != nil {
		faultHook(s)
	}
}
func writeY(path string, v any) error {
	if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		return e
	}
	b, e := yaml.Marshal(v)
	if e != nil {
		return e
	}
	tmp := path + ".tmp"
	f, e := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if e != nil {
		return e
	}
	_, e = f.Write(b)
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e != nil {
		return e
	}
	if ce != nil {
		return ce
	}
	return os.Rename(tmp, path)
}

// Source limits are independent of retained history. Hash history as streams.
func treeDigest(p string) (string, error) {
	st, e := os.Lstat(p)
	if e != nil {
		return "", e
	}
	if !st.IsDir() {
		return "", fail("path", "not a real tree")
	}
	entries := map[string]string{}
	seen := map[string]bool{}
	var total int64
	e = filepath.WalkDir(p, func(name string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == p {
			return nil
		}
		rel, _ := filepath.Rel(p, name)
		rel = filepath.ToSlash(rel)
		if !safePath(rel) || d.Type()&os.ModeSymlink != 0 {
			return fail("path", "unsafe tree entry")
		}
		key := strings.ToLower(rel)
		if seen[key] {
			return fail("path", "case collision")
		}
		seen[key] = true
		if len(seen) > 100000 {
			return fail("limit", "history tree exceeds 100000 entries")
		}
		if d.IsDir() {
			return nil
		}
		st, err := d.Info()
		if err != nil {
			return err
		}
		if !st.Mode().IsRegular() {
			return fail("path", "nonregular tree entry")
		}
		total += st.Size()
		if total > 1<<30 {
			return fail("limit", "retained artifact exceeds 1 GiB")
		}
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		h := sha256.New()
		_, err = io.Copy(h, io.LimitReader(f, (1<<30)+1))
		f.Close()
		if err != nil {
			return err
		}
		entries[rel] = hex.EncodeToString(h.Sum(nil))
		return nil
	})
	if e != nil {
		return "", e
	}
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	h.Write([]byte("SDP-model-sources/1\x00"))
	for _, k := range keys {
		var n [8]byte
		binary.BigEndian.PutUint64(n[:], uint64(len(k)))
		h.Write(n[:])
		h.Write([]byte(k))
		v, _ := hex.DecodeString(entries[k])
		h.Write(v)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
func journals(area string) ([]transaction, error) {
	entries, e := os.ReadDir(filepath.Join(area, operations))
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	out := []transaction{}
	for _, d := range entries {
		if !d.IsDir() {
			continue
		}
		if !uuidRE.MatchString(d.Name()) {
			return nil, fail("recovery", "unknown operation directory")
		}
		var t transaction
		b, e := readMetadata(filepath.Join(area, operations, d.Name(), "journal.yaml"))
		if os.IsNotExist(e) {
			return nil, fail("recovery", "incomplete staging operation "+d.Name())
		}
		if e != nil {
			return nil, e
		}
		if e = strict(b, &t); e != nil {
			return nil, e
		}
		if t.ID != d.Name() || !safePath(t.Target) || strings.Contains(t.Target, "/") || (!strings.HasPrefix(t.Target, "WORK--") && !strings.HasPrefix(t.Target, "CANDIDATE--") && !strings.HasPrefix(t.Target, "PROPOSAL--") && !strings.HasPrefix(t.Target, "RELEASE--")) {
			return nil, fail("recovery", "invalid journal target")
		}
		switch t.State {
		case "building", "prepared", "done", "aborted":
		default:
			return nil, fail("recovery", "unknown transaction state")
		}
		if (t.Expected != "" && !hashRE.MatchString(t.Expected)) || (t.State == "prepared" && !hashRE.MatchString(t.NewDigest)) {
			return nil, fail("recovery", "invalid transaction hash")
		}
		if t.State != "done" && t.State != "aborted" {
			out = append(out, t)
		}
	}
	return out, nil
}
func begin(area string) (func(), error) {
	abs, e := filepath.Abs(area)
	if e != nil {
		return nil, e
	}
	real, e := filepath.EvalSymlinks(abs)
	if e != nil {
		return nil, e
	}
	if real != abs {
		return nil, fail("path", "model area must not use symlinks")
	}
	unlock, e := lockArea(area)
	if e != nil {
		return nil, e
	}
	ts, e := journals(area)
	if e != nil {
		unlock()
		return nil, e
	}
	if len(ts) > 0 {
		unlock()
		return nil, fail("recovery", "run model recover "+ts[0].ID+" resume|abort")
	}
	return unlock, nil
}
func publish(area, target, source, expected string, build func(string) error) (string, error) {
	id := uuid()
	op := filepath.Join(area, operations, id)
	stage := filepath.Join(op, "stage", target)
	if e := os.MkdirAll(stage, 0700); e != nil {
		return "", e
	}
	t := transaction{ID: id, Target: target, Expected: expected, State: "building"}
	jp := filepath.Join(op, "journal.yaml")
	if e := writeY(jp, t); e != nil {
		return "", e
	}
	abandon := func(e error) (string, error) { t.State = "aborted"; _ = writeY(jp, t); return "", e }
	fault("building")
	if source != "" {
		if e := os.CopyFS(stage, os.DirFS(source)); e != nil {
			return abandon(e)
		}
	}
	if e := build(stage); e != nil {
		return abandon(e)
	}
	d, e := treeDigest(stage)
	if e != nil {
		return abandon(e)
	}
	t.NewDigest = d
	t.State = "prepared"
	if e = writeY(jp, t); e != nil {
		return "", e
	}
	fault("prepared")
	if e = finish(area, t); e != nil {
		return "", e
	}
	return filepath.Join(area, target), nil
}
func finish(area string, t transaction) error {
	op := filepath.Join(area, operations, t.ID)
	jp := filepath.Join(op, "journal.yaml")
	target := filepath.Join(area, t.Target)
	stage := filepath.Join(op, "stage", t.Target)
	backup := filepath.Join(op, "backup", t.Target)
	if t.NewDigest == "" {
		return fail("recovery", "incomplete build; abort required")
	}
	if _, e := os.Lstat(target); e == nil {
		d, e := treeDigest(target)
		if e != nil {
			return e
		}
		if d == t.NewDigest {
			t.State = "done"
			return writeY(jp, t)
		}
		if t.Expected == "" || d != t.Expected {
			return fail("stale", "target changed; recovery preserves all copies")
		}
		if _, e = os.Lstat(backup); !os.IsNotExist(e) {
			return fail("recovery", "backup already exists")
		}
		sd, e := treeDigest(stage)
		if e != nil {
			return e
		}
		if sd != t.NewDigest {
			return fail("integrity", "staged content corrupt")
		}
		if e = os.MkdirAll(filepath.Dir(backup), 0700); e != nil {
			return e
		}
		if e = os.Rename(target, backup); e != nil {
			return e
		}
		fault("backup")
	} else if !os.IsNotExist(e) {
		return e
	} else if t.Expected != "" {
		d, e := treeDigest(backup)
		if e != nil || d != t.Expected {
			return fail("integrity", "missing or changed recovery backup")
		}
	}
	sd, e := treeDigest(stage)
	if e != nil {
		return e
	}
	if sd != t.NewDigest {
		return fail("integrity", "staged content corrupt")
	}
	if e = os.Rename(stage, target); e != nil {
		return e
	}
	fault("installed")
	t.State = "done"
	return writeY(jp, t)
}
func Recover(area, id, action string) (Result, error) {
	if !uuidRE.MatchString(id) {
		return Result{}, fail("arguments", "invalid recovery ID")
	}
	unlock, e := lockArea(area)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	ts, e := journals(area)
	if e != nil {
		return Result{}, e
	}
	for _, t := range ts {
		if t.ID != id {
			continue
		}
		if action == "resume" {
			e = finish(area, t)
		} else if action == "abort" {
			target := filepath.Join(area, t.Target)
			backup := filepath.Join(area, operations, id, "backup", t.Target)
			if _, err := os.Lstat(target); os.IsNotExist(err) {
				if t.Expected != "" {
					d, err := treeDigest(backup)
					if err != nil || d != t.Expected {
						return Result{}, fail("integrity", "backup cannot restore")
					}
					e = os.Rename(backup, target)
				}
			} else if err != nil {
				e = err
			} else {
				d, err := treeDigest(target)
				if err != nil {
					return Result{}, err
				}
				if d != t.Expected {
					return Result{}, fail("recovery", "target already replaced/edited; resume or preserve manually")
				}
			}
			if e == nil {
				t.State = "aborted"
				e = writeY(filepath.Join(area, operations, id, "journal.yaml"), t)
			}
		} else {
			return Result{}, fail("arguments", "use resume or abort")
		}
		return Result{Schema: Schema, Operation: "recover", Status: action, Recovery: id}, e
	}
	return Result{}, fail("reference", fmt.Sprintf("pending operation %s not found", id))
}
