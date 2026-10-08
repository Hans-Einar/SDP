package parser

import (
	"regexp"
	"strings"
)

var commandFields = map[string]string{
	"label": "string", "icon": "string", "tooltip": "string", "toggle": "boolean", "checked": "boolean",
	"exclusive": "string", "key": "string", "context": "string", "callback": "reference", "effect": "string", "target": "string",
}
var interactions = map[string]widgetSchema{
	"command":   {"label", commandFields, "label"},
	"button":    {"label", buttonFields(), ""},
	"menu":      {"label", map[string]string{"label": "string", "mode": "string", "target": "string"}, "label"},
	"menuGroup": {"label", map[string]string{"label": "string"}, "label"},
	"item":      {"", map[string]string{"command": "string"}, "command"},
	"separator": {"", map[string]string{}, ""},
	"dialog":    {"label", map[string]string{"label": "string", "modal": "boolean", "callback": "reference"}, "label"},
}

func buttonFields() map[string]string {
	out := map[string]string{"command": "string"}
	for k, v := range commandFields {
		out[k] = v
	}
	return out
}
func boolArg(n *Instance, key string) bool {
	v, _ := n.Arguments[key].(Literal)
	b, _ := v.Value.(bool)
	return b
}
func hasArg(n *Instance, key string) bool { _, ok := n.Arguments[key]; return ok }

// IsCommandButton reports the explicit M2 opt-in, never inferring it from callback/icon/tooltip.
func IsCommandButton(n *Instance) bool {
	if n == nil || n.Profile != "sdui/0.3" || n.Kind != "widget" || n.Widget != "button" {
		return false
	}
	for _, key := range strings.Fields("command toggle checked exclusive key context target effect") {
		if hasArg(n, key) {
			return true
		}
	}
	return false
}
func isInteractionNode(n *Instance) bool {
	return n.Profile == "sdui/0.3" && ((n.Kind == "widget" && member(n.Widget, "command item")) ||
		(n.Kind == "composition" && member(n.Widget, "menu dialog")) || IsCommandButton(n))
}

// IsAuxiliary identifies declarations excluded from ordinary flow, not hidden content.
func IsAuxiliary(n *Instance) bool {
	if n == nil || n.Profile != "sdui/0.3" {
		return false
	}
	return n.Kind == "widget" && member(n.Widget, "command item separator") ||
		n.Kind == "composition" && (member(n.Widget, "menuGroup dialog") || n.Widget == "menu" && member(n.Argument("mode"), "context submenu"))
}
func validateInteractionSource(n *Node, profile string) {
	if profile != "sdui/0.3" {
		return
	}
	kind := val(n.Widget)
	if _, ok := interactions[kind]; !ok {
		return
	}
	body := member(kind, "menu menuGroup dialog")
	if body != (n.Kind == "composition") {
		fail("interaction-body", "Only menu/menuGroup/dialog require bodies", n.Span)
	}
	if !body && len(n.Rows) != 0 {
		fail("interaction-body", "Leaf interaction cannot have children", n.Span)
	}
	i := &Instance{Profile: profile, Kind: n.Kind, Widget: kind, Arguments: widgetArguments(n, profile), Span: n.Span}
	if member(kind, "command dialog") && n.Name == nil {
		fail("interaction-name", "Command/dialog requires a stable name", n.Span)
	}
	for _, key := range strings.Fields("label icon exclusive key target command") {
		if v, ok := i.Arguments[key].(Literal); ok && strings.TrimSpace(v.Value.(string)) == "" {
			fail("interaction-argument", "Nonempty "+key+" required", v.Span)
		}
	}
	if kind == "button" && !hasArg(i, "label") && !hasArg(i, "command") {
		fail("missing-argument", "Button needs label or command", n.Span)
	}
	if kind == "menu" {
		mode := i.Argument("mode")
		if mode == "" && !hasArg(i, "mode") {
			mode = "bar"
		}
		if !member(mode, "bar context submenu") {
			fail("menu-mode", "Unknown menu mode", n.Span)
		}
		if (mode == "context") != hasArg(i, "target") {
			fail("menu-target", "Only context menus require a target", n.Span)
		}
	}
	if icon := i.Argument("icon"); strings.Contains(icon, "://") {
		fail("interaction-icon", "Icon must be a symbolic resource, not a URL", n.Span)
	}
	if ref, ok := i.Arguments["callback"].(Reference); ok && (kind == "command" || kind == "dialog" || IsCommandButton(i)) && ref.Member != "invoke" {
		fail("interaction-callback", "M2 callbacks require @invoke", ref.Span)
	}
	if kind == "command" || IsCommandButton(i) {
		validateCommandSchema(i)
	}
	if body {
		for _, row := range n.Rows {
			for _, c := range row.Items {
				if c.Role != nil {
					fail("interaction-region", "Body assignments are names", c.Span)
				}
			}
		}
	}
}
func validateCommandSchema(n *Instance) {
	if hasArg(n, "command") {
		for key := range n.Arguments {
			if !member(key, "command label icon tooltip $scope") {
				fail("command-override", "Referring button cannot override "+key, n.Span)
			}
		}
		return
	}
	toggle := boolArg(n, "toggle")
	if (hasArg(n, "checked") || hasArg(n, "exclusive")) && !toggle {
		fail("command-toggle", "checked/exclusive require toggle=true", n.Span)
	}
	context := n.Argument("context")
	if context == "" && !hasArg(n, "context") {
		context = "none"
	}
	if !member(context, "none widget item") {
		fail("command-context", "Unknown context", n.Span)
	}
	effect := n.Argument("effect")
	if hasArg(n, "effect") {
		if !member(effect, "open accept cancel close") {
			fail("command-effect", "Unknown effect", n.Span)
		}
		if hasArg(n, "callback") || hasArg(n, "toggle") || context != "none" {
			fail("command-effect", "Effect forbids callback/toggle and requires context none", n.Span)
		}
	}
	if (context != "none" || effect == "open") != hasArg(n, "target") {
		fail("command-target", "Context/open requires target; otherwise target is forbidden", n.Span)
	}
	if hasArg(n, "key") {
		validateCommandKey(n.Argument("key"), n.Span)
	}
}

var commandKey = regexp.MustCompile(`^(?:[A-Z0-9]|F(?:[1-9]|1[0-2]))$`)

func validateCommandKey(key string, span Span) {
	parts := strings.Split(key, "+")
	last := parts[len(parts)-1]
	if !commandKey.MatchString(last) || len(parts) == 1 && len(last) == 1 {
		fail("command-key", "Expected modified letter/digit or F1-F12", span)
	}
	previous := -1
	primary, ctrl := false, false
	for _, mod := range parts[:len(parts)-1] {
		index := -1
		for i, m := range []string{"Primary", "Ctrl", "Alt", "Shift"} {
			if mod == m {
				index = i
			}
		}
		if index <= previous {
			fail("command-key", "Unknown, repeated or unordered modifier", span)
		}
		previous = index
		primary = primary || mod == "Primary"
		ctrl = ctrl || mod == "Ctrl"
	}
	if primary && ctrl {
		fail("command-key", "Primary and Ctrl cannot combine", span)
	}
}
func validateInteractionLayout(n *Instance) {
	if n.Profile != "sdui/0.3" {
		return
	}
	restricted := n.Widget == "command" || n.Widget == "menu" && n.Argument("mode") == "context"
	for key := range n.Layout {
		if n.Widget == "separator" || restricted && !member(key, "enabled visible") {
			fail("interaction-layout", "Unsupported declaration formatting "+key, n.Span)
		}
	}
}
func validateInteractionPlacement(n, parent *Instance) {
	menuParent := parent != nil && parent.Kind == "composition" && member(parent.Widget, "menu menuGroup")
	if n.Profile == "sdui/0.3" {
		if menuParent && !(n.Kind == "widget" && member(n.Widget, "item separator") || n.Kind == "composition" && member(n.Widget, "menu menuGroup")) {
			fail("menu-child", "Menus contain items, separators, groups or submenus", n.Span)
		}
		if parent != nil && member(n.Widget, "item separator menuGroup") && !menuParent {
			fail("menu-placement", "Menu child outside menu", n.Span)
		}
		if n.Widget == "menu" && (menuParent != (n.Argument("mode") == "submenu")) && parent != nil {
			fail("menu-placement", "Nested menus require mode=submenu; submenu requires menu parent", n.Span)
		}
		if n.Kind == "composition" && member(n.Widget, "menu menuGroup") {
			for _, row := range n.Rows {
				if len(row) != 1 {
					fail("menu-row", "One menu child per row required", n.Span)
				}
			}
		}
	}
	for _, row := range n.Rows {
		for _, child := range row {
			validateInteractionPlacement(child, n)
		}
	}
	for _, region := range n.Regions {
		validateInteractionPlacement(region.Node, n)
	}
}
