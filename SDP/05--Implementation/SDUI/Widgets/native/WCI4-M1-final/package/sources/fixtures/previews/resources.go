package previews

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// ShapeSVG is a direct closed-subset image, not a Markdown embedded image.
const ShapeSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 200"><title>Site output</title><desc>Red and blue sites with yellow marker</desc><rect x="0" y="0" width="400" height="200" fill="#f0f4f8"/><g transform="translate(10 10)"><rect x="0" y="0" width="140" height="140" fill="#d84040"/><circle cx="260" cy="70" r="65" fill="#2868c0"/><ellipse cx="180" cy="140" rx="25" ry="15" fill="#efb52b"/><path d="M 15 170 L 340 170" fill="none" stroke="#17324d" stroke-width="4"/><line x1="160" y1="10" x2="160" y2="110" stroke="#17324d" stroke-width="3"/><polyline points="175,10 190,30 205,10" fill="none" stroke="#238248" stroke-width="4"/><polygon points="345,110 370,150 330,150" fill="#238248" opacity="0.8"/></g></svg>`

func digest(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func ResourceSlot(slot string) (markdown.Resource, error) {
	s, w, h := ShapeSVG, 400.0, 200.0
	switch slot {
	case "good":
	case "alternate":
		s = strings.ReplaceAll(s, "#d84040", "#7c39b8")
	case "wide":
		s = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 100"><rect width="400" height="100" fill="#2868c0"/></svg>`
		h = 100
	case "tall":
		s = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 400"><rect width="100" height="400" fill="#d84040"/></svg>`
		w, h = 100, 400
	case "negative-origin":
		s = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="-40 -20 400 200"><rect x="-40" y="-20" width="400" height="200" fill="#238248"/></svg>`
	case "malformed":
		s = `<svg viewBox="0 0 400 200"><rect>`
	case "unsupported":
		s = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 200"><text x="10" y="20">Unsupported text</text></svg>`
	case "bad-dimensions":
		w = 399
	default:
		return markdown.Resource{}, fmt.Errorf("unknown resource slot %q", slot)
	}
	return markdown.Resource{SVG: []byte(s), Width: w, Height: h}, nil
}
func supplied(slot string) (markdown.PreparedSVG, error) {
	r, err := ResourceSlot(slot)
	return markdown.PreparedSVG{Source: parser.Reference{Module: "art", Object: "Chart", Member: "resource"}, ProviderID: "fixture-svg/1", SHA256: digest(r.SVG), Resource: r}, err
}

type countedRenderer struct {
	fixture    *Fixture
	path, mode string
}

func (r countedRenderer) Render(source string) (markdown.Resource, error) {
	f := r.fixture
	f.mu.Lock()
	f.rendererCalls[r.path]++
	count := f.rendererCalls[r.path]
	f.mu.Unlock()
	v, err := ResourceSlot("good")
	switch r.mode {
	case "error":
		err = fmt.Errorf("injected diagram renderer error")
	case "malformed":
		v, _ = ResourceSlot("malformed")
	case "unsafe":
		v.SVG = []byte(`<svg viewBox="0 0 400 200"><script>unsafe</script></svg>`)
	}
	f.mu.Lock()
	f.returns[r.path] = v.SVG
	f.mu.Unlock()
	event := map[string]any{"path": r.path, "calls": count, "mode": r.mode, "sourceSHA256": digest([]byte(source)), "outputSHA256": digest(v.SVG), "bytes": len(v.SVG), "width": v.Width, "height": v.Height}
	if err != nil {
		event["error"] = err.Error()
	}
	f.emit("provider-event", event)
	return v, err
}
