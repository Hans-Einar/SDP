//go:build desktop

// Native-driver acceptance probe; deliberately separate from the product CLI.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func main() {
	started := time.Now()
	a := app.NewWithID("no.sdp.sdui.acceptance")
	w := a.NewWindow("SDUI native acceptance")
	_, roots, err := parser.Compile(`sdui 0.2; page=[header="# Native SDUI",ok=button("OK"),edit=input("Value",value="12");"Resize preserves text and value"]*b {scale=1};`)
	must(err)
	v := fynehost.New(roots["page"])
	calls := 0
	value := ""
	v.Actions["page/ok"] = func(_, _ string) error { calls++; return nil }
	v.Actions["page/edit"] = func(_, text string) error { value = text; return nil }
	var lastError error
	v.OnStatus = func(err error) { lastError = err }
	w.SetPadded(false)
	w.SetContent(v.Container)
	w.Resize(fyne.NewSize(640, 400))
	w.Show()
	go func() {
		time.Sleep(400 * time.Millisecond)
		fyne.Do(func() {
			must(lastError)
			button := v.Controls["page/ok"].(*widget.Button)
			entry := v.Controls["page/edit"].(*widget.Entry)
			test.Tap(button)
			w.Canvas().Focus(button)
			w.Canvas().FocusNext()
			if w.Canvas().Focused() != entry {
				panic("Tab did not reach entry")
			}
			test.Type(entry, "3")
			if calls != 1 || value == "" {
				panic("native callback failed")
			}
			resize := time.Now()
			w.Resize(fyne.NewSize(800, 500))
			v.Container.Resize(fyne.NewSize(800, 500))
			duration := time.Since(resize)
			if entry.Text != value || calls != 1 {
				panic("resize state lost")
			}
			go func() {
				time.Sleep(350 * time.Millisecond)
				fyne.Do(func() {
					must(lastError)
					var mem runtime.MemStats
					runtime.ReadMemStats(&mem)
					status, err := os.ReadFile("/proc/self/status")
					must(err)
					must(os.WriteFile("/tmp/sdui-fyne-native-process.txt", status, 0600))
					json.NewEncoder(os.Stdout).Encode(map[string]any{"driver": "Fyne desktop/OpenGL", "startup_and_test_ms": time.Since(started).Milliseconds(), "resize_us": duration.Microseconds(), "go_heap_bytes": mem.Alloc, "callback_count": calls, "entry": value, "tab_focus": "entry"})
					go func() { time.Sleep(5 * time.Second); fyne.Do(func() { v.Close(); w.Close(); a.Quit() }) }()
				})
			}()
		})
	}()
	a.Run()
}
func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
}
