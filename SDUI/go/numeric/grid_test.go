package numeric

import (
	"math"
	"strings"
	"testing"
)

func grid(t *testing.T, a, b, d string) *Grid {
	t.Helper()
	g, e := NewGrid(a, b, d)
	if e != nil {
		t.Fatal(e)
	}
	return g
}
func TestExactTextVersusRoundedNumber(t *testing.T) {
	g := grid(t, "0", "1", "0.1")
	v, e := g.Parse("0.3")
	if e != nil || v != .3 {
		t.Fatal(v, e)
	}
	if _, e = g.Tick(.3); e != nil {
		t.Fatal(e)
	}
	if _, e = g.Parse("0.30000000000000000000001"); e == nil {
		t.Fatal("off-grid decimal rounded into admission")
	}
	if _, e = g.Tick(math.Nextafter(.3, 1)); e == nil {
		t.Fatal("typed tolerance")
	}
	for k := uint64(0); k <= g.LastTick(); k++ {
		v, e := g.At(k)
		if e != nil {
			t.Fatal(e)
		}
		back, e := g.Tick(v)
		if e != nil || back != k {
			t.Fatal(k, v, back, e)
		}
	}
}
func TestBoundedDecimalScan(t *testing.T) {
	for _, raw := range []string{"1e99999999999999999999999999999999999999", "1e-999999999999999999999999", "1e4097", "1e-4097", strings.Repeat("1", 32769), "01", "+1", "1.", ".1", "1e", "NaN", "Inf", " 1", "1_0", "1/2"} {
		if _, e := ExactInteger(raw); e == nil {
			t.Fatal("admitted", raw[:min(len(raw), 60)])
		}
	}
	for _, raw := range []string{"0e999999999999999999999999999999999999999", "-0e-999999999999999999999999", "0.000e9999"} {
		v, e := ExactInteger(raw)
		if e != nil || v != 0 {
			t.Fatal(raw, v, e)
		}
	}
	if _, e := decimal("1"+strings.Repeat("0", 4096)+"e-4096", "value"); e != nil {
		t.Fatal("stripped trailing compensation", e)
	}
	if _, e := decimal("0."+strings.Repeat("0", 4095)+"1", "value"); e != nil {
		t.Fatal("effective exponent boundary", e)
	}
}
func TestSafe53ExactNoCapAndFractionRejection(t *testing.T) {
	g := grid(t, "-9007199254740991", "9007199254740991", "1")
	if g.LastTick() != 18014398509481982 {
		t.Fatal(g.LastTick())
	}
	if e := g.ValidateSDL(); e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{"-9007199254740991", "9007199254740991"} {
		v, e := g.Parse(raw)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = Integer(v); e != nil {
			t.Fatal(e)
		}
	}
	for _, raw := range []string{"9007199254740991.1", "-9007199254740991.1", "9007199254740992", "1e-999"} {
		if _, e := ExactInteger(raw); e == nil {
			t.Fatal(raw)
		}
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), .1, 9007199254740992} {
		if _, e := Integer(v); e == nil {
			t.Fatal(v)
		}
	}
}
func TestRepresentabilityAndTickBounds(t *testing.T) {
	for _, tc := range [][3]string{{"0", "1", "0.0000000000000001"}, {"9007199254740992", "9007199254741000", "2"}, {"1e308", "1.00000000000001e308", "1e280"}, {"0", "1e-323", "4e-324"}, {"0", "1", "0"}, {"1", "1", "1"}, {"0", "1", "-1"}} {
		if _, e := NewGrid(tc[0], tc[1], tc[2]); e == nil {
			t.Fatal(tc)
		}
	}
	g := grid(t, "9007199254740992", "9007199254741004", "3")
	for k := uint64(0); k <= g.LastTick(); k++ {
		v, _ := g.At(k)
		got, e := g.Tick(v)
		if e != nil || got != k {
			t.Fatal(k, v, got, e)
		}
	}
	if e := g.ValidateSDL(); e == nil {
		t.Fatal("wide integral SDL grid admitted")
	}
	if _, e := g.At(g.LastTick() + 1); e == nil {
		t.Fatal("illegal tick")
	}
}
func TestGestureClampDoesNotChangeTextMembership(t *testing.T) {
	g := grid(t, "-1", "1", "0.3")
	k, e := g.Nearest(100)
	if e != nil || k != 6 {
		t.Fatal(k, e)
	}
	v, _ := g.At(k)
	if v != .8 {
		t.Fatal(v)
	}
	if _, e := g.Parse("1"); e == nil {
		t.Fatal("off-grid max accepted")
	}
	if _, e := g.Tick(1); e == nil {
		t.Fatal("typed max snapped")
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, e := g.Nearest(v); e == nil {
			t.Fatal(v)
		}
	}
}
func TestGridMethodsConcurrentAndImmutable(t *testing.T) {
	g := grid(t, "1000000000000.1", "1000000000010.1", "0.1")
	done := make(chan bool, 8)
	for i := 0; i < 8; i++ {
		go func() {
			for k := uint64(0); k <= g.LastTick(); k++ {
				v, e := g.At(k)
				if e != nil {
					t.Error(e)
				}
				got, e := g.Tick(v)
				if e != nil || got != k {
					t.Error(k, got, e)
				}
			}
			done <- true
		}()
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}
