// Package documents builds complete deterministic Markdown bundles in memory.
package documents

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"net/url"
	"path"
	"sort"
	"strings"
)

type Renderer interface {
	Render(context.Context, string) ([]byte, error)
	Identity() string
}
type Options struct {
	Navigator  bool
	Monolithic bool
	Project    string
	Viewpoints []string
	Renderer   Renderer
}
type Bundle struct {
	Files    map[string][]byte
	Manifest Manifest
}
type Manifest struct {
	Version  string              `json:"version"`
	Revision string              `json:"revision"`
	Renderer string              `json:"renderer"`
	Diagrams []viewpoint.Diagram `json:"diagrams"`
	Facts    []viewpoint.Fact    `json:"facts"`
	Gaps     []viewpoint.Fact    `json:"model_gaps"`
	Outputs  map[string]string   `json:"outputs"`
}

func Hash(b []byte) string           { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func (b *Bundle) Put(name, s string) { b.Files[name] = []byte(s) }
func (b *Bundle) Seal() {
	b.Manifest.Outputs = map[string]string{}
	for n, f := range b.Files {
		if n != "manifest.json" {
			b.Manifest.Outputs[n] = Hash(f)
		}
	}
	j, _ := json.MarshalIndent(b.Manifest, "", "  ")
	b.Files["manifest.json"] = append(j, '\n')
}
func Build(ctx context.Context, v *viewpoint.Views, o Options) (*Bundle, error) {
	selected := map[string]bool{}
	if len(o.Viewpoints) == 0 {
		for _, s := range viewpoint.Catalog {
			selected[s.ID] = true
		}
	} else {
		for _, id := range o.Viewpoints {
			if _, ok := viewpoint.SpecFor(id); !ok {
				return nil, fmt.Errorf("unknown viewpoint %s", id)
			}
			selected[id] = true
		}
	}
	b := &Bundle{Files: map[string][]byte{}, Manifest: Manifest{Version: viewpoint.Version, Revision: v.Revision, Diagrams: []viewpoint.Diagram{}, Facts: v.Facts, Gaps: v.Gaps}}
	if o.Renderer != nil {
		b.Manifest.Renderer = o.Renderer.Identity()
	}
	if o.Project == "" {
		o.Project = "local"
	}
	if strings.ContainsAny(o.Project, "/?#@:") {
		return nil, fmt.Errorf("invalid project ID")
	}
	nav := "# SDL — navigator\n\nRevisjon: `" + v.Revision + "`.\n\n[Oversikt](index.md)\n\n"
	full := "# SDL — genererte viewpoints\n\nStruktur og kildepåstander; ikke observert kjøring.\n\n"
	index := "# SDL — designoversikt\n\n[Navigator](navigator.md)\n\n"
	for i := 0; i < 6; i++ {
		level := fmt.Sprintf("A%d", i)
		title := []string{"Behov og forpliktelser", "Funksjonell hensikt", "System og containere", "Interne enheter", "Detaljdesign og kontrakter", "Realisering og bevis"}[i]
		index += "- [" + level + " — " + title + "](" + level + "/index.md)\n"
		content := "# " + level + " — " + title + "\n\n[Oversikt](../index.md) · [Typeinventar](inventory.md)\n\nNivåene klassifiserer visninger; objektenes eget nivå er uspesifisert.\n\n"
		for _, s := range viewpoint.Catalog {
			if selected[s.ID] && contains(s.Levels, level) {
				content += "- [" + s.ID + " — " + s.Title + "](../navigator.md#" + strings.ToLower(s.ID) + ")\n"
			}
		}
		b.Put(level+"/index.md", content)
		b.Put(level+"/inventory.md", inventory(v, level))
	}
	for _, s := range viewpoint.Catalog {
		if !selected[s.ID] {
			continue
		}
		nav += "## " + s.ID + "\n\n" + s.Title + ". " + s.Note + "\n\n"
		vpdir := "viewpoints/" + s.ID
		page := "# " + s.ID + " — " + s.Title + "\n\n[Navigator](../../navigator.md)\n\n" + s.Note + "\n\n"
		full += "## " + s.ID + " — " + s.Title + "\n\n" + s.Note + "\n\n"
		if o.Navigator {
			nav += "[Åpne ved behov](sdl-view://" + o.Project + "/" + s.ID + "?target=main&consumer=xfmd)\n\n"
			for _, d := range v.Diagrams {
				if strings.HasPrefix(d.ID, s.ID+"-") {
					nav += "- [" + d.Title + "](sdl-view://" + o.Project + "/" + s.ID + "?diagram=" + url.QueryEscape(d.ID) + "&target=main&consumer=xfmd)\n"
				}
			}
			nav += "\n"
			continue
		}
		nav += "[Åpne viewpoint](" + vpdir + "/index.md)\n\n"
		for _, d := range v.Diagrams {
			if !strings.HasPrefix(d.ID, s.ID+"-") {
				continue
			}
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			b.Manifest.Diagrams = append(b.Manifest.Diagrams, d)
			b.Put("diagrams/"+d.ID+".mmd", d.Mermaid())
			if o.Renderer != nil {
				svg, e := o.Renderer.Render(ctx, d.Mermaid())
				if e != nil {
					return nil, fmt.Errorf("%s: %w", d.ID, e)
				}
				b.Files["diagrams/"+d.ID+".svg"] = svg
			}
			page += "- [" + d.Title + "](" + d.ID + ".md)\n"
			b.Put(vpdir+"/"+d.ID+".md", "# "+d.Title+"\n\n[Viewpoint](index.md) · [Navigator](../../navigator.md)\n\nRevisjon: `"+v.Revision+"`.\n\n"+v.DiagramMarkdown(d, "../../diagrams/", o.Renderer != nil))
			full += v.DiagramMarkdown(d, "diagrams/", o.Renderer != nil)
		}
		page += "\n" + v.Tables(s.ID)
		full += v.Tables(s.ID)
		b.Put(vpdir+"/index.md", page)
	}
	if !o.Navigator {
		if selected["VP06"] {
			b.Put("implementation.md", v.Implementation())
			index += "\n[Utviklingsplan G-faser](implementation.md)\n"
		}
		if selected["VP08"] {
			j, _ := json.MarshalIndent(v.MessageSets, "", "  ")
			b.Files["message-sets.json"] = append(j, '\n')
		}
		if o.Monolithic {
			b.Put("viewpoints.md", full)
			index += "\n[Samlerapport](viewpoints.md)\n"
		}
	}
	b.Put("navigator.md", nav)
	b.Put("index.md", index+"\nRequirement, System og State støttes ikke av design-core 0.5. Mode er driftskontekst, ikke State.\n")
	b.Seal()
	if e := b.CheckLinks(); e != nil {
		return nil, e
	}
	return b, nil
}
func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
func inventory(v *viewpoint.Views, level string) string {
	s := "# Typeinventar — " + level + "\n\n[Opp](index.md)\n\nSamme modell-ID på alle nivåer; eget objektnivå er uspesifisert.\n\n"
	kinds := map[string][]string{}
	for n, k := range v.Kinds {
		kinds[k] = append(kinds[k], n)
	}
	ks := []string{}
	for k := range kinds {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		sort.Strings(kinds[k])
		s += fmt.Sprintf("## %s (%d)\n\n", k, len(kinds[k]))
		for _, n := range kinds[k] {
			s += "- <a id=\"" + n + "\"></a>" + n + "\n"
		}
		s += "\n"
	}
	return s
}

// LocalPath accepts a bundle-relative URL, never an absolute or escaping path.
func LocalPath(from, link string) (string, error) {
	u, e := url.Parse(link)
	if e != nil {
		return "", e
	}
	if u.Scheme != "" || u.Host != "" {
		return "", nil
	}
	p := path.Clean(path.Join(path.Dir(from), u.Path))
	if strings.HasPrefix(u.Path, "/") || p == ".." || strings.HasPrefix(p, "../") {
		return "", fmt.Errorf("escaping link %s", link)
	}
	return p, nil
}
