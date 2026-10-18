package fynehost

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	goruntime "runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/fyne-io/oksvg"
)

// prepareIcons owns immutable bytes before native publication and validates the
// concrete renderer formats, including resources in hidden/closed declarations.
func prepareIcons(root *parser.Instance, resources map[string]fyne.Resource) (map[string]fyne.Resource, error) {
	out := map[string]fyne.Resource{}
	var err error
	root.Walk(func(n *parser.Instance) {
		if err != nil {
			return
		}
		id := n.Argument("icon")
		if id == "" {
			return
		}
		if _, ok := out[id]; ok {
			return
		}
		r := resources[id]
		if r == nil {
			err = fmt.Errorf("icon-resource: missing %q at %s", id, n.Path)
			return
		}
		data := append([]byte(nil), r.Content()...)
		if len(data) == 0 || len(data) > 4<<20 {
			err = fmt.Errorf("icon-resource: invalid size for %q", id)
			return
		}
		name := r.Name()
		if strings.HasSuffix(strings.ToLower(name), ".svg") {
			var icon *oksvg.SvgIcon
			icon, err = oksvg.ReadIconStream(bytes.NewReader(data), oksvg.StrictErrorMode)
			if err == nil && (icon.ViewBox.W <= 0 || icon.ViewBox.H <= 0 || math.IsNaN(icon.ViewBox.W) || math.IsNaN(icon.ViewBox.H) || math.IsInf(icon.ViewBox.W, 0) || math.IsInf(icon.ViewBox.H, 0)) {
				err = fmt.Errorf("missing finite SVG dimensions")
			}
		} else {
			var config image.Config
			config, _, err = image.DecodeConfig(bytes.NewReader(data))
			if err == nil && (config.Width <= 0 || config.Height <= 0 || config.Width > 4096 || config.Height > 4096) {
				err = fmt.Errorf("invalid raster dimensions")
			}
		}
		if err != nil {
			err = fmt.Errorf("icon-resource: %q at %s: %w", id, n.Path, err)
			return
		}
		out[id] = frozenIcon(r, name, data)
	})
	return out, err
}

type commandKey struct {
	shortcut      *desktop.CustomShortcut
	path, surface string
}

func prepareKeys(root *parser.Instance) ([]commandKey, error) {
	identities, err := parser.ResolveInteractions(root)
	if err != nil {
		return nil, err
	}
	var keys []commandKey
	seen := map[string]string{}
	root.Walk(func(n *parser.Instance) {
		if err != nil {
			return
		}
		key := n.Argument("key")
		if key == "" {
			return
		}
		id := identities[n.Path]
		if id.Command != n.Path {
			err = fmt.Errorf("command-key: non-owner %s", n.Path)
			return
		}
		shortcut, e := nativeShortcut(key)
		if e != nil {
			err = fmt.Errorf("command-key: %s: %w", n.Path, e)
			return
		}
		identity := id.Dialog + "/" + shortcut.ShortcutName()
		if old := seen[identity]; old != "" {
			err = fmt.Errorf("command-key: duplicate %s at %s and %s", key, old, n.Path)
			return
		}
		seen[identity] = n.Path
		keys = append(keys, commandKey{shortcut: shortcut, path: n.Path, surface: id.Dialog})
	})
	return keys, err
}
func nativeShortcut(source string) (*desktop.CustomShortcut, error) {
	parts := strings.Split(source, "+")
	var modifier fyne.KeyModifier
	for _, part := range parts[:len(parts)-1] {
		var bit fyne.KeyModifier
		switch part {
		case "Primary":
			bit = fyne.KeyModifierControl
			if goruntime.GOOS == "darwin" {
				bit = fyne.KeyModifierSuper
			}
		case "Ctrl":
			bit = fyne.KeyModifierControl
		case "Alt":
			bit = fyne.KeyModifierAlt
		case "Shift":
			bit = fyne.KeyModifierShift
		default:
			return nil, fmt.Errorf("unknown modifier %q", part)
		}
		if modifier&bit != 0 {
			return nil, fmt.Errorf("duplicate normalized modifier")
		}
		modifier |= bit
	}
	key := fyne.KeyName(parts[len(parts)-1])
	if key == fyne.KeyF10 && modifier == fyne.KeyModifierShift || key == fyne.KeyF4 && modifier == fyne.KeyModifierAlt || goruntime.GOOS == "darwin" && modifier == fyne.KeyModifierSuper && (key == fyne.KeyQ || key == fyne.KeyW) {
		return nil, fmt.Errorf("host-reserved shortcut %s", source)
	}
	return &desktop.CustomShortcut{KeyName: key, Modifier: modifier}, nil
}

// Preserve the public theme-color contract as well as immutable bytes. Flattening
// a ThemedResource to StaticResource bakes the app color (often white/dark-mode)
// into our light component theme and makes an otherwise valid icon invisible.
type ownedThemedIcon struct {
	fyne.Resource
	color fyne.ThemeColorName
}

func (r ownedThemedIcon) ThemeColorName() fyne.ThemeColorName { return r.color }
func frozenIcon(source fyne.Resource, name string, data []byte) fyne.Resource {
	r := fyne.NewStaticResource(name, append([]byte(nil), data...))
	if themed, ok := source.(fyne.ThemedResource); ok {
		return ownedThemedIcon{r, themed.ThemeColorName()}
	}
	return r
}
