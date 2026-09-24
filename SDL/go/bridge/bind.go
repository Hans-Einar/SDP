// Package bridge resolves explicit SDUI callback/setHandle data against an
// already loaded SDL action runtime. It never opens ref paths or guesses bindings.
package bridge

import (
	"context"
	"fmt"
	"sort"
	"strings"

	uiparser "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

type Source struct {
	Widget  string
	Event   bool
	Literal *sdl.Value
	Context string
}
type Plan struct {
	Inputs                         map[string]Source
	OutputField                    string
	RevisionField, RevisionContext string
}
type Link struct {
	Widget, Symbol, Target string
	UISpan                 uiparser.Span
	SDLSpan                parser.Span
}
type Bridge struct {
	UI      *ui.Session
	Modules map[string]*sdl.Engine
	Plans   map[string]Plan
	Context map[string]sdl.Value
	Links   []Link
	ctx     context.Context
}

func Bind(ctx context.Context, session *ui.Session, document *uiparser.Document, modules map[string]*sdl.Engine, plans map[string]Plan, state map[string]sdl.Value) (*Bridge, error) {
	copiedModules := map[string]*sdl.Engine{}
	for k, v := range modules {
		copiedModules[k] = v
	}
	copiedPlans := map[string]Plan{}
	for k, p := range plans {
		inputs := map[string]Source{}
		for name, source := range p.Inputs {
			if source.Literal != nil {
				v := *source.Literal
				source.Literal = &v
			}
			inputs[name] = source
		}
		p.Inputs = inputs
		copiedPlans[k] = p
	}
	copiedState := map[string]sdl.Value{}
	for k, v := range state {
		copiedState[k] = v
	}
	b := &Bridge{UI: session, Modules: copiedModules, Plans: copiedPlans, Context: copiedState, ctx: ctx}
	if b.Context == nil {
		b.Context = map[string]sdl.Value{}
	}
	if err := b.Rebind(document); err != nil {
		return nil, err
	}
	return b, nil
}

// Rebind validates the complete binding set before installing any handlers.
func (b *Bridge) Rebind(document *uiparser.Document) error {
	refs := map[string]bool{}
	for _, r := range document.References {
		refs[r.Alias] = true
	}
	targets := map[string]string{}
	for _, c := range document.Connections {
		symbol := c.Module + "." + c.Object
		target := c.Definition + "/" + strings.Join(c.Path, "/")
		if targets[symbol] != "" {
			return fmt.Errorf("duplicate-handle: %s", symbol)
		}
		targets[symbol] = target
	}
	spans := map[string]uiparser.Span{}
	b.UI.SnapshotRoot().Walk(func(n *uiparser.Instance) { spans[n.Path] = n.Span })
	pending := map[ui.Handle]ui.Handler{}
	links := []Link{}
	used := map[string]bool{}
	for _, widget := range b.UI.Widgets() {
		ref := widget.Binding
		if ref.Module == "" {
			continue
		}
		symbol := ref.Module + "." + ref.Object
		engine := b.Modules[ref.Module]
		if !refs[ref.Module] || engine == nil {
			return fmt.Errorf("unbound-module: %s", ref.Module)
		}
		if ref.Member != "invoke" {
			return fmt.Errorf("callback-member: %s must use @invoke", symbol)
		}
		plan, ok := b.Plans[symbol]
		if !ok {
			return fmt.Errorf("missing-binding-plan: %s", symbol)
		}
		input, output, err := engine.Signature(ref.Object)
		if err != nil {
			return err
		}
		if len(input) != len(plan.Inputs) {
			return fmt.Errorf("binding-fields: %s", symbol)
		}
		for field, kind := range input {
			source, ok := plan.Inputs[field]
			if !ok {
				return fmt.Errorf("binding-field: %s.%s", symbol, field)
			}
			if err = b.validateSource(source, widget, kind); err != nil {
				return fmt.Errorf("binding-field %s.%s: %w", symbol, field, err)
			}
		}
		target := targets[symbol]
		if target == "" || plan.OutputField == "" || output[plan.OutputField] != parser.TextType {
			return fmt.Errorf("result-handle: %s needs a text result and explicit setHandle", symbol)
		}
		resultWidget, ok := b.UI.Widget(target)
		if !ok || resultWidget.Handle.Kind != "input" {
			return fmt.Errorf("result-widget: %s is not an input", target)
		}
		if plan.RevisionField != "" {
			if output[plan.RevisionField] != parser.IntegerType || plan.RevisionContext == "" {
				return fmt.Errorf("revision-binding: %s", symbol)
			}
		} else if plan.RevisionContext != "" {
			return fmt.Errorf("revision-binding: missing result field")
		}
		action, _ := engine.Action(ref.Object)
		links = append(links, Link{widget.Handle.Path, symbol, target, spans[widget.InstancePath], action.Span})
		used[symbol] = true
		pending[widget.Handle] = b.handler(engine, ref.Object, plan, target)
	}
	for symbol := range targets {
		if !used[symbol] {
			return fmt.Errorf("unbound-handle: %s has no matching callback in this session", symbol)
		}
	}
	for symbol := range b.Plans {
		if !used[symbol] {
			return fmt.Errorf("unused-binding-plan: %s", symbol)
		}
	}
	handles := []ui.Handle{}
	for h := range pending {
		handles = append(handles, h)
	}
	sort.Slice(handles, func(i, j int) bool { return handles[i].Path < handles[j].Path })
	for _, h := range handles {
		if err := b.UI.Bind(h, pending[h]); err != nil {
			return err
		}
	}
	b.Links = links
	return nil
}
