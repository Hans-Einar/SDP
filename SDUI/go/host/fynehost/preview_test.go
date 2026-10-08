package fynehost

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

const previewSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 200"><rect width="400" height="200" fill="#ff0000"/></svg>`
const previewSource = `sdui 0.3;ref: art "application-art";page=[picture=svg(art.Chart.@resource,label="Chart caption",description="Production by site",fallback="reject") {x=fill,y=fill};prose=markdown("# Prepared prose",description="Summary",fallback="label") {x=fill}] {x=fill,y=fill};`

func previewRequest(t *testing.T, source string, seq uint64) DocumentRequest {
	t.Helper()
	doc, err := parser.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	return DocumentRequest{Document: doc, Entry: "page", SessionID: "preview-test", Sequence: seq, SourceRevision: fmt.Sprint(seq), Mode: preparation.Prototype}
}
func bindPreview(r *DocumentRequest, path string, data []byte) {
	if r.SVGResources == nil {
		r.SVGResources = map[string]markdown.PreparedSVG{}
	}
	r.SVGResources[path] = markdown.PreparedSVG{Source: parser.Reference{Module: "art", Object: "Chart", Member: "resource"}, ProviderID: "fixture-art", SHA256: fmt.Sprintf("%x", sha256.Sum256(data)), Resource: markdown.Resource{SVG: data, Width: 400, Height: 200}}
}

func TestPreparedPreviewNativeOwnershipGeometryAndAccessibility(t *testing.T) {
	h := documentHost(t)
	r := previewRequest(t, previewSource, 1)
	data := []byte(previewSVG)
	bindPreview(&r, "page/picture", data)
	b, err := h.Prepare(r)
	if err != nil {
		t.Fatal(err)
	}
	data[0] = '!'
	if err = h.Commit(b); err != nil {
		t.Fatal(err)
	}
	c := b.Controls()["page/picture"].(*previewControl)
	if !bytes.Equal(c.image.Resource.Content(), []byte(previewSVG)) {
		t.Fatal("input buffer was retained")
	}
	copy := c.image.Resource.Content()
	copy[0] = '!'
	if c.image.Resource.Content()[0] != '<' {
		t.Fatal("native resource accessor leaked authoritative bytes")
	}
	if b.request.SVGResources != nil || b.request.MarkdownRenderers != nil {
		t.Fatal("request retained resource/provider inputs")
	}
	if c.frame.image.W <= 0 || c.frame.image.W/c.frame.image.H != 2 {
		t.Fatal("aspect ratio", c.frame)
	}
	if c.frame.caption.H <= 0 || c.frame.image.Intersect(c.frame.caption).H > 0 {
		t.Fatal("caption/image overlap", c.frame)
	}
	if c.frame.text == nil || strings.Contains(string(c.frame.text.Content()), "<image") {
		t.Fatal("missing direct glyph overlay / nested image")
	}
	if err := checkPreviewSVG(c.frame.text.Content()); err != nil {
		t.Fatal(err)
	}
	md := b.Controls()["page/prose"].(*previewControl)
	if md.frame.text != nil || md.frame.resource != nil {
		t.Fatal("Markdown duplicated background paint")
	}
	if b.nativeInventory(b.Session.SnapshotRoot())["page/prose"] != "" {
		t.Fatal("invented native Markdown inventory kind")
	}
	if b.nativeInventory(b.Session.SnapshotRoot())["page/picture"] != "svg" {
		t.Fatal("actual SVG missing exact inventory")
	}
	if !strings.Contains(fyne.Accessible(md).AccessibilityLabel(), "Summary") {
		t.Fatal("missing actual Markdown accessibility")
	}
	info := b.inspectPreviews()["page/picture"]
	if info["accessibleLabel"] != fyne.Accessible(c).AccessibilityLabel() || info["canvas"] != "main" {
		t.Fatal("inspector not actual native interface", info)
	}
	if err := h.Resize(layout.Size{W: 600, H: 450}); err != nil {
		t.Fatal(err)
	}
	if b.Controls()["page/picture"] != c || c.frame.image.W/c.frame.image.H != 2 {
		t.Fatal("resize replaced identity/fit")
	}
	b.Close()
	if c.image.Resource != nil || c.text.Resource != nil || b.previews != nil || b.previewResources != nil {
		t.Fatal("closed bundle retained preview resources")
	}
}

func TestPreparedPreviewFallbackIdentityAndPublicationFailure(t *testing.T) {
	h := documentHost(t)
	r := previewRequest(t, strings.Replace(previewSource, `fallback="reject"`, `fallback="label"`, 1), 1)
	// Missing SVG resource is an explicit measured label, even with a caption.
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	old := h.Current()
	c := old.Controls()["page/picture"].(*previewControl)
	if c.frame.resource != nil || c.frame.status.H <= 0 || c.frame.caption.H <= 0 {
		t.Fatal("fallback hidden by caption", c.frame)
	}
	if !strings.Contains(c.AccessibilityLabel(), "Preview unavailable: Production by site") {
		t.Fatal(c.AccessibilityLabel())
	}
	for _, cap := range old.capabilities() {
		if cap.ID == "svg-resource" || cap.ID == "mermaid-flowchart" {
			t.Fatal("label advertised rich readiness", cap)
		}
	}
	r = previewRequest(t, strings.Replace(previewSource, `fallback="reject"`, `fallback="label"`, 1), 2)
	bindPreview(&r, "page/picture", []byte(previewSVG))
	bad := r.SVGResources["page/picture"]
	bad.SHA256 = strings.Repeat("0", 64)
	r.SVGResources["page/picture"] = bad
	if _, err := h.Prepare(r); err == nil {
		t.Fatal("digest mismatch became fallback")
	}
	if h.Current() != old || old.closed {
		t.Fatal("failed candidate touched publication")
	}
	r = previewRequest(t, previewSource, 3)
	bindPreview(&r, "page/picture", []byte(previewSVG))
	stale := false
	r.Guard = func() error {
		if stale {
			return errors.New("resource revision changed")
		}
		return nil
	}
	next, err := h.Prepare(r)
	if err != nil {
		t.Fatal(err)
	}
	stale = true
	if err = h.Commit(next); err == nil || !next.closed || h.Current() != old {
		t.Fatal("stale resource publication", err)
	}
	if c.frame.status.H <= 0 || old.closed {
		t.Fatal("failed publication disposed previous content")
	}
}

type previewCountingRenderer struct {
	calls int
	bytes []byte
}

func (r *previewCountingRenderer) Render(string) (markdown.Resource, error) {
	r.calls++
	return markdown.Resource{SVG: r.bytes, Width: 400, Height: 200}, nil
}
func TestPreparedPreviewRendererOnlyDuringPreparation(t *testing.T) {
	h := documentHost(t)
	source := "sdui 0.3;page=[prose=markdown(\"Before\\n\\n```mermaid\\ngraph TD; A-->B\\n```\\n\\nAfter\",description=\"Diagram summary\",fallback=\"label\") {x=fill,y=fill}] {x=fill,y=fill};"
	r := previewRequest(t, source, 1)
	renderer := &previewCountingRenderer{bytes: []byte(previewSVG)}
	r.MarkdownRenderers = map[string]markdown.MarkdownRenderer{"page/prose": {ProviderID: "test-renderer", Revision: "one", Renderer: renderer}}
	b, err := h.Prepare(r)
	if err != nil {
		t.Fatal(err)
	}
	if renderer.calls != 1 {
		t.Fatal("preparation renderer count", renderer.calls)
	}
	renderer.bytes[0] = '!'
	if err := h.Commit(b); err != nil {
		t.Fatal(err)
	}
	for i := range 3 {
		if _, err := b.stage(b.Session.Snapshot()); err != nil {
			t.Fatal(err)
		}
		if err := h.Resize(layout.Size{W: float64(700 + i*10), H: 550}); err != nil {
			t.Fatal(err)
		}
	}
	if renderer.calls != 1 {
		t.Fatal("renderer invoked by gate/resize/commit", renderer.calls)
	}
	o, _ := b.previews.Outcome("page/prose")
	if o.Status != "partial" || len(o.Diagrams) != 1 || o.Diagrams[0].Resource != nil {
		t.Fatal("unsupported native diagram claimed rendered", o)
	}
	if !strings.Contains(string(b.presentation.background.Content()), "Preview unavailable") {
		t.Fatal("diagram fallback missing from background")
	}
	for _, cap := range b.capabilities() {
		if cap.ID == "mermaid-flowchart" {
			t.Fatal("false Mermaid capability")
		}
	}
}

func TestPreparedPreviewFailedTicketPreservesNativeCaption(t *testing.T) {
	h := documentHost(t)
	r := previewRequest(t, previewSource, 1)
	bindPreview(&r, "page/picture", []byte(previewSVG))
	reject := false
	r.PrepareResources = func(ui.Snapshot) error {
		if reject {
			return errors.New("prospective native resource failed")
		}
		return nil
	}
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	c := b.Controls()["page/picture"].(*previewControl)
	oldOverlay, oldPresentation := c.text.Resource, b.presentation
	w, _ := b.Session.Widget("page/picture")
	change := func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: ui.Label, Value: ui.Text("Replacement caption")}})
	}
	reject = true
	if err := h.Mutate(change); err == nil {
		t.Fatal("failed native ticket admitted")
	}
	if b.presentation != oldPresentation || c.text.Resource != oldOverlay || b.pending != nil {
		t.Fatal("failed ticket changed native resources")
	}
	current, _ := b.Session.Widget("page/picture")
	if current.Label != w.Label {
		t.Fatal("failed ticket published caption")
	}
	reject = false
	if err := h.Mutate(change); err != nil {
		t.Fatal(err)
	}
	if c.text.Resource == oldOverlay || !strings.Contains(string(c.text.Resource.Content()), "Replacement caption") {
		t.Fatal("accepted caption not published")
	}
	if c.image.Resource != b.previewResources["page/picture"] {
		t.Fatal("caption mutation replaced frozen resource")
	}
}

func TestPreparedPreviewScrollClipsWithoutRefitting(t *testing.T) {
	h := documentHost(t)
	source := `sdui 0.3;ref: art "application-art";page=[viewport=[picture=svg(art.Chart.@resource,description="Large image",fallback="reject") {scale-x=2,scale-y=2}] {x=fill,y=fill,overflow-x=scroll,overflow-y=scroll}] {x=fill,y=fill};`
	r := previewRequest(t, source, 1)
	bindPreview(&r, "page/viewport/picture", []byte(previewSVG))
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	c := b.Controls()["page/viewport/picture"].(*previewControl)
	before := c.frame.image
	if err := h.Mutate(func(s *ui.Session) error {
		return s.SetViewports(map[string]ui.ViewportState{"page/viewport": {X: 100, Y: 80}})
	}); err != nil {
		t.Fatal(err)
	}
	after := c.frame.image
	if before.W != after.W || before.H != after.H || after.X != before.X-100 || after.Y != before.Y-80 {
		t.Fatal("scroll refit artwork", before, after)
	}
	info := b.inspectPreviews()["page/viewport/picture"]
	clip := info["clip"].(layout.Rect)
	if clip.W >= after.W || clip.H >= c.frame.box.H {
		t.Fatal("missing ancestor clip", info)
	}
}

func TestPreparedPreviewClosedSurfaceInventoryAndReopen(t *testing.T) {
	h, w := lifecycleHost(t)
	source := `sdui 0.3;ref: art "application-art";page=[open=button("Open",effect="open",target="detail");detail=dialog("Preview detail",modal=false)[image=svg(art.Chart.@resource,description="Detail art",fallback="reject") {x=fill,y=fill};close=button("Close",effect="close")] {scale-x=0.6,scale-y=0.6}] {x=fill,y=fill};`
	r := previewRequest(t, source, 1)
	// Closed resources participate before any native window is opened.
	if _, err := h.Prepare(r); err == nil {
		t.Fatal("closed missing resource admitted")
	}
	r.Sequence++
	bindPreview(&r, "page/detail/image", []byte(previewSVG))
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	if b.inspectPreviews()["page/detail/image"]["mounted"] != false {
		t.Fatal("closed declaration pretends mounted")
	}
	open := b.Controls()["page/open"].(*commandButton)
	open.Tapped(&fyne.PointEvent{})
	native := b.surfaces["page/detail"]
	if native == nil || !native.shown {
		t.Fatal("preview surface did not open")
	}
	c := b.Controls()["page/detail/image"].(*previewControl)
	if c.image.Resource == nil {
		t.Fatal("opened preview lacks resource")
	}
	info := b.inspectPreviews()["page/detail/image"]
	if info["canvas"] != "page/detail" || info["accessibleLabel"] != c.AccessibilityLabel() {
		t.Fatal("wrong actual preview canvas", info)
	}
	b.Controls()["page/detail/close"].(*commandButton).Tapped(&fyne.PointEvent{})
	if b.inspectPreviews()["page/detail/image"]["mounted"] != false || c.image.Resource != nil {
		t.Fatal("closed surface retained native preview")
	}
	if b.previewResources["page/detail/image"] == nil {
		t.Fatal("surface closure disposed bundle declaration")
	}
	open.Tapped(&fyne.PointEvent{})
	if c.image.Resource == nil || b.Controls()["page/detail/image"] != c {
		t.Fatal("reopen lost resource/control identity")
	}
}
