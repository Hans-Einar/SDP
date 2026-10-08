//go:build desktop

// Native WCI3-M1 acceptance fixture. Product APIs never read this control channel.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/values"
	"os"
	"sync"
)

func main() {
	nonmodal := flag.Bool("nonmodal", false, "Use standard nonmodal form windows; source identities and protocol remain unchanged")
	requiredEmpty := flag.Bool("required-empty", false, "Start only main Mode required and empty; options and protocol remain unchanged")
	flag.Parse()
	var outputMu sync.Mutex
	log := func(event string, value any) {
		outputMu.Lock()
		defer outputMu.Unlock()
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"event": event, "data": value})
	}
	f, err := values.New()
	if err != nil {
		panic(err)
	}
	f.Log = log
	a := app.NewWithID("no.sdp.wci3.values")
	w := a.NewWindow(values.WindowTitle)
	w.SetPadded(false)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1100, H: 850})
	source := values.FixtureSource(*nonmodal)
	if *requiredEmpty {
		source = values.RequiredEmptySource(*nonmodal)
	}
	controls := &values.Controller{Fixture: f, Host: h, Source: source, Sequence: 1, Log: log}
	controls.Resize = func(size layout.Size) error {
		w.Resize(fyne.NewSize(float32(size.W), float32(size.H)))
		return h.Resize(size)
	}
	controls.InstallHooks()

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
	w.Resize(fyne.NewSize(1100, 850))
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
			log("closed", map[string]any{"sessionClosed": closed, "pending": controls.Pending(), "domain": f.Domain(), "actionCalls": calls(f)})
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
				result := map[string]any{"commandId": commandID, "command": line, "status": "ok", "candidateSequence": controls.Sequence, "actionCalls": calls(f), "bundleChanged": before != h.Current()}
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
		log("ready", map[string]any{"title": values.WindowTitle, "size": w.Canvas().Size(), "nonmodal": *nonmodal})
		h.OnChange(h.Current())
	})
	a.Run()
}

func calls(f *values.Fixture) map[string]int { a, _ := f.Counts(); return a }
