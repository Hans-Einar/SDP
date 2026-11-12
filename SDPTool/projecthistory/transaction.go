package projecthistory

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ErrClosed reports an operation attempted after a transaction released its lock.
var ErrClosed = errors.New("canonical history transaction closed")

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
	if t == nil || t.unlock == nil {
		return ErrClosed
	}
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

// Read returns validated exact ledger bytes under this transaction's held lock.
// It permits no reservation or a reservation belonging to owner.
func (t *Transaction) Read(owner string) (Snapshot, error) {
	if err := t.Check(owner); err != nil {
		return Snapshot{}, err
	}
	data, err := readBytes(t.path)
	if err != nil {
		return Snapshot{}, err
	}
	if _, err = validate(data); err != nil {
		return Snapshot{}, err
	}
	if err = t.Check(owner); err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Bytes: data, Revision: hash(data)}, nil
}

// BatchResult describes one atomic, append-only batch publication.
type BatchResult struct {
	Revision string
	EventIDs []string
	Appended bool
}

// AppendBatch publishes validated envelopes without reacquiring this transaction's
// lock. The caller retains responsibility for payload semantics, durable intent,
// other record effects and releasing any reservation after recovery is complete.
// Exact complete retries are accepted after later appends; partial overlap fails.
func (t *Transaction) AppendBatch(owner, expectedRevision string, events [][]byte) (BatchResult, error) {
	return t.appendBatch(owner, expectedRevision, events, nil)
}

func (t *Transaction) appendBatch(owner, expectedRevision string, events [][]byte, fault func(string) error) (BatchResult, error) {
	var out BatchResult
	if err := t.Check(owner); err != nil {
		return out, err
	}
	if !digest.MatchString(expectedRevision) || len(events) == 0 {
		return out, fmt.Errorf("invalid batch append request")
	}
	ids := make(map[string]bool, len(events))
	var suffix []byte
	for _, raw := range events {
		if len(raw) > MaxBytes {
			return out, fmt.Errorf("history capacity exceeded")
		}
		id, compact, err := envelope(raw)
		if err != nil {
			return out, err
		}
		if ids[id] {
			return out, fmt.Errorf("duplicate batch event %s", id)
		}
		ids[id] = true
		out.EventIDs = append(out.EventIDs, id)
		if len(compact)+1 > MaxBytes-len(suffix) {
			return out, fmt.Errorf("history capacity exceeded")
		}
		suffix = append(suffix, compact...)
		suffix = append(suffix, '\n')
	}
	snapshot, err := t.Read(owner)
	if err != nil {
		return out, err
	}
	out.Revision = snapshot.Revision
	existing, err := validate(snapshot.Bytes)
	if err != nil {
		return out, err
	}
	present := 0
	for _, raw := range events {
		id, compact, _ := envelope(raw)
		if old, ok := existing[id]; ok {
			if !bytes.Equal(old, compact) {
				return out, ErrConflict
			}
			present++
		}
	}
	if present != 0 {
		// A prior successful publication is exactly one contiguous compacted suffix
		// after the originally expected bytes, even if newer events follow it.
		at := bytes.Index(snapshot.Bytes, suffix)
		if present != len(events) || at < 0 || hash(snapshot.Bytes[:at]) != expectedRevision {
			return out, ErrConflict
		}
		return out, nil
	}
	if snapshot.Revision != expectedRevision {
		return out, ErrStale
	}
	if len(suffix) > MaxBytes-len(snapshot.Bytes) {
		return out, fmt.Errorf("history capacity exceeded")
	}
	next := append(append([]byte(nil), snapshot.Bytes...), suffix...)
	result, err := publish(t.path, snapshot.Bytes, next, Result{}, func(point string) error {
		if err := t.Check(owner); err != nil {
			return err
		}
		if fault != nil {
			return fault(point)
		}
		return nil
	})
	if result.Revision != "" {
		out.Revision = result.Revision
	}
	out.Appended = result.Appended
	return out, err
}
