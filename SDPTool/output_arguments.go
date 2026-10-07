package sdptool

import (
	"fmt"
	"strconv"
	"strings"
)

// Output selection is global, before or after the command. Preserve values of
// existing string options, even when those values happen to be "--json".
func outputArguments(args []string) ([]string, bool, error) {
	values := map[string]bool{}
	for _, name := range []string{"plan", "model", "output", "uri", "revision", "renderer", "entry", "viewer", "sdl-tool", "viewpoint", "artifact", "previous-artifact", "manifest", "plan-output", "apply", "resume", "release", "test-key", "notes", "task", "catalogue"} {
		values[name] = true
	}
	cleaned := make([]string, 0, len(args))
	jsonMode, releaseLog := false, false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			cleaned = append(cleaned, args[i:]...)
			break
		}
		name, value, equal := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		if strings.HasPrefix(arg, "-") && name == "json" {
			enabled := true
			var err error
			if equal {
				enabled, err = strconv.ParseBool(value)
			}
			if err != nil {
				return cleaned, jsonMode, fmt.Errorf("invalid --json value %q", value)
			}
			jsonMode = enabled
			continue
		}
		cleaned = append(cleaned, arg)
		if arg == "release-log" {
			releaseLog = true
		}
		if strings.HasPrefix(arg, "-") && !equal && (values[name] || name == "version" && releaseLog) && i+1 < len(args) {
			i++
			cleaned = append(cleaned, args[i])
		} else if arg == "preview" && i+1 < len(args) {
			// The required FILE positional can itself begin with a dash.
			i++
			cleaned = append(cleaned, args[i])
		}
	}
	return cleaned, jsonMode, nil
}
