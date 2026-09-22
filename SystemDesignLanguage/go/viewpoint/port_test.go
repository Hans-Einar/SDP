package viewpoint

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestFrozenProjectionPort(t *testing.T) {
	b, e := os.ReadFile("testdata/python-projection-source.design")
	if e != nil {
		t.Fatal(e)
	}
	v, e := New(string(b))
	if e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile("testdata/python-projections.json")
	if e != nil {
		t.Fatal(e)
	}
	var want map[string]any
	if e = json.Unmarshal(raw, &want); e != nil {
		t.Fatal(e)
	}
	actual := map[string]any{"source_sha256": v.Revision, "facts": v.Facts, "model_gaps": v.Gaps, "message_sets": v.MessageSets, "diagrams": v.Diagrams}
	out, _ := json.Marshal(actual)
	var got map[string]any
	json.Unmarshal(out, &got)
	for _, d := range got["diagrams"].([]any) {
		delete(d.(map[string]any), "title")
	}
	hashes := map[string]any{}
	for _, d := range v.Diagrams {
		hashes[d.ID] = fmt.Sprintf("%x", sha256.Sum256([]byte(d.Mermaid())))
	}
	got["mermaid_sha256"] = hashes
	for k, w := range want {
		if !reflect.DeepEqual(got[k], w) {
			os.WriteFile("/tmp/sdl-projection-got.json", out, 0600)
			t.Errorf("projection port differs: %s", k)
		}
	}
}
