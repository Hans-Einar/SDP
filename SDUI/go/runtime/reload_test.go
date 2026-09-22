package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"testing"
)

func root(t *testing.T, source string) *parser.Instance {
	t.Helper()
	_, roots, e := parser.Compile("sdui 0.2; " + source)
	if e != nil {
		t.Fatal(e)
	}
	return roots["page"]
}
func TestReloadIdentityDraftFocusAndRevocation(t *testing.T) {
	s := session(t)
	w, _ := s.Widget("page/edit")
	s.Draft(w.Handle, "draft")
	s.Focus(w.Handle)
	button, _ := s.Widget("page/go")
	calls := 0
	s.Bind(button.Handle, func(Event) ([]Update, error) { calls++; return nil, nil })
	oldEvent := Event{Handle: button.Handle, ModelRevision: 1, Sequence: 1, Kind: Activate}
	if e := s.Reload(root(t, `page=["new";< <edit=input("New label",value="new default")>,go=button("Go 2")>] {scale=1};`)); e != nil {
		t.Fatal(e)
	}
	now, _ := s.Widget("page/edit")
	if now.Handle != w.Handle || now.Value != "1" || now.Draft != "draft" || !now.Dirty || s.Focused() != "page/edit" || now.Label != "New label" {
		t.Fatal(now, s.Focused())
	}
	code(t, s.Dispatch(oldEvent), "stale-event")
	if calls != 0 {
		t.Fatal("reload executed callback")
	}
	oldEvent.ModelRevision = s.Revision
	if e := s.Dispatch(oldEvent); e != nil {
		t.Fatal(e)
	}
	if calls != 1 {
		t.Fatal(calls)
	}
	s.Reload(root(t, `page=[edit=button("Replacement")];`))
	replacement, _ := s.Widget("page/edit")
	if replacement.Handle.Generation == w.Handle.Generation || s.Focused() != "" {
		t.Fatal("type change kept handle/focus")
	}
	code(t, s.Draft(w.Handle, "late"), "stale-handle")
	code(t, s.Dispatch(oldEvent), "stale-handle")
	s.Close()
	code(t, s.Reload(root(t, `page=[];`)), "closed")
}
func TestChangedBindingDoesNotReuseHandler(t *testing.T) {
	r := root(t, `ref: s "x.sdl"; page=[go=button("Go",callback=s.first.@run)];`)
	s, e := New("test", r)
	if e != nil {
		t.Fatal(e)
	}
	w, _ := s.Widget("page/go")
	s.Bind(w.Handle, func(Event) ([]Update, error) { t.Fatal("old binding ran"); return nil, nil })
	if e = s.Reload(root(t, `ref: s "x.sdl"; page=[go=button("Go",callback=s.second.@run)];`)); e != nil {
		t.Fatal(e)
	}
	code(t, s.Dispatch(Event{Handle: w.Handle, Kind: Activate, ModelRevision: 2, Sequence: 1}), "unbound")
}
