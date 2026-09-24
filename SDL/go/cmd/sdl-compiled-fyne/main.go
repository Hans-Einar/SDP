//go:build desktop

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	model "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/generatedmodel"
	"os"
)

func main() {
	prototype, e := application.New(model.Program(), model.Document(), model.Root())
	must(e)
	desktop := app.NewWithID("no.sdp.compiled.apt-simulation")
	window := desktop.NewWindow("Generert SDL + SDUI · simulert EditAptCell")
	status := widget.NewLabel("Genererte Go-modeller · simulert domene")
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
