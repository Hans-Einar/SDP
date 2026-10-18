//go:build desktop

// Native WCI2-M2 acceptance fixture. Product APIs never read this control channel.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/commands"
	"os"
	"sync"
)

func main() {
	nonmodal := flag.Bool("nonmodal", false, "Use a standard nonmodal Settings window; source identities and protocol remain unchanged")
	flag.Parse()
	var outputMu sync.Mutex
	log := func(event string, value any) {
		outputMu.Lock()
		defer outputMu.Unlock()
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"event": event, "data": value})
	}
	f, err := commands.New()
	if err != nil {
		panic(err)
	}
	f.Log = log
	a := app.NewWithID("no.sdp.wci2.commands")
	w := a.NewWindow(commands.WindowTitle)
	w.SetPadded(false)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1100, H: 750})
	source := commands.FixtureSource(*nonmodal)
	controls := &commands.Controller{Fixture: f, Host: h, Source: source, Sequence: 1, Log: log}
	controls.Resize = func(size layout.Size) error {
		w.Resize(fyne.NewSize(float32(size.W), float32(size.H)))
		return h.Resize(size)
	}
	f.Conflict = func(path string) error {
		current := h.Current()
		if current == nil {
			return fmt.Errorf("no live bundle")
		}
		receiver, ok := current.Session.Widget(path)
		if !ok {
			return fmt.Errorf("missing conflict input %s", path)
		}
		return current.Session.Draft(receiver.Handle, "New draft during SDL action")
	}

	h.OnDialogResult = func(result ui.DialogResult) { log("dialog-result", result) }
	controls.ParentHide = func() error {
		if err := h.NativeParentHidden(); err != nil {
			return err
		}
		w.Hide()
		return nil
	}
	controls.ParentShow = func() error { w.Show(); return nil }
	h.OnStatus = func(e error) { log("error", e.Error()) }
	h.OnChange = func(*fynehost.Bundle) { log("state", controls.Inspect()) }
	request, e := f.Request(source, controls.Sequence)
	if e != nil {
		panic(e)
	}
	if e = h.Adopt(request); e != nil {
		panic(e)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(1100, 750))
	var closeOnce sync.Once
	closeFixture := func() error {
		if err := h.NativeParentClosed(); err != nil {
			return err
		}
		closeOnce.Do(func() {
			h.Close()
			f.Close()
			closed := true
			if b := h.Current(); b != nil {
				closed = b.Session.Closed()
			}
			log("closed", map[string]any{"sessionClosed": closed, "pending": f.Pending(), "actionCalls": f.Calls()})
			w.SetCloseIntercept(nil)
			w.Close()
			a.Quit()
		})
		return nil
	}
	controls.Close = closeFixture
	w.SetCloseIntercept(func() {
		if err := closeFixture(); err != nil {
			log("error", err.Error())
		}
	})
	var commandID uint64
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, 4096), 65536)
		for scanner.Scan() {
			line := scanner.Text()
			fyne.Do(func() {
				commandID++
				before := h.Current()
				e := controls.Run(line)
				result := map[string]any{"commandId": commandID, "command": line, "status": "ok", "candidateSequence": controls.Sequence, "actionCalls": f.Calls(), "bundleChanged": before != h.Current()}
				if b := h.Current(); b != nil {
					result["source"] = b.SourceRevision
					result["modelRevision"] = b.Session.Revision
				}
				if e != nil {
					result["status"] = "error"
					result["error"] = e.Error()
					log("error", e.Error())
					if h.Current() != nil && !h.Current().Session.Closed() {
						log("state", controls.Inspect())
					}
				} else {
					log("command", line)
				}
				log("command-result", result)
			})
		}
	}()
	w.Show()
	fyne.Do(func() {
		log("ready", map[string]any{"title": commands.WindowTitle, "size": w.Canvas().Size(), "nonmodal": *nonmodal})
		h.OnChange(h.Current())
	})
	a.Run()
}
