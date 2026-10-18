package governance

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"
)

// Decode rejects duplicate keys as well as unknown fields. This check precedes
// typed decoding so JSON cannot hide authority or revision fields by repetition.
func Decode(b []byte, out any) error {
	if len(b) > MaxRequest || !utf8.Valid(b) {
		return fail("invalid_request", "request exceeds 1 MiB or is not UTF-8")
	}
	if e := decode(b, out); e != nil {
		return e
	}
	var tree any
	_ = json.Unmarshal(b, &tree)
	if e := exactFields(tree, reflect.TypeOf(out)); e != nil {
		return e
	}
	var null func(any) bool
	null = func(v any) bool {
		switch x := v.(type) {
		case nil:
			return true
		case map[string]any:
			for _, v := range x {
				if null(v) {
					return true
				}
			}
		case []any:
			for _, v := range x {
				if null(v) {
					return true
				}
			}
		}
		return false
	}
	if null(tree) {
		return fail("invalid_request", "explicit null values are not supported")
	}
	return nil
}
func decode(b []byte, out any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var walk func(int) error
	walk = func(depth int) error {
		if depth > 32 {
			return fail("invalid_request", "JSON nesting limit")
		}
		t, e := d.Token()
		if e != nil {
			return e
		}
		v, ok := t.(json.Delim)
		if !ok {
			return nil
		}
		switch v {
		case '{':
			keys := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				if e != nil {
					return e
				}
				name, ok := k.(string)
				if !ok || keys[name] {
					return fail("invalid_request", "duplicate JSON key")
				}
				keys[name] = true
				if e = walk(depth + 1); e != nil {
					return e
				}
			}
		case '[':
			for d.More() {
				if e = walk(depth + 1); e != nil {
					return e
				}
			}
		default:
			return fail("invalid_request", "unexpected JSON delimiter")
		}
		_, e = d.Token()
		return e
	}
	if e := walk(0); e != nil {
		return fail("invalid_request", e.Error())
	}
	if _, e := d.Token(); e != io.EOF {
		return fail("invalid_request", "trailing JSON")
	}
	d = json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if e := d.Decode(out); e != nil {
		return fail("invalid_request", e.Error())
	}
	return nil
}

// Normalize through typed values first: omitted zero-valued fields have their
// explicit defaults. encoding/json sorts map keys recursively; arrays keep order.
func requestDigest(v any) string {
	b, _ := json.Marshal(v)
	var x any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	_ = d.Decode(&x)
	b, _ = json.Marshal(x)
	return hash(b)
}

// encoding/json accepts case-insensitive field aliases. The wire contract does
// not: schema and Schema must never be two ways to supply the same authority.
func exactFields(v any, typ reflect.Type) error {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	switch x := v.(type) {
	case map[string]any:
		if typ.Kind() != reflect.Struct {
			return nil
		}
		fields := map[string]reflect.Type{}
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			name := strings.Split(f.Tag.Get("json"), ",")[0]
			if name != "" && name != "-" {
				fields[name] = f.Type
			}
		}
		for name, value := range x {
			ft, ok := fields[name]
			if !ok {
				return fail("invalid_request", "unknown JSON field: "+name)
			}
			if e := exactFields(value, ft); e != nil {
				return e
			}
		}
	case []any:
		if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
			return nil
		}
		for _, value := range x {
			if e := exactFields(value, typ.Elem()); e != nil {
				return e
			}
		}
	}
	return nil
}
