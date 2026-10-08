package panes

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
)

// Controller changes fixture conditions and trusted application state only.
// Pointer/key user activation must come through the native input route.
type Controller struct {
	Fixture  *Fixture
	Host     *fynehost.DocumentHost
	Source   string
	Sequence uint64
	Resize   func(layout.Size) error
	Close    func()
	Log      func(string, any)
}

func (c *Controller) Inspect() map[string]any {
	state := map[string]any{}
	if b := c.Host.Current(); b != nil {
		state = b.Inspect()
	}
	state["actionCalls"] = c.Fixture.Calls()
	state["pending"] = c.Fixture.Pending()
	state["candidateSequence"] = c.Sequence
	state["sourceSHA256"] = fmt.Sprintf("%x", sha256.Sum256([]byte(c.Source)))
	state["actionsSHA256"] = fmt.Sprintf("%x", sha256.Sum256([]byte(Actions)))
	return state
}
func (c *Controller) emit(kind string, v any) {
	if c.Log != nil {
		c.Log(kind, v)
	}
}
func (c *Controller) Run(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return fmt.Errorf("empty fixture command")
	}
	need := func(n int) error {
		if len(parts) != n {
			return fmt.Errorf("invalid arguments for %s", parts[0])
		}
		return nil
	}
	switch parts[0] {
	case "state", "pending", "reload", "resource-error", "close":
		if err := need(1); err != nil {
			return err
		}
		switch parts[0] {
		case "state":
			c.emit("state", c.Inspect())
		case "pending":
			c.emit("pending", c.Fixture.Pending())
		case "reload":
			return c.reload("")
		case "resource-error":
			c.Fixture.FailResource()
		case "close":
			if c.Close == nil {
				return fmt.Errorf("close hook missing")
			}
			c.Close()
		}
		return nil
	case "action":
		if err := need(2); err != nil {
			return err
		}
		return c.Fixture.NextAction(parts[1])
	case "complete":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.Complete(parts[1], parts[2])
	case "cancel":
		if err := need(1); err != nil {
			return err
		}
		return c.Host.CancelCollection(TreePath)
	case "fail":
		if err := need(2); err != nil {
			return err
		}
		return c.reload(parts[1])
	case "resize":
		if err := need(3); err != nil {
			return err
		}
		w, err := finite(parts[1])
		if err != nil {
			return err
		}
		h, err := finite(parts[2])
		if err != nil {
			return err
		}
		if w <= 0 || h <= 0 || w > 32768 || h > 32768 {
			return fmt.Errorf("resize dimensions outside (0,32768]")
		}
		if c.Resize == nil {
			return fmt.Errorf("resize hook missing")
		}
		return c.Resize(layout.Size{W: w, H: h})
	case "select", "ratio", "collapse", "restore", "hide", "show", "disable", "enable":
		count := 2
		if parts[0] == "restore" {
			count = 1
		}
		if err := need(count); err != nil {
			return err
		}
		return c.Host.Mutate(func(s *ui.Session) error {
			path := SplitPath
			if parts[0] == "select" {
				path = TabsPath
			}
			switch parts[0] {
			case "hide", "show", "disable", "enable":
				switch parts[1] {
				case "workspace":
					path = TabsPath
				case "overview", "notes":
					path = TabsPath + "/" + parts[1]
				default:
					return fmt.Errorf("expected workspace, overview or notes")
				}
			}
			handle, ok := s.Pane(path)
			if !ok {
				return fmt.Errorf("missing pane %s", path)
			}
			switch parts[0] {
			case "select":
				return s.SelectPage(handle, parts[1])
			case "ratio":
				value, err := finite(parts[1])
				if err != nil {
					return err
				}
				return s.SetSplitProportion(handle, value)
			case "collapse":
				switch parts[1] {
				case "first":
					return s.CollapseSplit(handle, ui.SplitFirst)
				case "second":
					return s.CollapseSplit(handle, ui.SplitSecond)
				default:
					return fmt.Errorf("expected first or second")
				}
			case "restore":
				return s.RestoreSplit(handle)
			default:
				property := ui.Visible
				if parts[0] == "disable" || parts[0] == "enable" {
					property = ui.Enabled
				}
				return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: handle, Property: property, Value: ui.Bool(parts[0] == "show" || parts[0] == "enable")}})
			}
		})
	default:
		return fmt.Errorf("unknown fixture command %q", parts[0])
	}
}
func finite(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("expected finite number")
	}
	return v, nil
}
func (c *Controller) reload(failure string) error {
	switch failure {
	case "", "profile", "binding", "layout", "resource", "guard", "stale":
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
	r, err := c.Fixture.Request(source, c.Sequence)
	if err != nil {
		return err
	}
	switch failure {
	case "binding":
		r.Bind = func(s *ui.Session, d *parser.Document) error {
			_, err := bridge.Bind(context.Background(), s, d, nil, Plans(), nil)
			return err
		}
	case "resource":
		r.PrepareResources = func(ui.Snapshot) error { return fmt.Errorf("injected reload resource failure") }
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
		w, ok := s.Widget(PreviewPath)
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
