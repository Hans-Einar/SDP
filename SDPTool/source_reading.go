package sdptool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"net/url"
	"path/filepath"
	"strings"
)

func selectSource(ctx context.Context, p Project, m Model, entry string, o PreviewOptions) (Result, error) {
	u, e := url.Parse(o.URI)
	if e != nil {
		return Result{}, failure("selection", e)
	}
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil {
		return Result{}, failure("selection", e)
	}
	if u.Scheme != "sdl-view" || u.Host != p.Inventory.ProjectID || u.User != nil || u.Port() != "" || u.Fragment != "" || u.Path != "/VP11" || q.Get("projection") != "source" {
		return Result{}, failure("selection", fmt.Errorf("invalid source selection"))
	}
	for k, v := range q {
		if len(v) != 1 || !strings.Contains("|projection|source|target|consumer|", "|"+k+"|") {
			return Result{}, failure("selection", fmt.Errorf("invalid source query"))
		}
	}
	if (q.Get("target") != "" && q.Get("target") != "main") || (q.Get("consumer") != "" && q.Get("consumer") != "xfmd") || o.Output == "" || o.Revision == "" {
		return Result{}, failure("arguments", fmt.Errorf("source selection requires revision and output"))
	}
	source := q.Get("source")
	checked := true
	v, input, loadErr := loadModel(entry)
	var links []SourceLink
	copies := map[string][]byte{}
	revision := ""
	text := ""
	assets := map[string][]byte{}
	if loadErr == nil {
		revision = v.Revision
		if revision != o.Revision {
			return Result{}, failure("stale", fmt.Errorf("model changed; refresh"))
		}
		if e = input.Outside(o.Output); e != nil {
			return Result{}, failure("output", e)
		}
		links, copies, e = bindSourceLinks(v, entry)
		if e != nil {
			return Result{}, e
		}
		if source == "" {
			text = "# Source files\n\nThese files compose the selected checked model. The same physical file can belong to several distinct System contexts.\n\n"
			for _, s := range v.Sources {
				t := sourceTarget(p, m, revision, s.Path)
				text += fmt.Sprintf("- [%s](%s)\n", s.Path, t.URI)
			}
			text += "\n## Include relationships\n\n"
			for _, edge := range v.SourceEdges {
				text += fmt.Sprintf("- `%s` includes `%s`.\n", edge.From, edge.To)
			}
		} else {
			found := -1
			for i, s := range v.Sources {
				if s.Path == source {
					found = i
				}
			}
			if found < 0 {
				return Result{}, failure("selection", fmt.Errorf("source is outside selected model"))
			}
			raw := copies[links[found].Snapshot]
			syntax, err := parser.Parse(string(raw))
			if err != nil {
				return Result{}, failure("source", err)
			}
			text, assets["source.svg"] = documents.SourceReading(syntax, source, links[found].ID, checked)
		}
	} else {
		for _, known := range p.Sources {
			if known.ID == m.ID && known.State == "validated" {
				return Result{}, failure("stale", fmt.Errorf("checked model is no longer valid; refresh: %w", loadErr))
			}
		}
		if source != "" {
			return Result{}, failure("selection", fmt.Errorf("unassigned fragments cannot select external sources"))
		}
		raw, err := readSource(entry)
		if err != nil {
			return Result{}, err
		}
		revision = documents.Hash(raw)
		if revision != o.Revision {
			return Result{}, failure("stale", fmt.Errorf("source changed; refresh"))
		}
		syntax, err := parser.Parse(string(raw))
		if err != nil {
			return Result{}, failure("source", err)
		}
		if syntax.Header.Version != "0.6" {
			return Result{}, loadErr
		}
		for _, d := range syntax.Declarations {
			if d.Kind == "system" {
				return Result{}, loadErr
			}
		}
		if err = outside(o.Output, entry); err != nil {
			return Result{}, failure("output", err)
		}
		links = []SourceLink{{ID: "s0", Path: entry, Snapshot: "source-inputs/s0.design", Revision: revision}}
		copies[links[0].Snapshot] = raw
		text, assets["source.svg"] = documents.SourceReading(syntax, filepath.Base(entry), "s0", false)
	}
	if e = ctx.Err(); e != nil {
		return Result{}, e
	}
	output, e := filepath.Abs(o.Output)
	if e != nil {
		return Result{}, e
	}
	result := Result{SourceLinks: links, Schema: Version, Operation: "select", Source: entry, Profile: m.Profile, Revision: revision, Entry: filepath.Join(output, "entry.md"), Directory: output, Ownership: "caller-owned; remove directory after consumer release"}
	b := &documents.Bundle{Files: map[string][]byte{}, Manifest: documents.Manifest{Version: readingPresentation, Revision: revision}}
	b.Put("entry.md", text)
	for k, v := range copies {
		b.Files[k] = v
	}
	for k, v := range assets {
		b.Files[k] = v
	}
	raw, _ := json.MarshalIndent(result, "", "  ")
	b.Files["sdptool.json"] = raw
	b.Seal()
	size := 0
	for _, data := range b.Files {
		size += len(data)
	}
	if len(b.Files) > 128 || size > 32<<20 {
		return Result{}, failure("limit", fmt.Errorf("source bundle exceeds 128 files or 32 MiB"))
	}
	if e = b.CheckLinks(); e != nil {
		return Result{}, e
	}
	if input != nil {
		e = input.Fresh()
	} else {
		var raw []byte
		raw, e = readSource(entry)
		if e == nil && documents.Hash(raw) != revision {
			e = fmt.Errorf("source changed")
		}
	}
	if e != nil {
		return Result{}, failure("stale", e)
	}
	if e = b.Publish(output); e != nil {
		return Result{}, failure("output", e)
	}
	return result, nil
}
