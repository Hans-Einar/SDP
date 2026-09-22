package documents

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"os"
	"strings"
	"testing"
)

func TestClassDocuments(t *testing.T) {
	src, e := os.ReadFile("../examples/runtime-classes.sdl")
	if e != nil {
		t.Fatal(e)
	}
	a, e := ClassBundle(context.Background(), string(src), nil)
	if e != nil {
		t.Fatal(e)
	}
	b, e := ClassBundle(context.Background(), string(src), nil)
	if e != nil {
		t.Fatal(e)
	}
	for k, v := range a.Files {
		if !bytes.Equal(v, b.Files[k]) {
			t.Fatal("nondeterministic", k)
		}
	}
	d := a.Manifest.Diagrams[0]
	if len(d.Nodes) != 3 || len(d.Edges) != 2 || len(d.SourceFacts) != 8 || len(a.Manifest.Facts) != 8 {
		t.Fatal("source provenance", d)
	}
	for _, want := range []string{`class Runtime {`, `Runtime "0..* Client" o-- "1..* Modules" DomainModule`, `Runtime "1 Owner" *-- "0..* Widgets" Widget`} {
		if !strings.Contains(d.Mermaid(), want) {
			t.Fatal(want)
		}
	}
	exe := os.Getenv("MMDR")
	if exe == "" {
		t.Skip("set MMDR for class backend verification")
	}
	r, e := NewMmdr(exe)
	if e != nil {
		t.Fatal(e)
	}
	b, e = ClassBundle(context.Background(), string(src), r)
	if e != nil {
		t.Fatal(e)
	}
	dec := xml.NewDecoder(bytes.NewReader(b.Files["diagrams/CL01-classes.svg"]))
	var labels strings.Builder
	open, filled := 0, 0
	for {
		tok, e := dec.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal(e)
		}
		switch x := tok.(type) {
		case xml.CharData:
			labels.Write(x)
		case xml.StartElement:
			if x.Name.Local == "polygon" {
				attrs := map[string]string{}
				for _, a := range x.Attr {
					attrs[a.Name.Local] = a.Value
				}
				if attrs["points"] == "0,0 9,6 18,0 9,-6" {
					if attrs["fill"] == "none" {
						open++
					} else {
						filled++
					}
				}
			}
		}
	}
	if open != 1 || filled != 1 {
		t.Fatal("missing ownership markers", open, filled)
	}
	for _, want := range []string{"Runtime", "DomainModule", "Widget", "Client", "Modules", "Owner", "Widgets", "Revision", "Reload"} {
		if !strings.Contains(labels.String(), want) {
			t.Fatal("lost class text", want)
		}
	}
	if strings.Contains(labels.String(), "n_") || strings.Contains(labels.String(), "[\"") {
		t.Fatal("phantom alias nodes")
	}
}
