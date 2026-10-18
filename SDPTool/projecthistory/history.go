// Package projecthistory appends domain-validated events to canonical project
// history. It never owns assignment policy or stores a second authoritative log.
package projecthistory

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const MaxBytes = 16 << 20

var ErrBusy = errors.New("canonical history writer busy")
var ErrStale = errors.New("canonical history revision changed")
var ErrConflict = errors.New("history event ID reused with different content")
var ErrUnsupported = errors.New("canonical history mutation is supported only on Linux")
var eventID = regexp.MustCompile(`^EVT-[A-Z0-9]+(?:[.-][A-Z0-9]+)*$`)
var eventType = regexp.MustCompile(`^(?:x-[a-z0-9]+(?:[.-][a-z0-9]+)*:[a-z][a-z0-9]*(?:-[a-z0-9]+)*|(?:work|review|verification)-[a-z0-9]+(?:-[a-z0-9]+)*|release-(?:planned|version-selected|candidate-opened|verification-completed|approved|tag-created|published|yanked|migration-applied|notes-corrected))$`)
var releaseID = regexp.MustCompile(`^REL-(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-((?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(?:\.(?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
var digest = regexp.MustCompile(`^[a-f0-9]{64}$`)

type Snapshot struct {
	Bytes    []byte
	Revision string
}
type Result struct {
	Revision string
	EventID  string
	Appended bool
}

func hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }

// Safe paths and regular files are required. This does not defend against a
// malicious same-user process swapping ancestors between checks and use.
func safePath(path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for {
		st, e := os.Lstat(p)
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if e == nil && st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("history symlink: %s", p)
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
	return nil
}
func readBytes(path string) ([]byte, error) {
	if err := safePath(path); err != nil {
		return nil, err
	}
	st, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() || st.Size() > MaxBytes {
		return nil, fmt.Errorf("history must be a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, MaxBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > MaxBytes {
		return nil, fmt.Errorf("history exceeds %d bytes", MaxBytes)
	}
	return data, nil
}

// envelope checks exact generic envelope names and duplicate top-level keys.
// Payload semantics remain the domain's responsibility. Extension fields survive.
func envelope(raw []byte) (string, []byte, error) {
	if !utf8.Valid(raw) {
		return "", nil, fmt.Errorf("history event requires valid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return "", nil, fmt.Errorf("history event requires an object")
	}
	fields := map[string]json.RawMessage{}
	for d.More() {
		token, err = d.Token()
		if err != nil {
			return "", nil, err
		}
		name, ok := token.(string)
		if !ok {
			return "", nil, fmt.Errorf("invalid event key")
		}
		known := map[string]string{"schemaversion": "schemaVersion", "eventid": "eventId", "eventtype": "eventType", "occurredat": "occurredAt", "actor": "actor", "commit": "commit", "payload": "payload", "subjectid": "subjectId", "releaseid": "releaseId"}
		if canonical, exists := known[strings.ToLower(name)]; exists && name != canonical {
			return "", nil, fmt.Errorf("noncanonical envelope key %s", name)
		}
		if _, ok = fields[name]; ok {
			return "", nil, fmt.Errorf("duplicate event key %s", name)
		}
		var value json.RawMessage
		if err = d.Decode(&value); err != nil {
			return "", nil, err
		}
		fields[name] = value
	}
	if _, err = d.Token(); err != nil {
		return "", nil, err
	}
	if _, err = d.Token(); err != io.EOF {
		return "", nil, fmt.Errorf("trailing event input")
	}
	field := func(key string) string { var s string; _ = json.Unmarshal(fields[key], &s); return s }
	if field("schemaVersion") != "1.0" || !eventID.MatchString(field("eventId")) || !eventType.MatchString(field("eventType")) || field("actor") == "" {
		return "", nil, fmt.Errorf("invalid generic history envelope")
	}
	if _, err = time.Parse(time.RFC3339Nano, field("occurredAt")); err != nil {
		return "", nil, fmt.Errorf("invalid event time")
	}
	if v, ok := fields["commit"]; !ok || (string(v) != "null" && field("commit") == "") {
		return "", nil, fmt.Errorf("invalid commit reference")
	}
	if _, ok := fields["subjectId"]; ok && field("subjectId") == "" {
		return "", nil, fmt.Errorf("invalid subject reference")
	}
	if _, ok := fields["releaseId"]; ok && !releaseID.MatchString(field("releaseId")) {
		return "", nil, fmt.Errorf("invalid release reference")
	}
	var payload map[string]json.RawMessage
	if err = json.Unmarshal(fields["payload"], &payload); err != nil || payload == nil {
		return "", nil, fmt.Errorf("event payload requires object")
	}
	var compact bytes.Buffer
	if err = json.Compact(&compact, raw); err != nil {
		return "", nil, err
	}
	return field("eventId"), compact.Bytes(), nil
}
func validate(data []byte) (map[string][]byte, error) {
	if len(data) > 0 && data[len(data)-1] != '\n' {
		return nil, fmt.Errorf("incomplete history tail; explicit reconciliation required")
	}
	events := map[string][]byte{}
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		id, compact, err := envelope(line)
		if err != nil {
			return nil, err
		}
		if _, ok := events[id]; ok {
			return nil, fmt.Errorf("duplicate history event %s", id)
		}
		events[id] = compact
	}
	return events, nil
}

// Validate checks a captured generic history stream without filesystem I/O.
func Validate(data []byte) error {
	if len(data) > MaxBytes {
		return fmt.Errorf("history capacity exceeded")
	}
	_, err := validate(data)
	return err
}

func Read(path string) (Snapshot, error) {
	data, err := readBytes(path)
	if err != nil {
		return Snapshot{}, err
	}
	if _, err = validate(data); err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Bytes: data, Revision: hash(data)}, nil
}

// Append preserves every previous byte. An exact event retry is idempotent.
// A post-rename error can mean committed-but-unacknowledged; retry identical input.
func Append(path, expectedRevision string, event []byte) (Result, error) {
	return appendEvent(path, expectedRevision, event, nil)
}
func appendEvent(path, expectedRevision string, event []byte, fault func(string) error) (Result, error) {
	var result Result
	if !digest.MatchString(expectedRevision) || len(event) > MaxBytes {
		return result, fmt.Errorf("invalid append request")
	}
	id, compact, err := envelope(event)
	if err != nil {
		return result, err
	}
	result.EventID = id
	path, err = filepath.Abs(path)
	if err != nil {
		return result, err
	}
	if err = safePath(path); err != nil {
		return result, err
	}
	unlock, err := lock(path)
	if err != nil {
		return result, err
	}
	defer unlock()
	original, err := readBytes(path)
	if err != nil {
		return result, err
	}
	events, err := validate(original)
	if err != nil {
		return result, err
	}
	result.Revision = hash(original)
	if prior, ok := events[id]; ok {
		if !bytes.Equal(prior, compact) {
			return result, ErrConflict
		}
		return result, nil
	}
	if result.Revision != expectedRevision {
		return result, ErrStale
	}
	next := make([]byte, 0, len(original)+len(compact)+1)
	next = append(next, original...)
	next = append(next, compact...)
	next = append(next, '\n')
	if len(next) > MaxBytes {
		return result, fmt.Errorf("history capacity exceeded")
	}
	dir := filepath.Dir(path)
	stage, err := os.CreateTemp(dir, ".ledger-stage-")
	if err != nil {
		return result, err
	}
	defer os.Remove(stage.Name())
	// Preserve an existing ledger's permissions. New ledgers are private by default.
	if st, e := os.Stat(path); e == nil {
		if err = stage.Chmod(st.Mode().Perm()); err != nil {
			stage.Close()
			return result, err
		}
	}
	if _, err = stage.Write(next); err != nil {
		stage.Close()
		return result, err
	}
	if err = stage.Sync(); err != nil {
		stage.Close()
		return result, err
	}
	if err = stage.Close(); err != nil {
		return result, err
	}
	if fault != nil {
		if err = fault("before-publish"); err != nil {
			return result, err
		}
	}
	current, err := readBytes(path)
	if err != nil {
		return result, err
	}
	if !bytes.Equal(current, original) {
		return result, ErrStale
	}
	if err = os.Rename(stage.Name(), path); err != nil {
		return result, err
	}
	result.Revision = hash(next)
	result.Appended = true
	if fault != nil {
		if err = fault("after-publish"); err != nil {
			return result, err
		}
	}
	f, err := os.Open(dir)
	if err != nil {
		return result, err
	}
	err = f.Sync()
	closeErr := f.Close()
	if err != nil {
		return result, err
	}
	if closeErr != nil {
		return result, closeErr
	}
	if fault != nil {
		if err = fault("after-sync"); err != nil {
			return result, err
		}
	}
	return result, nil
}
