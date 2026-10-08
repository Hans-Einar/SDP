//go:build desktop

// Native WCI2-M1 acceptance fixture. Product APIs never read this control channel.
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
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/panes"
	"os"
	"strings"
	"sync"
)

func main() {
	vertical := flag.Bool("vertical", false, "Use a vertical split; fixture protocol and identities remain unchanged")
	flag.Parse()
	var outputMu sync.Mutex
	log := func(event string, value any) {
		outputMu.Lock()
		defer outputMu.Unlock()
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"event": event, "data": value})
	}
	f, err := panes.New()
	if err != nil {
		panic(err)
	}
	f.Log = log
	a := app.NewWithID("no.sdp.wci2.panes")
	w := a.NewWindow(panes.WindowTitle)
	w.SetPadded(false)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1000, H: 650})
	source := panes.Source
	if *vertical {
		source = strings.Replace(source, `axis="horizontal"`, `axis="vertical"`, 1)
	}
	controls := &panes.Controller{Fixture: f, Host: h, Source: source, Sequence: 1, Log: log}
	controls.Resize = func(size layout.Size) error {
		w.Resize(fyne.NewSize(float32(size.W), float32(size.H)))
		return h.Resize(size)
	}
	f.Conflict = func() error {
		current := h.Current()
		if current == nil {
			return fmt.Errorf("no live bundle")
		}
		receiver, ok := current.Session.Widget(panes.PreviewPath)
		if !ok {
			return fmt.Errorf("missing preview")
		}
		return current.Session.Draft(receiver.Handle, "New draft during Page action")
	}

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
		log("ready", map[string]any{"title": panes.WindowTitle, "size": w.Canvas().Size()})
		h.OnChange(h.Current())
	})
	a.Run()
}
