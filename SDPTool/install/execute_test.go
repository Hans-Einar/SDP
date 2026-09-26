package install

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestApplyRepeatAndKnownUpgrade(t *testing.T) {
	root := t.TempDir()
	p := preview(t, root, fixture())
	r, e := (Executor{}).Apply(context.Background(), root, p)
	if e != nil {
		t.Fatal(r, e)
	}
	if r.Status != "completed" {
		t.Fatal(r)
	}
	b, e := os.ReadFile(filepath.Join(root, ReceiptPath))
	if e != nil {
		t.Fatal(e)
	}
	var receipt Receipt
	if e = Decode(b, MetadataLimit, &receipt); e != nil {
		t.Fatal(e)
	}
	if e = ValidateReceipt(receipt); e != nil {
		t.Fatal(e)
	}
	q, e := Preview(Options{root, "upgrade", p.Release.Path, "", "", true})
	if e != nil || !q.NoChange {
		t.Fatal(q.Actions, e)
	}
	before, e := Inspect(root, nil)
	if e != nil {
		t.Fatal(e)
	}
	r, e = (Executor{}).Apply(context.Background(), root, q)
	if e != nil || r.Status != "no-change" {
		t.Fatal(r, e)
	}
	after, _ := Inspect(root, nil)
	if !SameSnapshot(before.Snapshot, after.Snapshot) {
		t.Fatal("repeat changed bytes")
	}
	d := fixture()
	d.Release = "dev-next"
	d.UpgradesFrom = []string{p.Release.SHA256}
	d.Files[0].Content = []byte("new managed instructions")
	h := Hash(d.Files[0].Content)
	d.Files[0].SHA256 = &h
	q, e = Preview(Options{root, "upgrade", artifact(t, d), p.Release.Path, "", true})
	if e != nil {
		t.Fatal(e)
	}
	r, e = (Executor{}).Apply(context.Background(), root, q)
	if e != nil {
		t.Fatal(r, e)
	}
	pending, e := Pending(root)
	if e != nil || len(pending) != 0 {
		t.Fatal(pending, e)
	}
}
func TestApplyBindsInputsRootAndDrift(t *testing.T) {
	root := t.TempDir()
	p := preview(t, root, fixture())
	put(t, root, "unexpected.md", []byte("edit"))
	if _, e := (Executor{}).Apply(context.Background(), root, p); e == nil {
		t.Fatal("drift accepted")
	}
	if _, e := os.Stat(filepath.Join(root, Operations)); !os.IsNotExist(e) {
		t.Fatal("journal created on drift")
	}
	os.Remove(filepath.Join(root, "unexpected.md"))
	put(t, filepath.Dir(p.Release.Path), filepath.Base(p.Release.Path), []byte("{}"))
	if _, e := (Executor{}).Apply(context.Background(), root, p); e == nil {
		t.Fatal("input drift accepted")
	}
	if _, e := (Executor{}).Apply(context.Background(), t.TempDir(), p); e == nil {
		t.Fatal("wrong root")
	}
}
func TestLockAndLegacyPending(t *testing.T) {
	root := t.TempDir()
	p := preview(t, root, fixture())
	unlock, e := lock(root)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = (Executor{}).Apply(context.Background(), root, p); e == nil {
		t.Fatal("lock ignored")
	}
	unlock()
	put(t, root, Operations+"/install-0123456789abcdef01234567/journal.json", []byte(`{"schemaVersion":"2.0","operationId":"install-0123456789abcdef01234567","status":"failed"}`))
	if _, e = Preview(Options{root, "install", p.Release.Path, "", "", true}); e == nil {
		t.Fatal("legacy pending ignored")
	}
}

func TestUnstartedPublicationFailureIsNotResumable(t *testing.T) {
	root := t.TempDir()
	p := preview(t, root, fixture())
	x := Executor{Fault: func(name string, index int) error {
		if name == "publication" {
			return os.ErrPermission
		}
		return nil
	}}
	r, e := x.Apply(context.Background(), root, p)
	if e == nil || r.OperationID != "" {
		t.Fatal(r, e)
	}
	failure, ok := e.(*Error)
	if !ok || failure.Exit != 4 {
		t.Fatal(e)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatal("unstarted failure changed project")
	}
	if _, e = (Executor{}).Apply(context.Background(), root, p); e != nil {
		t.Fatal("retry failed", e)
	}
}
