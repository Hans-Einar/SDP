package collections

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

// NestedSource keeps stable collection/preview paths. Body and each collection
// are nested scroll owners; the fixed footer preview is a separate sibling.
// Relative sizing is finite: collections overflow their body's height by 60%.
var NestedSource = strings.Replace(strings.Replace(Source,
	`y=fill,overflow-x=scroll,overflow-y=scroll`,
	`scale-y=1.6,overflow-x=scroll,overflow-y=scroll`, 2),
	`> {x=fill,y=fill};`, `> {x=fill,y=fill,overflow-y=scroll};`, 1)

const WindowTitle = "SDUI WCI1 Collections"

// Controller is a fixture-only control channel. Product hosts do not parse these
// commands. All methods run on the UI owner goroutine; providers remain explicit
// barrier-controlled requests through Fixture.Complete.
type Controller struct {
	Fixture  *Fixture
	Host     *fynehost.DocumentHost
	Source   string
	Empty    bool
	Sequence uint64
	Resize   func(layout.Size) error
	Capture  func(string) error
	Close    func()
	Log      func(string, any)
}

func (c *Controller) Inspect() map[string]any {
	b := c.Host.Current()
	if b == nil {
		return map[string]any{"actionCalls": c.Fixture.Calls(), "pending": c.Fixture.Pending()}
	}
	state := b.Inspect()
	state["actionCalls"] = c.Fixture.Calls()
	state["pending"] = c.Fixture.Pending()
	state["candidateSequence"] = c.Sequence
	state["sourceSHA256"] = fmt.Sprintf("%x", sha256.Sum256([]byte(c.Source)))
	state["actionsSHA256"] = fmt.Sprintf("%x", sha256.Sum256([]byte(Actions)))
	state["nested"] = c.Source == NestedSource
	return state
}
func (c *Controller) emit(event string, value any) {
	if c.Log != nil {
		c.Log(event, value)
	}
}

// Run rejects malformed commands before changing fixture/runtime state. A
// successful complete command releases one named provider barrier; the later
// load-return/state logs establish delivery, not the command acknowledgement.
func (c *Controller) Run(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return fmt.Errorf("empty fixture command")
	}
	need := func(count int) error {
		if len(parts) != count {
			return fmt.Errorf("invalid arguments for %s", parts[0])
		}
		return nil
	}
	switch parts[0] {
	case "state", "pending":
		if err := need(1); err != nil {
			return err
		}
		if parts[0] == "state" {
			c.emit("state", c.Inspect())
		} else {
			c.emit("pending", c.Fixture.Pending())
		}
		return nil
	case "reload", "bad-reload":
		if err := need(1); err != nil {
			return err
		}
		kind := ""
		if parts[0] == "bad-reload" {
			kind = "layout"
		}
		return c.reload(kind)
	case "fail":
		if err := need(2); err != nil {
			return err
		}
		return c.reload(parts[1])
	case "complete":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.Complete(parts[1], parts[2])
	case "cancel":
		if err := need(2); err != nil {
			return err
		}
		path, err := c.path(parts[1])
		if err != nil {
			return err
		}
		return c.Host.CancelCollection(path)
	case "resize":
		if err := need(3); err != nil {
			return err
		}
		w, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return err
		}
		h, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return err
		}
		if math.IsNaN(w) || math.IsNaN(h) || math.IsInf(w, 0) || math.IsInf(h, 0) || w <= 0 || h <= 0 || w > 32768 || h > 32768 {
			return fmt.Errorf("resize dimensions must be finite in (0,32768]")
		}
		if c.Resize == nil {
			return fmt.Errorf("fixture resize callback missing")
		}
		return c.Resize(layout.Size{W: w, H: h})
	case "shrink", "reorder", "reset":
		count := -1
		if parts[0] == "shrink" {
			if err := need(3); err != nil {
				return err
			}
			n, err := strconv.Atoi(parts[2])
			if err != nil || n < 0 {
				return fmt.Errorf("shrink requires nonnegative count")
			}
			count = n
		} else if err := need(2); err != nil {
			return err
		}
		return c.replace(parts[1], parts[0], count)
	case "hide", "show", "disable", "enable":
		if err := need(2); err != nil {
			return err
		}
		path, err := c.path(parts[1])
		if err != nil {
			return err
		}
		return c.Host.Mutate(func(s *ui.Session) error {
			w, ok := s.Widget(path)
			if !ok {
				return fmt.Errorf("missing fixture widget %s", path)
			}
			property := ui.Visible
			if parts[0] == "disable" || parts[0] == "enable" {
				property = ui.Enabled
			}
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: property, Value: ui.Bool(parts[0] == "show" || parts[0] == "enable")}})
		})
	case "capture":
		if err := need(2); err != nil {
			return err
		}
		if c.Capture == nil {
			return fmt.Errorf("fixture capture callback missing")
		}
		return c.Capture(parts[1])
	case "close":
		if err := need(1); err != nil {
			return err
		}
		if c.Close == nil {
			return fmt.Errorf("fixture close callback missing")
		}
		c.Close()
		return nil
	default:
		return fmt.Errorf("unknown fixture command %q", parts[0])
	}
}
func (c *Controller) path(name string) (string, error) {
	switch name {
	case "tree":
		name = "page/body/nav"
	case "list":
		name = "page/body/entries"
	case "preview":
		name = "page/footer/preview"
	}
	if c.Host.Current() == nil {
		return "", fmt.Errorf("no published fixture")
	}
	if _, ok := c.Host.Current().Session.Widget(name); !ok {
		return "", fmt.Errorf("unknown fixture widget %q", name)
	}
	return name, nil
}
func (c *Controller) replace(name, operation string, count int) error {
	path, err := c.path(name)
	if err != nil {
		return err
	}
	return c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget(path)
		state, ok := s.Collection(w.Handle)
		if !ok {
			return fmt.Errorf("not a collection: %s", path)
		}
		data := state.Data
		switch operation {
		case "shrink":
			if count > len(data.Items) {
				return fmt.Errorf("shrink count exceeds current item count")
			}
			data.Items = data.Items[:count]
		case "reorder":
			for i, j := 0, len(data.Items)-1; i < j; i, j = i+1, j-1 {
				data.Items[i], data.Items[j] = data.Items[j], data.Items[i]
			}
		case "reset":
			data = c.Fixture.Providers(false)[state.InstancePath].Initial
		}
		target, err := s.Target(w.Handle, "")
		if err != nil {
			return err
		}
		return s.ReplaceCollection(target, data)
	})
}
func (c *Controller) reload(failure string) error {
	switch failure {
	case "", "profile", "provider", "binding", "layout", "resource", "guard", "stale":
	default:
		return fmt.Errorf("unknown failure stage %q", failure)
	}
	c.Sequence++
	source := c.Source
	if failure == "profile" {
		source = strings.Replace(source, "sdui 0.3;", "sdui 0.9;", 1)
	}
	if failure == "layout" {
		source = strings.Replace(source, `] {x=fill,y=fill};`, `] {scale-x=0.001,y=fill};`, 1)
	}
	r, err := c.Fixture.Request(source, c.Sequence, c.Empty)
	if err != nil {
		return err
	}
	switch failure {
	case "provider":
		p := r.Providers["page/body/nav"]
		p.Epoch++
		p.RootLoaded = true
		p.Initial = ui.CollectionData{Items: []ui.CollectionItem{{Kind: ui.Row, Label: "Missing identity", ChildrenLoaded: true}}}
		r.Providers["page/body/nav"] = p
	case "binding":
		r.Bind = func(s *ui.Session, d *parser.Document) error {
			_, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{}, map[string]bridge.Plan{"nav.Activate": {Inputs: map[string]bridge.Source{"ItemId": {EventField: bridge.CollectionItemID}}, OutputField: "Preview"}}, nil)
			return err
		}
	case "resource":
		r.PrepareResources = func(ui.Snapshot) error { return fmt.Errorf("injected fixture resource preparation failure") }
	case "guard":
		original := r.Guard
		checks := 0
		r.Guard = func() error {
			if err := original(); err != nil {
				return err
			}
			checks++
			if checks == 2 {
				return fmt.Errorf("injected final publication guard failure")
			}
			return nil
		}
	}
	if failure != "stale" {
		return c.Host.Adopt(r)
	}
	candidate, err := c.Host.Prepare(r)
	if err != nil {
		return err
	}
	if err = c.Host.Mutate(func(s *ui.Session) error {
		w, ok := s.Widget("page/footer/preview")
		if !ok {
			return fmt.Errorf("missing preview")
		}
		return s.Draft(w.Handle, "New draft during detached preparation")
	}); err != nil {
		candidate.Close()
		return err
	}
	return c.Host.Commit(candidate)
}
