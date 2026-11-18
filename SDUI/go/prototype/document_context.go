package prototype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"unicode/utf8"
)

const contextLimit = 1 << 20
const svgLimit = 4 << 20
const resourcesLimit = 32 << 20

type contextOption struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
}
type contextSVG struct {
	File   string  `json:"file"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
type contextViewport struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
type documentContext struct {
	Schema   string                     `json:"schema"`
	Viewport *contextViewport           `json:"viewport"`
	Choices  map[string][]contextOption `json:"choices"`
	SVG      map[string]contextSVG      `json:"svg"`
}

// Inspect tokens before decoding: encoding/json otherwise silently replaces
// duplicate keys (including duplicate exact-path fixture registrations).
func jsonValue(d *json.Decoder, depth int) error {
	if depth > 32 {
		return fmt.Errorf("context: JSON nesting exceeds 32")
	}
	token, err := d.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			token, err := d.Token()
			if err != nil {
				return err
			}
			key, ok := token.(string)
			if !ok || seen[key] {
				return fmt.Errorf("context: duplicate/invalid JSON key %q", key)
			}
			seen[key] = true
			if err := jsonValue(d, depth+1); err != nil {
				return err
			}
		}
	case '[':
		for d.More() {
			if err := jsonValue(d, depth+1); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("context: unexpected JSON delimiter")
	}
	_, err = d.Token()
	return err
}
func objectKeys(data json.RawMessage, allowed []string, required []string) (map[string]json.RawMessage, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil || fields == nil {
		return nil, fmt.Errorf("context: expected object")
	}
	for key := range fields {
		found := false
		for _, name := range allowed {
			found = found || name == key
		}
		if !found {
			return nil, fmt.Errorf("context: unknown key %q", key)
		}
	}
	for _, key := range required {
		v, ok := fields[key]
		if !ok || bytes.Equal(bytes.TrimSpace(v), []byte("null")) {
			return nil, fmt.Errorf("context: missing %s", key)
		}
	}
	return fields, nil
}
func parseDocumentContext(data []byte) (documentContext, error) {
	var c documentContext
	if !utf8.Valid(data) {
		return c, fmt.Errorf("context: invalid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	if err := jsonValue(d, 0); err != nil {
		return c, err
	}
	if _, err := d.Token(); err != io.EOF {
		return c, fmt.Errorf("context: trailing JSON data")
	}
	fields, err := objectKeys(data, []string{"schema", "viewport", "choices", "svg"}, []string{"schema"})
	if err != nil {
		return c, err
	}
	if raw, ok := fields["viewport"]; ok {
		if _, err := objectKeys(raw, []string{"width", "height"}, []string{"width", "height"}); err != nil {
			return c, err
		}
	}
	for _, family := range []string{"choices", "svg"} {
		raw, ok := fields[family]
		if !ok {
			continue
		}
		var paths map[string]json.RawMessage
		if err := json.Unmarshal(raw, &paths); err != nil || paths == nil {
			return c, fmt.Errorf("context: %s must be an object", family)
		}
		for path, value := range paths {
			if path == "" {
				return c, fmt.Errorf("context: empty instance path")
			}
			if family == "svg" {
				if _, err := objectKeys(value, []string{"file", "width", "height"}, []string{"file", "width", "height"}); err != nil {
					return c, err
				}
			} else {
				var options []json.RawMessage
				if err := json.Unmarshal(value, &options); err != nil || options == nil {
					return c, fmt.Errorf("context: choices at %s must be an array", path)
				}
				for _, option := range options {
					if _, err := objectKeys(option, []string{"id", "label", "enabled"}, []string{"id", "label", "enabled"}); err != nil {
						return c, err
					}
				}
			}
		}
	}
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&c); err != nil {
		return c, err
	}
	if c.Schema != "sdui-prototype-context/1" {
		return c, fmt.Errorf("context: unsupported schema %q", c.Schema)
	}
	if c.Viewport != nil && (!positive(c.Viewport.Width) || !positive(c.Viewport.Height) || c.Viewport.Width > 32768 || c.Viewport.Height > 32768) {
		return c, fmt.Errorf("context: viewport must have finite positive dimensions at most 32768")
	}
	for path, svg := range c.SVG {
		if !positive(svg.Width) || !positive(svg.Height) || svg.Width > 32768 || svg.Height > 32768 {
			return c, fmt.Errorf("context: invalid SVG dimensions at %s", path)
		}
	}
	return c, nil
}
func positive(v float64) bool { return v > 0 && !math.IsInf(v, 0) && !math.IsNaN(v) }
