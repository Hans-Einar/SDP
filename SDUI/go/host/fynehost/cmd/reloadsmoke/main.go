//go:build desktop

// Native file-watch/reload acceptance, including draft preservation and no replay.
package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
	uiruntime "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	dir, err := os.MkdirTemp("", "sdui-native-reload-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "page.sdui")
	must(os.WriteFile(path, []byte(`sdui 0.2; page=[header="# Before reload",edit=input("Value",value="12")] {scale=1};`), 0600))
	first := reload.Read(path, "page", 1, nil)
	must(first.Err)
	s, err := uiruntime.New("native-test", first.Root)
	must(err)
	a := app.NewWithID("no.sdp.sdui.reload-acceptance")
	w := a.NewWindow("SDUI native hot reload")
	v, err := fynehost.NewRuntime(s, w.Canvas())
	must(err)
	calls := 0
	initial, _ := s.Widget("page/edit")
	must(s.Bind(initial.Handle, func(e uiruntime.Event) ([]uiruntime.Update, error) {
		calls++
		current, _ := s.Widget(e.Handle.Path)
		return []uiruntime.Update{{Handle: e.Handle, Property: uiruntime.AcceptedValue, Value: e.Value, ExpectedValueRevision: current.ValueRevision, AcceptDraft: true}}, nil
	}))
	w.SetPadded(false)
	w.SetContent(v.Container)
	w.Resize(fyne.NewSize(640, 400))
	w.Show()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	phase := 0
	draft := ""
	old := v.View.Controls[initial.InstancePath].(*fynehost.Input)
	go func() {
		for c := range reload.Watch(ctx, path, "page", 20*time.Millisecond, nil) {
			fyne.Do(func() {
				if c.Hash == first.Hash {
					return
				}
				switch phase {
				case 0:
					if v.Adopt(c) == nil {
						panic("invalid source accepted")
					}
					state, _ := s.Widget("page/edit")
					if state.Draft != draft || s.Revision != 1 || calls != 0 {
						panic("invalid save changed state")
					}
					phase = 1
					replacement := path + ".tmp"
					must(os.WriteFile(replacement, []byte(`sdui 0.2; page=[header="# Reloaded without losing draft",<edit=input("New label",value="new default")>;"Runtime callback accepted exactly once"] {scale=1};`), 0600))
					must(os.Rename(replacement, path))
				case 1:
					must(v.Adopt(c))
					state, _ := s.Widget("page/edit")
					current := v.View.Controls[state.InstancePath].(*fynehost.Input)
					if current.Text != draft || w.Canvas().Focused() != current || calls != 0 {
						panic("reload lost draft/focus or replayed")
					}
					old.OnSubmitted(draft)
					if calls != 0 {
						panic("retired callback executed")
					}
					current.OnSubmitted(draft)
					state, _ = s.Widget("page/edit")
					if calls != 1 || state.Dirty {
						panic("commit failed")
					}
					phase = 2
					cancel()
					json.NewEncoder(os.Stdout).Encode(map[string]any{"driver": "native OpenGL", "model_revision": s.Revision, "invalid_source": "last-good retained", "draft": state.Draft, "focus": "preserved", "callback_count": calls, "retired_callback": "rejected"})
					go func() { time.Sleep(5 * time.Second); fyne.Do(func() { v.Close(); w.Close(); a.Quit() }) }()
				}
			})
		}
	}()
	go func() {
		time.Sleep(400 * time.Millisecond)
		fyne.Do(func() {
			w.Canvas().Focus(old)
			test.Type(old, "draft")
			state, _ := s.Widget("page/edit")
			draft = state.Draft
			must(os.WriteFile(path, []byte(`sdui 0.2; page=[`), 0600))
		})
	}()
	a.Run()
}
