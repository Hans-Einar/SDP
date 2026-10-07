package preparation_test

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"testing"
)

func request(t *testing.T) preparation.Request {
	t.Helper()
	d, _, e := parser.Compile(`sdui 0.2; ref: api "unopened.sdl"; page=[ok=button("Go",callback=api.Go.@invoke)];`)
	if e != nil {
		t.Fatal(e)
	}
	return preparation.Request{Document: d, Entry: "page", SessionID: "candidate", SourceRevision: "revision", Mode: preparation.Prototype, Capabilities: admission.Capabilities(), ValidateLayout: func(n *parser.Instance) error { return admission.Layout(n, layout.Size{W: 1280, H: 800}) }}
}
func TestPreparationFailures(t *testing.T) {
	for _, kind := range []string{"document", "entry", "profile", "mode", "revision", "layout", "geometry", "missing", "unbound", "error", "closed", "stale"} {
		t.Run(kind, func(t *testing.T) {
			r := request(t)
			var captured *ui.Session
			switch kind {
			case "document":
				r.Document = nil
			case "entry":
				r.Entry = "missing"
			case "profile":
				r.Document.Profile = "sdui/0.3"
			case "mode":
				r.Mode = ""
			case "revision":
				r.SourceRevision = ""
			case "layout":
				r.ValidateLayout = nil
			case "geometry":
				r.ValidateLayout = func(*parser.Instance) error { return errors.New("bad geometry") }
			default:
				r.Mode = preparation.Connected
			}
			if kind != "missing" {
				r.Bind = func(s *ui.Session, _ *parser.Document) error {
					captured = s
					switch kind {
					case "error":
						return errors.New("bad binding")
					case "closed":
						s.Close()
					case "stale":
						return s.InvalidateEvents()
					}
					return nil
				}
			}
			c, e := preparation.Prepare(r)
			if e == nil || c != nil {
				t.Fatal("accepted invalid candidate")
			}
			if captured != nil && !captured.Closed() {
				t.Fatal("failed candidate not disposed")
			}
			if r.Mode != preparation.Connected && captured != nil {
				t.Fatal("preflight called binding")
			}
		})
	}
}
func TestPrototypeOwnershipAndAdmission(t *testing.T) {
	r := request(t)
	calls := 0
	r.Bind = func(*ui.Session, *parser.Document) error { calls++; return errors.New("must not bind") }
	c, e := preparation.Prepare(r)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	if calls != 0 || c.Unbound != 1 {
		t.Fatal(c.Unbound, calls)
	}
	if c.Admit("revision") != nil || c.Admit("changed") == nil {
		t.Fatal("source admission")
	}
	if e = c.Session.InvalidateEvents(); e != nil {
		t.Fatal(e)
	}
	if c.Admit("revision") == nil {
		t.Fatal("stale admitted")
	}
	c.Close()
	if c.Admit("revision") == nil {
		t.Fatal("closed admitted")
	}
}
func TestCapabilitiesInspectHiddenNodesAndPreserveOrigin(t *testing.T) {
	for _, kind := range []string{"version", "widget", "node", "scroll-x", "scroll-y", "provider", "host"} {
		t.Run(kind, func(t *testing.T) {
			r := request(t)
			roots, e := parser.Normalize(r.Document)
			if e != nil {
				t.Fatal(e)
			}
			root := roots["page"]
			n := root.Rows[0][0]
			n.Layout["visible"] = false
			caps := admission.Capabilities()
			switch kind {
			case "version":
				for i := range caps {
					if caps[i].Dimension == preparation.Widget && caps[i].ID == "button" {
						caps[i].Major = 2
					}
				}
			case "widget":
				n.Widget = "tree"
				caps = append(caps, preparation.Capability{Dimension: preparation.Widget, ID: "tree", Major: 1})
			case "node":
				n.Kind = "tabs"
			case "scroll-x":
				n.Layout["overflow-x"] = "scroll"
			case "scroll-y":
				n.Layout["overflow-y"] = "scroll"
			case "provider":
				n.Widget = "svg"
				for i := range caps {
					if caps[i].Dimension == preparation.Provider {
						caps[i].Major = 2
					}
				}
			case "host":
				for i := range caps {
					if caps[i].Dimension == preparation.Host {
						caps[i].Major = 2
					}
				}
			}
			e = preparation.Check(r.Document.Profile, root, caps)
			var d *preparation.Diagnostic
			if !errors.As(e, &d) || d.Path != n.Path || d.Span != n.Span {
				t.Fatalf("lost diagnostic origin: %#v", e)
			}
		})
	}
}
func TestConnectedRequiresAdapterWithoutCallbacks(t *testing.T) {
	r := request(t)
	d, _, e := parser.Compile(`sdui 0.2; page=[];`)
	if e != nil {
		t.Fatal(e)
	}
	r.Document = d
	r.Mode = preparation.Connected
	if _, e = preparation.Prepare(r); e == nil {
		t.Fatal("inferred connected readiness")
	}
}

func TestDocumentDeterminesCallbackAndDetachedRoot(t *testing.T) {
	r := request(t)
	c, e := preparation.Prepare(r)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	// Editing the caller's AST after preparation cannot change detached bindings.
	r.Document.Definitions[0].Root.Rows[0].Items[0].Arguments = nil
	w, ok := c.Session.Widget("page/ok")
	if !ok || w.Binding.Module != "api" || w.Binding.Object != "Go" {
		t.Fatal("lost document/root consistency", w)
	}
}

func TestMeasuredGeometryFailureDoesNotReachBinder(t *testing.T) {
	r := request(t)
	d, _, e := parser.Compile(`sdui 0.2; page=[button("Too wide") {scale-x=10}] {scale=1};`)
	if e != nil {
		t.Fatal(e)
	}
	r.Document = d
	r.Mode = preparation.Connected
	calls := 0
	r.Bind = func(*ui.Session, *parser.Document) error { calls++; return nil }
	if c, e := preparation.Prepare(r); e == nil || c != nil || calls != 0 {
		t.Fatal("failed geometry reached binding", e, calls)
	}
}

func TestSelectedConnectionsAndUnboundHandle(t *testing.T) {
	r := request(t)
	d, _, e := parser.Compile(`sdui 0.2; ref: api "unused.sdl"; page=[edit=input("Edit")]; other=[edit=input("Other")]; api.One.setHandle(page.edit); api.Two.setHandle(other.edit);`)
	if e != nil {
		t.Fatal(e)
	}
	r.Document = d
	c, e := preparation.Prepare(r)
	if e != nil {
		t.Fatal(e)
	}
	if c.Unbound != 1 {
		t.Fatal("counted another entry")
	}
	c.Close()
	r.Mode = preparation.Connected
	r.Bind = func(_ *ui.Session, selected *parser.Document) error {
		if len(selected.Connections) != 1 || selected.Connections[0].Definition != "page" {
			t.Fatal("unselected connections sent to binder")
		}
		return nil
	}
	if c, e = preparation.Prepare(r); e == nil || c != nil {
		t.Fatal("unmatched connection passed readiness")
	}
}
