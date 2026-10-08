package values

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/numeric"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
)

type Controller struct {
	Fixture                       *Fixture
	Host                          *fynehost.DocumentHost
	Source                        string
	Sequence                      uint64
	Resize                        func(layout.Size) error
	ParentHide, ParentShow, Close func() error
	Log                           func(string, any)
}

func (c *Controller) Inspect() map[string]any {
	state := map[string]any{}
	if b := c.Host.Current(); b != nil {
		state = b.Inspect()
	}
	state["actionCalls"], state["changeCounts"] = c.Fixture.Counts()
	state["domain"] = c.Fixture.Domain()
	state["pending"] = c.Pending()
	state["candidateSequence"] = c.Sequence
	state["sourceSHA256"] = fmt.Sprintf("%x", sha256.Sum256([]byte(c.Source)))
	state["actionsSHA256"] = fmt.Sprintf("%x", sha256.Sum256([]byte(Actions)))
	return state
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
	case "state", "reload", "parent-hide", "parent-show", "close":
		if err := need(1); err != nil {
			return err
		}
		switch p[0] {
		case "state":
			if c.Log != nil {
				c.Log("state", c.Inspect())
			}
			return nil
		case "reload":
			return c.reload("")
		case "parent-hide":
			if c.ParentHide != nil {
				return c.ParentHide()
			}
		case "parent-show":
			if c.ParentShow != nil {
				return c.ParentShow()
			}
		case "close":
			if c.Close != nil {
				return c.Close()
			}
		}
		return fmt.Errorf("missing lifecycle hook")
	case "action":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.NextAction(p[1], p[2])
	case "accept":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.NextAccept(p[1], p[2])
	case "observe":
		if err := need(3); err != nil {
			return err
		}
		path, ok := Targets[p[1]]
		if !ok || p[2] != "draft-conflict" {
			return fmt.Errorf("expected observe TARGET draft-conflict")
		}
		c.Fixture.mu.Lock()
		c.Fixture.observers[path] = p[2]
		c.Fixture.mu.Unlock()
		return nil
	case "options":
		if err := need(3); err != nil {
			return err
		}
		if p[1] != "Mode" && p[1] != "FormMode" {
			return fmt.Errorf("expected Mode or FormMode")
		}
		switch p[2] {
		case "initial", "removed", "disabled", "reordered":
		default:
			return fmt.Errorf("invalid options condition")
		}
		return c.Host.Mutate(func(s *ui.Session) error { return c.replace(s, Targets[p[1]], OptionSet(p[2])) })
	case "set", "enable", "disable", "readonly", "writable":
		n := 2
		if p[0] == "set" {
			n = 3
		}
		if err := need(n); err != nil {
			return err
		}
		path, ok := Targets[p[1]]
		if !ok {
			return fmt.Errorf("unknown target alias %q", p[1])
		}
		return c.Host.Mutate(func(s *ui.Session) error {
			w, ok := s.Widget(path)
			if !ok {
				return fmt.Errorf("missing field %s", path)
			}
			field, ok := s.Field(w.Handle)
			if !ok {
				return fmt.Errorf("missing field state")
			}
			update := ui.Update{Handle: w.Handle, Property: ui.Enabled, Value: ui.Bool(p[0] == "enable")}
			if p[0] == "readonly" || p[0] == "writable" {
				update.Property = ui.ReadOnly
				update.ExpectedValueRevision = field.Target.ValueRevision
				update.ExpectedDraftRevision = field.Target.DraftRevision
				update.ExpectedOptionGeneration = field.Target.OptionGeneration
				update.Value = ui.Bool(p[0] == "readonly")
			}
			if p[0] == "set" {
				value, err := programmaticValue(field, p[2])
				if err != nil {
					return err
				}
				update = ui.Update{Handle: w.Handle, Property: ui.AcceptedValue, Value: value, ExpectedValueRevision: field.Target.ValueRevision, ExpectedDraftRevision: field.Target.DraftRevision, ExpectedOptionGeneration: field.Target.OptionGeneration}
			}
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{update})
		})
	case "fail":
		if err := need(2); err != nil {
			return err
		}
		return c.reload(p[1])
	case "resize":
		if err := need(3); err != nil {
			return err
		}
		w, err := dimension(p[1])
		if err != nil {
			return err
		}
		h, err := dimension(p[2])
		if err != nil {
			return err
		}
		if c.Resize == nil {
			return fmt.Errorf("missing resize hook")
		}
		return c.Resize(layout.Size{W: w, H: h})
	}
	return fmt.Errorf("unknown fixture command %q", p[0])
}
func programmaticValue(field ui.FieldState, raw string) (ui.Value, error) {
	switch field.Target.Handle.Kind {
	case "checkbox":
		if raw == "true" || raw == "false" {
			return ui.Bool(raw == "true"), nil
		}
		return ui.Value{}, fmt.Errorf("expected true or false")
	case "select":
		if raw == "empty" {
			raw = ""
		}
		return ui.Choice(ui.ItemID(raw)), nil
	case "slider", "number":
		if field.Numeric == nil {
			return ui.Value{}, fmt.Errorf("missing numeric constraints")
		}
		n := field.Numeric
		grid, err := numeric.NewGrid(n.Min, n.Max, n.Step)
		if err != nil {
			return ui.Value{}, err
		}
		v, err := grid.Parse(raw)
		return ui.Numeric(v), err
	}
	return ui.Value{}, fmt.Errorf("unsupported fixture target")
}
func dimension(raw string) (float64, error) {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 || v > 32768 {
		return 0, fmt.Errorf("expected finite dimension in (0,32768]")
	}
	return v, nil
}
func (f *Fixture) NextAction(name, mode string) error {
	switch name {
	case "AcceptFlag", "AcceptLevel", "AcceptMode", "AcceptCount", "ChildCount", "Load", "SaveText":
	default:
		return fmt.Errorf("unknown action %q", name)
	}
	switch mode {
	case "success", "error", "malformed", "draft-conflict":
	case "options-conflict":
		if name != "AcceptMode" {
			return fmt.Errorf("options-conflict requires AcceptMode")
		}
	default:
		return fmt.Errorf("unknown action outcome")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.outcomes[name] = mode
	return nil
}
func (f *Fixture) NextAccept(form, mode string) error {
	if form != "mixed" && form != "text" {
		return fmt.Errorf("expected mixed or text")
	}
	switch mode {
	case "true", "false", "error", "draft-conflict":
	case "options-conflict":
		if form != "mixed" {
			return fmt.Errorf("options-conflict requires mixed")
		}
	case "malformed":
		if form != "text" {
			return fmt.Errorf("malformed requires text")
		}
	default:
		return fmt.Errorf("unknown accept outcome")
	}
	if form == "text" {
		form = "TextSave"
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.outcomes[form] = mode
	return nil
}
func (c *Controller) replace(s *ui.Session, path string, options []ui.ChoiceOption) error {
	w, ok := s.Widget(path)
	if !ok {
		return fmt.Errorf("missing choice")
	}
	field, ok := s.Field(w.Handle)
	if !ok {
		return fmt.Errorf("missing choice state")
	}
	if err := s.ReplaceChoices(w.Handle, field.Target.OptionGeneration, options); err != nil {
		return err
	}
	c.Fixture.storeChoices(path, options)
	return nil
}

// Hooks run inside actual callbacks on the UI owner; no control command invokes them.
func (c *Controller) InstallHooks() {
	c.Fixture.Replace = func(path string) error { return c.replace(c.Host.Current().Session, path, c.Fixture.Choices()[path]) }
	c.Fixture.Conflict = func(path string) error {
		s := c.Host.Current().Session
		w, ok := s.Widget(path)
		if !ok {
			return fmt.Errorf("missing conflict target")
		}
		if w.Handle.Kind == "input" {
			return s.Draft(w.Handle, w.Draft+" newer")
		}
		field, ok := s.Field(w.Handle)
		if !ok {
			return fmt.Errorf("missing conflict field")
		}
		switch w.Handle.Kind {
		case "checkbox":
			_, err := s.EditField(w.Handle, s.Revision, ui.Bool(!field.Proposed.Bool))
			return err
		case "select":
			for _, o := range field.Options {
				if o.Enabled && o.ID != field.Proposed.OptionID {
					target, err := s.Option(w.Handle, o.ID)
					if err != nil {
						return err
					}
					_, err = s.ChooseOption(target)
					return err
				}
			}
			return fmt.Errorf("no alternate option")
		case "number", "slider":
			n := field.Numeric
			grid, err := numeric.NewGrid(n.Min, n.Max, n.Step)
			if err != nil {
				return err
			}
			tick, err := grid.Tick(field.Proposed.Number)
			if err != nil {
				return err
			}
			if tick < grid.LastTick() {
				tick++
			} else if tick > 0 {
				tick--
			}
			_, err = s.EditTick(w.Handle, s.Revision, tick)
			return err
		}
		return fmt.Errorf("unsupported conflict field")
	}
}
func (c *Controller) reload(failure string) error {
	switch failure {
	case "", "profile", "binding", "resource", "guard", "stale":
	default:
		return fmt.Errorf("unknown reload failure")
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
		r.PrepareResources = func(ui.Snapshot) error { return fmt.Errorf("injected resources failure") }
	case "guard":
		original := r.Guard
		checks := 0
		r.Guard = func() error {
			if err := original(); err != nil {
				return err
			}
			checks++
			if checks == 2 {
				return fmt.Errorf("injected final guard failure")
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
		w, _ := s.Widget(PreviewPath)
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: ui.Label, Value: ui.Text("Changed during preparation")}})
	}); err != nil {
		candidate.Close()
		return err
	}
	return c.Host.Commit(candidate)
}

// Pending reports the runtime's actual outstanding collection requests, including
// after host teardown; this fixture does not create a synthetic provider queue.
func (c *Controller) Pending() []ui.LoadRequest {
	out := []ui.LoadRequest{}
	if b := c.Host.Current(); b != nil {
		for _, collection := range b.Session.Snapshot().Collections {
			if collection.Request != nil {
				out = append(out, *collection.Request)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target.Handle.Path < out[j].Target.Handle.Path })
	return out
}
