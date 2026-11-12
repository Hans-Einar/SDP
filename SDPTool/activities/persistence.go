package activities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"os"
	"path/filepath"
)

type image struct {
	Path   string `json:"path"`
	Before []byte `json:"before"`
	After  []byte `json:"after"`
}
type intent struct {
	Pending     bool     `json:"pending,omitempty"`
	Schema      string   `json:"schema"`
	ProjectID   string   `json:"projectId"`
	OperationID string   `json:"operationId"`
	Digest      string   `json:"digest"`
	Expected    string   `json:"expected"`
	Files       []image  `json:"files"`
	Events      [][]byte `json:"events"`
	Result      Result   `json:"result"`
}

func (s *Service) intentPath(id string) string { return "Activities/.operations/" + id + ".json" }
func (s *Service) load(id string) (intent, error) {
	var in intent
	if !uuidRE.MatchString(id) {
		return in, fmt.Errorf("invalid operation UUID")
	}
	b, e := s.read(s.intentPath(id))
	if e != nil {
		return in, e
	}
	e = Decode(b, &in)
	if e == nil && (in.Schema != "sdp-activity-operation/0.2" || in.ProjectID != s.ProjectID || in.OperationID != id || !hashRE.MatchString(in.Digest) || !hashRE.MatchString(in.Expected) || len(in.Events) == 0) {
		e = fmt.Errorf("invalid operation intent")
	}
	return in, e
}
func (s *Service) write(rel string, data []byte) error {
	p, e := s.path(rel)
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		return e
	}
	if _, e = s.path(rel); e != nil {
		return e
	}
	f, e := os.CreateTemp(filepath.Dir(p), ".activity-stage-")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	if _, e = f.Write(data); e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e != nil {
		return e
	}
	if ce != nil {
		return ce
	}
	if e = os.Rename(f.Name(), p); e != nil {
		return e
	}
	d, e := os.Open(filepath.Dir(p))
	if e != nil {
		return e
	}
	defer d.Close()
	if e = d.Sync(); e != nil {
		return e
	}
	root, _ := filepath.Abs(s.Area)
	for q := filepath.Dir(filepath.Dir(p)); ; q = filepath.Dir(q) {
		if len(q) < len(root) {
			break
		}
		ancestor, e := os.Open(q)
		if e != nil {
			return e
		}
		e = ancestor.Sync()
		ce := ancestor.Close()
		if e != nil {
			return e
		}
		if ce != nil {
			return ce
		}
		if q == root || filepath.Dir(q) == q {
			break
		}
	}
	return nil
}
func (s *Service) point(p string) error {
	if s.fault != nil {
		return s.fault(p)
	}
	return nil
}
func (s *Service) apply(tx *projecthistory.Transaction, in intent) (Result, error) {
	r := in.Result
	owner := "activity:" + in.OperationID
	snap, e := tx.Read(owner)
	if e != nil {
		return r, e
	}
	found := 0
	for _, line := range bytes.Split(snap.Bytes, []byte{'\n'}) {
		var head struct {
			EventID string `json:"eventId"`
		}
		_ = json.Unmarshal(line, &head)
		for _, raw := range in.Events {
			var want Event
			_ = json.Unmarshal(raw, &want)
			if head.EventID == want.EventID {
				found++
			}
		}
	}
	if found > 0 {
		if _, e = tx.AppendBatch(owner, in.Expected, in.Events); e != nil {
			return r, e
		}
		if in.Pending {
			view, ve := s.ReadAdmission(tx, owner)
			if ve != nil {
				return r, ve
			}
			if view.Pending != nil && view.Pending.OperationID == in.OperationID {
				if e = s.ClearPendingLifecycle(tx, owner, in.OperationID); e != nil {
					return r, e
				}
			}
		}
		if e = tx.Release(owner); e != nil {
			return r, e
		}
		return r, nil
	}
	if snap.Revision != in.Expected {
		return r, projecthistory.ErrStale
	}
	if e := tx.Reserve(owner); e != nil {
		return r, e
	}
	if e := s.point("after-reserve"); e != nil {
		return r, e
	}
	for i, f := range in.Files {
		b, e := s.read(f.Path)
		absent := os.IsNotExist(e)
		if e != nil && !absent {
			return r, e
		}
		if bytes.Equal(b, f.After) && !absent {
			continue
		}
		if !(f.Before == nil && absent) && !(f.Before != nil && !absent && bytes.Equal(b, f.Before)) {
			return r, fmt.Errorf("file preimage conflict: %s", f.Path)
		}
		if e = s.write(f.Path, f.After); e != nil {
			return r, e
		}
		if e = s.point(fmt.Sprintf("after-file-%d", i)); e != nil {
			return r, e
		}
	}
	next := append([]byte{}, snap.Bytes...)
	for _, ev := range in.Events {
		next = append(next, ev...)
		next = append(next, '\n')
	}
	if _, e = Replay(s.Area, s.ProjectID, projecthistory.Snapshot{Bytes: next, Revision: digest(next)}); e != nil {
		return r, e
	}
	if e = s.point("before-history"); e != nil {
		return r, e
	}
	_, e = tx.AppendBatch(owner, in.Expected, in.Events)
	if e != nil {
		return r, e
	}
	if e = s.point("after-history"); e != nil {
		return r, e
	}
	if in.Pending {
		if e = s.ClearPendingLifecycle(tx, owner, in.OperationID); e != nil {
			return r, e
		}
	}
	if e = tx.Release(owner); e != nil {
		return r, e
	}
	return r, s.point("after-release")
}

// Recover requires trusted host authority just like creation; operation files are not a separate ledger.
func (s *Service) Recover(operationID string, caller Caller) (Result, error) {
	if e := s.validateCaller(caller); e != nil {
		return Result{}, e
	}
	p, e := s.path("ProjectManagement/Ledger.ndjson")
	if e != nil {
		return Result{}, e
	}
	tx, e := projecthistory.BeginTransaction(p)
	if e != nil {
		return Result{}, e
	}
	defer tx.Close()
	in, e := s.load(operationID)
	if os.IsNotExist(e) {
		view, ve := s.ReadAdmission(tx, "activity:"+operationID)
		if ve != nil {
			return Result{}, ve
		}
		if view.Pending != nil && view.Pending.OperationID == operationID {
			return Result{OperationID: operationID, ActivityID: view.Pending.ActivityID, Pending: true, HistoryRevision: view.Catalog.HistoryRevision}, nil
		}
	}
	if e != nil {
		return Result{}, e
	}
	r, e := s.apply(tx, in)
	r.Recovered = true
	return r, e
}
