package parser

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestPreviewSourceNormalizationReuseAndNoIO(t *testing.T) {
	const source = `sdui 0.3; ref: art "/does/not/exist.sdl";
Text=markdown("å🙂\n# Heading",description="  Reading  ",fallback="label");
Main=[left=Text;right=Text {visible=false};image=svg(art.Image.@resource,label="Caption",description="Figure",fallback="reject");old=svg(art.Old.@anything);bare="Bare"];`
	doc, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	before, _ := json.Marshal(Data(doc))
	node := doc.Definitions[0].Root
	if node.Kind != "widget" || val(node.Widget) != "markdown" || node.Text != nil {
		t.Fatal("source AST rewritten", node)
	}
	root := roots["Main"]
	left, right := root.Rows[0][0], root.Rows[1][0]
	for _, n := range []*Instance{left, right} {
		p, err := PreviewOptions(n)
		if err != nil || !p.Explicit || p.Description != "  Reading  " || p.Fallback != "label" || n.Kind != "markdown" || n.Widget != "" || n.Text != "å🙂\n# Heading" || len(n.Arguments) != 3 || len(n.Uses) != 1 {
			t.Fatal(n, p, err)
		}
	}
	if left.Path == right.Path || left.Uses[0].Span == right.Uses[0].Span || left.Span != node.Span {
		t.Fatal("reuse identity/span lost")
	}
	if _, err := ResolveInteractions(root); err != nil {
		t.Fatal(err)
	}
	for _, n := range root.Rows[3:] {
		p, err := PreviewOptions(n[0])
		if err != nil || p.Explicit {
			t.Fatal("legacy opted in", p, err)
		}
	}
	again, err := Normalize(doc)
	if err != nil || !reflect.DeepEqual(root, again["Main"]) {
		t.Fatal("normalization unstable", err)
	}
	left.Arguments["description"] = Literal{Kind: "string", Value: "Changed"}
	after, _ := json.Marshal(Data(doc))
	if string(before) != string(after) || right.Argument("description") != "  Reading  " {
		t.Fatal("arguments alias")
	}
}
func TestPreviewClosedSchemasAndLegacyBoundary(t *testing.T) {
	bad := []string{
		`markdown("x")`, `markdown("x",description="D")`, `markdown("x",fallback="label")`,
		`markdown("x",description=" ",fallback="label")`, `markdown("x",description="D",fallback="LABEL")`,
		`markdown("x",description=false,fallback="label")`, `markdown(1,description="D",fallback="label")`,
		`markdown("x",description="D",fallback=null)`, `markdown(text="x","D",fallback="label")`,
		`markdown("x",description="D",fallback="label",text="again")`,
		`markdown("x",description="D",fallback="label",callback=art.Do.@invoke)`,
		`markdown("x",description="D",fallback="label")[]`,
		`svg(art.Image.@resource,description="D")`, `svg(art.Image.@resource,fallback="label")`,
		`svg(art.Image.@invoke,description="D",fallback="label")`,
		`svg(missing.Image.@resource,description="D",fallback="label")`,
		`svg(art.Image.@resource,description="D",fallback="label",label=false)`,
		`svg(art.Image.@resource,description="D",fallback="label",description="again")`,
		`svg(art.Image.@resource,description="D",fallback="label")[]`,
	}
	for _, body := range bad {
		t.Run(body, func(t *testing.T) {
			if _, _, err := Compile(`sdui 0.3; ref: art "unopened"; Main=[p=` + body + `];`); err == nil {
				t.Fatal("accepted")
			}
		})
	}
	for _, body := range []string{`markdown("",description="D",fallback="label")`, `svg(art.Image.@resource,description="D",fallback="reject")`} {
		if _, _, err := Compile(`sdui 0.2; ref: art "unopened"; Main=[p=` + body + `];`); err == nil {
			t.Fatal(".2 accepted", body)
		}
	}
	for _, profile := range []string{"0.2", "0.3"} {
		for _, member := range []string{"resource", "anything"} {
			_, roots, err := Compile(`sdui ` + profile + `; ref: art "unopened"; Main=[p=svg(art.Image.@` + member + `,label="D");m="Plain"];`)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range roots["Main"].Rows {
				p, err := PreviewOptions(r[0])
				if err != nil || p.Explicit {
					t.Fatal(p, err)
				}
			}
		}
	}
	// Unsupported HTML/images are provider content decisions, not grammar errors.
	for _, text := range []string{"", "<script>x</script>", "![image](remote)", strings.Repeat("å", 16384)} {
		if _, _, err := Compile(`sdui 0.3; Main=[markdown(` + strconv.Quote(text) + `,description="D",fallback="label")];`); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range []string{`"x",description=` + strconv.Quote(strings.Repeat("x", 4097)) + `,fallback="label"`, strconv.Quote(strings.Repeat("x", 32769)) + `,description="D",fallback="label"`} {
		if _, _, err := Compile(`sdui 0.3; Main=[markdown(` + args + `)];`); err == nil {
			t.Fatal("limit not enforced")
		}
	}
}
func TestPreviewConstructedValidationAndStrictHiddenRoot(t *testing.T) {
	for _, mutate := range []func(*Instance){
		func(n *Instance) { n.Text = "different" }, func(n *Instance) { n.Widget = "markdown" },
		func(n *Instance) { n.Arguments["description"] = Literal{Kind: "string", Value: 1} },
		func(n *Instance) { n.Arguments["description"] = Literal{Kind: "string", Value: string([]byte{255})} },
		func(n *Instance) { n.Arguments["text"] = Reference{Module: "art", Object: "X", Member: "resource"} },
		func(n *Instance) { n.Profile = "" }, func(n *Instance) { n.Rows = [][]*Instance{{{Kind: "markdown"}}} },
	} {
		_, roots, err := Compile(`sdui 0.3; Main=[p=markdown("x",description="D",fallback="label") {visible=false}];`)
		if err != nil {
			t.Fatal(err)
		}
		n := roots["Main"].Rows[0][0]
		mutate(n)
		if _, err := PreviewOptions(n); err == nil {
			t.Fatal("malformed preview accepted")
		}
		if _, err := ResolveInteractions(roots["Main"]); err == nil {
			t.Fatal("hidden malformed preview accepted")
		}
	}
	_, roots, err := Compile(`sdui 0.3; ref: art "unopened"; Main=[p=svg(art.X.@resource,description="D",fallback="label")];`)
	if err != nil {
		t.Fatal(err)
	}
	n := roots["Main"].Rows[0][0]
	n.Arguments["source"] = Reference{Module: "bad/path", Object: "X", Member: "resource", Span: n.Span}
	_, err = PreviewOptions(n)
	var d *Diagnostic
	if !errors.As(err, &d) || d.Span != n.Span || !strings.Contains(d.Message, n.Path) {
		t.Fatal(err)
	}
	if _, err := PreviewOptions(nil); err == nil {
		t.Fatal("nil accepted")
	}
}
func TestPreviewMarkdownFormattingAndWidgetRoles(t *testing.T) {
	for _, body := range []string{`"Text"`, `markdown("Text",description="D",fallback="label")`} {
		for _, suffix := range []string{` {align-x=start}`, ` {visible=false}`} {
			if _, _, err := Compile(`sdui 0.3; M=` + body + `; Main=[p=M` + suffix + `];`); err != nil {
				t.Fatal(body, suffix, err)
			}
		}
	}
	for _, tail := range []string{`app.Save.setHandle(Main.p);`, ``} {
		src := `sdui 0.3; ref: app "unopened"; Main=[p=markdown("x",description="D",fallback="label")`
		if tail == "" {
			src += `;cmd=command("C",context="field",target="p",callback=app.Save.@invoke)`
		}
		src += `];` + tail
		_, roots, err := Compile(src)
		if err == nil {
			_, err = ResolveInteractions(roots["Main"])
		}
		if err == nil {
			t.Fatal("Markdown admitted as interactive field", src)
		}
	}
}
