package runtime

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"strings"
	"testing"
	"time"
)

func TestReloadWaitsAndPreservesState(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	state := 0
	r := registry(func(_ context.Context, v Record) (Record, error) { close(entered); <-release; state++; return v, nil })
	r["GoOther"] = Binding{Input: parser.RecordType{"Value": parser.TextType}, Output: parser.RecordType{"Value": parser.TextType}, Call: func(_ context.Context, v Record) (Record, error) { state++; v["Value"] = Text("new"); return v, nil }}
	e, err := New(program(t), r)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := e.Execute(context.Background(), Request{"Echo", 1, 1, Record{"Value": Text("old")}})
		done <- err
	}()
	<-entered
	source, err := os.ReadFile("../examples/echo.sdl")
	if err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(string(source), "GoEcho", "GoOther", 1)
	reloaded := make(chan error, 1)
	go func() { reloaded <- e.ReloadSource(changed) }()
	select {
	case err := <-reloaded:
		t.Fatal("reload crossed active invocation", err)
	case <-time.After(30 * time.Millisecond):
	}
	close(release)
	if err = <-done; err != nil {
		t.Fatal(err)
	}
	if err = <-reloaded; err != nil {
		t.Fatal(err)
	}
	if e.Revision() != 2 || state != 1 {
		t.Fatal("state reset or replay")
	}
	if _, err = e.Execute(context.Background(), Request{"Echo", 1, 2, Record{"Value": Text("late")}}); err == nil {
		t.Fatal("old model command accepted")
	}
	result, err := e.Execute(context.Background(), Request{"Echo", 2, 2, Record{"Value": Text("input")}})
	if err != nil || result.Output["Value"].Text != "new" || state != 2 {
		t.Fatal(result, err, state)
	}
	if err = e.ReloadSource("bad SDL"); err == nil || e.Revision() != 2 {
		t.Fatal("invalid reload changed revision")
	}
	if err = e.ReloadSource(strings.Replace(changed, "GoOther", "MissingGo", 1)); err == nil || e.Revision() != 2 {
		t.Fatal("unbound reload changed revision")
	}
}
