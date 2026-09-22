// Package runtime invokes explicitly registered Go handlers for action-core.
// It owns model revisions and call validation, not domain algorithms or UI state.
package runtime

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

type Value struct {
	Type    parser.ScalarType `json:"type"`
	Text    string            `json:"text,omitempty"`
	Integer int64             `json:"integer,omitempty"`
	Boolean bool              `json:"boolean,omitempty"`
}

func Text(v string) Value   { return Value{Type: parser.TextType, Text: v} }
func Integer(v int64) Value { return Value{Type: parser.IntegerType, Integer: v} }
func Boolean(v bool) Value  { return Value{Type: parser.BooleanType, Boolean: v} }

type Record map[string]Value

func cloneRecord(r Record) Record {
	out := Record{}
	for k, v := range r {
		out[k] = v
	}
	return out
}
func checkRecord(record Record, typ parser.RecordType) error {
	if len(record) != len(typ) {
		return fmt.Errorf("record-fields: expected exactly %d fields", len(typ))
	}
	for name, kind := range typ {
		v, ok := record[name]
		if !ok || v.Type != kind {
			return fmt.Errorf("record-type: %s requires %s", name, kind)
		}
		switch kind {
		case parser.TextType:
			if !utf8.ValidString(v.Text) || strings.ContainsRune(v.Text, 0) || len(v.Text) > 32768 || v.Integer != 0 || v.Boolean {
				return fmt.Errorf("record-value: invalid text %s", name)
			}
		case parser.IntegerType:
			if v.Text != "" || v.Boolean {
				return fmt.Errorf("record-value: invalid integer %s", name)
			}
		case parser.BooleanType:
			if v.Text != "" || v.Integer != 0 {
				return fmt.Errorf("record-value: invalid boolean %s", name)
			}
		default:
			return fmt.Errorf("record-type: unsupported %s", kind)
		}
	}
	return nil
}
func cloneType(t parser.RecordType) parser.RecordType {
	c := parser.RecordType{}
	for k, v := range t {
		c[k] = v
	}
	return c
}
func sameType(a, b parser.RecordType) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// ValidateRecord exposes the same closed scalar contract to explicit adapters.
func ValidateRecord(record Record, typ parser.RecordType) error { return checkRecord(record, typ) }
