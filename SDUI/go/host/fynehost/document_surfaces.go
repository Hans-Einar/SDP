package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
	"math"
	"sort"
	"strings"
)

type canvasPresentation struct {
	geometry   *layout.SnapshotLayout
	size       layout.Size
	background fyne.Resource
	objects    []fyne.CanvasObject
	native     *nativeSurface
}
type nativeSurface struct {
	target                                  ui.SurfaceTarget
	path                                    string
	window                                  fyne.Window
	modal                                   *dialog.CustomDialog
	canvas                                  fyne.Canvas
	content                                 *fyne.Container
	chrome                                  *fyne.Container
	status                                  *widget.Label
	anchor                                  *surfaceFocus
	image                                   *canvas.Image
	owner                                   *Bundle
	size                                    layout.Size
	shown, retiring, resizing, modalContent bool
}

func (s *nativeSurface) MinSize([]fyne.CanvasObject) fyne.Size {
	if s.modalContent {
		return fyne.NewSize(float32(s.size.W), float32(s.size.H))
	}
	return fyne.NewSize(1, 1)
}
func (s *nativeSurface) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	if !s.shown || s.retiring || s.resizing || s.owner.muted || s.modal != nil {
		return
	}
	next := layout.Size{W: float64(size.Width), H: float64(size.Height)}
	if next == s.size || !validDocumentSize(next) {
		return
	}
	b := s.owner
	old := b.surfaceSizes[s.path]
	b.surfaceSizes[s.path] = next
	s.resizing = true
	defer func() { s.resizing = false }()
	if err := b.Session.SetViewports(nil); err != nil {
		b.surfaceSizes[s.path] = old
		b.owner.status(err)
		return
	}
	b.owner.after()
}
func windowForCanvas(c fyne.Canvas) fyne.Window {
	for _, w := range fyne.CurrentApp().Driver().AllWindows() {
		if w.Canvas() == c {
			return w
		}
	}
	return nil
}
func surfacePath(snapshot ui.Snapshot, path string) string {
	owner := ""
	for p := range snapshot.Surfaces {
		if strings.HasPrefix(path, p+"/") && len(p) > len(owner) {
			owner = p
		}
	}
	return owner
}

// Source nesting is not native ownership: a root-sibling dialog may be opened
// from any live surface. Order exact opening tokens, never source path lengths.
func orderedSurfacePaths(snapshot ui.Snapshot) ([]string, error) {
	paths := []string{}
	owners := map[ui.SurfaceTarget]string{}
	for path, state := range snapshot.Surfaces {
		if state.Open {
			paths = append(paths, path)
			owners[state.Target] = path
		}
	}
	sort.Strings(paths)
	order := []string{}
	visited := map[string]uint8{}
	var visit func(string) error
	visit = func(path string) error {
		if visited[path] == 2 {
			return nil
		}
		if visited[path] == 1 {
			return fmt.Errorf("surface-native: cyclic parent opening at %s", path)
		}
		visited[path] = 1
		state := snapshot.Surfaces[path]
		if state.ParentSurface != nil {
			parent, ok := owners[*state.ParentSurface]
			if !ok {
				return fmt.Errorf("surface-native: missing exact parent opening for %s", path)
			}
			if err := visit(parent); err != nil {
				return err
			}
		}
		visited[path] = 2
		order = append(order, path)
		return nil
	}
	for _, path := range paths {
		if err := visit(path); err != nil {
			return nil, err
		}
	}
	return order, nil
}
func surfacePaths(snapshot ui.Snapshot) []string {
	// Accepted snapshots have already passed orderedSurfacePaths in the gate.
	paths, _ := orderedSurfacePaths(snapshot)
	return paths
}
func exactSurfacePath(snapshot ui.Snapshot, target ui.SurfaceTarget) string {
	for path, state := range snapshot.Surfaces {
		if state.Open && state.Target == target {
			return path
		}
	}
	return ""
}
func (b *Bundle) canvasSizes(snapshot ui.Snapshot, engine *layout.Engine) (map[string]layout.Size, error) {
	sizes := map[string]layout.Size{}
	paths, err := orderedSurfacePaths(snapshot)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		state := snapshot.Surfaces[path]
		if old := b.surfaces[path]; old != nil && old.target == state.Target && !state.Modal {
			sizes[path] = b.surfaceSizes[path]
			continue
		}
		reference := b.size
		if state.ParentSurface != nil {
			parentPath := exactSurfacePath(snapshot, *state.ParentSurface)
			var ok bool
			reference, ok = sizes[parentPath]
			if !ok {
				return nil, fmt.Errorf("surface-native: missing prepared parent size for %s", path)
			}
		}
		size, err := engine.SurfaceSize(snapshot, path, reference)
		if err != nil {
			return nil, err
		}
		// Concrete canvas adapter MinSize is 1x1, including an empty dialog.
		size.W = math.Max(1, size.W)
		size.H = math.Max(1, size.H)
		sizes[path] = size
		if state.Modal {
			// Probe actual chrome against the prospective parent, including an
			// already-open dialog. This detached object is never shown.
			parent := windowForCanvas(b.owner.canvas)
			if parent == nil {
				return nil, fmt.Errorf("surface-native: missing parent window for %s", path)
			}
			measure := &nativeSurface{size: size, modalContent: true}
			body := container.NewCenter(container.New(measure))
			status := widget.NewLabel(" ")
			status.Truncation = fyne.TextTruncateEllipsis
			chrome := container.NewBorder(status, nil, nil, nil, body)
			minimum := dialog.NewCustomWithoutButtons(state.Label, chrome, parent).MinSize()
			if float64(minimum.Width) > reference.W || float64(minimum.Height) > reference.H {
				return nil, fmt.Errorf("surface-native: dialog chrome/content exceeds parent at %s", path)
			}
		}
	}
	return sizes, nil
}
func (b *Bundle) nativeInventory(root *parser.Instance) map[string]string {
	all := b.view.nativeControls()
	out := map[string]string{}
	root.Walk(func(n *parser.Instance) {
		if kind, ok := all[n.Path]; ok {
			out[n.Path] = kind
		}
		if n.Widget == "menu" || n.Widget == "dialog" {
			out[n.Path] = n.Widget
		}
	})
	return out
}
func (b *Bundle) prepareCanvases(p *nativePresentation, all *layout.CanvasLayout, sizes map[string]layout.Size, provider svg.ContentRenderer) error {
	p.canvases = map[string]*canvasPresentation{}
	paths, err := orderedSurfacePaths(p.snapshot)
	if err != nil {
		return err
	}
	for _, path := range paths {
		state := p.snapshot.Surfaces[path]
		g := all.Surfaces[path]
		size := sizes[path]
		background, err := svg.Render(g.Root, svg.Options{Width: size.W, Height: size.H, SkipControls: true, NativeControls: b.nativeInventory(g.Root.Instance), InteractionRoot: p.snapshot.Root, Content: provider})
		if err != nil {
			b.discardCanvases(p)
			return err
		}
		native := b.surfaces[path]
		if native == nil || native.target != state.Target {
			native = &nativeSurface{target: state.Target, path: path, owner: b, size: size, modalContent: state.Modal}
			native.image = canvas.NewImageFromResource(fyne.NewStaticResource("surface.svg", []byte(background)))
			native.image.FillMode = canvas.ImageFillStretch
			native.content = container.New(native, native.image)
			native.status = widget.NewLabel(" ")
			native.status.Truncation = fyne.TextTruncateEllipsis
			var body fyne.CanvasObject = native.content
			if state.Modal {
				body = container.NewCenter(native.content)
			}
			native.chrome = container.NewBorder(native.status, nil, nil, nil, body)
			native.anchor = newSurfaceFocus(b, native)
			if state.Modal {
				var parent fyne.Window
				parentSize := b.size
				if state.ParentSurface == nil {
					parent = windowForCanvas(b.owner.canvas)
				} else {
					parentPath := exactSurfacePath(p.snapshot, *state.ParentSurface)
					frame := p.canvases[parentPath]
					if frame == nil || frame.native.target != *state.ParentSurface {
						b.discardCanvases(p)
						return fmt.Errorf("surface-native: missing prepared exact parent for %s", path)
					}
					parent = windowForCanvas(frame.native.canvas)
					parentSize = frame.size
				}
				if parent == nil {
					b.discardCanvases(p)
					return fmt.Errorf("surface-native: missing parent window for %s", path)
				}
				native.canvas = parent.Canvas()
				native.modal = dialog.NewCustomWithoutButtons(state.Label, native.chrome, parent)
				min := native.modal.MinSize()
				if float64(min.Width) > parentSize.W || float64(min.Height) > parentSize.H {
					native.destroy()
					b.discardCanvases(p)
					return fmt.Errorf("surface-native: dialog chrome/content exceeds parent at %s", path)
				}
				native.modal.SetOnClosed(func() {
					if native.shown && !native.retiring && !b.muted {
						b.dispatchDialog(native.target, ui.Close)
					}
				})
			} else {
				native.window = fyne.CurrentApp().NewWindow(state.Label)
				native.canvas = native.window.Canvas()
				native.window.SetPadded(false)
				native.window.SetContent(native.chrome)
				native.window.SetOnClosed(native.nativeClosed)
				native.window.SetCloseIntercept(native.nativeCloseRequested)
			}
		}
		backs, thumbs := b.viewportObjects(p.snapshot, g)
		objects := []fyne.CanvasObject{native.image}
		objects = append(objects, backs...)
		readingOrder(p.snapshot.Root, func(n *parser.Instance) {
			if surfacePath(p.snapshot, n.Path) == path {
				if c := b.view.controls[n.Path]; c != nil {
					objects = append(objects, c.clip)
				}
			}
		})
		objects = append(objects, thumbs...)
		if state.Focused == "" {
			objects = append(objects, native.anchor)
		}
		p.canvases[path] = &canvasPresentation{geometry: g, size: size, background: fyne.NewStaticResource("surface.svg", []byte(background)), objects: objects, native: native}
	}
	return nil
}
func (b *Bundle) discardCanvases(p *nativePresentation) {
	if p == nil {
		return
	}
	for path, frame := range p.canvases {
		if b.surfaces[path] != frame.native {
			frame.native.destroy()
		}
	}
}
func (s *nativeSurface) destroy() {
	if s.retiring {
		return
	}
	s.retiring = true
	if s.modal != nil {
		s.modal.Hide()
	}
	if s.window != nil {
		s.window.SetCloseIntercept(nil)
		s.window.Close()
	}
}
func (b *Bundle) syncSurfaces() {
	if b.Session == nil {
		return
	}
	muted := b.muted
	b.muted = true
	defer func() { b.muted = muted }()
	snapshot := b.Session.Snapshot()
	for path, native := range b.surfaces {
		state, ok := snapshot.Surfaces[path]
		if !ok || !state.Open || state.Target != native.target {
			if native.window != nil {
				b.removeKeys(native.canvas, path)
			}
			native.destroy()
			b.retireSurfacePreviews(native)
			delete(b.surfaces, path)
			delete(b.surfaceSizes, path)
		}
	}
	if b.presentation == nil {
		return
	}
	for _, path := range surfacePaths(snapshot) {
		frame := b.presentation.canvases[path]
		if frame == nil || frame.native.target != snapshot.Surfaces[path].Target {
			continue
		}
		native := frame.native
		b.surfaces[path] = native
		b.surfaceSizes[path] = frame.size
		sizeChanged := native.size != frame.size
		native.size = frame.size
		native.anchor.Resize(fyne.NewSize(float32(frame.size.W), float32(frame.size.H)))
		native.content.Objects = frame.objects
		native.image.Resource = frame.background
		native.image.Resize(fyne.NewSize(float32(frame.size.W), float32(frame.size.H)))
		native.image.Refresh()
		native.content.Resize(fyne.NewSize(float32(frame.size.W), float32(frame.size.H)))
		native.content.Refresh()
		message := snapshot.Surfaces[path].Message
		if snapshot.Surfaces[path].AcceptBlocked {
			message = fmt.Sprintf("Accept blocked (%s): %s", snapshot.Surfaces[path].Domain, message)
		}
		if message == "" {
			message = " "
		}
		native.status.SetText(strings.Join(strings.Fields(message), " "))
		if sizeChanged && native.shown && native.modal != nil {
			native.modal.Refresh()
			native.modal.Resize(native.modal.MinSize())
		}
		if !native.shown {
			native.shown = true
			b.installKeys(native.canvas, path)
			if native.window != nil {
				// Ask the actual Border chrome for the outer size, including its
				// padding. Source/shared size describes the content only.
				native.modalContent = true
				outer := native.chrome.MinSize()
				native.window.Resize(outer)
				native.modalContent = false
				native.window.Show()
			}
			if native.modal != nil {
				native.modal.Resize(native.modal.MinSize())
				native.modal.Show()
			}
			b.owner.status(b.Session.ConfirmSurfacePublication(native.target))
		}
	}
}
func (b *Bundle) drainDialogResults() {
	for _, result := range b.Session.DrainDialogResults() {
		if b.owner.OnDialogResult != nil {
			func() {
				defer func() {
					if p := recover(); p != nil {
						b.owner.status(fmt.Errorf("dialog observer panic: %v", p))
					}
				}()
				b.owner.OnDialogResult(result)
			}()
		}
	}
}
func (h *DocumentHost) NativeParentHidden() error { return h.revokeNativeParent("parent-hidden") }
func (h *DocumentHost) NativeParentClosed() error { return h.revokeNativeParent("parent-closed") }
func (h *DocumentHost) revokeNativeParent(reason string) error {
	b := h.current
	if b == nil || b.closed {
		return nil
	}
	err := b.Session.RevokeSurfaces(ui.Handle{}, reason)
	h.status(err)
	b.syncSurfaces()
	b.drainDialogResults()
	h.reconcileLoads(b)
	return err
}

func (b *Bundle) geometryFor(path string) *layout.SnapshotLayout {
	if b.presentation == nil {
		return nil
	}
	if owner := surfacePath(b.presentation.snapshot, path); owner != "" {
		if frame := b.presentation.canvases[owner]; frame != nil {
			return frame.geometry
		}
	}
	return b.presentation.geometry
}
func (b *Bundle) canvasFor(path string) fyne.Canvas {
	if b.presentation != nil {
		if owner := surfacePath(b.presentation.snapshot, path); owner != "" {
			if native := b.surfaces[owner]; native != nil {
				return native.canvas
			}
		}
	}
	return b.owner.canvas
}

func (native *nativeSurface) nativeClosed() {
	b := native.owner
	if native.retiring || b.closed || b.owner.current != b {
		return
	}
	// Fyne Window.Close bypasses the interceptor; real teardown cannot be gated.
	native.retiring = true
	b.owner.status(b.Session.RevokeSurface(native.target, "parent-closed"))
	b.syncSurfaces()
	b.owner.restoreFocus(b)
	b.drainDialogResults()
	b.owner.reconcileLoads(b)
}

// WM/decoration close is native owner loss, not a source Close interaction.
// Revoke the exact opening (and descendants) before destroying any native window;
// neither a child modal input guard nor a failing presentation gate can veto it.
func (native *nativeSurface) nativeCloseRequested() {
	b := native.owner
	if !native.shown || native.retiring || b.closed || b.owner.current != b {
		return
	}
	snapshot := b.Session.Snapshot()
	before := snapshot.ActiveSurface
	activeOwned := false
	for current := before; current != nil; {
		if *current == native.target {
			activeOwned = true
			break
		}
		path := exactSurfacePath(snapshot, *current)
		if path == "" {
			break
		}
		current = snapshot.Surfaces[path].ParentSurface
	}
	if err := b.Session.RevokeSurface(native.target, "parent-closed"); err != nil {
		b.owner.status(err)
		return
	}
	b.syncSurfaces()
	b.owner.restoreFocus(b)
	b.syncMenus()
	b.owner.reconcileLoads(b)
	b.drainDialogResults()
	if activeOwned {
		b.focusAfterGesture(before)
	}
	if b.owner.current == b && !b.closed && !b.owner.closed && b.owner.OnChange != nil {
		b.owner.OnChange(b)
	}
}
