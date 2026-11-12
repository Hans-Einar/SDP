package documents

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

func TestReadingBundlesPreserveContextAndSeparateRegisters(t *testing.T) {
	v := model(t)
	v.ReadingProject = "trial"
	for _, spec := range viewpoint.Catalog {
		queries := []viewpoint.Query{{Viewpoint: spec.ID, Projection: "reading"}}
		for _, c := range v.ReadingCollections(spec.ID) {
			queries = append(queries, viewpoint.Query{Viewpoint: spec.ID, Projection: "reading", Collection: c.ID})
			if len(c.Names) > 0 {
				queries = append(queries, viewpoint.Query{Viewpoint: spec.ID, Projection: "reading", Focus: c.Names[0]})
			}
		}
		for _, q := range queries {
			b, e := Selected(context.Background(), v, q, nil)
			if e != nil {
				t.Fatal(q, e)
			}
			text := string(b.Files["entry.md"])
			if strings.Contains(text, v.Revision) || strings.Contains(text, "## Complete fact registry") || strings.Contains(text, "## Declaration register") {
				t.Fatal("identity/register in reading page", q)
			}
			for _, f := range v.Facts {
				if strings.Contains(text, f.S("id")) {
					t.Fatal("visible fact ID", q)
				}
			}
			var saved viewpoint.Query
			if e = json.Unmarshal(b.Files["selection.json"], &saved); e != nil || saved.Viewpoint != q.Viewpoint || saved.Collection != q.Collection || saved.Focus != q.Focus {
				t.Fatal("context lost", saved, e)
			}
			if e = b.CheckLinks(); e != nil {
				t.Fatal(e)
			}
		}
	}
	for _, reg := range []string{"facts", "declarations"} {
		b, e := Selected(context.Background(), v, viewpoint.Query{Viewpoint: "VP11", Projection: "reading", Register: reg}, nil)
		if e != nil {
			t.Fatal(e)
		}
		text := string(b.Files["entry.md"])
		if strings.Contains(text, "| ID |") || strings.Contains(text, v.Revision) || !strings.Contains(text, "complete checked model") {
			t.Fatal(text)
		}
	}
}
