package prototype

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// DocumentInput owns a local fixture snapshot. Guard performs bounded I/O; call
// it off the UI thread before dispatching a snapshot to the native owner.
// The maps/bytes must be treated as immutable after loading.
type DocumentInput struct {
	Candidate       reload.Candidate
	Choices         map[string][]ui.ChoiceOption
	SVGResources    map[string]markdown.PreparedSVG
	Size            layout.Size
	ContextRevision string
	guard           func() error
}

func (d *DocumentInput) Guard() error {
	if d == nil || d.guard == nil {
		return fmt.Errorf("stale: missing document snapshot")
	}
	return d.guard()
}
func digest(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }
func boundedFile(f *os.File, limit int) ([]byte, error) {
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("context: regular file required")
	}
	if info.Size() > int64(limit) {
		return nil, fmt.Errorf("context: file exceeds %d bytes", limit)
	}
	data, err := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err == nil && len(data) > limit {
		err = fmt.Errorf("context: file exceeds %d bytes", limit)
	}
	return data, err
}
func readBounded(path string, limit int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return boundedFile(f, limit)
}
func localResource(root *os.Root, path string) ([]byte, error) {
	// OpenRoot also enforces containment during symlink resolution, rather than
	// relying on a check followed by an unrestricted open vulnerable to replacement.
	if !filepath.IsLocal(path) || strings.Contains(path, "\\") {
		return nil, fmt.Errorf("context: resource must be a relative local path: %s", path)
	}
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == ".." {
			return nil, fmt.Errorf("context: resource traversal: %s", path)
		}
	}
	f, err := root.Open(path)
	if err != nil {
		return nil, err
	}
	return boundedFile(f, svgLimit)
}

type fixtureSnapshot struct {
	context   documentContext
	resources map[string][]byte
	revision  string
}

func readFixtures(path string) (fixtureSnapshot, error) {
	s := fixtureSnapshot{resources: map[string][]byte{}}
	if path == "" {
		return s, nil
	}
	data, err := readBounded(path, contextLimit)
	if err != nil {
		return s, err
	}
	s.context, err = parseDocumentContext(data)
	if err != nil {
		return s, err
	}
	root, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		return s, err
	}
	defer root.Close()
	paths := make([]string, 0, len(s.context.SVG))
	for path := range s.context.SVG {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	hash := sha256.New()
	fmt.Fprintf(hash, "context:%s\n", digest(data))
	total := 0
	for _, path := range paths {
		spec := s.context.SVG[path]
		resource, err := localResource(root, spec.File)
		if err != nil {
			return s, fmt.Errorf("context: SVG %s: %w", path, err)
		}
		total += len(resource)
		if total > resourcesLimit {
			return s, fmt.Errorf("context: SVG inventory exceeds 32 MiB")
		}
		s.resources[path] = resource
		fmt.Fprintf(hash, "%q:%q:%s\n", path, spec.File, digest(resource))
	}
	s.revision = fmt.Sprintf("%x", hash.Sum(nil))
	return s, nil
}

// LoadDocument reads only explicitly named source, context and local resources.
// It validates fixture identities, including descendants of hidden panes.
func LoadDocument(source, entry, revision, contextPath string) (*DocumentInput, error) {
	source, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	if contextPath != "" {
		contextPath, err = filepath.Abs(contextPath)
		if err != nil {
			return nil, err
		}
	}
	c := reload.Read(source, entry, 1, nil)
	if c.Err != nil {
		return nil, c.Err
	}
	if revision != "" && revision != c.Hash {
		return nil, fmt.Errorf("stale: SDUI source changed; regenerate the preview")
	}
	var unavailable error
	c.Root.Walk(func(n *parser.Instance) {
		if unavailable == nil && n.Kind == "widget" && (n.Widget == "tree" || n.Widget == "list") {
			unavailable = fmt.Errorf("unsupported-provider: standalone prototype requires an application-supplied collection provider: %w", &preparation.Diagnostic{Capability: preparation.Capability{Dimension: preparation.Provider, ID: "collection-data", Major: 1}, Path: n.Path, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...)})
		}
	})
	if unavailable != nil {
		return nil, unavailable
	}
	fixture, err := readFixtures(contextPath)
	if err != nil {
		return nil, err
	}
	d := &DocumentInput{Candidate: c, Size: layout.Size{W: 1280, H: 800}, ContextRevision: fixture.revision, Choices: map[string][]ui.ChoiceOption{}, SVGResources: map[string]markdown.PreparedSVG{}}
	if fixture.context.Viewport != nil {
		d.Size = layout.Size{W: fixture.context.Viewport.Width, H: fixture.context.Viewport.Height}
	}
	nodes := map[string]*parser.Instance{}
	c.Root.Walk(func(n *parser.Instance) { nodes[n.Path] = n })
	for path, options := range fixture.context.Choices {
		n := nodes[path]
		if n == nil || n.Kind != "widget" || n.Widget != "select" {
			return nil, fmt.Errorf("context: unknown/wrong-kind choices path %s", path)
		}
		d.Choices[path] = make([]ui.ChoiceOption, 0, len(options))
		for _, o := range options {
			d.Choices[path] = append(d.Choices[path], ui.ChoiceOption{ID: ui.ItemID(o.ID), Label: o.Label, Enabled: o.Enabled})
		}
	}
	for path, spec := range fixture.context.SVG {
		n := nodes[path]
		if n == nil || n.Kind != "widget" || n.Widget != "svg" {
			return nil, fmt.Errorf("context: unknown/wrong-kind SVG path %s", path)
		}
		policy, e := parser.PreviewOptions(n)
		if e != nil {
			return nil, e
		}
		if !policy.Explicit {
			return nil, fmt.Errorf("context: SVG path requires an explicit resource preview: %s", path)
		}
		ref, ok := n.Arguments["source"].(parser.Reference)
		if !ok {
			return nil, fmt.Errorf("context: SVG source reference required at %s", path)
		}
		resource, e := markdown.ValidateSVGResource(markdown.Resource{SVG: fixture.resources[path], Width: spec.Width, Height: spec.Height})
		if e != nil {
			return nil, fmt.Errorf("context: SVG %s: %w", path, e)
		}
		d.SVGResources[path] = markdown.PreparedSVG{Source: ref, ProviderID: "sdui-prototype-context", SHA256: digest(resource.SVG), Resource: resource}
	}
	for path, n := range nodes {
		if n.Widget == "select" {
			if _, ok := d.Choices[path]; !ok {
				return nil, fmt.Errorf("choice-options: missing %s", path)
			}
		}
	}
	// Reuse runtime's option/value validation without measuring or publishing GUI.
	session, e := ui.New("prototype-fixtures", c.Root)
	if e != nil {
		return nil, e
	}
	defer session.Close()
	if err = session.BindChoices(d.Choices); err != nil {
		return nil, err
	}
	sourceHash, contextHash := c.Hash, fixture.revision
	d.guard = func() error {
		data, e := readBounded(source, parser.MaxBytes)
		if e != nil {
			return e
		}
		if digest(data) != sourceHash {
			return fmt.Errorf("stale: SDUI source snapshot changed")
		}
		current, e := readFixtures(contextPath)
		if e != nil {
			return e
		}
		if current.revision != contextHash {
			return fmt.Errorf("stale: prototype context/resource snapshot changed")
		}
		return nil
	}
	if err = d.Guard(); err != nil {
		return nil, err
	}
	return d, nil
}
