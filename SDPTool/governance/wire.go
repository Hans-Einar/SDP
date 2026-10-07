package governance

import (
	_ "embed"
	"encoding/json"
)

//go:embed schema.json
var wireSchema []byte

// WireDefinition shares the exact application schema with transport adapters.
// Each call returns an independent JSON value.
func WireDefinition(name string) map[string]any {
	var schema map[string]any
	_ = json.Unmarshal(wireSchema, &schema)
	defs := schema["$defs"].(map[string]any)
	out := defs[name].(map[string]any)
	// The root definition must not also occur inside its own $defs map.
	copy := map[string]any{}
	for k, v := range out {
		copy[k] = v
	}
	copy["$defs"] = defs
	return copy
}
