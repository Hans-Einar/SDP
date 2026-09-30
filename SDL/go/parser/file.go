package parser

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// File is immutable syntax. Accessors return copies; no binding context is cached.
type File struct {
	name, text string
	model      *Model
}

func (f *File) TokenCount() int { return f.model.tokenCount }
func (f *File) Name() string    { return f.name }
func (f *File) Text() string    { return f.text }
func CopyModel(m *Model) *Model {
	b, _ := json.Marshal(m)
	var out Model
	_ = json.Unmarshal(b, &out)
	return &out
}
func (f *File) Model() *Model {
	m := CopyModel(f.model)
	if m.Header.Version == "0.6" {
		tagSpans(reflect.ValueOf(m), f.name)
	}
	return m
}
func tagSpans(v reflect.Value, name string) {
	if v.Kind() == reflect.Pointer {
		if !v.IsNil() {
			tagSpans(v.Elem(), name)
		}
		return
	}
	if v.Type() == reflect.TypeOf(Span{}) {
		if v.FieldByName("Line").Int() > 0 {
			v.FieldByName("Source").SetString(name)
		}
		return
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			tagSpans(v.Field(i), name)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			tagSpans(v.Index(i), name)
		}
	}
}

// SyntaxCache reuses parsing by bytes/profile implementation, never System bindings.
// Bounded eviction affects performance only. A zero-value cache is ready to use.
type SyntaxCache struct {
	mu           sync.Mutex
	models       map[[32]byte]*Model
	hits, misses uint64
}

func (c *SyntaxCache) Stats() (uint64, uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hits, c.misses
}
func (c *SyntaxCache) Parse(name, text string) (*File, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := sha256.Sum256([]byte("sdl-file-0.6/1\x00" + text))
	if m := c.models[key]; m != nil {
		c.hits++
		return &File{name, text, m}, nil
	}
	m, e := Parse(text)
	if e != nil {
		d := e.(Diagnostic)
		d.Span.Source = name
		return nil, d
	}
	if len(c.models) >= 128 {
		c.models = nil
	}
	if c.models == nil {
		c.models = map[[32]byte]*Model{}
	}
	c.models[key] = m
	c.misses++
	return &File{name, text, m}, nil
}
func (f *File) Inspect() map[string]any {
	m := f.Model()
	deps := []Include{}
	deps = append(deps, m.Includes...)
	for _, s := range m.Statements {
		if s.Path != "" {
			deps = append(deps, Include{s.Path + ".design", s.PathSpan})
		}
	}
	return map[string]any{"schema": "sdl-fragment/1", "state": "parsed-context-required", "validated": false, "file": f.name, "ast": Data(m), "dependencies": Data(deps), "note": fmt.Sprintf("Syntax only (%s); supply an explicit System entry for resolution", m.Header.Version)}
}
