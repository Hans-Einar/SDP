package snapshot

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestSelectedStateBundle(t *testing.T) {
	src, e := os.ReadFile("../examples/edit-apt-cell.sdui")
	if e != nil {
		t.Fatal(e)
	}
	value := "440"
	label := "Accepted"
	disabled := false
	o := Options{Entry: "page", Width: 1000, Height: 650, State: map[string]WidgetState{"page/value": {Value: &value}, "page/apply": {Label: &label, Enabled: &disabled}}}
	a, e := Build(context.Background(), string(src), o)
	if e != nil {
		t.Fatal(e)
	}
	b, e := Build(context.Background(), string(src), o)
	if e != nil {
		t.Fatal(e)
	}
	for k, v := range a.Files {
		if !bytes.Equal(v, b.Files[k]) {
			t.Fatal("nonreproducible", k)
		}
	}
	if !strings.Contains(string(a.Files["structure.md"]), "440") || !strings.Contains(string(a.Files["structure.md"]), "Accepted") {
		t.Fatal("snapshot not applied")
	}
	var provenance map[string]any
	if e = json.Unmarshal(a.Files["provenance.json"], &provenance); e != nil || provenance["tool"] != Version {
		t.Fatal("missing provenance", e)
	}
	if !strings.Contains(string(a.Files["ui.svg"]), `data-path="page/`) {
		t.Fatal("missing measured geometry")
	}
	for _, state := range []map[string]WidgetState{{"missing": {Value: &value}}, {"page/apply": {Value: &value}}} {
		o.State = state
		if _, e := Build(context.Background(), string(src), o); e == nil {
			t.Fatal("invalid state accepted")
		}
	}
	o.State = nil
	o.Width = -1
	if _, e := Build(context.Background(), string(src), o); e == nil {
		t.Fatal("invalid viewport")
	}
}
func TestDesignConsumer(t *testing.T) {
	src, e := os.ReadFile("../examples/edit-apt-cell.sdui")
	if e != nil {
		t.Fatal(e)
	}
	design, e := os.ReadFile("../../../SDUI/design/architecture.design")
	if e != nil {
		t.Fatal(e)
	}
	b, e := Build(context.Background(), string(src), Options{Entry: "page", Width: 1000, Height: 650, Design: string(design)})
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(b.Files["design/navigator.md"]), "sdl-view://design/VP") {
		t.Fatal("not consuming SDL projector")
	}
	for name := range b.Files {
		if strings.HasPrefix(name, "design/diagrams/") {
			t.Fatal("navigator eagerly rendered detail")
		}
	}
}
