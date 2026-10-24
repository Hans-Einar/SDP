package projecthistory

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Transaction shares Append's canonical ledger lock. A durable reservation blocks
// cooperating appenders between interrupted multi-record effects. It contains an
// opaque owner identity, not executable input. Callers must preserve file preimages.
type Transaction struct {
	path   string
	unlock func()
}

func BeginTransaction(path string) (*Transaction, error) {
	absolute, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	if e = safePath(absolute); e != nil {
		return nil, e
	}
	unlock, e := lock(absolute)
	if e != nil {
		return nil, e
	}
	return &Transaction{absolute, unlock}, nil
}
func (t *Transaction) Close() {
	if t.unlock != nil {
		t.unlock()
		t.unlock = nil
	}
}
func pendingOwner(path string) (string, error) {
	marker := path + ".pg-pending"
	if e := safePath(marker); e != nil {
		return "", e
	}
	info, e := os.Lstat(marker)
	if os.IsNotExist(e) {
		return "", nil
	}
	if e != nil {
		return "", e
	}
	if !info.Mode().IsRegular() || info.Size() > 8192 {
		return "", ErrBusy
	}
	f, e := os.Open(marker)
	if os.IsNotExist(e) {
		return "", nil
	}
	if e != nil {
		return "", e
	}
	defer f.Close()
	info, e = f.Stat()
	if e != nil {
		return "", e
	}
	if !info.Mode().IsRegular() {
		return "", ErrBusy
	}
	b, e := io.ReadAll(io.LimitReader(f, 8193))
	if os.IsNotExist(e) {
		return "", nil
	}
	if e != nil {
		return "", e
	}
	if len(b) == 0 || len(b) > 8192 {
		return "", fmt.Errorf("invalid pending canonical transaction")
	}
	return string(b), nil
}
func (t *Transaction) Check(owner string) error {
	pending, e := pendingOwner(t.path)
	if e != nil {
		return e
	}
	if pending != "" && pending != owner {
		return ErrBusy
	}
	return nil
}
func syncTransactionDir(path string) error {
	f, e := os.Open(filepath.Dir(path))
	if e != nil {
		return e
	}
	defer f.Close()
	return f.Sync()
}
func (t *Transaction) Reserve(owner string) error {
	if strings.TrimSpace(owner) == "" || len(owner) > 8192 {
		return fmt.Errorf("invalid transaction owner")
	}
	if e := t.Check(owner); e != nil {
		return e
	}
	marker := t.path + ".pg-pending"
	old, e := pendingOwner(t.path)
	if e != nil {
		return e
	}
	if old == owner {
		return nil
	}
	f, e := os.CreateTemp(filepath.Dir(marker), ".pg-claim-")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	_, e = f.WriteString(owner)
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
	if e = os.Link(f.Name(), marker); e != nil {
		return e
	}
	return syncTransactionDir(marker)
}
func (t *Transaction) Release(owner string) error {
	if e := t.Check(owner); e != nil {
		return e
	}
	marker := t.path + ".pg-pending"
	if e := os.Remove(marker); e != nil && !os.IsNotExist(e) {
		return e
	}
	return syncTransactionDir(marker)
}
