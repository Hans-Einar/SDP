//go:build desktop

// Native WCI1 acceptance fixture. Product APIs never read this control channel.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/collections"
	"image/png"
	"os"
	"sync"
)

func main() {
	empty := flag.Bool("empty", false, "Start tree with unloaded root")
	nested := flag.Bool("nested", false, "Use nested body and collection viewports")
	flag.Parse()
	var outputMu sync.Mutex
	log := func(event string, value any) {
		outputMu.Lock()
		defer outputMu.Unlock()
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"event": event, "data": value})
	}
	f, err := collections.New()
	if err != nil {
		panic(err)
	}
	f.Log = log
	a := app.NewWithID("no.sdp.wci1.collections")
	w := a.NewWindow(collections.WindowTitle)
	w.SetPadded(false)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1000, H: 650})
	source := collections.Source
	if *nested {
		source = collections.NestedSource
	}
	controls := &collections.Controller{Fixture: f, Host: h, Source: source, Empty: *empty, Sequence: 1, Log: log}
	controls.Resize = func(size layout.Size) error {
		w.Resize(fyne.NewSize(float32(size.W), float32(size.H)))
		return h.Resize(size)
	}
	controls.Capture = func(path string) error {
		file, e := os.Create(path)
		if e != nil {
			return e
		}
		defer file.Close()
		return png.Encode(file, w.Canvas().Capture())
	}

	h.OnStatus = func(e error) { log("error", e.Error()) }
	h.OnChange = func(*fynehost.Bundle) { log("state", controls.Inspect()) }
	request, e := f.Request(source, controls.Sequence, *empty)
	if e != nil {
		panic(e)
	}
	if e = h.Adopt(request); e != nil {
		panic(e)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(1000, 650))
	var closeOnce sync.Once
	closeFixture := func() {
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
	}
	controls.Close = closeFixture
	w.SetCloseIntercept(closeFixture)
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
		log("ready", map[string]any{"title": collections.WindowTitle, "size": w.Canvas().Size(), "nested": *nested, "empty": *empty})
		h.OnChange(h.Current())
	})
	a.Run()
}
