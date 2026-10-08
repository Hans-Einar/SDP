// Temporary Architect probe. No SDUI/SDL product imports or persistence.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var logSequence uint64

func record(event string, data map[string]any) {
	logSequence++
	if data == nil {
		data = map[string]any{}
	}
	data["sequence"], data["event"] = logSequence, event
	data["time"] = time.Now().UTC().Format(time.RFC3339Nano)
	if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
		panic(err)
	}
}

type entry struct {
	widget.Entry
	id     string
	escape func()
}

func newEntry(id, text string) *entry {
	e := &entry{id: id}
	e.ExtendBaseWidget(e)
	e.SetText(text)
	e.OnChanged = func(text string) {
		record("draft", map[string]any{"entry": id, "bytes": len(text), "sha256": fmt.Sprintf("%x", sha256.Sum256([]byte(text)))})
	}
	return e
}

func (e *entry) FocusGained() {
	e.Entry.FocusGained()
	record("focus-gained", map[string]any{"entry": e.id})
}
func (e *entry) FocusLost() { e.Entry.FocusLost(); record("focus-lost", map[string]any{"entry": e.id}) }
func (e *entry) TypedKey(k *fyne.KeyEvent) {
	if k.Name == fyne.KeyEscape && e.escape != nil {
		e.escape()
		return
	}
	e.Entry.TypedKey(k)
}

type owner struct {
	app                     fyne.App
	parent                  fyne.Window
	input                   *entry
	child                   *surface
	opener                  fyne.Focusable
	live                    bool
	generation, nextOpening uint64
}
type surface struct {
	owner                     *owner
	window                    fyne.Window
	input                     *entry
	parentGeneration, opening uint64
	closed                    bool
	accepted                  string
	commits                   int
}

func (o *owner) open() {
	if !o.live {
		return
	}
	if o.child != nil {
		o.child.window.RequestFocus()
		record("existing-opening-focused", nil)
		return
	}
	o.nextOpening++
	c := &surface{owner: o, parentGeneration: o.generation, opening: o.nextOpening, accepted: "child accepted"}
	c.window = o.app.NewWindow(fmt.Sprintf("WCI2 child %d", c.opening))
	c.input = newEntry(fmt.Sprintf("child-%d", c.opening), c.accepted)
	c.input.OnSubmitted = func(text string) {
		if !c.current() {
			record("stale-commit-rejected", nil)
			return
		}
		c.accepted, c.commits = text, c.commits+1
		record("local-field-commit", map[string]any{"opening": c.opening, "commits": c.commits})
	}
	c.input.escape = func() {
		if !c.current() {
			return
		}
		if c.input.Text != c.accepted {
			c.input.SetText(c.accepted)
			record("draft-reverted", nil)
			return
		}
		c.finish("cancel", "escape", true, true)
	}
	c.window.SetContent(container.NewVBox(widget.NewLabel("Synthetic draft. Enter commits locally; Escape reverts then cancels."), c.input,
		widget.NewButton("Log snapshot", func() { o.snapshot("child-snapshot") }),
		widget.NewButton("Cancel child", func() { c.finish("cancel", "button", true, true) })))
	c.window.SetCloseIntercept(func() { c.finish("close", "window-chrome", true, false) })
	c.window.SetOnClosed(func() { c.finish("close", "native-close-fallback", false, false) })
	c.window.Resize(fyne.NewSize(620, 210))
	o.child = c // publish logical identity before native callbacks can run
	record("surface-published", map[string]any{"opening": c.opening, "parentGeneration": o.generation})
	c.window.Show()
	c.window.Canvas().Focus(c.input)
}

func (c *surface) current() bool {
	return !c.closed && c.owner.live && c.owner.child == c && c.owner.generation == c.parentGeneration
}

func (c *surface) finish(kind, reason string, closeNative, restore bool) {
	if c.closed {
		return
	}
	c.closed = true // revoke before Close/OnClosed reentrancy; emit exactly one result
	if c.owner.child == c {
		c.owner.child = nil
	}
	record("surface-result", map[string]any{"opening": c.opening, "parentGeneration": c.parentGeneration,
		"kind": kind, "reason": reason, "dirtyDiscarded": c.input.Text != c.accepted, "localCommits": c.commits})
	if closeNative {
		c.window.Close()
	}
	// Only an explicit gesture within the child asks for focus. Parent-driven closure never raises it.
	if restore && c.owner.live && c.owner.generation == c.parentGeneration {
		c.owner.parent.Canvas().Focus(c.owner.opener)
		c.owner.parent.RequestFocus()
		record("opener-focus-requested", nil) // request, not proof the window manager granted it
	}
}

func (o *owner) disposeChild(reason string) {
	if o.child != nil {
		o.child.finish("close", reason, true, false)
	}
}
func (o *owner) closeParent() {
	if !o.live {
		return
	}
	o.live = false
	o.disposeChild("parent-closed")
	record("parent-close", nil)
	o.parent.Close()
}
func (o *owner) snapshot(event string) {
	d := map[string]any{"parentGeneration": o.generation, "childOpen": o.child != nil,
		"parentFocus": fmt.Sprintf("%T", o.parent.Canvas().Focused()), "parentDraftBytes": len(o.input.Text)}
	if c := o.child; c != nil {
		d["opening"], d["childFocus"], d["childDirty"] = c.opening, fmt.Sprintf("%T", c.window.Canvas().Focused()), c.input.Text != c.accepted
	}
	record(event, d) // canvas focus is NOT a native active-window oracle
}

func main() {
	a := app.NewWithID("org.sdp.wci2.surface-probe")
	o := &owner{app: a, live: true, generation: 1}
	o.parent = a.NewWindow("WCI2 parent")
	o.input = newEntry("parent", "parent accepted")
	open := widget.NewButton("Open or focus child", o.open)
	o.opener = open
	o.parent.SetContent(container.NewVBox(widget.NewLabel("Temporary native adapter probe. No SDL or persistence."), o.input, open,
		widget.NewButton("Parent action (child must remain open)", func() { o.snapshot("parent-action") }),
		widget.NewButton("Close child from parent", func() { o.disposeChild("parent-command"); o.snapshot("parent-remains") }),
		widget.NewButton("Simulate reload", func() { o.disposeChild("reload"); o.generation++; o.snapshot("reload") }),
		widget.NewButton("Hide parent for one second", func() {
			o.disposeChild("parent-hidden")
			o.parent.Hide()
			record("parent-hidden", nil)
			go func() {
				time.Sleep(time.Second)
				fyne.Do(func() {
					if o.live {
						o.parent.Show()
						record("parent-shown", nil)
					}
				})
			}()
		}),
		widget.NewButton("Close parent and exit", o.closeParent)))
	o.parent.SetCloseIntercept(o.closeParent)
	o.parent.SetOnClosed(func() { o.live = false; o.disposeChild("parent-closed"); record("parent-onclosed", nil) })
	o.parent.Resize(fyne.NewSize(660, 420))
	record("probe-start", map[string]any{"fyne": "v2.8.1", "purpose": "manual native experiment, not SDUI acceptance"})
	o.parent.ShowAndRun()
}
