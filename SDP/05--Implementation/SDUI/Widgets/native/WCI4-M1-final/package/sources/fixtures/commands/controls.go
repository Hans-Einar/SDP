package commands

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

type Controller struct {
	Fixture                *Fixture
	Host                   *fynehost.DocumentHost
	Source                 string
	Sequence               uint64
	Resize                 func(layout.Size) error
	ParentHide, ParentShow func() error
	Close                  func() error
	Log                    func(string, any)
}

func (c *Controller) Inspect() map[string]any {
	state := map[string]any{}
	if b := c.Host.Current(); b != nil {
		state = b.Inspect()
	}
	state["actionCalls"] = c.Fixture.Calls()
	state["domain"] = c.Fixture.Domain()
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
	p := strings.Fields(line)
	if len(p) == 0 {
		return fmt.Errorf("empty fixture command")
	}
	need := func(n int) error {
		if len(p) != n {
			return fmt.Errorf("invalid arguments for %s", p[0])
		}
		return nil
	}
	switch p[0] {
	case "state", "pending", "reload", "replace-items", "cancel-load", "parent-hide", "parent-show", "close":
		if err := need(1); err != nil {
			return err
		}
		switch p[0] {
		case "state":
			c.emit("state", c.Inspect())
		case "pending":
			c.emit("pending", c.Fixture.Pending())
		case "reload":
			return c.reload("")
		case "replace-items":
			return c.Host.Mutate(func(s *ui.Session) error {
				w, ok := s.Widget(ItemsPath)
				if !ok {
					return fmt.Errorf("missing items")
				}
				target, err := s.Target(w.Handle, "")
				if err != nil {
					return err
				}
				return s.ReplaceCollection(target, c.Fixture.Providers()[ItemsPath].Initial)
			})
		case "cancel-load":
			return c.Host.CancelCollection(ItemsPath)
		case "parent-hide":
			if c.ParentHide == nil {
				return fmt.Errorf("parent hide hook missing")
			}
			return c.ParentHide()
		case "parent-show":
			if c.ParentShow == nil {
				return fmt.Errorf("parent show hook missing")
			}
			return c.ParentShow()
		case "close":
			if c.Close == nil {
				return fmt.Errorf("parent close hook missing")
			}
			return c.Close()
		}
		return nil
	case "action":
		if err := need(3); err != nil {
			return err
		}
		if p[1] == "Save" {
			return fmt.Errorf("use accept for Save outcome")
		}
		return c.Fixture.NextAction(p[1], p[2])
	case "accept":
		if err := need(2); err != nil {
			return err
		}
		return c.Fixture.NextAction("Save", p[1])
	case "complete":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.Complete(p[1], p[2])
	case "fail":
		if err := need(2); err != nil {
			return err
		}
		return c.reload(p[1])
	case "enable", "disable":
		if err := need(2); err != nil {
			return err
		}
		path := ""
		switch p[1] {
		case "Run":
			path = "page/run"
		case "Flag":
			path = "page/flag"
		default:
			return fmt.Errorf("expected Run or Flag")
		}
		return c.Host.Mutate(func(s *ui.Session) error {
			h, ok := s.Command(path)
			if !ok {
				return fmt.Errorf("missing command %s", path)
			}
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: h, Property: ui.Enabled, Value: ui.Bool(p[0] == "enable")}})
		})
	case "resize":
		if err := need(3); err != nil {
			return err
		}
		w, err := boundedSize(p[1])
		if err != nil {
			return err
		}
		h, err := boundedSize(p[2])
		if err != nil {
			return err
		}
		if c.Resize == nil {
			return fmt.Errorf("resize hook missing")
		}
		return c.Resize(layout.Size{W: w, H: h})
	default:
		return fmt.Errorf("unknown fixture command %q", p[0])
	}
}
func boundedSize(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 || v > 32768 {
		return 0, fmt.Errorf("expected finite dimension in (0,32768]")
	}
	return v, nil
}
func (c *Controller) reload(failure string) error {
	switch failure {
	case "", "profile", "binding", "resource", "guard", "stale":
	default:
		return fmt.Errorf("unknown reload failure %q", failure)
	}
	c.Sequence++
	source := c.Source
	if failure == "profile" {
		source = strings.Replace(source, "sdui 0.3;", "sdui 0.9;", 1)
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
				return fmt.Errorf("injected final reload guard failure")
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
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: ui.Label, Value: ui.Text("Changed during preparation")}})
	}); err != nil {
		candidate.Close()
		return err
	}
	return c.Host.Commit(candidate)
}
