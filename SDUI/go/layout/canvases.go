package layout

import (
	"sort"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// CanvasLayout keeps independent canvas coordinates. Paths remain globally unique
// so one PresentationGate can validate all offsets/split bounds atomically.
type CanvasLayout struct {
	Main     *SnapshotLayout
	Surfaces map[string]*SnapshotLayout
}

func (g *CanvasLayout) PresentationState() runtime.PresentationState {
	out := runtime.PresentationState{Viewports: map[string]runtime.ViewportState{}, Splits: map[string]runtime.SplitGeometry{}}
	merge := func(canvas *SnapshotLayout) {
		if canvas == nil {
			return
		}
		p := canvas.PresentationState()
		for k, v := range p.Viewports {
			out.Viewports[k] = v
		}
		for k, v := range p.Splits {
			out.Splits[k] = v
		}
	}
	merge(g.Main)
	for _, s := range g.Surfaces {
		merge(s)
	}
	return out
}

type surfaceNode struct {
	node *parser.Instance
	font float64
}

func surfaceNodes(root *parser.Instance) (map[string]surfaceNode, error) {
	profile, err := parser.EffectiveProfile(root)
	if err != nil {
		return nil, err
	}
	out := map[string]surfaceNode{}
	paths := map[string]bool{}
	var visit func(*parser.Instance, float64) error
	visit = func(n *parser.Instance, font float64) error {
		if paths[n.Path] {
			return diag(n, "canvas-path", "Duplicate normalized path")
		}
		paths[n.Path] = true
		font = number(n, "font", font)
		if dialog(n) {
			out[n.Path] = surfaceNode{n, font}
		}
		for _, row := range allRows(n) {
			for _, c := range row {
				if err := visit(c, font); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if profile == "sdui/0.3" {
		if err := visit(root, 14); err != nil {
			return nil, err
		}
	}
	return out, nil
}
func validCanvas(size Size) bool {
	return finite(size.W) && finite(size.H) && size.W > 0 && size.H > 0 && size.W <= 32768 && size.H <= 32768
}
func (e *Engine) snapshotRun(snapshot runtime.Snapshot) *Engine {
	run := &Engine{Measure: e.Measure, fieldState: snapshot.Fields, requested: snapshot.Viewports, viewports: map[string]Viewport{}, tabState: snapshot.Tabs, splitState: snapshot.Splits}
	if run.Measure == nil {
		run.Measure = TextMetrics{}
	}
	return run
}
func snapshotResult(run *Engine, root *Box) *SnapshotLayout {
	return &SnapshotLayout{Root: root, Fields: run.fields, Viewports: run.viewports, Tabs: run.tabs, Splits: run.splits}
}
func checkOffsets(snapshot runtime.Snapshot) error {
	for _, offset := range snapshot.Viewports {
		if !validOffset(offset) {
			return &parser.Diagnostic{Code: "viewport-offset", Message: "Viewport offsets must be finite and nonnegative"}
		}
	}
	return nil
}

// SurfaceSize resolves opening-only root sizing against the caller's correct
// parent reference. The prospective snapshot must already identify an open dialog.
func (e *Engine) SurfaceSize(snapshot runtime.Snapshot, path string, reference Size) (Size, error) {
	nodes, err := surfaceNodes(snapshot.Root)
	if err != nil {
		return Size{}, err
	}
	n, ok := nodes[path]
	if !ok {
		return Size{}, diag(snapshot.Root, "surface-layout", "Unknown dialog path")
	}
	state, ok := snapshot.Surfaces[path]
	if !ok || !state.Open || !visible(n.node) {
		return Size{}, diag(n.node, "surface-layout", "Dialog must be prospectively open and visible")
	}
	if !validCanvas(reference) {
		return Size{}, diag(n.node, "surface-layout", "Surface reference must be finite and in (0,32768]")
	}
	if err = checkOffsets(snapshot); err != nil {
		return Size{}, err
	}
	if err = validateScrollOwners(snapshot.Root); err != nil {
		return Size{}, err
	}
	run := e.snapshotRun(snapshot)
	run.profile = "sdui/0.3"
	run.surfaceRoot = n.node
	size, err := run.desired(n.node, reference, reference, assigned{}, n.font)
	if err != nil {
		return Size{}, err
	}
	if !finiteExtent(size) || size.W > 32768 || size.H > 32768 {
		return Size{}, diag(n.node, "surface-layout", "Natural opening size must be finite and in [0,32768]")
	}
	return size, nil
}

// LayoutCanvases measures the entire prospective presentation or fails without a
// result. Existing nonmodal sizes are supplied explicitly, never derived again
// from a resized parent. Closed-surface sizes can remain in the caller's cache.
func (e *Engine) LayoutCanvases(snapshot runtime.Snapshot, mainSize Size, surfaceSizes map[string]Size) (*CanvasLayout, error) {
	nodes, err := surfaceNodes(snapshot.Root)
	if err != nil {
		return nil, err
	}
	if !validCanvas(mainSize) {
		return nil, diag(snapshot.Root, "viewport", "Viewport must be finite and in (0,32768]")
	}
	for path := range surfaceSizes {
		if _, ok := nodes[path]; !ok {
			return nil, diag(snapshot.Root, "surface-layout", "Unknown surface size path: "+path)
		}
	}
	paths := []string{}
	for path, state := range snapshot.Surfaces {
		n, ok := nodes[path]
		if !ok {
			return nil, diag(snapshot.Root, "surface-layout", "Unknown surface state path: "+path)
		}
		if !state.Open {
			continue
		}
		if !visible(n.node) {
			return nil, diag(n.node, "surface-layout", "Open surface is inactive")
		}
		if size, ok := surfaceSizes[path]; !ok || !validCanvas(size) {
			return nil, diag(n.node, "surface-layout", "Open surface requires a finite accepted canvas size in (0,32768]")
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	if err = checkOffsets(snapshot); err != nil {
		return nil, err
	}
	mainRun := e.snapshotRun(snapshot)
	mainRoot, err := mainRun.layout(snapshot.Root, mainSize)
	if err != nil {
		return nil, err
	}
	main := snapshotResult(mainRun, mainRoot)
	out := &CanvasLayout{Main: main, Surfaces: map[string]*SnapshotLayout{}}
	// A single invocation retains the same whole-model operation budget even
	// though each canvas starts fresh clipping/viewport ancestry and metric caches.
	operations := mainRun.operations
	for _, path := range paths {
		n := nodes[path]
		size := surfaceSizes[path]
		run := e.snapshotRun(snapshot)
		run.profile = "sdui/0.3"
		run.surfaceRoot = n.node
		run.operations = operations
		rect := Rect{W: size.W, H: size.H}
		root, err := run.arrange(n.node, rect, size, rect, n.font, true, "")
		if err != nil {
			return nil, err
		}
		out.Surfaces[path] = snapshotResult(run, root)
		operations = run.operations
	}
	// No path may occur in two canvas geometries, even for malformed hand-built
	// input. This also protects aggregate map assembly from silent replacement.
	seen := map[string]bool{}
	canvases := []*SnapshotLayout{main}
	for _, path := range paths {
		canvases = append(canvases, out.Surfaces[path])
	}
	for _, canvas := range canvases {
		var duplicate string
		canvas.Root.Walk(func(b *Box) {
			if seen[b.Path] {
				duplicate = b.Path
			}
			seen[b.Path] = true
		})
		if duplicate != "" {
			return nil, diag(snapshot.Root, "canvas-path", "Geometry appears in multiple canvases: "+duplicate)
		}
	}
	return out, nil
}
