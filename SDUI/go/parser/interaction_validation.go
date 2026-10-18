package parser

import "strings"

// Validate the serialized interaction schema as well as parsed source. Generated
// constructors and detached runtime roots use the same strict selected-entry seam.
func validateNormalizedInteraction(n *Instance) {
	if n.Profile != "sdui/0.3" {
		return
	}
	schema, ok := interactions[n.Widget]
	if !ok || n.Kind != "widget" && n.Kind != "composition" {
		return
	}
	if (member(n.Widget, "menu menuGroup dialog")) != (n.Kind == "composition") {
		fail("interaction-body", "Invalid normalized interaction kind", n.Span)
	}
	if len(n.Regions) != 0 || n.Kind == "widget" && len(n.Rows) != 0 {
		fail("interaction-body", "Interaction body has invalid children/regions", n.Span)
	}
	for key, value := range n.Arguments {
		expected, known := schema.fields[key]
		if key == "$scope" && isInteractionNode(n) {
			expected = "string"
			known = true
		}
		if !known {
			fail("interaction-argument", n.Path+": unknown argument "+key, n.Span)
		}
		switch v := value.(type) {
		case Literal:
			valid := v.Kind == expected
			switch expected {
			case "string":
				_, ok := v.Value.(string)
				valid = valid && ok
			case "boolean":
				_, ok := v.Value.(bool)
				valid = valid && ok
			default:
				valid = false
			}
			if !valid {
				fail("interaction-argument", n.Path+": wrong type for "+key, v.Span)
			}
		case Reference:
			if expected != "reference" {
				fail("interaction-argument", n.Path+": unexpected reference", v.Span)
			}
			if (isInteractionNode(n) || n.Widget == "dialog") && v.Member != "invoke" {
				fail("interaction-callback", "M2 callbacks require @invoke", v.Span)
			}
		default:
			fail("interaction-argument", n.Path+": invalid argument "+key, n.Span)
		}
	}
	if schema.required != "" && !hasArg(n, schema.required) {
		fail("interaction-argument", "Missing "+schema.required, n.Span)
	}
	for _, key := range strings.Fields("label icon exclusive key target command") {
		if hasArg(n, key) && strings.TrimSpace(n.Argument(key)) == "" {
			fail("interaction-argument", "Nonempty "+key+" required", n.Span)
		}
	}
	if icon := n.Argument("icon"); strings.Contains(icon, "://") {
		fail("interaction-icon", "Icon must be a symbolic resource, not a URL", n.Span)
	}
	if n.Widget == "button" && !hasArg(n, "label") && !hasArg(n, "command") {
		fail("interaction-argument", "Button needs label or command", n.Span)
	}
	if member(n.Widget, "command dialog") {
		parts := strings.Split(n.Path, "/")
		if strings.HasPrefix(parts[len(parts)-1], "$") {
			fail("interaction-name", "Command/dialog requires a stable name", n.Span)
		}
	}
	if n.Widget == "menu" {
		mode := n.Argument("mode")
		if mode == "" && !hasArg(n, "mode") {
			mode = "bar"
		}
		if !member(mode, "bar context submenu") || (mode == "context") != hasArg(n, "target") {
			fail("menu-mode", "Invalid menu mode/target", n.Span)
		}
	}
	if n.Widget == "command" || IsCommandButton(n) {
		validateCommandSchema(n)
	}
	validateInteractionLayout(n)
}
