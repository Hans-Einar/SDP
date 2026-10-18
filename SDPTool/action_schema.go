package sdptool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
)

// Schemas and strict decoders share the actual request types, including domain requests.
func actionSchema(t reflect.Type) *InputSchema {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return &InputSchema{Type: "string"}
	case reflect.Struct:
		no := false
		s := &InputSchema{Type: "object", Properties: map[string]*InputSchema{}, AdditionalProperties: &no}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := strings.Split(f.Tag.Get("json"), ",")
			name := tag[0]
			if name == "" || name == "-" {
				continue
			}
			v := actionSchema(f.Type)
			v.PathBase = f.Tag.Get("pathBase")
			if f.Tag.Get("nonempty") == "true" {
				n := 1
				v.MinLength = &n
			}
			s.Properties[name] = v
			if len(tag) < 2 || tag[1] != "omitempty" {
				s.Required = append(s.Required, name)
			}
		}
		return s
	default:
		panic("unsupported action schema type: " + t.String())
	}
}
func validateActionShape(raw json.RawMessage, s *InputSchema) error {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return fmt.Errorf("null is not allowed")
	}
	if s.Type == "string" {
		var str string
		if err := json.Unmarshal(raw, &str); err != nil {
			return fmt.Errorf("string required")
		}
		n := utf8.RuneCountInString(str)
		if s.MinLength != nil && n < *s.MinLength {
			return fmt.Errorf("string shorter than minimum")
		}
		if s.MaxLength != nil && n > *s.MaxLength {
			return fmt.Errorf("string longer than maximum")
		}
		if len(s.Enum) > 0 {
			ok := false
			for _, v := range s.Enum {
				ok = ok || str == v
			}
			if !ok {
				return fmt.Errorf("unsupported value")
			}
		}
	}
	if s.Type == "object" || len(s.Required) > 0 {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
			return fmt.Errorf("object required")
		}
		for _, key := range s.Required {
			if _, ok := fields[key]; !ok {
				return fmt.Errorf("missing field %s", key)
			}
		}
		for key, v := range fields {
			prop, ok := s.Properties[key]
			if !ok {
				if s.AdditionalProperties != nil && !*s.AdditionalProperties {
					return fmt.Errorf("unknown or noncanonical field %s", key)
				}
				continue
			}
			if err := validateActionShape(v, prop); err != nil {
				return fmt.Errorf("%s: %w", key, err)
			}
		}
	}
	if len(s.OneOf) > 0 {
		matches := 0
		for _, variant := range s.OneOf {
			if validateActionShape(raw, variant) == nil {
				matches++
			}
		}
		if matches != 1 {
			return fmt.Errorf("exactly one alternative is required")
		}
	}
	return nil
}
func decodeAction[T any](raw json.RawMessage, s *InputSchema) (any, error) {
	var v T
	if err := blueprintstate.Strict(raw, &v); err != nil {
		return nil, err
	}
	if err := validateActionShape(raw, s); err != nil {
		return nil, err
	}
	return v, nil
}
