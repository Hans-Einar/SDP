package runtime

import (
	"sort"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type Session struct {
	commands                map[string]*CommandState
	presentations           map[string]*CommandPresentation
	menus                   map[string]*MenuState
	surfaces                map[string]*SurfaceState
	aux                     map[string]*Widget
	ownerSurface, menuOwner map[string]string
	menuTargets             map[string]Handle
	rootOwner               Handle
	surfaceEpoch            uint64
	publishedSurfaces       map[SurfaceTarget]bool
	dialogResults           []DialogResult
	activeSurface           *SurfaceTarget
	mainFocus               string
	accepting               *SurfaceTarget
	check                   func(*parser.Instance) error
	stateCheck              StateGate
	presentationCheck       PresentationGate
	presentationPrepare     PresentationPrepare
	panes                   map[string]*Widget
	tabs                    map[string]*TabsState
	splits                  map[string]*SplitState
	intent, active          map[string]activity
	interactions            map[string]InteractionHandler
	interacting             bool
	strictSplit             string
	StateRevision           uint64
	collections             map[string]*collection
	viewports               map[string]ViewportState
	viewportHandles         map[string]Handle
	ID                      string
	Revision, BatchRevision uint64
	root                    *parser.Instance
	widgets                 map[string]*Widget
	handlers                map[string]Handler
	focused                 string
	generation, sequence    uint64
	closed                  bool
}

func New(id string, root *parser.Instance) (*Session, error) {
	if id == "" || root == nil || root.Kind != "frame" {
		return nil, fault("session", "Session ID and frame root are required")
	}
	if _, err := parser.EffectiveProfile(root); err != nil {
		return nil, err
	}
	s := &Session{ID: id, Revision: 1, StateRevision: 1, collections: map[string]*collection{}, viewports: map[string]ViewportState{}, viewportHandles: map[string]Handle{}, root: clone(root), widgets: map[string]*Widget{}, handlers: map[string]Handler{}}
	if err := s.register(s.root, true, true); err != nil {
		return nil, err
	}
	if err := s.initPanes(); err != nil {
		return nil, err
	}
	if err := s.initCommands(); err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Session) register(n *parser.Instance, enabled, visible bool) error {
	parentEnabled, parentVisible := enabled, visible
	if n.Layout["overflow-x"] == "scroll" || n.Layout["overflow-y"] == "scroll" {
		s.generation++
		s.viewportHandles[n.Path] = Handle{s.ID, n.Path, s.generation, "viewport:" + n.Kind + ":" + n.Widget}
		s.viewports[n.Path] = ViewportState{}
	}
	enabled = enabled && n.Layout["enabled"] != false
	visible = visible && n.Layout["visible"] != false
	if n.Kind == "widget" && n.Widget != "command" && n.Widget != "item" && n.Widget != "separator" {
		path := public(n.Path)
		if s.widgets[path] != nil {
			return fault("ambiguous-instance", path)
		}
		s.generation++
		w := &Widget{ancestorEnabled: parentEnabled, ancestorVisible: parentVisible, Handle: Handle{s.ID, path, s.generation, n.Widget}, InstancePath: n.Path, Label: n.Argument("label"), Value: n.Argument("value"), Draft: n.Argument("value"), ValueRevision: 1, DraftRevision: 1, Enabled: enabled, Visible: visible}
		if n.Widget == "input" {
			w.Label = n.Argument("text")
		}
		if ref, ok := n.Arguments["callback"].(parser.Reference); ok {
			w.Binding = ref
		}
		s.widgets[path] = w
	}
	for _, r := range n.Regions {
		if err := s.register(r.Node, enabled, visible); err != nil {
			return err
		}
	}
	for _, row := range n.Rows {
		for _, c := range row {
			if err := s.register(c, enabled, visible); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Session) Widgets() []Widget {
	keys := []string{}
	for k := range s.widgets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []Widget{}
	for _, k := range keys {
		out = append(out, *s.widgets[k])
	}
	return out
}
func (s *Session) Widget(path string) (Widget, bool) {
	w, ok := s.widgets[path]
	if !ok {
		return Widget{}, false
	}
	return *w, true
}
func (s *Session) Bind(handle Handle, handler Handler) error {
	if _, e := s.lookup(handle); e != nil {
		return e
	}
	if s.presentations[handle.Path] != nil {
		return fault("binding", "Command buttons use BindInteraction")
	}
	if handler == nil {
		return fault("binding", "Nil callback")
	}
	s.handlers[handle.Path] = handler
	s.StateRevision++
	return nil
}
func (s *Session) Draft(handle Handle, value string) error {
	w, e := s.lookup(handle)
	if e != nil {
		return e
	}
	if w.Handle.Kind != "input" || !s.inputAllowed(w) {
		return fault("draft", "Input is not editable")
	}
	if len(value) > 32768 {
		return fault("value-limit", "Draft exceeds 32768 bytes")
	}
	n := s.copyState()
	w = n.widgets[handle.Path]
	w.Draft = value
	w.DraftRevision++
	w.Dirty = w.Draft != w.Value
	return s.publish(n)
}
func (s *Session) Revert(handle Handle) error {
	w, e := s.lookup(handle)
	if e != nil {
		return e
	}
	if w.Handle.Kind != "input" {
		return fault("widget-type", "Revert requires input")
	}
	n := s.copyState()
	w = n.widgets[handle.Path]
	w.Draft = w.Value
	w.Dirty = false
	w.DraftRevision++
	return s.publish(n)
}
func (s *Session) Focus(handle Handle) error {
	w, e := s.lookupControl(handle)
	if e != nil {
		return e
	}
	if w.Handle.Kind == "page" || s.aux[handle.Path] != nil {
		return fault("focus", "Page content has no independent focus stop")
	}
	if !s.inputAllowed(w) {
		return fault("focus", "Widget is not focusable")
	}
	n := s.copyState()
	n.focused = handle.Path
	n.rememberFocus(handle.Path)
	return s.publish(n)
}
func (s *Session) Focused() string { return s.focused }
func (s *Session) Close() {
	if s.closed {
		return
	}
	_ = s.RevokeSurfaces(Handle{}, "dispose")
	for _, c := range s.collections {
		c.cancel()
	}
	s.closed = true
	s.StateRevision++
	s.handlers = map[string]Handler{}
	s.interactions = map[string]InteractionHandler{}
	s.presentationPrepare = nil
	s.focused = ""
	s.Revision++
}

// CheckWith installs a pure presentation gate. Failed geometry never publishes
// partial state; callbacks are never invoked by the gate.
func (s *Session) CheckWith(check func(*parser.Instance) error) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	if check != nil {
		if err := check(s.SnapshotRoot()); err != nil {
			return err
		}
	}
	s.check = check
	s.StateRevision++
	return nil
}

// InvalidateEvents advances the published binding/model epoch without changing
// widget state. Hosts rebuild native event closures after an SDL model change.
func (s *Session) InvalidateEvents() error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	n := s.copyState()
	for _, c := range n.collections {
		c.cancel()
	}
	n.Revision++
	return s.publish(n)
}

func (s *Session) Closed() bool { return s.closed }

// HasBinding reports an installed handler for this exact live widget identity.
// It does not execute the handler or attest to an external module signature.
func (s *Session) HasBinding(handle Handle) bool {
	if _, err := s.lookup(handle); err != nil {
		return false
	}
	return s.handlers[handle.Path] != nil
}
