package documents

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"reflect"
	"strings"
	"testing"
)

func TestSingleBundleMatchesFullDiagram(t *testing.T) {
	v := model(t)
	all, e := Build(context.Background(), v, Options{})
	if e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{"VP02-roots", "VP08-SelectedViewOpened"} {
		single, e := Selected(context.Background(), v, viewpoint.Query{Viewpoint: id[:4], Diagram: id}, nil)
		if e != nil {
			t.Fatal(e)
		}
		if !reflect.DeepEqual(single.Files["diagrams/"+id+".mmd"], all.Files["diagrams/"+id+".mmd"]) {
			t.Fatal("different projection", id)
		}
		if len(single.Manifest.Diagrams) != 1 {
			t.Fatal("unselected diagrams generated")
		}
		if single.Manifest.Revision != all.Manifest.Revision {
			t.Fatal("different revision")
		}
	}
}

func TestNavigatorSelectedLinks(t *testing.T) {
	b, e := Build(context.Background(), model(t), Options{Navigator: true, Viewpoints: []string{"VP02"}})
	if e != nil {
		t.Fatal(e)
	}
	s := string(b.Files["navigator.md"])
	if !strings.Contains(s, "diagram=VP02-roots") || strings.Contains(s, "/VP08?") {
		t.Fatal("wrong navigator selection")
	}
}
