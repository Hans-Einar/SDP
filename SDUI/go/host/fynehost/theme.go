package fynehost

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/image/font/gofont/goregular"
)

var regular = fyne.NewStaticResource("Go-Regular.ttf", goregular.TTF)

type componentTheme struct{ font float32 }

func (t componentTheme) Font(fyne.TextStyle) fyne.Resource       { return regular }
func (t componentTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (t componentTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, theme.VariantLight)
}
func (t componentTheme) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNameText:
		return t.font
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameInnerPadding:
		return 6
	}
	return theme.DefaultTheme().Size(n)
}
