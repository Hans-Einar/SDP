// Package markdown implements the bounded SDUI Markdown presentation profile.
// Parsing is independent of diagram rendering and never loads URLs or executes HTML.
package markdown

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Block struct {
	Text    string
	Scale   float64
	Strong  bool
	Code    bool
	Diagram string
	Table   [][]string
}
type Diagram struct{ ID, Source string }
type Document struct {
	Blocks   []Block
	Diagrams []Diagram
}

func Parse(source string) (*Document, error) {
	if len(source) > 32768 {
		return nil, fmt.Errorf("markdown-limit: widget exceeds 32768 bytes")
	}
	raw := []byte(source)
	tree := goldmark.New(goldmark.WithExtensions(extension.Table)).Parser().Parse(text.NewReader(raw))
	doc := &Document{}
	var plain func(ast.Node) string
	plain = func(n ast.Node) string {
		var b strings.Builder
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			switch x := c.(type) {
			case *ast.Text:
				value := x.Segment.Value(raw)
				if !x.IsRaw() {
					value = util.ResolveEntityNames(util.ResolveNumericReferences(util.UnescapePunctuations(value)))
				}
				b.Write(value)
				if x.SoftLineBreak() || x.HardLineBreak() {
					b.WriteByte('\n')
				}
			case *ast.String:
				b.Write(x.Value)
			default:
				b.WriteString(plain(c))
			}
		}
		return b.String()
	}
	lines := func(n ast.Node) string {
		var b strings.Builder
		for i := 0; i < n.Lines().Len(); i++ {
			s := n.Lines().At(i)
			b.Write(s.Value(raw))
		}
		return b.String()
	}
	err := ast.Walk(tree, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		switch n.(type) {
		case *ast.HTMLBlock, *ast.RawHTML, *ast.Image:
			return ast.WalkStop, fmt.Errorf("markdown-unsupported: HTML and images require a registered resource provider")
		}
		if len(doc.Blocks) > 256 {
			return ast.WalkStop, fmt.Errorf("markdown-limit: more than 256 blocks")
		}
		switch x := n.(type) {
		case *ast.Heading:
			doc.Blocks = append(doc.Blocks, Block{Text: plain(n), Scale: 1 + float64(7-x.Level)*.12, Strong: true})
			return ast.WalkSkipChildren, nil
		case *ast.Paragraph, *ast.TextBlock:
			t := plain(n)
			if _, ok := n.Parent().(*ast.ListItem); ok {
				t = "• " + t
			}
			doc.Blocks = append(doc.Blocks, Block{Text: t, Scale: 1})
		case *ast.FencedCodeBlock:
			source := lines(n)
			if string(x.Language(raw)) == "mermaid" {
				hash := sha256.Sum256([]byte(source))
				id := fmt.Sprintf("mermaid-%x", hash[:12])
				doc.Diagrams = append(doc.Diagrams, Diagram{id, source})
				doc.Blocks = append(doc.Blocks, Block{Diagram: id, Scale: 1})
				if len(doc.Diagrams) > 8 {
					return ast.WalkStop, fmt.Errorf("markdown-limit: more than eight diagrams")
				}
			} else {
				doc.Blocks = append(doc.Blocks, Block{Text: strings.TrimSuffix(source, "\n"), Scale: 1, Code: true})
			}
			return ast.WalkSkipChildren, nil
		case *ast.CodeBlock:
			doc.Blocks = append(doc.Blocks, Block{Text: strings.TrimSuffix(lines(n), "\n"), Scale: 1, Code: true})
			return ast.WalkSkipChildren, nil
		case *extast.Table:
			rows := [][]string{}
			for r := n.FirstChild(); r != nil; r = r.NextSibling() {
				cells := []string{}
				for c := r.FirstChild(); c != nil; c = c.NextSibling() {
					cells = append(cells, plain(c))
				}
				rows = append(rows, cells)
			}
			doc.Blocks = append(doc.Blocks, Block{Table: rows, Scale: 1})
			return ast.WalkSkipChildren, nil
		case *ast.ThematicBreak:
			doc.Blocks = append(doc.Blocks, Block{Text: "────────────────", Scale: 1})
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	// Check inline resources even when heading/table traversal skipped children.
	err = ast.Walk(tree, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if enter {
			switch n.(type) {
			case *ast.HTMLBlock, *ast.RawHTML, *ast.Image:
				return ast.WalkStop, fmt.Errorf("markdown-unsupported: HTML and images are not enabled")
			}
		}
		return ast.WalkContinue, nil
	})
	return doc, err
}
