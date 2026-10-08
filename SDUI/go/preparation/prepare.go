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

// Request accepts a bounded profile document; adapters are synchronous and must
// treat their document and layout inputs as read-only.
type Request struct {
	Document                         *parser.Document
	Previous                         *ui.Session
	Providers                        map[string]ui.CollectionProvider
	Choices                          map[string][]ui.ChoiceOption
	ValidateState                    ui.StateGate
	ValidatePresentation             ui.PresentationGate
	PreparePresentation              ui.PresentationPrepare
	Entry, SessionID, SourceRevision string
	Capabilities                     Capabilities
	ValidateLayout                   func(*parser.Instance) error
	Mode                             Mode
	Bind                             Binder
}

// Candidate owns an unmounted session. The caller must Close abandoned candidates.
// Admission is only a final synchronous revision check, not a native publication.
type Candidate struct {
	Session                      *ui.Session
	Unbound                      int
	sourceRevision               string
	modelRevision                uint64
	stateRevision                uint64
	previous                     *ui.Session
	previousState, previousModel uint64
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
	if c.previous != nil && (c.previous.Closed() || c.previous.StateRevision != c.previousState || c.previous.Revision != c.previousModel) {
		return fmt.Errorf("stale: published state changed during preparation")
	}
	if sourceRevision != c.sourceRevision || c.Session.Revision != c.modelRevision || c.Session.StateRevision != c.stateRevision {
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
	if r.ValidateState != nil && r.ValidatePresentation != nil {
		return nil, fmt.Errorf("preparation: choose one presentation gate")
	}
	if r.ValidateLayout == nil && r.ValidateState == nil && r.ValidatePresentation == nil {
		return nil, fmt.Errorf("preparation: layout validator required")
	}
	needsChoices := len(r.Choices) > 0
	root.Walk(func(n *parser.Instance) {
		if n.Widget == "select" {
			needsChoices = true
		}
	})
	var s *ui.Session
	var oldState, oldModel uint64
	if r.Previous != nil {
		oldState, oldModel = r.Previous.StateRevision, r.Previous.Revision
		if needsChoices {
			s, err = r.Previous.SuccessorWithChoices(root, r.Providers, r.Choices)
		} else {
			s, err = r.Previous.Successor(root, r.Providers)
		}
	} else {
		s, err = ui.New(r.SessionID, root)
		if err == nil {
			err = s.BindProviders(r.Providers)
			if err == nil && needsChoices {
				err = s.BindChoices(r.Choices)
			}
		}
	}
	if err != nil {
		if s != nil {
			s.Close()
		}
		return nil, err
	}
	c := &Candidate{Session: s, sourceRevision: r.SourceRevision, modelRevision: s.Revision, previous: r.Previous, previousState: oldState, previousModel: oldModel}
	success := false
	defer func() {
		if !success {
			c.Close()
		}
	}()
	// Options must be supplied even for hidden/closed controls. Positive generation
	// distinguishes an admitted empty inventory from an absent provider.
	for path, field := range s.Snapshot().Fields {
		if field.Target.Handle.Kind == "select" && field.Target.OptionGeneration == 0 {
			return nil, fmt.Errorf("choice-options: missing %s", path)
		}
	}
	if r.ValidateLayout != nil {
		if err = r.ValidateLayout(s.SnapshotRoot()); err != nil {
			return nil, err
		}
	}
	if r.ValidateState != nil {
		if err = s.CheckStateWith(r.ValidateState); err != nil {
			return nil, err
		}
	}
	if r.ValidatePresentation != nil {
		if err = s.CheckPresentationWith(r.ValidatePresentation); err != nil {
			return nil, err
		}
	}
	for path, p := range r.Providers {
		if p.Load != nil {
			found := false
			for _, cap := range r.Capabilities {
				if cap == (Capability{Provider, "collection-load", 1}) {
					found = true
				}
			}
			if !found {
				return nil, fmt.Errorf("capability: provider collection-load/1 at %s", path)
			}
		}
	}
	c.stateRevision = s.StateRevision
	selected := *r.Document
	selected.Connections = nil
	for _, con := range r.Document.Connections {
		if con.Definition == r.Entry {
			selected.Connections = append(selected.Connections, con)
		}
	}
	c.Unbound = len(selected.Connections)
	for _, w := range append(s.Widgets(), s.CallbackOwners()...) {
		if w.Binding.Module != "" {
			c.Unbound++
		}
	}
	for _, command := range s.Snapshot().Commands {
		if command.Effect == "" && !command.Toggle && command.Binding.Module == "" {
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
		c.stateRevision = s.StateRevision // Installing the prepared handlers advances state.
		if err = c.Admit(r.SourceRevision); err != nil {
			return nil, err
		}
		symbols := map[string]bool{}
		for _, w := range append(s.Widgets(), s.CallbackOwners()...) {
			if w.Binding.Module != "" {
				if !(s.HasBinding(w.Handle) || s.HasInteractionBinding(w.Handle)) {
					return nil, fmt.Errorf("unbound: %s", w.Handle.Path)
				}
				symbols[w.Binding.Module+"."+w.Binding.Object] = true
			}
		}
		for _, con := range selected.Connections {
			target := con.Definition + "/" + strings.Join(con.Path, "/")
			w, ok := s.Widget(target)
			if !ok || !fieldReceiver(w.Handle.Kind) || !symbols[con.Module+"."+con.Object] {
				return nil, fmt.Errorf("unbound-handle: %s", target)
			}
		}
		for _, command := range s.Snapshot().Commands {
			if command.Effect == "" && !command.Toggle && !s.HasInteractionBinding(command.Handle) {
				return nil, fmt.Errorf("unbound-command: %s", command.Handle.Path)
			}
		}
		c.Unbound = 0
	}
	// The binder and all readiness postchecks precede native candidate resources.
	// Pure layout/presentation validation still precedes application binding.
	if r.PreparePresentation != nil {
		if err = s.PreparePresentationWith(r.PreparePresentation); err != nil {
			return nil, err
		}
	}
	c.stateRevision = s.StateRevision
	if err = c.Admit(r.SourceRevision); err != nil {
		return nil, err
	}
	success = true
	return c, nil
}

func fieldReceiver(kind string) bool {
	switch kind {
	case "input", "checkbox", "slider", "select", "number":
		return true
	}
	return false
}
