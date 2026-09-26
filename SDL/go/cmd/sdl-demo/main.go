//go:build desktop

package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	uireload "github.com/Hans-Einar/SDP/SDUI/go/reload"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	sdlreload "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/reload"
	"os"
	"time"
)

func main() {
	actions := flag.String("actions", "examples/edit-apt-cell.sdl", "Explicit SDL source")
	source := flag.String("ui", "examples/edit-apt-cell.sdui", "SDUI source")
	watch := flag.Bool("watch", true, "Reload SDL and SDUI source changes")
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
	sourceErrors := map[string]error{}
	view.OnStatus = func(e error) {
		if e == nil {
			for _, language := range []string{"SDUI", "SDL"} {
				if sourceErrors[language] != nil {
					e = fmt.Errorf("%s: %w", language, sourceErrors[language])
					break
				}
			}
		}
		if e != nil {
			status.SetText(e.Error())
		} else {
			value, revision, _ := prototype.Domain.Snapshot()
			status.SetText(fmt.Sprintf("Simulert domene · verdi %s · revisjon %d · SDUI r%d · SDL r%d", value, revision, prototype.UI.Revision, prototype.SDL.Revision()))
		}
	}
	window.SetContent(container.NewBorder(nil, status, nil, nil, view.Container))
	window.Resize(fyne.NewSize(1000, 650))
	ctx, cancel := context.WithCancel(context.Background())
	window.SetOnClosed(func() { cancel(); view.Close(); prototype.SDL.Close() })
	if *watch {
		watchModels(ctx, *actions, *source, string(a), string(s), prototype, view, func(language string, err error) { sourceErrors[language] = err; view.OnStatus(nil) })
	}
	window.ShowAndRun()
}
func must(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}

func watchModels(ctx context.Context, actions, source, initialActions, initialUI string, prototype *application.App, view *fynehost.RuntimeView, report func(string, error)) {
	go func() {
		for c := range uireload.Watch(ctx, source, "page", 250*time.Millisecond, nil) {
			fyne.Do(func() {
				if c.Sequence == 1 && c.Hash == fmt.Sprintf("%x", sha256.Sum256([]byte(initialUI))) {
					return
				}
				if c.Err == nil {
					c.Err = prototype.ValidateUI(c.Document, c.Root)
				}
				err := view.Adopt(c)
				if err == nil {
					err = prototype.RebindUI(c.Document)
				}
				report("SDUI", err)
			})
		}
	}()
	go func() {
		var latest uint64
		for c := range sdlreload.Watch(ctx, actions, 250*time.Millisecond) {
			fyne.Do(func() {
				if c.Sequence <= latest {
					return
				}
				latest = c.Sequence
				if c.Sequence == 1 && c.Hash == fmt.Sprintf("%x", sha256.Sum256([]byte(initialActions))) {
					return
				}
				err := c.Err
				if err == nil {
					err = prototype.ReloadSDL(c.Program)
				}
				if err == nil {
					err = view.RefreshModel()
				}
				report("SDL", err)
			})
		}
	}()
}
