//go:build desktop

package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	"os"
)

func main() {
	actions := flag.String("actions", "examples/edit-apt-cell.sdl", "Explicit SDL source")
	source := flag.String("ui", "examples/edit-apt-cell.sdui", "SDUI source")
	flag.Parse()
	a, e := os.ReadFile(*actions)
	must(e)
	s, e := os.ReadFile(*source)
	must(e)
	prototype, e := application.Load(string(a), string(s))
	must(e)
	desktop := app.NewWithID("no.sdp.sdl.apt-simulation")
	window := desktop.NewWindow("SDL + SDUI · simulert EditAptCell")
	status := widget.NewLabel("Simulert domene · ingen Ponsse-runtime")
	view, e := fynehost.NewRuntime(prototype.UI, window.Canvas())
	must(e)
	view.OnStatus = func(e error) {
		if e != nil {
			status.SetText(e.Error())
		} else {
			value, revision, _ := prototype.Domain.Snapshot()
			status.SetText(fmt.Sprintf("Simulert domene · verdi %s · revisjon %d", value, revision))
		}
	}
	window.SetContent(container.NewBorder(nil, status, nil, nil, view.Container))
	window.Resize(fyne.NewSize(1000, 650))
	window.SetOnClosed(func() { view.Close(); prototype.SDL.Close() })
	window.ShowAndRun()
}
func must(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
