package text

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Controller struct {
	Fixture                       *Fixture
	Host                          *fynehost.DocumentHost
	Source                        string
	Sequence                      uint64
	Resize                        func(layout.Size) error
	ParentHide, ParentShow, Close func() error
	Log                           func(string, any)
	rejectEdit                    string // UI-owner condition consumed by actual prospective draft change.
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
		return fmt.Errorf("empty command")
	}
	need := func(n int) error {
		if len(p) != n {
			return fmt.Errorf("invalid arguments for %s", p[0])
		}
		return nil
	}
	target := func() (string, error) {
		path, ok := Targets[p[1]]
		if !ok {
			return "", fmt.Errorf("unknown alias %q", p[1])
		}
		return path, nil
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
			return c.reload(c.Source, "")
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
		if err := need(2); err != nil {
			return err
		}
		return c.Fixture.NextAction("TextSave", p[1])
	case "load":
		if err := need(2); err != nil {
			return err
		}
		v, err := Slot(p[1])
		if err != nil {
			return err
		}
		c.Fixture.mu.Lock()
		c.Fixture.load = v
		c.Fixture.mu.Unlock()
		return nil
	case "reject-edit":
		if err := need(2); err != nil {
			return err
		}
		path, err := target()
		if err != nil {
			return err
		}
		c.rejectEdit = path
		return nil
	case "set", "enable", "disable", "readonly", "writable":
		n := 2
		if p[0] == "set" {
			n = 3
		}
		if err := need(n); err != nil {
			return err
		}
		path, err := target()
		if err != nil {
			return err
		}
		return c.Host.Mutate(func(s *ui.Session) error {
			w, ok := s.Widget(path)
			if !ok {
				return fmt.Errorf("missing widget %s", path)
			}
			field, ok := s.Field(w.Handle)
			if !ok || field.Input == nil {
				return fmt.Errorf("missing extended field")
			}
			u := ui.Update{Handle: w.Handle, Property: ui.Enabled, Value: ui.Bool(p[0] == "enable")}
			if p[0] == "readonly" || p[0] == "writable" {
				u.Property = ui.ReadOnly
				u.Value = ui.Bool(p[0] == "readonly")
				u.ExpectedValueRevision = field.Target.ValueRevision
				u.ExpectedDraftRevision = field.Target.DraftRevision
			}
			if p[0] == "set" {
				value := w.Draft
				same := p[2] == "same"
				if !same {
					var err error
					value, err = Slot(p[2])
					if err != nil {
						return err
					}
				}
				u = ui.Update{Handle: w.Handle, Property: ui.AcceptedValue, Value: ui.Text(value), ExpectedValueRevision: field.Target.ValueRevision, ExpectedDraftRevision: field.Target.DraftRevision, AcceptDraft: same}
			}
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{u})
		})
	case "constraint":
		if err := need(3); err != nil {
			return err
		}
		path, err := target()
		if err != nil {
			return err
		}
		source, err := constrain(c.Source, path, p[2])
		if err != nil {
			return err
		}
		return c.reload(source, "")
	case "fail":
		if err := need(2); err != nil {
			return err
		}
		return c.reload(c.Source, p[1])
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
func (f *Fixture) NextAction(name, mode string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.calls[name]; !ok {
		return fmt.Errorf("unknown SDL action %q", name)
	}
	switch mode {
	case "success", "error", "malformed", "draft-conflict":
	case "true", "false":
		if name != "TextSave" {
			return fmt.Errorf("accept outcome requires TextSave")
		}
	case "non-echo":
		if name == "TextSave" {
			return fmt.Errorf("non-echo requires text result")
		}
	default:
		return fmt.Errorf("unknown action outcome")
	}
	f.outcomes[name] = mode
	return nil
}
func dimension(raw string) (float64, error) {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 || v > 32768 {
		return 0, fmt.Errorf("expected finite dimension in (0,32768]")
	}
	return v, nil
}
func constrain(source, path, mode string) (string, error) {
	key, value := "", ""
	switch mode {
	case "required":
		key, value = "required", "true"
	case "optional":
		key, value = "required", "false"
	case "single":
		key, value = "multiline", "false"
	case "multi":
		key, value = "multiline", "true"
	default:
		return "", fmt.Errorf("unknown constraint")
	}
	name := path[strings.LastIndex(path, "/")+1:]
	start := strings.Index(source, name+"=input(")
	if start < 0 {
		return "", fmt.Errorf("missing source input")
	}
	end := strings.Index(source[start:], ") {x=fill}")
	if end < 0 {
		return "", fmt.Errorf("missing fixture input terminator")
	}
	end += start
	call := source[start:end]
	re := regexp.MustCompile("," + key + "=(true|false)")
	if re.MatchString(call) {
		call = re.ReplaceAllString(call, ","+key+"="+value)
	} else {
		call += "," + key + "=" + value
	}
	return source[:start] + call + source[end:], nil
}

// These hooks run only while actual callbacks or the runtime publication gate run.
func (c *Controller) InstallHooks() {
	c.Fixture.Conflict = func(path string) error {
		s := c.Host.Current().Session
		w, ok := s.Widget(path)
		if !ok {
			return fmt.Errorf("missing conflict target")
		}
		_, err := s.EditField(w.Handle, s.Revision, ui.Text(w.Draft+" newer"))
		return err
	}
	c.Fixture.PrepareResources = func(proposed ui.Snapshot) error {
		if c.rejectEdit == "" || c.Host.Current() == nil {
			return nil
		}
		path := c.rejectEdit
		before, ok := c.Host.Current().Session.Snapshot().Fields[path]
		after, aok := proposed.Fields[path]
		if !ok || !aok || before.Target.Handle != after.Target.Handle || before.Target.ModelRevision != after.Target.ModelRevision || before.Target.DraftRevision == after.Target.DraftRevision || before.Proposed == after.Proposed || before.Accepted != after.Accepted {
			return nil
		}
		c.rejectEdit = ""
		if c.Log != nil {
			c.Log("diagnostic", map[string]any{"kind": "reject-edit", "path": path})
		}
		return fmt.Errorf("injected native edit publication rejection: %s", path)
	}
}
func (c *Controller) reload(source, failure string) error {
	switch failure {
	case "", "profile", "binding", "resource", "guard", "stale":
	default:
		return fmt.Errorf("unknown reload failure")
	}
	c.Sequence++
	candidateSource := source
	if failure == "profile" {
		candidateSource = strings.Replace(source, "sdui 0.3;", "sdui 0.9;", 1)
	}
	r, err := c.Fixture.Request(candidateSource, c.Sequence)
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
		err = c.Host.Adopt(r)
	} else {
		return c.stale(r)
	}
	if err == nil {
		c.Source = source
	}
	return err
}
func (c *Controller) stale(r fynehost.DocumentRequest) error {
	candidate, err := c.Host.Prepare(r)
	if err != nil {
		return err
	}
	if err = c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget(Targets["Preview"])
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: ui.Label, Value: ui.Text("Changed during preparation")}})
	}); err != nil {
		candidate.Close()
		return err
	}
	return c.Host.Commit(candidate)
}
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
