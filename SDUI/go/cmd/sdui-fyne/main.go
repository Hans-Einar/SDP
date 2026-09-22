//go:build desktop

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
	uiruntime "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func main() {
	entry := flag.String("entry", "page", "Root frame")
	watch := flag.Bool("watch", true, "Reload changed SDUI source")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage: sdui-fyne [-entry page] [-watch=true] source.sdui")
		os.Exit(2)
	}
	initial := reload.Read(flag.Arg(0), *entry, 1, nil)
	if initial.Err != nil {
		fmt.Fprintln(os.Stderr, initial.Err)
		os.Exit(2)
	}
	session, err := uiruntime.New(fmt.Sprintf("ui-%d", time.Now().UnixNano()), initial.Root)
	if err != nil {
		panic(err)
	}
	a := app.NewWithID("no.sdp.sdui.prototype")
	w := a.NewWindow("SDUI · native prototype")
	status := widget.NewLabel("Lokal prototype · ingen SDL-runtime")
	view, err := fynehost.NewRuntime(session, w.Canvas())
	if err != nil {
		panic(err)
	}
	view.OnStatus = func(err error) {
		if err != nil {
			status.SetText(err.Error())
		} else {
			status.SetText("OK · modellrevisjon " + fmt.Sprint(session.Revision))
		}
	}
	registerPrototype(session, *entry)
	w.SetContent(container.NewBorder(nil, status, nil, nil, view.Container))
	w.Resize(fyne.NewSize(1280, 800))
	ctx, cancel := context.WithCancel(context.Background())
	w.SetOnClosed(func() { cancel(); view.Close() })
	if *watch {
		go func() {
			for candidate := range reload.Watch(ctx, flag.Arg(0), *entry, 250*time.Millisecond, nil) {
				fyne.Do(func() {
					if candidate.Hash == initial.Hash && candidate.Sequence == 1 {
						return
					}
					if err := view.Adopt(candidate); err == nil {
						registerPrototype(session, *entry)
					}
				})
			}
		}()
	}
	w.ShowAndRun()
}

// Explicit local prototype bindings. Symbolic SDL callbacks remain unbound.
func registerPrototype(s *uiruntime.Session, entry string) {
	for _, w := range s.Widgets() {
		if w.Binding.Module != "" {
			continue
		}
		switch {
		case w.Handle.Kind == "input":
			_ = s.Bind(w.Handle, func(e uiruntime.Event) ([]uiruntime.Update, error) {
				current, _ := s.Widget(e.Handle.Path)
				return []uiruntime.Update{{Handle: e.Handle, Property: uiruntime.AcceptedValue, Value: e.Value, ExpectedValueRevision: current.ValueRevision, AcceptDraft: true}}, nil
			})
		case w.Handle.Path == entry+"/ok":
			_ = s.Bind(w.Handle, func(e uiruntime.Event) ([]uiruntime.Update, error) {
				fmt.Println("local callback:", e.Handle.Path)
				return nil, nil
			})
		}
	}
}
