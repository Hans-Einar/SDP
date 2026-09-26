package application

import (
	"fmt"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Exercise drives the same public input/event path in source and generated mode.
// The reports are simulation evidence, not production domain behavior.
func Exercise(app *App, values []string) (map[string]any, error) {
	reports := []any{}
	for i, value := range values {
		edit, ok := app.UI.Widget("page/value")
		if !ok {
			return nil, fmt.Errorf("missing prototype control")
		}
		button, ok := app.UI.Widget("page/apply")
		if !ok {
			return nil, fmt.Errorf("missing prototype control")
		}
		if e := app.UI.Draft(edit.Handle, value); e != nil {
			return nil, e
		}
		e := app.UI.Dispatch(ui.Event{Handle: button.Handle, ModelRevision: app.UI.Revision, Sequence: uint64(i + 1), Kind: ui.Activate})
		report := map[string]any{"sequence": i + 1, "draft": value, "accepted": e == nil}
		if e != nil {
			report["error"] = e.Error()
		}
		reports = append(reports, report)
	}
	value, revision, trace := app.Domain.Snapshot()
	return map[string]any{"simulation": true, "events": reports, "domain_value": value, "domain_revision": revision, "domain_trace": trace, "links": app.Bindings.Links, "widgets": app.UI.Widgets()}, nil
}
