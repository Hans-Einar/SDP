package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestDeferredCLIAndJSONActions(t *testing.T) {
	root := deferredFixture(t)
	ctx := context.Background()
	p := compact(t, root)
	var model, revision string
	for _, s := range p.Sources {
		if s.State == "validated" {
			model, revision = s.ID, s.Revision
		}
	}
	cases := []struct {
		args   []string
		action string
		params any
	}{
		{[]string{root, "discover", "--navigation", "lazy-sdl", "--json"}, "project.discover", map[string]any{"navigation": "lazy-sdl"}},
		{[]string{root, "tree", "--model", model, "--revision", revision, "--model-only", "--json"}, "project.model-tree", ModelTreeOptions{Model: model, Revision: revision}},
	}
	for _, c := range cases {
		var out, errs bytes.Buffer
		if code := Run(ctx, c.args, &out, &errs); code != 0 {
			t.Fatal(code, errs.String())
		}
		reply := Invoke(ctx, bytes.NewReader(requestBytes(t, root, c.action, c.params)), ActionCaller{})
		if reply.ExitCode != 0 {
			t.Fatal(reply)
		}
		var a, b any
		json.Unmarshal(out.Bytes(), &a)
		raw, _ := json.Marshal(reply.Result)
		json.Unmarshal(raw, &b)
		if !reflect.DeepEqual(a, b) {
			t.Fatal("CLI/invoke mismatch", c.action)
		}
	}
	catalog, e := Actions(ctx, root)
	if e != nil {
		t.Fatal(e)
	}
	found := 0
	for _, a := range catalog.Actions {
		switch a.ID {
		case "project.discover":
			if !reflect.DeepEqual(a.InputSchema.Properties["navigation"].Enum, []string{"eager", "lazy-sdl"}) {
				t.Fatal(a)
			}
			found++
		case "project.model-tree":
			if a.Mutation != "none" || len(a.InputSchema.Required) != 2 {
				t.Fatal(a)
			}
			found++
		}
	}
	if found != 2 {
		t.Fatal("capability not advertised")
	}
	for _, args := range [][]string{{root, "discover", "--navigation", "unknown", "--json"}, {root, "discover", "extra", "--json"}, {root, "tree", "--model-only", "--json"}, {root, "tree", "--model", model, "--revision", revision, "--json"}} {
		var o, e bytes.Buffer
		if Run(ctx, args, &o, &e) == 0 || o.Len() != 0 {
			t.Fatal("invalid CLI accepted", args, o.String())
		}
	}
	for _, bad := range []struct {
		action string
		params any
	}{{"project.discover", map[string]any{"navigation": "unknown"}}, {"project.discover", map[string]any{"navigation": true}}, {"project.model-tree", map[string]any{"model": model}}, {"project.model-tree", ModelTreeOptions{}}, {"project.model-tree", ModelTreeOptions{Model: model, Revision: "old"}}} {
		v := Invoke(ctx, bytes.NewReader(requestBytes(t, root, bad.action, bad.params)), ActionCaller{})
		if v.ExitCode == 0 || v.Result != nil {
			t.Fatal("invalid JSON request accepted", bad, v)
		}
	}
}
