package main

import (
	"encoding/json"
	"flag"
	"fmt"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	"os"
	"strings"
)

func main() {
	actions := flag.String("actions", "examples/edit-apt-cell.sdl", "Explicit SDL source")
	source := flag.String("ui", "examples/edit-apt-cell.sdui", "SDUI source")
	values := flag.String("values", "430,invalid,440", "Ordered simulated drafts")
	flag.Parse()
	a, e := os.ReadFile(*actions)
	must(e)
	s, e := os.ReadFile(*source)
	must(e)
	app, e := application.Load(string(a), string(s))
	must(e)
	defer app.Close()
	reports := []any{}
	for i, value := range strings.Split(*values, ",") {
		edit, _ := app.UI.Widget("page/value")
		button, _ := app.UI.Widget("page/apply")
		must(app.UI.Draft(edit.Handle, value))
		e = app.UI.Dispatch(ui.Event{Handle: button.Handle, ModelRevision: app.UI.Revision, Sequence: uint64(i + 1), Kind: ui.Activate})
		report := map[string]any{"sequence": i + 1, "draft": value, "accepted": e == nil}
		if e != nil {
			report["error"] = e.Error()
		}
		reports = append(reports, report)
	}
	value, revision, trace := app.Domain.Snapshot()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	must(enc.Encode(map[string]any{"simulation": true, "events": reports, "domain_value": value, "domain_revision": revision, "domain_trace": trace, "links": app.Bindings.Links, "widgets": app.UI.Widgets()}))
}
func must(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
