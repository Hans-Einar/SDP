package codegen

import (
	"math"
	"strings"
	"testing"
)

func TestLiteralBoundaries(t *testing.T) {
	for _, v := range []any{func() {}, math.Inf(1), map[int]int{1: 2}, make(chan int)} {
		if _, e := Literal(v, nil); e == nil {
			t.Fatal("accepted", v)
		}
	}
	a, e := Literal(map[string]any{"z": float64(2), "a": "$(echo nope)\n\""}, nil)
	if e != nil {
		t.Fatal(e)
	}
	b, e := Literal(map[string]any{"a": "$(echo nope)\n\"", "z": float64(2)}, nil)
	if e != nil || a != b || !strings.Contains(a, "float64(2)") {
		t.Fatal("unstable/untagged values", a, b, e)
	}
}
