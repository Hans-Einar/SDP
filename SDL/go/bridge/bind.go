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

type EventField string

// CollectionItemID extracts only the validated collection Activate item identity.
const CollectionItemID EventField = "collection.item-id"

// Tab selectors extract stable direct-page IDs from a validated ActivatePage.
const (
	TabPageID            EventField = "tab.page-id"
	TabPreviousPageID    EventField = "tab.previous-page-id"
	CommandContextItemID EventField = "command.context-item-id"
	CommandChecked       EventField = "command.checked"
	DialogFieldValue     EventField = "dialog.field-value"
)

// ResultMode is closed. TextResult is the zero/default legacy text receiver mode.
type ResultMode string

const (
	TextResult         ResultMode = ""
	DialogAcceptResult ResultMode = "dialog-accept"
)

type Source struct {
	EventField    EventField
	FieldPath     string
	resolvedField ui.Handle
	Widget        string
	Event         bool
	Literal       *sdl.Value
	Context       string
}
type Plan struct {
	ResultMode                     ResultMode
	AcceptField, MessageField      string
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
		copiedPlans[k] = copyPlan(p)
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
	interactions := map[ui.Handle]ui.InteractionHandler{}
	links := []Link{}
	used := map[string]bool{}
	owners := append(b.UI.Widgets(), b.UI.CallbackOwners()...)
	seen := map[ui.Handle]bool{}
	for _, widget := range owners {
		ref := widget.Binding
		if ref.Module == "" {
			continue
		}
		if seen[widget.Handle] {
			continue
		}
		seen[widget.Handle] = true
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
		plan = copyPlan(plan) // Resolve owned captures separately for each reused callback owner.
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
			if source.EventField == DialogFieldValue {
				source.resolvedField, err = b.UI.DialogField(widget.Handle, source.FieldPath)
				if err != nil {
					return fmt.Errorf("binding-field %s.%s: %w", symbol, field, err)
				}
				plan.Inputs[field] = source
			}
		}
		target := targets[symbol]
		if err = b.validateResult(plan, widget, target, output); err != nil {
			return fmt.Errorf("binding-result %s: %w", symbol, err)
		}

		action, _ := engine.Action(ref.Object)
		links = append(links, Link{widget.Handle.Path, symbol, target, spans[widget.InstancePath], action.Span})
		used[symbol] = true
		_, command := b.UI.CommandState(widget.Handle)
		if widget.Handle.Kind == "tabs" || widget.Handle.Kind == "dialog" || command {
			interactions[widget.Handle] = b.interactionHandler(engine, ref.Object, plan, target)
		} else {
			pending[widget.Handle] = b.handler(engine, ref.Object, plan, target)
		}
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
	for h := range interactions {
		handles = append(handles, h)
	}
	sort.Slice(handles, func(i, j int) bool { return handles[i].Path < handles[j].Path })
	for _, h := range handles {
		if handler, ok := interactions[h]; ok {
			if err := b.UI.BindInteraction(h, handler); err != nil {
				return err
			}
		} else if err := b.UI.Bind(h, pending[h]); err != nil {
			return err
		}
	}
	b.Links = links
	return nil
}

func copyPlan(p Plan) Plan {
	inputs := make(map[string]Source, len(p.Inputs))
	for name, source := range p.Inputs {
		source.resolvedField = ui.Handle{}
		if source.Literal != nil {
			v := *source.Literal
			source.Literal = &v
		}
		inputs[name] = source
	}
	p.Inputs = inputs
	return p
}
