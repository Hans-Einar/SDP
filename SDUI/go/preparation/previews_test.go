package preparation_test

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

const preparedPreviewSource = `sdui 0.3;ref: art "resource-module";page=[closed=dialog("Unopened")[image=svg(art.Chart.@resource,description="Original description",fallback="label")]];`

func previewPreparation(t *testing.T) (preparation.Request, *markdown.Previews) {
	t.Helper()
	doc, roots, err := parser.Compile(preparedPreviewSource)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 10"><rect width="20" height="10" fill="red"/></svg>`)
	resources := map[string]markdown.PreparedSVG{"page/closed/image": {Source: parser.Reference{Module: "art", Object: "Chart", Member: "resource"}, ProviderID: "resource", SHA256: fmt.Sprintf("%x", sha256.Sum256(data)), Resource: markdown.Resource{SVG: data, Width: 20, Height: 10}}}
	p, err := markdown.PreparePreviews(roots["page"], resources, nil, markdown.PreviewBackend{SVG: func(markdown.Resource) error { return nil }})
	if err != nil {
		t.Fatal(err)
	}
	caps := append(admission.TextCapabilities(), preparation.Capabilities{
		{Dimension: preparation.Widget, ID: "svg", Major: 1},
		{Dimension: preparation.Provider, ID: "svg-resource", Major: 1},
		{Dimension: preparation.Layout, ID: "preview-resource", Major: 1},
		{Dimension: preparation.Host, ID: "svg-resource", Major: 1},
	}...)
	return preparation.Request{Document: doc, Entry: "page", SessionID: "candidate", SourceRevision: "source", Mode: preparation.Prototype, Capabilities: caps, Previews: p, ValidateLayout: func(*parser.Instance) error { return nil }}, p
}

func TestPreviewFingerprintPrecedesCapsSessionAndGates(t *testing.T) {
	r, _ := previewPreparation(t)
	doc, err := parser.Parse(strings.Replace(preparedPreviewSource, "Original description", "Foreign description", 1))
	if err != nil {
		t.Fatal(err)
	}
	r.Document = doc
	r.Capabilities = nil // Root mismatch must not be disguised as missing frontend.
	calls := 0
	r.ValidateLayout = func(*parser.Instance) error { calls++; return nil }
	r.Bind = func(*ui.Session, *parser.Document) error { calls++; return nil }
	if c, err := preparation.Prepare(r); err == nil || c != nil || calls != 0 {
		t.Fatal("foreign prepared root reached gates/session", err, calls)
	} else {
		var cap *preparation.Diagnostic
		if errors.As(err, &cap) {
			t.Fatal("capabilities checked before fingerprint", err)
		}
	}
}

func TestHiddenPreviewExactNativeCapabilityFacts(t *testing.T) {
	for _, removed := range []preparation.Capability{{preparation.Widget, "svg", 1}, {preparation.Provider, "svg-resource", 1}, {preparation.Layout, "preview-resource", 1}, {preparation.Host, "svg-resource", 1}} {
		r, _ := previewPreparation(t)
		var caps preparation.Capabilities
		for _, c := range r.Capabilities {
			if c != removed {
				caps = append(caps, c)
			}
		}
		r.Capabilities = caps
		c, err := preparation.Prepare(r)
		var diagnostic *preparation.Diagnostic
		if c != nil || !errors.As(err, &diagnostic) || diagnostic.Capability != removed || diagnostic.Path != "page/closed/image" {
			t.Fatal("hidden capability missing", removed, err)
		}
	}
	r, _ := previewPreparation(t)
	r.Previews = nil
	if _, err := preparation.Prepare(r); err == nil {
		t.Fatal("rich capability booleans admitted missing prepared inventory")
	}
	r, _ = previewPreparation(t)
	c, err := preparation.Prepare(r)
	if err != nil {
		t.Fatal(err)
	}
	c.Close()
}
