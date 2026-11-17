//go:build desktop

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	fynetest "fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"github.com/Hans-Einar/SDP/SDUI/go/prototype"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestDocumentWatcherIncludesResourcesAndRejectsInvalidSnapshots(t *testing.T) {
	dir := t.TempDir()
	write := func(name, data string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	source := write("page.sdui", `sdui 0.3; ref: art "unopened"; Main=[s=svg(art.X.@resource,description="Image",fallback="reject")];`)
	contextPath := write("context.json", `{"schema":"sdui-prototype-context/1","svg":{"Main/s":{"file":"image.svg","width":100,"height":100}}}`)
	write("image.svg", `<svg viewBox="0 0 100 100"><rect x="0" y="0" width="20" height="20"/></svg>`)
	initial, err := prototype.LoadDocument(source, "Main", "", contextPath)
	if err != nil {
		t.Fatal(err)
	}
	type update struct {
		input      *prototype.DocumentInput
		generation uint64
		err        error
	}
	updates := make(chan update, 8)
	var observed atomic.Uint64
	observed.Store(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		watchDocuments(ctx, source, "Main", contextPath, initial, &observed, func(input *prototype.DocumentInput, generation uint64, err error) {
			updates <- update{input, generation, err}
		})
	}()
	next := func() update {
		t.Helper()
		select {
		case u := <-updates:
			return u
		case <-time.After(3 * time.Second):
			t.Fatal("no watcher update")
			return update{}
		}
	}
	write("image.svg", `<svg viewBox="0 0 100 100"><circle cx="10" cy="10" r="5"/></svg>`)
	resource := next()
	if resource.err != nil || resource.generation != 2 || resource.input.ContextRevision == initial.ContextRevision || resource.input.Candidate.Hash != initial.Candidate.Hash {
		t.Fatal(resource)
	}
	// Context-only edits also invalidate the queued resource generation.
	write("context.json", `{"schema":"sdui-prototype-context/1","viewport":{"width":1680,"height":1150},"svg":{"Main/s":{"file":"image.svg","width":100,"height":100}}}`)
	fixture := next()
	if fixture.err != nil || fixture.generation != 3 || fixture.input.Size.W != 1680 || observed.Load() == resource.generation {
		t.Fatal(fixture)
	}
	write("page.sdui", `broken source`)
	invalid := next()
	if invalid.err == nil || invalid.input != nil || invalid.generation != 4 {
		t.Fatal(invalid)
	}
	// The last valid snapshot's pure model/resources remain available to its host;
	// no invalid candidate is substituted, and the old disk guard now rejects.
	if !strings.HasPrefix(fixture.input.Candidate.Document.Profile, "sdui/") || fixture.input.Guard() == nil {
		t.Fatal("invalid reload replaced last candidate or left stale guard valid")
	}
	write("page.sdui", `sdui 0.3; ref: art "unopened"; Main=[s=svg(art.X.@resource,description="Image",fallback="reject")];`)
	restored := next()
	if restored.err != nil || restored.generation != 5 || restored.input.ContextRevision != fixture.input.ContextRevision {
		t.Fatal(restored)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("watcher did not stop")
	}
}

// This harness uses Fyne's test driver only in a test binary. Production
// preflight neither imports it nor claims this as native desktop geometry proof.
func TestDocumentNativeAdoptionRejectsStaleContextAndKeepsLastBundle(t *testing.T) {
	a := fynetest.NewApp()
	defer a.Quit()
	w := a.NewWindow("prototype test")
	defer w.Close()
	dir := t.TempDir()
	source := filepath.Join(dir, "page.sdui")
	if err := os.WriteFile(source, []byte(`sdui 0.3; Main=[t=tabs("Tasks")[one=page("One")[c=select("Choice",value="one")];two=page("Two")[flag=checkbox("Flag")]]];`), 0600); err != nil {
		t.Fatal(err)
	}
	contextPath := filepath.Join(dir, "context.json")
	if err := os.WriteFile(contextPath, []byte(`{"schema":"sdui-prototype-context/1","choices":{"Main/t/one/c":[{"id":"one","label":"One","enabled":true}]}}`), 0600); err != nil {
		t.Fatal(err)
	}
	input, err := prototype.LoadDocument(source, "Main", "", contextPath)
	if err != nil {
		t.Fatal(err)
	}
	host := fynehost.NewDocumentHost(w.Canvas(), input.Size)
	defer host.Close()
	var generation atomic.Uint64
	generation.Store(1)
	request := func(sequence uint64) fynehost.DocumentRequest {
		return fynehost.DocumentRequest{Document: input.Candidate.Document, Entry: "Main", SessionID: "test-prototype", Sequence: sequence, SourceRevision: input.Candidate.Hash, Mode: preparation.Prototype, Choices: input.Choices, Guard: func() error {
			if generation.Load() != sequence {
				return fmt.Errorf("stale: context generation")
			}
			return nil
		}}
	}
	if err = host.Adopt(request(1)); err != nil {
		t.Fatal(err)
	}
	original := host.Current()
	if err = host.Mutate(func(s *ui.Session) error { h, _ := s.Pane("Main/t"); return s.SelectPage(h, "two") }); err != nil {
		t.Fatal(err)
	}
	generation.Store(2)
	candidate, err := host.Prepare(request(2))
	if err != nil {
		t.Fatal(err)
	}
	generation.Store(3)
	if err = host.Commit(candidate); err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatal(err)
	}
	if host.Current() != original || original.Session.Closed() || original.Session.Snapshot().Tabs["Main/t"].Selected != "two" {
		t.Fatal("failed stale commit changed published bundle/state")
	}
	invalid := request(3)
	invalid.Choices = nil
	if err = host.Adopt(invalid); err == nil {
		t.Fatal("native host accepted missing hidden choice provider")
	}
	if host.Current() != original || original.Session.Closed() {
		t.Fatal("invalid candidate retired last valid bundle")
	}
	if err = host.Adopt(request(3)); err != nil {
		t.Fatal(err)
	}
	if host.Current() == original || !original.Session.Closed() || host.Current().Session.Snapshot().Tabs["Main/t"].Selected != "two" {
		t.Fatal("valid reload did not preserve compatible tab state")
	}
}
