// Package codegen produces Go model constructors, with no runtime source parsing.
package codegen

import (
	"fmt"
	"go/token"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Literal emits bounded, deterministic typed data. Only explicitly registered
// packages and exported data fields are accepted; executable values are rejected.
func Literal(value any, packages map[string]string) (string, error) {
	var b strings.Builder
	nodes := 0
	var typ func(reflect.Type) (string, error)
	typ = func(t reflect.Type) (string, error) {
		if t.Name() != "" {
			if t.PkgPath() == "" {
				return t.Name(), nil
			}
			a, ok := packages[t.PkgPath()]
			if !ok {
				return "", fmt.Errorf("unregistered type %s", t)
			}
			return a + "." + t.Name(), nil
		}
		switch t.Kind() {
		case reflect.Pointer:
			s, e := typ(t.Elem())
			return "*" + s, e
		case reflect.Slice:
			s, e := typ(t.Elem())
			return "[]" + s, e
		case reflect.Map:
			k, e := typ(t.Key())
			if e != nil {
				return "", e
			}
			v, e := typ(t.Elem())
			return "map[" + k + "]" + v, e
		case reflect.Interface:
			if t.NumMethod() == 0 {
				return "any", nil
			}
		}
		return "", fmt.Errorf("unsupported Go type %s", t)
	}
	var emit func(reflect.Value, int) error
	emit = func(v reflect.Value, depth int) error {
		nodes++
		if depth > 512 || nodes > 1000000 || b.Len() > 32<<20 {
			return fmt.Errorf("Go literal resource limit")
		}
		if !v.IsValid() {
			b.WriteString("nil")
			return nil
		}
		if v.Kind() == reflect.Interface {
			if v.IsNil() {
				b.WriteString("nil")
				return nil
			}
			return emit(v.Elem(), depth+1)
		}
		t, e := typ(v.Type())
		if e != nil {
			return e
		}
		switch v.Kind() {
		case reflect.Pointer:
			if v.IsNil() {
				b.WriteString("nil")
				return nil
			}
			b.WriteString("func() " + t + " { v := ")
			if e := emit(v.Elem(), depth+1); e != nil {
				return e
			}
			b.WriteString("; return &v }()")
		case reflect.Struct:
			b.WriteString(t + "{")
			for i := 0; i < v.NumField(); i++ {
				f := v.Type().Field(i)
				if !f.IsExported() {
					return fmt.Errorf("unexported data field %s", f.Name)
				}
				b.WriteString(f.Name + ":")
				if e := emit(v.Field(i), depth+1); e != nil {
					return e
				}
				b.WriteString(",")
			}
			b.WriteString("}")
		case reflect.Map:
			if v.IsNil() {
				b.WriteString("nil")
				return nil
			}
			if v.Type().Key().Kind() != reflect.String {
				return fmt.Errorf("only string map keys supported")
			}
			keys := v.MapKeys()
			sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
			b.WriteString(t + "{")
			for _, k := range keys {
				b.WriteString(strconv.Quote(k.String()) + ":")
				if e := emit(v.MapIndex(k), depth+1); e != nil {
					return e
				}
				b.WriteString(",")
			}
			b.WriteString("}")
		case reflect.Slice:
			if v.IsNil() {
				b.WriteString("nil")
				return nil
			}
			b.WriteString(t + "{")
			for i := 0; i < v.Len(); i++ {
				if e := emit(v.Index(i), depth+1); e != nil {
					return e
				}
				b.WriteString(",")
			}
			b.WriteString("}")
		case reflect.String:
			b.WriteString(t + "(" + strconv.Quote(v.String()) + ")")
		case reflect.Bool:
			b.WriteString(t + "(" + strconv.FormatBool(v.Bool()) + ")")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b.WriteString(t + "(" + strconv.FormatInt(v.Int(), 10) + ")")
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b.WriteString(t + "(" + strconv.FormatUint(v.Uint(), 10) + ")")
		case reflect.Float32, reflect.Float64:
			f := v.Float()
			if math.IsNaN(f) || math.IsInf(f, 0) {
				return fmt.Errorf("nonfinite literal")
			}
			b.WriteString(t + "(" + strconv.FormatFloat(f, 'g', -1, v.Type().Bits()) + ")")
		default:
			return fmt.Errorf("unsupported literal %s", v.Kind())
		}
		return nil
	}
	if e := emit(reflect.ValueOf(value), 0); e != nil {
		return "", e
	}
	return b.String(), nil
}
func ValidPackage(name string) bool { return token.IsIdentifier(name) && name != "_" }
