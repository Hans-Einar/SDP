package svg

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

type preparedPreviewProbe struct {
	nodes           map[string]string
	checks, renders int
}

func (p *preparedPreviewProbe) CheckPreview(n *parser.Instance) error {
	p.checks++
	if p.nodes[n.Path] != n.Text+"|"+n.Argument("description")+"|"+n.Argument("fallback") {
		return &parser.Diagnostic{Code: "preview-unprepared", Message: n.Path + ": missing or mismatched frozen outcome", Span: n.Span}
	}
	return nil
}
func (p *preparedPreviewProbe) Render(out *strings.Builder, b *layout.Box) error {
	p.renders++
	out.WriteString("<!-- prepared prose -->")
	return nil
}
func previewRoot(t *testing.T, body string) *parser.Instance {
	t.Helper()
	_, roots, err := parser.Compile(`sdui 0.3; ref: art "unopened"; Main=[` + body + `];`)
	if err != nil {
		t.Fatal(err)
	}
	return roots["Main"]
}
func TestExplicitPreviewPublicSVGRejectsBeforeRendererOrGeometry(t *testing.T) {
	for _, kind := range []string{`svg(art.Image.@resource,description="D",fallback=`, `markdown("# Text",description="D",fallback=`} {
		for _, policy := range []string{"label", "reject"} {
			for _, hidden := range []string{"", ` {visible=false}`} {
				root := previewRoot(t, `p=`+kind+`"`+policy+`")`+hidden)
				node := root.Rows[0][0]
				probe := &preparedPreviewProbe{}
				out, err := Render(&layout.Box{Instance: root, Path: root.Path}, Options{Content: probe})
				var d *parser.Diagnostic
				if out != "" || !errors.As(err, &d) || d.Code != "unsupported-resource-export" || d.Span != node.Span || !strings.Contains(d.Message, node.Path) || probe.checks != 0 || probe.renders != 0 {
					t.Fatal(out, err, probe)
				}
			}
		}
	}
}
func TestNativePreviewExactInventoryAndFrozenMarkdownCheck(t *testing.T) {
	root := previewRoot(t, `image=svg(art.Image.@resource,description="D",fallback="label");a=markdown("same",description="A",fallback="label");b=markdown("same",description="B",fallback="reject") {visible=false};old=svg(art.Old.@resource,label="Legacy")`)
	probe := &preparedPreviewProbe{nodes: map[string]string{"Main/a": "same|A|label", "Main/b": "same|B|reject"}}
	options := Options{SkipControls: true, Content: probe, NativeControls: map[string]string{"Main/image": "svg"}}
	for i := 0; i < 3; i++ {
		if err := Check(root, options); err != nil {
			t.Fatal(err)
		}
	}
	if probe.checks != 6 || probe.renders != 0 {
		t.Fatal("preflight rendered or skipped hidden content", probe)
	}
	for _, inv := range []map[string]string{nil, {"Main/image": "input"}, {"Main/image": "svg", "Main/old": "svg"}, {"Main/image": "svg", "Main/a": "markdown"}} {
		invalid := options
		invalid.NativeControls = inv
		if err := Check(root, invalid); err == nil {
			t.Fatal("invalid inventory", inv)
		}
	}
	invalid := options
	invalid.Content = nil
	if err := Check(root, invalid); err == nil {
		t.Fatal("missing prepared content accepted")
	}
	delete(probe.nodes, "Main/b")
	if err := Check(root, options); err == nil {
		t.Fatal("hidden missing outcome accepted")
	}
	probe.nodes["Main/b"] = "same|WRONG|reject"
	if err := Check(root, options); err == nil {
		t.Fatal("policy/description mismatch accepted")
	}
	probe.nodes["Main/b"] = "same|B|reject"
	// Nonzero geometry verifies omission of the native image, retention of the legacy
	// placeholder, and use of the already-prepared Markdown renderer.
	box := &layout.Box{Instance: root, Path: root.Path, Rect: layout.Rect{W: 400, H: 400}, Clip: layout.Rect{W: 400, H: 400}}
	for _, row := range root.Rows {
		n := row[0]
		box.Children = append(box.Children, &layout.Box{Instance: n, Path: n.Path, Rect: layout.Rect{W: 300, H: 60}, Clip: box.Clip, Font: 14, Enabled: true})
	}
	out, err := Render(box, options)
	if err != nil || !strings.Contains(out, "Legacy") || strings.Contains(out, ">D<") || probe.renders != 2 {
		t.Fatal(out, err, probe)
	}
}

func TestAllFamilyPublicSVGSupportBoundary(t *testing.T) {
	cases := []struct{ body, code string }{
		{`p=button("B")`, ``}, {`p=input("I",value="old")`, ``}, {`p=svg(art.X.@draw)`, ``}, {`p="Bare Markdown"`, ``},
		{`p=tree("Tree")`, `unsupported-collection-export`}, {`p=list("List")`, `unsupported-collection-export`},
		{`p=tabs("T")[one=page("One")[]]`, `unsupported-pane-export`},
		{`p=split("horizontal")[left=[];right=[]]`, `unsupported-pane-export`},
		{`c=command("C");p=button(command="c")`, `unsupported-interaction-export`},
		{`p=button("B",toggle=false)`, `unsupported-interaction-export`},
		{`p=button("B",icon="symbol")`, `unsupported-interaction-export`},
		{`p=menu("M")[menuGroup("Group")[separator()]]`, `unsupported-interaction-export`},
		{`p=dialog("D")[]`, `unsupported-interaction-export`},
		{`p=dialog("D",modal=false)[]`, `unsupported-interaction-export`},
		{`p=checkbox("C")`, `unsupported-value-export`}, {`p=select("S")`, `unsupported-value-export`},
		{`p=number("N",min=0,max=1,step=1,value=0)`, `unsupported-value-export`},
		{`p=slider("S",min=0,max=1,step=1,value=0)`, `unsupported-value-export`},
		{`p=input("I",placeholder="")`, `unsupported-text-export`},
		{`p=input("I",multiline=true)`, `unsupported-text-export`},
		{`p=svg(art.X.@resource,description="D",fallback="label")`, `unsupported-resource-export`},
		{`p=markdown("Text",description="D",fallback="reject")`, `unsupported-resource-export`},
	}
	for _, tc := range cases {
		t.Run(tc.body, func(t *testing.T) {
			root := previewRoot(t, tc.body)
			err := Check(root, Options{})
			if tc.code == "" {
				if err != nil {
					t.Fatal("legacy public export regressed", err)
				}
				return
			}
			var d *parser.Diagnostic
			if !errors.As(err, &d) || d.Code != tc.code || d.Span.Line == 0 || !strings.Contains(d.Message, "Main/") {
				t.Fatal(tc.code, err)
			}
		})
	}
	root := previewRoot(t, `p="Text"`)
	root.Rows[0][0].Kind = "widget"
	root.Rows[0][0].Widget = "future"
	root.Rows[0][0].Layout["visible"] = false
	if err := Check(root, Options{SkipControls: true, NativeControls: map[string]string{"Main/p": "future"}}); err == nil {
		t.Fatal("unknown hidden widget omitted")
	}
}
