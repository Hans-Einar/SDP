package fynehost

import (
	"context"
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
)

// DocumentRequest is owned by an application composition root. Guard rechecks
// source/provider/SDL identities without performing I/O or changing live state.
type DocumentRequest struct {
	Document                         *parser.Document
	Entry, SessionID, SourceRevision string
	Sequence                         uint64
	Providers                        map[string]ui.CollectionProvider
	Mode                             preparation.Mode
	Bind                             preparation.Binder
	Guard                            func() error
	// PrepareResources validates/prepares resources using read-only prospective state.
	// Its viewport offsets are the measured effective values, not raw requests.
	PrepareResources func(ui.Snapshot) error
}
type nativePresentation struct {
	objects    []fyne.CanvasObject
	snapshot   ui.Snapshot
	geometry   *layout.SnapshotLayout
	background fyne.Resource
}
type providerFlight struct {
	request ui.LoadRequest
	cancel  context.CancelFunc
}
type Bundle struct {
	Session        *ui.Session
	Document       *parser.Document
	SourceRevision string
	Sequence       uint64
	prepared       *preparation.Candidate
	request        DocumentRequest
	owner          *DocumentHost
	old            *Bundle
	// size is the last accepted presentation size (or the detached candidate size).
	// It may differ from the physical host size after a rejected resize.
	size                     layout.Size
	view                     *View
	presentation, pending    *nativePresentation
	children                 []fyne.CanvasObject
	flights                  map[string]providerFlight
	closed, published, muted bool
}

// DocumentHost owns exactly one published bundle. All methods except provider
// functions run on the Fyne owner goroutine. Legacy RuntimeView remains separate.
type DocumentHost struct {
	Container        *fyne.Container
	canvas           fyne.Canvas
	current          *Bundle
	size             layout.Size // Actual window size, including rejected small sizes.
	latest           uint64
	latestSource     string
	closed, resizing bool
	OnChange         func(*Bundle)
	OnStatus         func(error)
	// Post marshals provider completion to the owner goroutine; defaults to fyne.Do.
	Post func(func())
}

func NewDocumentHost(canvas fyne.Canvas, size layout.Size) *DocumentHost {
	h := &DocumentHost{canvas: canvas, size: size, Post: fyne.Do}
	h.Container = container.New(h)
	return h
}
func (h *DocumentHost) Current() *Bundle                      { return h.current }
func (h *DocumentHost) MinSize([]fyne.CanvasObject) fyne.Size { return fyne.NewSize(1, 1) }
func (h *DocumentHost) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	if h.closed || h.resizing || size.Width <= 0 || size.Height <= 0 {
		return
	}
	next := layout.Size{W: float64(size.Width), H: float64(size.Height)}
	if next != h.size {
		h.status(h.Resize(next))
	}
	if h.current != nil {
		h.current.view.Container.Resize(size)
	}
}
func (h *DocumentHost) status(err error) {
	if err != nil && h.OnStatus != nil {
		h.OnStatus(err)
	}
}
func (h *DocumentHost) Prepare(r DocumentRequest) (*Bundle, error) {
	if h.closed {
		return nil, fmt.Errorf("closed: document host")
	}
	if r.Sequence == 0 || r.Sequence < h.latest {
		return nil, fmt.Errorf("stale: source sequence")
	}
	if !validDocumentSize(h.size) {
		return nil, fmt.Errorf("viewport: intended native size required")
	}
	h.latest, h.latestSource = r.Sequence, r.SourceRevision
	b := &Bundle{owner: h, request: r, old: h.current, size: h.size, Document: r.Document, SourceRevision: r.SourceRevision, Sequence: r.Sequence, flights: map[string]providerFlight{}}
	req := preparation.Request{Document: r.Document, Entry: r.Entry, SessionID: r.SessionID, SourceRevision: r.SourceRevision, Mode: r.Mode, Bind: r.Bind, Providers: r.Providers, Capabilities: admission.PaneCapabilities(), ValidatePresentation: b.stage, PreparePresentation: b.preparePresentation}
	if h.current != nil {
		req.Previous = h.current.Session
	}
	prepared, err := preparation.Prepare(req)
	if err != nil {
		b.Close()
		return nil, err
	}
	b.prepared = prepared
	b.Session = prepared.Session
	if r.Guard != nil {
		if err = r.Guard(); err != nil {
			b.Close()
			return nil, err
		}
	}
	b.connect()
	// Native resources and geometry are installed while detached. Commit uses
	// this ready child and never re-runs measurement, parsing or binding.
	b.apply()
	b.children = []fyne.CanvasObject{b.view.Container}
	return b, nil
}

// stage is a pure prospective geometry probe. It cannot promote native state.
func (b *Bundle) stage(snapshot ui.Snapshot) (ui.PresentationState, error) {
	measured, _, err := b.measure(snapshot)
	if err != nil {
		return ui.PresentationState{}, err
	}
	return measured.PresentationState(), nil
}
func (b *Bundle) measure(snapshot ui.Snapshot) (*layout.SnapshotLayout, *markdown.Provider, error) {
	if b.closed {
		return nil, nil, fmt.Errorf("closed: bundle")
	}
	if err := preparation.Check(b.Document.Profile, snapshot.Root, admission.PaneCapabilities()); err != nil {
		return nil, nil, err
	}
	provider, err := markdown.Prepare(snapshot.Root, nil)
	if err != nil {
		return nil, nil, err
	}
	metrics := &collectionMeasure{snapshot: snapshot, markdown: provider}
	measured, err := (&layout.Engine{Measure: metrics}).LayoutSnapshot(snapshot, b.size)
	return measured, provider, err
}

// preparePresentation owns one private ticket. Runtime alone promotes it after
// final state publication; sequence consumption and failed probes cannot do so.
func (b *Bundle) preparePresentation(snapshot ui.Snapshot) (ui.PresentationTicket, error) {
	measured, provider, err := b.measure(snapshot)
	if err != nil {
		return ui.PresentationTicket{}, err
	}
	size, source := b.size, b.SourceRevision
	if b.view == nil {
		b.view = New(snapshot.Root)
		b.view.prepared = true
		snapshot.Root.Walk(func(n *parser.Instance) {
			switch {
			case n.Widget == "tree" || n.Widget == "list":
				b.view.addCollection(n.Path, newCollectionControl(b.owner, b, n.Path))
			case n.Kind == "composition" && n.Widget == "tabs":
				b.view.addPane(n.Path, newPaneHeader())
			case n.Kind == "composition" && n.Widget == "split":
				b.view.addPane(n.Path, newPaneDivider())
			}
		})
	}
	if b.request.PrepareResources != nil {
		if err = b.request.PrepareResources(snapshot); err != nil {
			return ui.PresentationTicket{}, err
		}
	}
	background, err := svg.Render(measured.Root, svg.Options{Width: size.W, Height: size.H, SkipControls: true, NativeControls: b.view.nativeControls(), Content: provider})
	if err != nil {
		return ui.PresentationTicket{}, err
	}
	backs, thumbs := b.viewportObjects(snapshot, measured)
	objects := []fyne.CanvasObject{b.view.image}
	objects = append(objects, backs...)
	readingOrder(snapshot.Root, func(n *parser.Instance) {
		if c := b.view.controls[n.Path]; c != nil {
			objects = append(objects, c.clip)
		}
	})
	objects = append(objects, thumbs...)
	prepared := &nativePresentation{snapshot: snapshot, geometry: measured, background: fyne.NewStaticResource("sdui.svg", []byte(background)), objects: objects}
	if b.closed || b.size != size || b.SourceRevision != source || b.published && b.owner.current != b {
		return ui.PresentationTicket{}, fmt.Errorf("stale: prepared native bundle")
	}
	return ui.PresentationTicket{Publish: func() { b.pending = prepared }, Discard: func() { prepared = nil }}, nil
}
func (h *DocumentHost) Commit(b *Bundle) error {
	if h.closed || b == nil || b.closed || b.published || b.owner != h {
		return fmt.Errorf("publication: invalid candidate")
	}
	fail := func(err error) error { b.Close(); return err }
	if h.current != b.old || h.latest != b.Sequence || h.latestSource != b.SourceRevision || h.size != b.size {
		return fail(fmt.Errorf("stale: publication identity/size"))
	}
	if err := b.prepared.Admit(b.SourceRevision); err != nil {
		return fail(err)
	}
	if b.request.Guard != nil {
		if err := b.request.Guard(); err != nil {
			return fail(err)
		}
	}
	// Non-yielding, preallocated model/child swap. Toolkit repaint cannot enter
	// old handlers: they compare bundle identity before dispatching.
	old := h.current
	b.muted = true
	b.published = true
	h.current = b
	h.Container.Objects = b.children
	h.Container.Refresh()
	b.muted = false
	if old != nil {
		old.Close()
	}
	// Establish native focus before after() exposes the new bundle to observers.
	// Recheck after its final local loading presentation as well; this is a no-op
	// when the prepared focused control was preserved by native synchronization.
	h.restoreFocus(b)
	h.after()
	h.restoreFocus(b)
	return nil
}

func (h *DocumentHost) restoreFocus(b *Bundle) {
	if h.current != b || b.closed || h.closed || h.canvas == nil {
		return
	}
	path := ""
	if w, ok := b.Session.Widget(b.Session.Focused()); ok && w.Enabled && w.Visible {
		path = w.InstancePath
	}
	for p, t := range b.Session.Snapshot().Tabs {
		if t.Handle.Path == b.Session.Focused() && t.Enabled && t.Visible {
			path = p
		}
	}
	for p, t := range b.Session.Snapshot().Splits {
		if t.Handle.Path == b.Session.Focused() && t.Enabled && t.Visible {
			path = p
		}
	}
	obj, ok := b.view.Controls[path].(fyne.Focusable)
	if !ok || h.canvas.Focused() == obj {
		return
	}
	// The new control and wrapper renderers were prepared before publication.
	// FocusGained may draw the marker, but must not synthesize a runtime event.
	muted := b.muted
	b.muted = true
	defer func() { b.muted = muted }()
	h.canvas.Focus(obj)
}

func (h *DocumentHost) Adopt(r DocumentRequest) error {
	b, err := h.Prepare(r)
	if err != nil {
		return err
	}
	return h.Commit(b)
}
func (b *Bundle) Close() {
	if b == nil || b.closed {
		return
	}
	b.closed = true
	if b.Session != nil {
		b.Session.Close()
	} // revoke token acceptance before cancellation
	for _, f := range b.flights {
		f.cancel()
	}
	b.flights = nil
	if b.view != nil {
		b.view.Close()
	}
	if b.prepared != nil {
		b.prepared.Close()
	}
}
func (h *DocumentHost) Close() {
	if h.closed {
		return
	}
	h.closed = true
	if h.current != nil {
		h.current.Close()
	}
}
func (h *DocumentHost) Resize(size layout.Size) error {
	if h.closed {
		return fmt.Errorf("closed: document host")
	}
	if !validDocumentSize(size) {
		return fmt.Errorf("viewport: finite positive native size required")
	}
	h.resizing = true
	defer func() { h.resizing = false }()
	h.size = size
	if h.current == nil {
		return nil
	}
	b := h.current
	previousSize := b.size
	b.size = size
	if err := b.Session.SetViewports(nil); err != nil {
		// Preserve the last valid measurement basis for cancellation and completion.
		// The canvas still clips that presentation to the actual physical window.
		b.size = previousSize
		return err
	}
	h.after()
	return nil
}

func validDocumentSize(size layout.Size) bool {
	return size.W > 0 && size.H > 0 && !math.IsNaN(size.W) && !math.IsNaN(size.H) && !math.IsInf(size.W, 0) && !math.IsInf(size.H, 0)
}

// Mutate applies an application operation through the Session's state gate,
// then synchronizes the accepted presentation and provider lifecycle.
func (h *DocumentHost) Mutate(fn func(*ui.Session) error) error {
	if h.closed || h.current == nil {
		return fmt.Errorf("closed: no published bundle")
	}
	b := h.current
	before := b.Session.StateRevision
	err := fn(b.Session)
	// A handler may publish a valid reentrant update before Dispatch rejects its
	// stale result. Reconcile that accepted state even though the outer call fails.
	if h.current == b && !b.closed && (err == nil || b.Session.StateRevision != before) {
		h.after()
	}
	h.status(err)
	return err
}
func (b *Bundle) Geometry() *layout.SnapshotLayout {
	if b.presentation == nil {
		return nil
	}
	return b.presentation.geometry
}
func (b *Bundle) Controls() map[string]fyne.CanvasObject { return b.view.Controls }
