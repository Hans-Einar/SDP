package collections

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func controllerFixture(t *testing.T, source string, empty bool, log func(string, any)) (*Controller, chan func()) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow(WindowTitle)
	f, err := New()
	if err != nil {
		t.Fatal(err)
	}
	f.Log = log
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1000, H: 650})
	posts := make(chan func(), 32)
	h.Post = func(fn func()) { posts <- fn }
	c := &Controller{Fixture: f, Host: h, Source: source, Empty: empty, Sequence: 1, Resize: h.Resize, Log: log}
	c.Close = func() { h.Close(); f.Close() }
	t.Cleanup(func() { c.Close(); w.Close() })
	r, err := f.Request(source, 1, empty)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	return c, posts
}
func collectionState(t *testing.T, c *Controller, path string) ui.CollectionState {
	t.Helper()
	state, ok := c.Host.Current().Session.Snapshot().Collections[path]
	if !ok {
		t.Fatalf("missing %s", path)
	}
	return state
}

func TestNestedFixtureGeometryMutationAndIdentity(t *testing.T) {
	c, _ := controllerFixture(t, NestedSource, false, nil)
	g := c.Host.Current().Geometry()
	if g.Viewports["page/body"].Maximum.Y <= 0 || g.Viewports["page/body/entries"].Parent != "page/body" {
		t.Fatalf("not nested scroll geometry: %+v", g.Viewports)
	}
	state := collectionState(t, c, "page/body/entries")
	target, err := c.Host.Current().Session.Target(state.Handle, "entry-098")
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Host.Mutate(func(s *ui.Session) error { return s.SelectItem(target) }); err != nil {
		t.Fatal(err)
	}
	if err = c.Host.Mutate(func(s *ui.Session) error {
		return s.SetViewports(map[string]ui.ViewportState{"page/body/entries": {Y: 10000}, "page/body": {Y: 50}})
	}); err != nil {
		t.Fatal(err)
	}
	if err = c.Run("reorder list"); err != nil {
		t.Fatal(err)
	}
	current := collectionState(t, c, "page/body/entries")
	if current.Selected != "entry-098" || current.Data.Items[0].ID != "entry-099" || current.Generation != state.Generation+1 {
		t.Fatalf("reorder lost identity: %+v", current)
	}
	if err = c.Host.Current().Session.ValidateCollectionTarget(target); err == nil {
		t.Fatal("old target survived reorder")
	}
	if err = c.Run("shrink list 1"); err != nil {
		t.Fatal(err)
	}
	current = collectionState(t, c, "page/body/entries")
	if len(current.Data.Items) != 1 || current.Selected != "" || c.Host.Current().Session.Snapshot().Viewports["page/body/entries"].Y != 0 {
		t.Fatal("shrink selection/clamp")
	}
	if err = c.Run("reset list"); err != nil {
		t.Fatal(err)
	}
	if len(collectionState(t, c, "page/body/entries").Data.Items) != 100 {
		t.Fatal("reset")
	}
	if err = c.Run("hide tree"); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Host.Current().Geometry().Viewports["page/body/nav"]; ok {
		t.Fatal("hidden viewport remains")
	}
	if err = c.Run("show tree"); err != nil {
		t.Fatal(err)
	}
	if err = c.Run("disable tree"); err != nil {
		t.Fatal(err)
	}
	w, _ := c.Host.Current().Session.Widget("page/body/nav")
	if w.Enabled {
		t.Fatal("disable missing")
	}
	if err = c.Run("enable tree"); err != nil {
		t.Fatal(err)
	}
	if err = c.Run("resize 900 600"); err != nil {
		t.Fatal(err)
	}
	if c.Host.Current().Geometry().Root.Rect.W != 900 {
		t.Fatal("resize geometry")
	}
	if c.Fixture.Calls() != 0 {
		t.Fatal("fixture mutation invoked SDL action")
	}
	inspected := c.Inspect()
	if inspected["rows"] == nil || inspected["viewports"] == nil || inspected["nested"] != true {
		t.Fatal("missing coordinate logs")
	}
}

func TestFixtureFailuresPreservePublishedBundle(t *testing.T) {
	c, _ := controllerFixture(t, Source, false, nil)
	old := c.Host.Current()
	for _, stage := range []string{"profile", "provider", "binding", "layout", "resource", "guard"} {
		t.Run(stage, func(t *testing.T) {
			before := old.Session.Snapshot()
			geometry := old.Geometry()
			controls := old.Controls()["page/body/nav"]
			if err := c.Run("fail " + stage); err == nil {
				t.Fatalf("%s injection succeeded", stage)
			}
			if c.Host.Current() != old || !reflect.DeepEqual(before, old.Session.Snapshot()) || old.Geometry() != geometry || old.Controls()["page/body/nav"] != controls {
				t.Fatalf("%s mutated live bundle", stage)
			}
		})
	}
	if err := c.Run("fail stale"); err == nil {
		t.Fatal("stale candidate published")
	}
	if c.Host.Current() != old {
		t.Fatal("stale failure changed bundle")
	}
	preview, _ := old.Session.Widget("page/footer/preview")
	if preview.Draft != "New draft during detached preparation" {
		t.Fatal("stale failure overwrote newer draft")
	}
	if err := c.Run("reload"); err != nil {
		t.Fatal(err)
	}
	if c.Host.Current() == old || !old.Session.Closed() {
		t.Fatal("valid reload did not retire old bundle")
	}
	if c.Fixture.Calls() != 0 {
		t.Fatal("failure/preparation invoked SDL action")
	}
}

func TestFixtureCommandValidationAndTooSmallGeometry(t *testing.T) {
	c, _ := controllerFixture(t, Source, false, nil)
	before := c.Host.Current().Session.Snapshot()
	for _, line := range []string{"", "no-such-command", "resize NaN 5", "resize 0 5", "resize 400 Inf", "resize 400 500 extra", "shrink list -1", "shrink list 101", "hide missing", "fail unknown", "cancel", "reorder", "state extra"} {
		if err := c.Run(line); err == nil {
			t.Errorf("malformed command succeeded: %s", line)
		}
		if !reflect.DeepEqual(before, c.Host.Current().Session.Snapshot()) {
			t.Fatalf("bad command changed state: %s", line)
		}
	}
	geometry := c.Host.Current().Geometry()
	if err := c.Run("resize 1 1"); err == nil {
		t.Fatal("too-small geometry accepted")
	}
	if !reflect.DeepEqual(before, c.Host.Current().Session.Snapshot()) || geometry != c.Host.Current().Geometry() {
		t.Fatal("too-small resize lost last valid presentation")
	}
	if err := c.Run("resize 1000 650"); err != nil {
		t.Fatal(err)
	}
}

func TestExplicitCompletionBarriersAndLateResults(t *testing.T) {
	loads := make(chan string, 8)
	c, posts := controllerFixture(t, Source, true, func(event string, value any) {
		if event == "load" {
			loads <- value.(map[string]any)["key"].(string)
		}
	})
	first := <-loads
	if len(c.Fixture.Pending()) != 1 {
		t.Fatal("missing started barrier")
	}
	if err := c.Run("complete " + first + " nonsense"); err == nil {
		t.Fatal("unknown outcome accepted")
	}
	if len(c.Fixture.Pending()) != 1 {
		t.Fatal("bad outcome consumed request")
	}
	if err := c.Run("cancel tree"); err != nil {
		t.Fatal(err)
	}
	canceled := collectionState(t, c, "page/body/nav")
	if err := c.Run("complete " + first + " success"); err != nil {
		t.Fatal(err)
	}
	(<-posts)()
	after := collectionState(t, c, "page/body/nav")
	if !reflect.DeepEqual(canceled, after) {
		t.Fatal("late completion changed canceled state")
	}
	if err := c.Host.Mutate(func(s *ui.Session) error {
		target, e := s.Target(after.Handle, "")
		if e != nil {
			return e
		}
		return s.RetryItem(target)
	}); err != nil {
		t.Fatal(err)
	}
	second := <-loads
	if second == first {
		t.Fatal("request key reused")
	}
	if err := c.Run("complete " + second + " error"); err != nil {
		t.Fatal(err)
	}
	(<-posts)()
	failed := collectionState(t, c, "page/body/nav")
	if failed.Status[""].Phase != ui.LoadError {
		t.Fatal("error completion not visible")
	}
	if err := c.Host.Mutate(func(s *ui.Session) error {
		target, e := s.Target(failed.Handle, "")
		if e != nil {
			return e
		}
		return s.RetryItem(target)
	}); err != nil {
		t.Fatal(err)
	}
	third := <-loads
	if err := c.Run("complete " + third + " empty"); err != nil {
		t.Fatal(err)
	}
	(<-posts)()
	current := collectionState(t, c, "page/body/nav")
	if !current.RootLoaded || current.Request != nil || len(current.Data.Items) != 0 || current.AutoLoadPending {
		t.Fatal("empty success restarted or was not accepted")
	}
	if c.Fixture.Calls() != 0 {
		t.Fatal("provider lifecycle invoked SDL")
	}
}
