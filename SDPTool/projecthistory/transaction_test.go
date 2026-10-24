package projecthistory

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestTransactionReservationSurvivesOwnerExit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Ledger.ndjson")
	if e := os.WriteFile(path, nil, 0600); e != nil {
		t.Fatal(e)
	}
	tx, e := BeginTransaction(path)
	if e != nil {
		t.Fatal(e)
	}
	if e = tx.Reserve("store\noperation\n"); e != nil {
		t.Fatal(e)
	}
	if _, e = BeginTransaction(path); !errors.Is(e, ErrBusy) {
		t.Fatal("concurrent writer acquired history", e)
	}
	tx.Close()
	if _, e = Read(path); !errors.Is(e, ErrBusy) {
		t.Fatal("reader missed interrupted transaction", e)
	}
	next, e := BeginTransaction(path)
	if e != nil {
		t.Fatal(e)
	}
	defer next.Close()
	if e = next.Check("different"); !errors.Is(e, ErrBusy) {
		t.Fatal("foreign owner accepted", e)
	}
	if e = next.Release("different"); !errors.Is(e, ErrBusy) {
		t.Fatal("foreign reservation released", e)
	}
	if e = next.Reserve("store\noperation\n"); e != nil {
		t.Fatal(e)
	}
	if e = next.Release("store\noperation\n"); e != nil {
		t.Fatal(e)
	}
	if _, e = Read(path); e != nil {
		t.Fatal(e)
	}
}
