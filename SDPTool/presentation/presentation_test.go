package presentation

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestTreeAndMachineBoundary(t *testing.T) {
	raw := json.RawMessage(`{"schema":"sdptool/0.2","operation":"tree","project":"xfmd","roots":["sdl","kb","ui"],"nodes":[{"id":"sdl","label":"SDL","state":"absent"},{"id":"kb","label":"KanBan","state":"validated","children":["backlog"]},{"id":"backlog","label":"backlog","state":"available","children":["card"]},{"id":"card","label":"#003--Question.md","state":"available"},{"id":"ui","label":"SDUI","state":"absent"}],"expansionDepthLimit":8}`)
	want := "xfmd\n├── SDL [absent]\n├── KanBan [validated]\n│   └── backlog [available]\n│       └── #003--Question.md [available]\n└── SDUI [absent]\n"
	for _, machine := range []bool{false, true} {
		var out bytes.Buffer
		if err := Default().WriteJSON(&out, raw, machine); err != nil {
			t.Fatal(err)
		}
		expected := want
		if machine {
			expected = string(raw) + "\n"
		}
		if out.String() != expected {
			t.Fatalf("machine=%v: %q", machine, out.String())
		}
	}
}

func TestUnknownAndExtension(t *testing.T) {
	raw := json.RawMessage(`{ "schema":"new/1", "operation":"future", "value":9007199254740993 }`)
	r := Default()
	var out bytes.Buffer
	if err := r.WriteJSON(&out, raw, false); err != nil || out.String() != string(raw)+"\n" {
		t.Fatalf("fallback: %v %s", err, &out)
	}
	r[Key{"new/1", "future"}] = func(w io.Writer, b json.RawMessage) error { _, e := io.WriteString(w, "extension\n"); return e }
	out.Reset()
	if err := r.WriteJSON(&out, raw, false); err != nil || out.String() != "extension\n" {
		t.Fatalf("extension: %v %s", err, &out)
	}
	out.Reset()
	if err := r.WriteJSON(&out, raw, true); err != nil || out.String() != string(raw)+"\n" {
		t.Fatalf("bypass: %v %s", err, &out)
	}
}

func TestTreeLimitsReferencesAndControls(t *testing.T) {
	for _, tc := range []struct{ name, extra, want string }{
		{"cycle", ``, "[cycle]"},
		{"reference", `,"reference":"elsewhere"`, "→ elsewhere"},
		{"external", `,"externalReference":"KB-EXT-001"`, "→ external: KB-EXT-001"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw := json.RawMessage(`{"schema":"sdptool/0.2","operation":"tree","project":"test","roots":["a"],"nodes":[{"id":"a","label":"x\u001b\ny","children":["a"]` + tc.extra + `}]}`)
			var out bytes.Buffer
			if err := Default().WriteJSON(&out, raw, false); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), tc.want) || strings.Contains(out.String(), "\x1b") {
				t.Fatal(out.String())
			}
		})
	}
	var out bytes.Buffer
	raw := json.RawMessage(`{"schema":"sdptool/0.2","operation":"tree","project":"x","expansionDepthLimit":1,"roots":["a"],"nodes":[{"id":"a","label":"A","children":["b"]},{"id":"b","label":"B"}]}`)
	if err := Default().WriteJSON(&out, raw, false); err != nil || !strings.Contains(out.String(), "[depth limit]") {
		t.Fatalf("%v %s", err, &out)
	}
	for _, raw := range []string{
		`{"schema":"sdptool/0.2","operation":"tree","roots":["missing"]}`,
		`{"schema":"sdptool/0.2","operation":"tree","nodes":[{"id":"a"},{"id":"a"}]}`,
	} {
		out.Reset()
		if err := Default().WriteJSON(&out, []byte(raw), false); err == nil || out.Len() != 0 {
			t.Fatalf("invalid tree: %v %s", err, &out)
		}
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, errors.New("broken pipe") }
func TestWriterFailure(t *testing.T) {
	for _, machine := range []bool{false, true} {
		if err := Default().Write(brokenWriter{}, map[string]string{"schema": "sdptool/0.2", "operation": "version"}, machine); err == nil {
			t.Fatal("missing write error")
		}
	}
}
