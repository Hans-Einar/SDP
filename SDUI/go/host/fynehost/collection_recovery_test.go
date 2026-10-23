package fynehost

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestLongProviderErrorFitsDefaultOverflowAndRetriesOnce(t *testing.T) {
	for _, depth := range []int{0, 14} {
		t.Run(fmt.Sprint(depth), func(t *testing.T) {
			h := documentHost(t)
			if err := h.Resize(layout.Size{W: 300, H: 500}); err != nil {
				t.Fatal(err)
			}
			posts := make(chan func(), 8)
			h.Post = func(f func()) { posts <- f }
			calls := make(chan ui.LoadRequest, 8)
			replies := make(chan error, 8)
			message := strings.Repeat("bounded provider diagnostic ", 40)
			p := ui.CollectionProvider{ID: "recovery", Epoch: 1, RootLoaded: depth > 0, Load: func(ctx context.Context, r ui.LoadRequest) (ui.CollectionData, error) {
				calls <- r
				select {
				case err := <-replies:
					return ui.CollectionData{}, err
				case <-ctx.Done():
					return ui.CollectionData{}, ctx.Err()
				}
			}}
			parent := ui.ItemID("")
			for i := 0; i < depth; i++ {
				id := ui.ItemID(fmt.Sprint(i))
				p.Initial.Items = append(p.Initial.Items, ui.CollectionItem{ID: id, Parent: parent, Kind: ui.Group, Label: "x", HasChildren: true, ChildrenLoaded: i < depth-1})
				parent = id
			}
			req := documentRequest(t, 1, p)
			d, err := parser.Parse(strings.Replace(collectionSource, "overflow-x=scroll,", "", 1))
			if err != nil {
				t.Fatal(err)
			}
			req.Document = d
			if err = h.Adopt(req); err != nil {
				t.Fatal(err)
			}
			b := h.Current()
			c := b.Controls()["page/tree"].(*CollectionControl)
			for i := 0; i < depth; i++ {
				target, err := b.Session.Target(c.state.Handle, ui.ItemID(fmt.Sprint(i)))
				if err != nil {
					t.Fatal(err)
				}
				if err = h.Mutate(func(s *ui.Session) error { return s.ExpandItem(target) }); err != nil {
					t.Fatal(err)
				}
			}
			var first ui.LoadRequest
			select {
			case first = <-calls:
			case <-time.After(3 * time.Second):
				t.Fatal("load did not start")
			}
			replies <- errors.New(message)
			select {
			case f := <-posts:
				f()
			case <-time.After(3 * time.Second):
				t.Fatal("completion not delivered")
			}
			state, _ := b.Session.Collection(c.state.Handle)
			if state.Request != nil || state.Status[parent].Phase != ui.LoadError {
				t.Fatal("provider error stranded loading", state.Request, state.Status)
			}
			if state.Status[parent].Error != message {
				t.Fatal("runtime diagnostic was truncated")
			}
			found := false
			for _, row := range c.rows {
				if row.Status && row.Recovery {
					label, indent := fittedRow(row, state, c.font, c.viewport.Rect.W)
					if !strings.Contains(label, "Retry") || !strings.Contains(label, "R") {
						t.Fatal("recovery affordance lost", label)
					}
					if float64(nativeText(label, c.font).Width)+indent+8 > c.viewport.Rect.W {
						t.Fatal("status exceeds admitted width", label, indent)
					}
					if depth > 0 && indent >= float64(row.Depth)*18 {
						t.Fatal("deep status indentation was not capped")
					}
					found = true
				}
			}
			if !found {
				t.Fatal("no fitted recovery row")
			}
			if depth > 0 {
				target, _ := b.Session.Target(c.state.Handle, parent)
				if err = h.Mutate(func(s *ui.Session) error { return s.FocusItem(target) }); err != nil {
					t.Fatal(err)
				}
			}
			c.TypedRune('r')
			c.TypedRune('r')
			var next ui.LoadRequest
			select {
			case next = <-calls:
			case <-time.After(3 * time.Second):
				t.Fatal("Retry did not start a fresh load")
			}
			if next.RequestID == first.RequestID {
				t.Fatal("Retry reused request")
			}
			state, _ = b.Session.Collection(c.state.Handle)
			if state.Request == nil || state.Request.RequestID != next.RequestID {
				t.Fatal("duplicate Retry replaced request")
			}
			replies <- nil
			select {
			case f := <-posts:
				f()
			case <-time.After(3 * time.Second):
				t.Fatal("retry completion not delivered")
			}
			select {
			case extra := <-calls:
				t.Fatal("duplicate Retry started load", extra)
			default:
			}
			target, _ := b.Session.Target(c.state.Handle, "")
			bad := ui.CollectionData{Items: []ui.CollectionItem{{ID: "wide", Kind: ui.Row, Label: message, ChildrenLoaded: true}}}
			if err = h.Mutate(func(s *ui.Session) error { return s.ReplaceCollection(target, bad) }); err == nil {
				t.Fatal("wide data row bypassed source overflow-x=error")
			}
			// A data row still uses its full label and depth; no source overflow policy is relaxed.
			dataRow := collectionRow{Item: ui.CollectionItem{ID: "wide", Kind: ui.Row}, Text: message, Depth: depth}
			label, indent := fittedRow(dataRow, state, c.font, c.viewport.Rect.W)
			if label != rowText(dataRow, state) || indent != float64(depth)*18 {
				t.Fatal("data row was fitted")
			}
		})
	}
}

func TestRecoveryLabelFitsMinimumViewport(t *testing.T) {
	documentHost(t)
	for _, font := range []float64{14, 20, 28} {
		for _, depth := range []int{0, 100} {
			for _, phase := range []ui.LoadPhase{ui.LoadError, ui.Canceled} {
				width := collectionMinimum(font, false, true).W - collectionGutter
				state := ui.CollectionState{Status: map[ui.ItemID]ui.LoadStatus{"": {Phase: phase, Error: strings.Repeat("diagnostic", 100)}}}
				row := collectionRow{Text: "Retry — R: " + state.Status[""].Error, Depth: depth, Status: true, Recovery: true}
				text, indent := fittedRow(row, state, font, width)
				prefix := "Retry/R"
				if phase == ui.Canceled {
					prefix = "Load/R"
				}
				if !strings.HasPrefix(text, prefix) || indent+float64(nativeText(text, font).Width)+8 > width {
					t.Fatal("minimum viewport lost recovery label", font, depth, phase, text, indent, width)
				}
			}
		}
	}
}
