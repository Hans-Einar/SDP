package previews

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Controller runs on the UI owner. Its commands arrange application conditions;
// they never dispatch events or simulate a pointer/key/scroll/dialog gesture.
type Controller struct {
	Fixture                       *Fixture
	Host                          *fynehost.DocumentHost
	Sequence                      uint64
	Resize                        func(layout.Size) error
	ParentHide, ParentShow, Close func() error
	Log                           func(string, any)
	held                          *fynehost.Bundle
	failTicket                    bool
}

func (c *Controller) InstallHooks() {
	c.Fixture.PrepareResources = func(ui.Snapshot) error {
		if c.failTicket {
			c.failTicket = false
			return fmt.Errorf("injected preview resource-ticket failure")
		}
		return nil
	}
}
func (c *Controller) Pending() []ui.LoadRequest {
	out := []ui.LoadRequest{}
	if b := c.Host.Current(); b != nil {
		for _, v := range b.Session.Snapshot().Collections {
			if v.Request != nil {
				out = append(out, *v.Request)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target.Handle.Path < out[j].Target.Handle.Path })
	return out
}
func (c *Controller) Inspect() map[string]any {
	out := map[string]any{}
	b := c.Host.Current()
	if b != nil {
		out = b.Inspect()
	}
	out["actionCalls"], out["changeCounts"], out["rendererCalls"] = c.Fixture.Counts()
	out["domain"] = c.Fixture.Domain()
	out["pending"] = c.Pending()
	out["candidateSequence"] = c.Sequence
	out["candidateHeld"] = c.held != nil
	out["sourceSHA256"] = digest([]byte(c.Fixture.Source()))
	out["actionsSHA256"] = digest([]byte(Actions))
	f := c.Fixture
	f.mu.Lock()
	defer f.mu.Unlock()
	out["stagedRevision"] = f.revision
	out["sourceRequests"] = len(f.requests)
	if b != nil {
		out["publishedIdentity"] = f.requests[b.Sequence]
	}
	if c.held != nil {
		out["heldIdentity"] = f.requests[c.held.Sequence]
	}
	loans := map[string]string{}
	for p, v := range f.loans {
		loans[p] = digest(v)
	}
	out["loanSHA256"] = loans
	return out
}
func (c *Controller) Prepare() error {
	if c.held != nil {
		return fmt.Errorf("candidate already held; publish or abandon first")
	}
	c.Sequence++
	r, err := c.Fixture.Request(c.Sequence)
	if err != nil {
		return err
	}
	c.held, err = c.Host.Prepare(r)
	if err == nil {
		a, b, renderers := c.Fixture.Counts()
		if c.Log != nil {
			c.Log("prepared", map[string]any{"candidateSequence": c.Sequence, "sourceSHA256": digest([]byte(c.Fixture.Source())), "actionCalls": a, "changeCounts": b, "rendererBaseline": renderers})
		}
	}
	return err
}
func (c *Controller) Publish() error {
	if c.held == nil {
		return fmt.Errorf("no held candidate")
	}
	b := c.held
	c.held = nil
	return c.Host.Commit(b)
}
func (c *Controller) Abandon() {
	if c.held != nil {
		c.held.Close()
		c.held = nil
	}
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
	switch p[0] {
	case "state", "reload", "prepare", "publish", "abandon", "parent-hide", "parent-show", "close":
		if err := need(1); err != nil {
			return err
		}
		switch p[0] {
		case "state":
			if c.Log != nil {
				c.Log("state", c.Inspect())
			}
			return nil
		case "prepare":
			return c.Prepare()
		case "publish":
			return c.Publish()
		case "reload":
			if err := c.Prepare(); err != nil {
				return err
			}
			return c.Publish()
		case "abandon":
			c.Abandon()
			return nil
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
	case "resource", "renderer", "renderer-revision", "policy", "document", "binding":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.configure(p[0], p[1], p[2])
	case "mutate-input", "mutate-return":
		if err := need(2); err != nil {
			return err
		}
		path, ok := Targets[p[1]]
		if !ok {
			return fmt.Errorf("unknown alias %q", p[1])
		}
		f := c.Fixture
		f.mu.Lock()
		defer f.mu.Unlock()
		buffers := f.loans
		if p[0] == "mutate-return" {
			buffers = f.returns
		}
		v := buffers[path]
		if len(v) == 0 {
			return fmt.Errorf("no observed caller buffer for %s", p[1])
		}
		v[0] = '!'
		return nil
	case "action":
		if err := need(3); err != nil {
			return err
		}
		return c.Fixture.NextAction(p[1], p[2])
	case "accept":
		if err := need(2); err != nil {
			return err
		}
		return c.Fixture.NextAction("AcceptDetail", p[1])
	case "fail":
		if err := need(2); err != nil {
			return err
		}
		if p[1] != "resource-ticket" {
			return fmt.Errorf("unknown failure")
		}
		c.failTicket = true
		return nil
	case "caption":
		if err := need(3); err != nil {
			return err
		}
		if !isSVG(p[1]) {
			return fmt.Errorf("caption requires SVG alias")
		}
		caption := ""
		switch p[2] {
		case "short":
			caption = "Changed caption"
		case "long":
			caption = "A deliberately long optional chart caption whose complete accessibility description remains independent"
		case "empty":
		default:
			return fmt.Errorf("unknown caption slot")
		}
		return c.Host.Mutate(func(s *ui.Session) error {
			w, ok := s.Widget(Targets[p[1]])
			if !ok {
				return fmt.Errorf("missing SVG")
			}
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: ui.Label, Value: ui.Text(caption)}})
		})
	case "resize":
		if err := need(3); err != nil {
			return err
		}
		w, e := dimension(p[1])
		if e != nil {
			return e
		}
		h, e := dimension(p[2])
		if e != nil {
			return e
		}
		if c.Resize == nil {
			return fmt.Errorf("missing resize hook")
		}
		return c.Resize(layout.Size{W: w, H: h})
	}
	return fmt.Errorf("unknown fixture command %q", p[0])
}
func dimension(s string) (float64, error) {
	v, e := strconv.ParseFloat(s, 64)
	if e != nil || math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 || v > 32768 {
		return 0, fmt.Errorf("expected finite dimension in (0,32768]")
	}
	return v, nil
}
func isSVG(alias string) bool { return alias == "Hero" || alias == "ReportSVG" || alias == "DetailSVG" }
func isMD(alias string) bool {
	return alias == "Prose" || alias == "Tail" || alias == "ReportA" || alias == "ReportB" || alias == "DetailDoc"
}
func isReport(alias string) bool { return alias == "ReportA" || alias == "ReportB" }
func (f *Fixture) configure(kind, alias, value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch kind {
	case "resource":
		if !isSVG(alias) {
			return fmt.Errorf("resource requires SVG alias")
		}
		if value != "missing" {
			if _, err := ResourceSlot(value); err != nil {
				return err
			}
		}
		f.resources[alias] = value
	case "renderer":
		if !isReport(alias) {
			return fmt.Errorf("renderer requires ReportA or ReportB")
		}
		switch value {
		case "good", "error", "malformed", "unsafe", "missing":
		default:
			return fmt.Errorf("unknown renderer mode")
		}
		f.rendererModes[alias] = value
	case "renderer-revision":
		if !isReport(alias) || (value != "r1" && value != "r2") {
			return fmt.Errorf("renderer revision requires ReportA/B and r1/r2")
		}
		f.rendererRevisions[alias] = value
	case "policy":
		if (!isSVG(alias) && !isMD(alias)) || (value != "label" && value != "reject") {
			return fmt.Errorf("invalid preview policy")
		}
		f.policies[alias] = value
	case "document":
		if !isMD(alias) {
			return fmt.Errorf("document requires Markdown alias")
		}
		switch value {
		case "prose":
			f.documents[alias] = ProseText
		case "diagram":
			f.documents[alias] = DiagramText
		case "html":
			f.documents[alias] = "<div>Unsupported HTML</div>"
		case "image":
			f.documents[alias] = "![Unsupported image](not-fetched.png)"
		default:
			return fmt.Errorf("unknown document slot")
		}
	case "binding":
		valid := false
		if isSVG(alias) && f.resources[alias] != "missing" {
			switch value {
			case "digest", "source", "provider", "extra", "wrong-kind":
				valid = true
			}
		}
		if isReport(alias) && f.rendererModes[alias] != "missing" {
			switch value {
			case "provider", "revision", "nil", "extra", "wrong-kind", "unused":
				valid = true
			}
		}
		if !valid {
			return fmt.Errorf("invalid binding fault/target")
		}
		f.faults[alias] = value
	default:
		return fmt.Errorf("unknown condition")
	}
	f.revision++
	return nil
}
