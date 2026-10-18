package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math"
)

// Tooltip paint lives in the existing canvas content, not OverlayStack: a new
// top overlay would steal menu/focus/input. It owns no event or model state.
func (b *Bundle) hideTooltip() {
	if b.tip != nil && b.tipParent != nil {
		b.tipParent.Remove(b.tip)
	}
	b.tip = nil
	b.tipParent = nil
}
func (b *Bundle) showTooltip(path string, button *commandButton, show bool) {
	b.hideTooltip()
	if !show || !b.livePane() || button.tooltip == "" {
		return
	}
	parent := b.view.Container
	if owner := surfacePath(b.presentation.snapshot, path); owner != "" {
		if surface := b.surfaces[owner]; surface != nil {
			parent = surface.content
		}
	}
	label := widget.NewLabel(button.tooltip)
	background := canvas.NewRectangle(theme.BackgroundColor())
	background.StrokeColor = theme.ForegroundColor()
	background.StrokeWidth = 1
	box := container.NewStack(background, container.NewPadded(label))
	natural := box.MinSize()
	width := float32(math.Min(float64(parent.Size().Width), float64(natural.Width)))
	height := natural.Height
	label.Truncation = fyne.TextTruncateEllipsis
	position := fyne.CurrentApp().Driver().AbsolutePositionForObject(button).Subtract(fyne.CurrentApp().Driver().AbsolutePositionForObject(parent)).Add(fyne.NewPos(0, button.Size().Height))
	if position.X+width > parent.Size().Width {
		position.X = parent.Size().Width - width
	}
	if position.Y+height > parent.Size().Height {
		position.Y -= button.Size().Height + height
	}
	if position.X < 0 {
		position.X = 0
	}
	if position.Y < 0 {
		position.Y = 0
	}
	box.Resize(fyne.NewSize(width, height))
	box.Move(position)
	b.tip, b.tipParent = box, parent
	parent.Add(box)
}
