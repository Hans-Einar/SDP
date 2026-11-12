package projecthistory

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func batchTransaction(t *testing.T) (string, *Transaction, Snapshot) {
	t.Helper()
	requireLinux(t)
	path := filepath.Join(t.TempDir(), "Ledger.ndjson")
	original := append([]byte("\n  "), event(1)...)
	original = append(original, '\n', '\n')
	if err := os.WriteFile(path, original, 0640); err != nil {
		t.Fatal(err)
	}
	tx, err := BeginTransaction(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(tx.Close)
	snapshot, err := tx.Read("")
	if err != nil {
		t.Fatal(err)
	}
	return path, tx, snapshot
}

func TestTransactionBatchReservationAndRetry(t *testing.T) {
	path, tx, before := batchTransaction(t)
	const owner = "activity\noperation\n"
	if err := tx.Reserve(owner); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Read(""); !errors.Is(err, ErrBusy) {
		t.Fatal(err)
	}
	if _, err := tx.AppendBatch("foreign", before.Revision, [][]byte{event(2)}); !errors.Is(err, ErrBusy) {
		t.Fatal(err)
	}
	if _, err := Read(path); !errors.Is(err, ErrBusy) {
		t.Fatal(err)
	}
	if _, err := Append(path, before.Revision, event(2)); !errors.Is(err, ErrBusy) {
		t.Fatal(err)
	}
	batch := [][]byte{event(2), event(3)}
	result, err := tx.AppendBatch(owner, before.Revision, batch)
	if err != nil || !result.Appended || len(result.EventIDs) != 2 {
		t.Fatal(result, err)
	}
	after, err := tx.Read(owner)
	if err != nil || !bytes.HasPrefix(after.Bytes, before.Bytes) || after.Revision != result.Revision {
		t.Fatal(after, err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0640 {
		t.Fatal(info, err)
	}
	if _, err = Read(path); !errors.Is(err, ErrBusy) {
		t.Fatal(err)
	}
	tx.Close()
	tx, err = BeginTransaction(path)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Close()
	newer, err := tx.AppendBatch(owner, after.Revision, [][]byte{event(4)})
	if err != nil || !newer.Appended {
		t.Fatal(newer, err)
	}
	retry, err := tx.AppendBatch(owner, before.Revision, batch)
	if err != nil || retry.Appended || retry.Revision != newer.Revision {
		t.Fatal(retry, err)
	}
	if err = tx.Release(owner); err != nil {
		t.Fatal(err)
	}
	if _, err = Read(path); err != nil {
		t.Fatal(err)
	}
}

func TestTransactionBatchRejectsInvalidPublication(t *testing.T) {
	cases := []struct {
		name   string
		events [][]byte
		want   error
	}{
		{"duplicate", [][]byte{event(2), event(2)}, nil},
		{"invalid envelope", [][]byte{[]byte(`{"eventId":"EVT-X"}`)}, nil},
		{"partial retry", [][]byte{event(1), event(2)}, ErrConflict},
		{"different content", [][]byte{bytes.Replace(event(1), []byte(`"payload":{}`), []byte(`"payload":{"changed":true}`), 1)}, ErrConflict},
		{"empty", nil, nil},
		{"capacity", [][]byte{bytes.Repeat([]byte(" "), MaxBytes+1)}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path, tx, before := batchTransaction(t)
			_, err := tx.AppendBatch("", before.Revision, tc.events)
			if err == nil || tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
			actual, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(actual, before.Bytes) {
				t.Fatal("failed batch changed history", err)
			}
		})
	}
	_, tx, before := batchTransaction(t)
	if _, err := tx.AppendBatch("", hash(nil), [][]byte{event(2)}); !errors.Is(err, ErrStale) {
		t.Fatal(err)
	}
	result, err := tx.AppendBatch("", before.Revision, [][]byte{event(2), event(3)})
	if err != nil {
		t.Fatal(err)
	}
	for _, batch := range [][][]byte{{event(3), event(2)}, {event(1), event(3)}} {
		if _, err = tx.AppendBatch("", before.Revision, batch); !errors.Is(err, ErrConflict) {
			t.Fatal(err)
		}
	}
	if _, err = tx.AppendBatch("", result.Revision, [][]byte{event(2), event(3)}); !errors.Is(err, ErrConflict) {
		t.Fatal(err)
	}
}

func TestTransactionBatchFaultAndInterference(t *testing.T) {
	for _, point := range []string{"before-publish", "after-publish", "after-sync"} {
		t.Run(point, func(t *testing.T) {
			path, tx, before := batchTransaction(t)
			injected := errors.New("interrupted")
			result, err := tx.appendBatch("", before.Revision, [][]byte{event(2), event(3)}, func(p string) error {
				if p == point {
					return injected
				}
				return nil
			})
			if !errors.Is(err, injected) || result.Appended != (point != "before-publish") {
				t.Fatal(result, err)
			}
			actual, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if point == "before-publish" && !bytes.Equal(actual, before.Bytes) {
				t.Fatal("published before boundary")
			}
			retry, err := tx.AppendBatch("", before.Revision, [][]byte{event(2), event(3)})
			if err != nil || retry.Appended != (point == "before-publish") {
				t.Fatal(retry, err)
			}
		})
	}
	path, tx, before := batchTransaction(t)
	external := append(append([]byte(nil), before.Bytes...), '\n')
	_, err := tx.appendBatch("", before.Revision, [][]byte{event(2)}, func(point string) error {
		if point == "before-publish" {
			return os.WriteFile(path, external, 0600)
		}
		return nil
	})
	if !errors.Is(err, ErrStale) {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(actual, external) {
		t.Fatal("external edit lost", err)
	}
}

func TestTransactionReadRejectsMalformedHistoryAndClosedMethods(t *testing.T) {
	path, tx, before := batchTransaction(t)
	if err := os.WriteFile(path, []byte("incomplete"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Read(""); err == nil {
		t.Fatal("malformed history accepted")
	}
	tx.Close()
	for name, call := range map[string]func() error{
		"check":   func() error { return tx.Check("") },
		"reserve": func() error { return tx.Reserve("owner") },
		"release": func() error { return tx.Release("owner") },
		"read":    func() error { _, err := tx.Read(""); return err },
		"append":  func() error { _, err := tx.AppendBatch("", before.Revision, [][]byte{event(2)}); return err },
	} {
		if err := call(); !errors.Is(err, ErrClosed) {
			t.Fatalf("%s: %v", name, err)
		}
	}
	emptyPath := filepath.Join(t.TempDir(), "Ledger.ndjson")
	fresh, err := BeginTransaction(emptyPath)
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	snapshot, err := fresh.Read("")
	if err != nil || len(snapshot.Bytes) != 0 || snapshot.Revision != hash(nil) {
		t.Fatal(snapshot, err)
	}
	if _, err = fresh.AppendBatch("", snapshot.Revision, [][]byte{event(1)}); err != nil {
		t.Fatal(err)
	}
}
