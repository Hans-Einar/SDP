package preparation

import (
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type Mode string

const (
	Prototype Mode = "prototype"
	Connected Mode = "connected"
)

// Bind must validate all selected declarations against loaded modules and typed
// plans, install handlers on the detached session, and never execute an action.
type Binder func(*ui.Session, *parser.Document) error

// Request accepts a validated 0.2 document; adapters are synchronous and must
// treat their document and layout inputs as read-only.
type Request struct {
	Document                         *parser.Document
	Entry, SessionID, SourceRevision string
	Capabilities                     Capabilities
	ValidateLayout                   func(*parser.Instance) error
	Mode                             Mode
	Bind                             Binder
}

// Candidate owns an unmounted session. The caller must Close abandoned candidates.
// Admission is only a final synchronous revision check, not a native publication.
type Candidate struct {
	Session        *ui.Session
	Unbound        int
	sourceRevision string
	modelRevision  uint64
}

func (c *Candidate) Close() {
	if c != nil && c.Session != nil {
		c.Session.Close()
	}
}
func (c *Candidate) Admit(sourceRevision string) error {
	if c == nil || c.Session == nil || c.Session.Closed() {
		return fmt.Errorf("closed: preparation candidate")
	}
	if sourceRevision != c.sourceRevision || c.Session.Revision != c.modelRevision {
		return fmt.Errorf("stale: preparation candidate")
	}
	return nil
}
func Prepare(r Request) (*Candidate, error) {
	if r.Mode != Prototype && r.Mode != Connected {
		return nil, fmt.Errorf("preparation: explicit mode required")
	}
	if r.Document == nil {
		return nil, fmt.Errorf("preparation: document required")
	}
	if r.SourceRevision == "" {
		return nil, fmt.Errorf("preparation: source revision required")
	}
	// Normalize internally: callers cannot supply an unrelated normalized root.
	roots, err := parser.Normalize(r.Document)
	if err != nil {
		return nil, err
	}
	root := roots[r.Entry]
	if root == nil || root.Kind != "frame" {
		return nil, fmt.Errorf("preparation: entry %q must be a frame", r.Entry)
	}
	if err = Check(r.Document.Profile, root, r.Capabilities); err != nil {
		return nil, err
	}
	if r.ValidateLayout == nil {
		return nil, fmt.Errorf("preparation: layout validator required")
	}
	if err = r.ValidateLayout(root); err != nil {
		return nil, err
	}
	s, err := ui.New(r.SessionID, root)
	if err != nil {
		return nil, err
	}
	c := &Candidate{Session: s, sourceRevision: r.SourceRevision, modelRevision: s.Revision}
	success := false
	defer func() {
		if !success {
			c.Close()
		}
	}()
	selected := *r.Document
	selected.Connections = nil
	for _, con := range r.Document.Connections {
		if con.Definition == r.Entry {
			selected.Connections = append(selected.Connections, con)
		}
	}
	c.Unbound = len(selected.Connections)
	for _, w := range s.Widgets() {
		if w.Binding.Module != "" {
			c.Unbound++
		}
	}
	if r.Mode == Connected {
		if r.Bind == nil {
			return nil, fmt.Errorf("preparation: connected binding adapter required")
		}
		if err = r.Bind(s, &selected); err != nil {
			return nil, err
		}
		if err = c.Admit(r.SourceRevision); err != nil {
			return nil, err
		}
		symbols := map[string]bool{}
		for _, w := range s.Widgets() {
			if w.Binding.Module != "" {
				if !s.HasBinding(w.Handle) {
					return nil, fmt.Errorf("unbound: %s", w.Handle.Path)
				}
				symbols[w.Binding.Module+"."+w.Binding.Object] = true
			}
		}
		for _, con := range selected.Connections {
			target := con.Definition + "/" + strings.Join(con.Path, "/")
			w, ok := s.Widget(target)
			if !ok || w.Handle.Kind != "input" || !symbols[con.Module+"."+con.Object] {
				return nil, fmt.Errorf("unbound-handle: %s", target)
			}
		}
		c.Unbound = 0
	}
	success = true
	return c, nil
}
