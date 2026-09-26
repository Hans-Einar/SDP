//go:build desktop

package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	"os"
	"time"
)

func must(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	a, e := os.ReadFile("examples/edit-apt-cell.sdl")
	must(e)
	u, e := os.ReadFile("examples/edit-apt-cell.sdui")
	must(e)
	prototype, e := application.Load(string(a), string(u))
	must(e)
	desktop := app.NewWithID("no.sdp.sdl.native-binding-probe")
	w := desktop.NewWindow("SDL native binding acceptance")
	view, e := fynehost.NewRuntime(prototype.UI, w.Canvas())
	must(e)
	var lastError error
	view.OnStatus = func(e error) { lastError = e }
	w.SetPadded(false)
	w.SetContent(view.Container)
	w.Resize(fyne.NewSize(1000, 600))
	w.Show()
	go func() {
		time.Sleep(400 * time.Millisecond)
		fyne.Do(func() {
			must(lastError)
			entry := view.View.Controls["page/value"].(*fynehost.Input)
			button := view.View.Controls["page/apply"].(*widget.Button)
			w.Canvas().Focus(entry)
			entry.SetText("")
			test.Type(entry, "430")
			test.Tap(button)
			must(lastError)
			value, revision, trace := prototype.Domain.Snapshot()
			if value != "430" || entry.Text != "430" || revision != 1 || len(trace) != 1 {
				panic("native SDL action failed")
			}
			entry.SetText("invalid")
			test.Tap(button)
			if lastError == nil {
				panic("invalid simulated edit accepted")
			}
			rejection := lastError.Error()
			entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
			value, revision, trace = prototype.Domain.Snapshot()
			if value != "430" || revision != 1 || entry.Text != "430" || len(trace) != 2 || trace[1].Accepted {
				panic("rejected action or revert changed domain truth")
			}
			json.NewEncoder(os.Stdout).Encode(map[string]any{"driver": "native OpenGL", "simulation": true, "accepted_value": value, "domain_revision": revision, "domain_calls": len(trace), "rejection": rejection, "source_links": len(prototype.Bindings.Links)})
			go func() {
				time.Sleep(5 * time.Second)
				fyne.Do(func() { view.Close(); prototype.SDL.Close(); w.Close(); desktop.Quit() })
			}()
		})
	}()
	desktop.Run()
}
