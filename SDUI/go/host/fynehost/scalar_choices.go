package fynehost

import (
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type choicePopup struct {
	target ui.FieldTarget

	native  *nativeMenu
	targets []ui.OptionTarget
	options []ui.ChoiceOption
}

func (c *scalarControl) openChoices() {
	if !c.editable() || c.popup != nil {
		return
	}
	b := c.bundle
	f, ok := b.Session.Field(c.state.Target.Handle)
	if !ok {
		return
	}
	popup := &choicePopup{target: f.Target, options: append([]ui.ChoiceOption(nil), f.Options...)}
	model := fyne.NewMenu("")
	for _, option := range popup.options {
		target := ui.OptionTarget{Handle: f.Target.Handle, ModelRevision: f.Target.ModelRevision, OptionGeneration: f.Target.OptionGeneration, OptionID: option.ID}
		popup.targets = append(popup.targets, target)
		item := fyne.NewMenuItem(option.Label, func() {
			if !c.editable() {
				return
			}
			b.owner.Mutate(func(s *ui.Session) error {
				change, err := s.ChooseOption(target)
				if err != nil {
					return err
				}
				return c.dispatch(s, change.Field.Target)
			})
		})
		item.Disabled = !option.Enabled
		item.Checked = option.ID == f.Proposed.OptionID
		model.Items = append(model.Items, item)
	}
	cv := b.canvasFor(c.path)
	popup.native = newNativeMenu(model, cv, func() {
		if c.popup != popup {
			return
		}
		c.popup = nil
		if b.livePane() {
			b.owner.restoreFocus(b)
			if b.owner.OnChange != nil {
				b.owner.OnChange(b)
			}
		}
	})
	popup.native.nativeTransition = func(fn func()) { muted := b.muted; b.muted = true; defer func() { b.muted = muted }(); fn() }
	c.popup = popup
	muted := b.muted
	b.muted = true
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c.body).Add(fyne.NewPos(0, c.body.Size().Height))
	popup.native.show(pos)
	b.muted = muted
	if b.owner.OnChange != nil {
		b.owner.OnChange(b)
	}
}
func (b *Bundle) syncChoices() {
	for _, obj := range b.view.Controls {
		c, ok := obj.(*scalarControl)
		if !ok || c.popup == nil {
			continue
		}
		popup := c.popup
		if popup.native.scoped {
			continue
		}
		field, found := b.Session.Field(c.state.Target.Handle)
		w, widgetOK := b.Session.Widget(c.state.Target.Handle.Path)
		stale := !found || !widgetOK || !w.Visible || !w.Enabled || field.ReadOnly || field.Target.ModelRevision != c.state.Target.ModelRevision
		stale = stale || field.Target.OptionGeneration != popup.target.OptionGeneration || field.Target.Handle != popup.target.Handle
		if stale {
			popup.native.close()
		}
	}
}
func (b *Bundle) closeChoices() {
	if b.view == nil {
		return
	}
	for _, obj := range b.view.Controls {
		if c, ok := obj.(*scalarControl); ok && c.popup != nil {
			c.popup.native.close()
		}
	}
}
func (b *Bundle) inspectFields() (map[string]map[string]any, map[string]map[string]any) {
	fields := map[string]map[string]any{}
	choices := map[string]map[string]any{}
	for path, obj := range b.view.Controls {
		c, ok := obj.(*scalarControl)
		if !ok {
			continue
		}
		cv := b.canvasFor(path)
		id, title := b.inspectionCanvas(cv)
		holder := b.view.controls[path]
		clip := inspectionRect(holder.clip).Intersect(inspectionBounds(cv))
		rect := func(o fyne.CanvasObject) layout.Rect {
			if o == nil {
				return layout.Rect{}
			}
			return inspectionRect(o)
		}
		label := rect(c.label)
		if c.check != nil {
			label = rect(c.body)
		}
		parts := map[string]any{"kind": c.kind, "canvas": id, "title": title, "label": label, "control": rect(c.body), "feedback": rect(c.feedback), "clip": clip, "visible": holder.clip.Visible() && clip.W > 0 && clip.H > 0, "entry": layout.Rect{}, "increment": layout.Rect{}, "decrement": layout.Rect{}, "slidertrack": layout.Rect{}, "sliderthumb": layout.Rect{}}
		if c.entry != nil {
			parts["entry"] = rect(c.entry)
			parts["increment"] = rect(c.increment)
			parts["decrement"] = rect(c.decrement)
		}
		if c.slider != nil && c.slider.renderer != nil {
			objects := c.slider.renderer.Objects()
			if len(objects) >= 3 {
				parts["slidertrack"] = rect(objects[0])
				parts["sliderthumb"] = rect(objects[2])
			}
			parts["orientation"] = "horizontal"
		}
		parts["labelText"] = c.label.Text
		fields[path] = parts
		if popup := c.popup; popup != nil && !popup.native.closed && !popup.native.dismissed {
			items := []map[string]any{}
			bounds := inspectionRect(popup.native.popup).Intersect(inspectionBounds(cv))
			for i, object := range popup.native.popup.Items {
				if i >= len(popup.options) {
					break
				}
				option := popup.options[i]
				r := rect(object)
				items = append(items, map[string]any{"id": option.ID, "label": option.Label, "enabled": option.Enabled, "rect": r, "clip": r.Intersect(bounds), "target": popup.targets[i]})
			}
			choices[path] = map[string]any{"canvas": id, "title": title, "generation": popup.target.OptionGeneration, "items": items}
		}
	}
	return fields, choices
}
