package sdptool

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

// ActivityView renders producer-owned content into a revision-bound private bundle.
func ActivityView(ctx context.Context, p Project, uri, revision, output string) (Result, error) {
	if output == "" {
		return Result{}, failure("arguments", fmt.Errorf("Activity view requires output"))
	}
	if e := ctx.Err(); e != nil {
		return Result{}, e
	}
	if e := outside(output, p.Area); e != nil {
		return Result{}, failure("output", e)
	}
	out, e := physical(output)
	if e != nil {
		return Result{}, e
	}
	// A generated view must never replace or pollute the canonical source tree.
	if e = outside(p.Area, out); e != nil {
		return Result{}, failure("output", fmt.Errorf("Activity output must be outside SDP"))
	}
	page, e := renderActivityView(p, uri)
	if e != nil {
		return Result{}, e
	}
	if revision != "" && revision != page.Revision {
		return Result{}, failure("stale", fmt.Errorf("Activity sources changed; refresh inventory"))
	}
	result := Result{Schema: Version, Operation: "activity-view", Source: p.Area, Profile: "activity/1", Revision: page.Revision, Entry: filepath.Join(out, "entry.md"), Directory: out, Ownership: "caller-owned; remove directory after consumer release"}
	bundle := documents.Bundle{Files: map[string][]byte{}, Manifest: documents.Manifest{Version: "activity/1", Revision: page.Revision}}
	bundle.Put("entry.md", page.Markdown)
	for name, data := range page.Assets {
		bundle.Files[name] = data
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	bundle.Files["sdptool.json"] = append(data, '\n')
	bundle.Seal()
	total := 0
	for _, b := range bundle.Files {
		total += len(b)
	}
	if len(bundle.Files) > 2048 || total > 32<<20 {
		return Result{}, failure("limit", fmt.Errorf("Activity view exceeds 2048 files or 32 MiB; select a smaller view"))
	}
	if e = ctx.Err(); e != nil {
		return Result{}, e
	}
	fresh, e := activityCatalog(p)
	if e != nil {
		return Result{}, e
	}
	if fresh.Revision != page.Revision {
		return Result{}, failure("stale", fmt.Errorf("Activity sources changed while rendering"))
	}
	if e = bundle.Publish(out); e != nil {
		return Result{}, failure("output", e)
	}
	return result, nil
}
