package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/snapshot"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	source := flag.String("ui", "", "SDUI file")
	entry := flag.String("entry", "page", "Entry definition")
	state := flag.String("state", "", "Optional JSON widget state")
	design := flag.String("design", "", "Optional validated SDL design model")
	output := flag.String("output", "", "Document directory")
	renderer := flag.String("renderer", "", "Optional registered mmdr")
	w := flag.Float64("width", 1000, "Logical export width")
	h := flag.Float64("height", 650, "Logical export height")
	flag.Parse()
	if *source == "" || *output == "" || flag.NArg() != 0 {
		check(fmt.Errorf("require -ui and -output"))
	}
	target, e := filepath.Abs(*output)
	check(e)
	read := func(name string) []byte {
		p, e := filepath.Abs(name)
		check(e)
		rel, e := filepath.Rel(target, p)
		check(e)
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			check(fmt.Errorf("source must be outside output"))
		}
		f, e := os.Open(p)
		check(e)
		defer f.Close()
		b, e := io.ReadAll(io.LimitReader(f, (2<<20)+1))
		check(e)
		if len(b) > 2<<20 {
			check(fmt.Errorf("input limit"))
		}
		return b
	}
	o := snapshot.Options{Entry: *entry, Width: *w, Height: *h, State: map[string]snapshot.WidgetState{}}
	if *state != "" {
		d := json.NewDecoder(bytes.NewReader(read(*state)))
		d.DisallowUnknownFields()
		check(d.Decode(&o.State))
		var extra any
		if d.Decode(&extra) != io.EOF {
			check(fmt.Errorf("trailing state JSON"))
		}
	}
	if *design != "" {
		o.Design = string(read(*design))
	}
	if *renderer != "" {
		p, e := filepath.Abs(*renderer)
		check(e)
		f, e := os.ReadFile(p)
		check(e)
		o.Renderer = markdown.Mmdr{Executable: p}
		o.RendererID = "mmdr:" + documents.Hash(f)
	}
	b, e := snapshot.Build(context.Background(), string(read(*source)), o)
	check(e)
	check(b.Publish(target))
	fmt.Println(filepath.Join(target, "entry.md"))
}
func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
