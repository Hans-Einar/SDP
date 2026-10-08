package numeric

import (
	"math"
	"math/big"
)

// Grid is immutable. Methods never return its mutable exact rational operands.
type Grid struct {
	min, max, step *big.Rat
	last           uint64
	integer        bool
}

func NewGrid(min, max, step string) (*Grid, error) {
	a, err := decimal(min, "min")
	if err != nil {
		return nil, err
	}
	b, err := decimal(max, "max")
	if err != nil {
		return nil, err
	}
	d, err := decimal(step, "step")
	if err != nil {
		return nil, err
	}
	if a.Cmp(b) >= 0 {
		return nil, invalid("number-range", "max")
	}
	if d.Sign() <= 0 {
		return nil, invalid("number-step", "step")
	}
	q := new(big.Rat).Quo(new(big.Rat).Sub(b, a), d)
	n := new(big.Int).Quo(q.Num(), q.Denom())
	if !n.IsUint64() {
		return nil, invalid("number-ticks", "step")
	}
	g := &Grid{a, b, d, n.Uint64(), safe(a) && safe(b) && safe(d)}
	if !g.integer {
		if g.last > 1<<26 {
			return nil, invalid("number-ticks", "step")
		}
		for _, endpoint := range []*big.Rat{a, b} {
			f, _ := endpoint.Float64()
			for _, dir := range []float64{math.Inf(-1), math.Inf(1)} {
				next := math.Nextafter(f, dir)
				if math.IsInf(next, 0) {
					return nil, invalid("number-spacing", "step")
				}
				spacing := new(big.Rat).Sub(new(big.Rat).SetFloat64(next), new(big.Rat).SetFloat64(f))
				spacing.Abs(spacing)
				if d.Cmp(spacing) <= 0 {
					return nil, invalid("number-spacing", "step")
				}
			}
		}
	}
	return g, nil
}
func (g *Grid) LastTick() uint64 { return g.last }
func (g *Grid) point(k uint64) *big.Rat {
	return new(big.Rat).Add(g.min, new(big.Rat).Mul(g.step, new(big.Rat).SetInt(new(big.Int).SetUint64(k))))
}
func (g *Grid) At(k uint64) (float64, error) {
	if k > g.last {
		return 0, invalid("number-range", "tick")
	}
	f, _ := g.point(k).Float64()
	return f, nil
}
func (g *Grid) Parse(raw string) (float64, error) {
	r, err := decimal(raw, "value")
	if err != nil {
		return 0, err
	}
	if r.Cmp(g.min) < 0 || r.Cmp(g.max) > 0 {
		return 0, invalid("number-range", "value")
	}
	q := new(big.Rat).Quo(new(big.Rat).Sub(r, g.min), g.step)
	if !q.IsInt() {
		return 0, invalid("number-grid", "value")
	}
	f, _ := r.Float64()
	return f, nil
}

// nearest clamps exact quotient and rounds a half tick upward. This is a gesture
// convention only; Parse and Tick never silently snap an invalid value.
func (g *Grid) nearest(value float64) (uint64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, invalid("number-finite", "value")
	}
	r := new(big.Rat).SetFloat64(value)
	q := new(big.Rat).Quo(new(big.Rat).Sub(r, g.min), g.step)
	if q.Sign() <= 0 {
		return 0, nil
	}
	last := new(big.Rat).SetInt(new(big.Int).SetUint64(g.last))
	if q.Cmp(last) >= 0 {
		return g.last, nil
	}
	num := new(big.Int).Add(new(big.Int).Lsh(new(big.Int).Set(q.Num()), 1), q.Denom())
	den := new(big.Int).Lsh(new(big.Int).Set(q.Denom()), 1)
	return num.Quo(num, den).Uint64(), nil
}
func (g *Grid) Nearest(value float64) (uint64, error) { return g.nearest(value) }
func (g *Grid) Tick(value float64) (uint64, error) {
	k, err := g.nearest(value)
	if err != nil {
		return 0, err
	}
	f, _ := g.At(k)
	if f != value {
		return 0, invalid("number-grid", "value")
	}
	return k, nil
}
func (g *Grid) ValidateSDL() error {
	if !g.integer {
		return invalid("number-integer", "constraints")
	}
	return nil
}

// Text returns the exact finite decimal of a legal tick, for editable number
// drafts. Formatting the rounded binary64 instead would lose grid provenance.
func (g *Grid) Text(k uint64) (string, error) {
	if k > g.last {
		return "", invalid("number-range", "tick")
	}
	r := g.point(k)
	den := new(big.Int).Set(r.Denom())
	two, five := big.NewInt(2), big.NewInt(5)
	a, b := 0, 0
	rem := new(big.Int)
	for rem.Mod(den, two).Sign() == 0 {
		den.Quo(den, two)
		a++
	}
	for rem.Mod(den, five).Sign() == 0 {
		den.Quo(den, five)
		b++
	}
	scale := a
	if b > scale {
		scale = b
	}
	return r.FloatString(scale), nil
}
