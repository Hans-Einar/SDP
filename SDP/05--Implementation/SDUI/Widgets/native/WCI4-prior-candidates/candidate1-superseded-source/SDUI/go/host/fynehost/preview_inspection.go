package fynehost

import (
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
)

// Accessible values are queried from the actual mounted adapter, independently
// from the immutable provider identity and diagnostic inventory.
func (b *Bundle) inspectPreviews() map[string]map[string]any {
	out := map[string]map[string]any{}
	for path, outcome := range b.previewOutcomes {
		record := map[string]any{
			"kind": outcome.Kind, "status": outcome.Status, "description": outcome.Description,
			"diagnostic": outcome.Diagnostic, "providerID": outcome.ProviderID,
			"revision": outcome.Revision, "sha256": outcome.SHA256, "visible": false, "mounted": false,
		}
		out[path] = record
		if b.closed || b.view == nil || b.presentation == nil {
			continue
		}
		control := b.view.controls[path]
		if control == nil {
			continue
		}
		c, ok := control.widget.(*previewControl)
		if !ok {
			continue
		}
		frame, mounted := b.presentation.previews[path]
		if !mounted {
			continue
		}
		cnv := b.canvasFor(path)
		canvas, title := b.inspectionCanvas(cnv)
		origin := b.inspectionOrigin(path)
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(control.clip)
		size := control.clip.Size()
		clip := layout.Rect{X: float64(position.X), Y: float64(position.Y), W: float64(size.Width), H: float64(size.Height)}.Intersect(inspectionBounds(cnv))
		accessible := fyne.Accessible(c)
		record["mounted"], record["canvas"], record["title"] = true, canvas, title
		record["rect"], record["clip"] = inspectionTranslate(frame.box, origin), clip
		for name, rect := range map[string]layout.Rect{"image": frame.image, "caption": frame.caption, "statusRect": frame.status} {
			if rect.W > 0 && rect.H > 0 {
				rect = inspectionTranslate(rect, origin)
			}
			record[name] = rect
		}
		record["visible"] = c.Visible() && control.clip.Visible() && clip.W > 0 && clip.H > 0
		record["accessibleLabel"], record["accessibleRole"] = accessible.AccessibilityLabel(), accessible.AccessibilityRole()
	}
	return out
}
