package fynehost

import (
	"fmt"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestInspectionM2SurfaceCanvasProjection(t *testing.T) {
	for _, modal := range []bool{true, false} {
		t.Run(fmt.Sprint(modal), func(t *testing.T) {
			h, w := lifecycleHost(t)
			w.SetTitle("Inspector main")
			source := fmt.Sprintf(`sdui 0.3; page=[open=button("Open",effect="open",target="child");child=dialog("Inspector child",modal=%t)[field=input("Field",value="saved") {x=fill}] {scale-x=0.4,scale-y=0.5}] {x=fill,y=fill};`, modal)
			doc, err := parser.Parse(source)
			lifecycleOK(t, err)
			lifecycleOK(t, h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "inspection", Sequence: 1, SourceRevision: "inspect", Mode: preparation.Prototype}))
			w.SetContent(h.Container)
			w.Resize(fyne.NewSize(800, 500))
			w.Show()
			b := h.Current()
			closed := b.Inspect()
			if len(closed["surfaces"].(map[string]map[string]any)) != 0 || closed["controls"].(map[string]map[string]any)["page/child/field"]["visible"] != false {
				t.Fatal("closed surface reported as native visible", closed)
			}
			opener := b.Session.Snapshot().Presentations["page/open"].Handle
			lifecycleOK(t, b.invokeCommand(opener, "button"))
			surface := b.surfaces["page/child"]
			if surface == nil {
				t.Fatal("surface not published")
			}
			before := b.Session.Snapshot()
			inspected := b.Inspect()
			row := inspected["surfaces"].(map[string]map[string]any)["page/child"]
			control := inspected["controls"].(map[string]map[string]any)["page/child/field"]
			title, canvas := "Inspector main", "main"
			if !modal {
				title, canvas = "Inspector child", "page/child"
			}
			if row["title"] != title || row["canvas"] != canvas || row["modal"] != modal || row["target"] != surface.target || row["visible"] != true {
				t.Fatal("surface identity/visibility", row)
			}
			if control["title"] != title || control["canvas"] != canvas || control["surface"] != "page/child" || control["visible"] != true {
				t.Fatal("flat control canvas identity", control)
			}
			body := row["rect"].(layout.Rect)
			clip := control["clip"].(layout.Rect)
			if body.Y <= 0 || clip.Y < body.Y || clip.W <= 0 || clip.H <= 0 || clip.Intersect(body) != clip {
				t.Fatal("body chrome offset/clipping lost", body, clip)
			}
			c := b.view.controls["page/child/field"]
			nativePos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c.clip)
			if clip.X != float64(nativePos.X) || clip.Y != float64(nativePos.Y) {
				t.Fatal("control is not in actual canvas coordinates", clip, nativePos)
			}
			if !reflect.DeepEqual(before, b.Session.Snapshot()) {
				t.Fatal("inspection changed runtime")
			}
			row["title"] = "tampered"
			control["visible"] = false
			next := b.Inspect()
			if next["surfaces"].(map[string]map[string]any)["page/child"]["title"] != title || next["controls"].(map[string]map[string]any)["page/child/field"]["visible"] != true {
				t.Fatal("inspection maps alias live state")
			}
			lifecycleOK(t, h.Mutate(func(s *ui.Session) error { return s.CloseSurface(surface.target, "cancel") }))
			after := b.Inspect()
			if len(after["surfaces"].(map[string]map[string]any)) != 0 || after["controls"].(map[string]map[string]any)["page/child/field"]["visible"] != false {
				t.Fatal("closed native geometry retained", after)
			}
		})
	}
}

func TestInspectionM2EmbeddedMainOrigin(t *testing.T) {
	h, w := lifecycleHost(t)
	lifecycleOK(t, h.Adopt(paneRequest(t)))
	w.SetContent(container.NewWithoutLayout(h.Container))
	w.Resize(fyne.NewSize(1000, 700))
	h.Container.Resize(fyne.NewSize(800, 500))
	h.Container.Move(fyne.NewPos(30, 20))
	w.Show()
	b := h.Current()
	inspected := b.Inspect()
	tabs := inspected["tabs"].(map[string]map[string]any)["page/panes/tabs"]
	expected := b.Geometry().Tabs["page/panes/tabs"].Header
	expected.X += 30
	expected.Y += 20
	if tabs["header"] != expected || tabs["nativeSelected"] != "one" || tabs["title"] != w.Title() {
		t.Fatal("embedded header projection/selection", tabs, expected)
	}
	for _, page := range tabs["pages"].([]map[string]any) {
		if r := page["clip"].(layout.Rect); r.W <= 0 || r.H <= 0 {
			t.Fatal("native page clip lost", page)
		}
	}
	control := inspected["controls"].(map[string]map[string]any)["page/panes/tabs/one/edit"]
	if r := control["clip"].(layout.Rect); r.X < 30 || r.Y < 20 || control["title"] != w.Title() {
		t.Fatal("main origin/title lost", control)
	}
}

func TestInspectionM2NativeMenuPathsAndClips(t *testing.T) {
	app := test.NewTempApp(t)
	w := app.NewWindow("Inspector menu")
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(600, 400))
	w.Show()
	disabled := fyne.NewMenuItem("Disabled", func() { t.Fatal("inspection dispatched action") })
	disabled.Disabled = true
	leaf := fyne.NewMenuItem("Checked", func() { t.Fatal("inspection dispatched child action") })
	leaf.Checked = true
	nested := &fyne.MenuItem{Label: "Nested", ChildMenu: fyne.NewMenu("Nested", leaf)}
	native := newNativeMenu(fyne.NewMenu("Root", disabled, fyne.NewMenuItemSeparator(), nested), w.Canvas(), func() {})
	native.show(fyne.NewPos(70, 40))
	defer native.close()
	b := &Bundle{owner: &DocumentHost{canvas: w.Canvas()}, surfaces: map[string]*nativeSurface{}, menus: map[string]*documentMenu{
		"page/menu": {native: native, path: "page/menu", rows: map[string]string{"/0": "page/menu/disabled", "/1": "page/menu/separator", "/2": "page/menu/nested", "/2/0": "page/menu/nested/leaf"}},
	}}
	root := b.inspectMenus()["page/menu"]
	if root["title"] != w.Title() || len(root["items"].([]map[string]any)) != 3 {
		t.Fatal("hidden submenu included or title missing", root)
	}
	native.popup.Items[2].(desktop.Hoverable).MouseIn(&desktop.MouseEvent{})
	items := b.inspectMenus()["page/menu"]["items"].([]map[string]any)
	if len(items) != 4 {
		t.Fatal("visible submenu missing", items)
	}
	if items[0]["enabled"] != false || items[1]["separator"] != true || items[1]["enabled"] != false || items[2]["submenu"] != true || items[3]["path"] != "page/menu/nested/leaf" || items[3]["checked"] != true {
		t.Fatal("native model/source mapping", items)
	}
	objects := append([]fyne.CanvasObject{}, native.popup.Items...)
	objects = append(objects, native.popup.Items[2].(nativeMenuItemChildren).Child().Items[0])
	for i, item := range items {
		r := item["clip"].(layout.Rect)
		if r.W <= 0 || r.H <= 0 {
			t.Fatal("native item clip empty", item)
		}
		point := fyne.NewPos(float32(r.X+r.W/2), float32(r.Y+r.H/2))
		if hit := native.hit(point); hit != objects[i] {
			t.Fatal("inspection center differs from native hit", item, hit)
		}
	}
	// Native menu scrolling/clipping must not invent reachable off-canvas rows.
	many := make([]*fyne.MenuItem, 30)
	for i := range many {
		many[i] = fyne.NewMenuItem(fmt.Sprint(i), func() {})
	}
	native.close()
	long := newNativeMenu(fyne.NewMenu("Long", many...), w.Canvas(), func() {})
	long.show(fyne.NewPos(80, 30))
	defer long.close()
	b.menus["page/menu"].native = long
	clipped := b.inspectMenus()["page/menu"]["items"].([]map[string]any)
	empty := 0
	for _, item := range clipped {
		r := item["clip"].(layout.Rect)
		if r.H == 0 {
			empty++
		}
		if r.H > 0 && (r.Y < 0 || r.Y+r.H > 400) {
			t.Fatal("menu clip outside canvas", item)
		}
	}
	if empty == 0 {
		t.Fatal("long menu omitted clipping evidence")
	}
	long.close()
	if len(b.inspectMenus()) != 0 {
		t.Fatal("dismissed menu remains inspectable")
	}
}
