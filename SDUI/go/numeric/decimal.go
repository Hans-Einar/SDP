// Package numeric implements bounded exact decimal grids, without UI or parser dependencies.
package numeric

import (
	"math"
	"math/big"
	"strings"
)

const SafeInteger int64 = 9007199254740991
const maxText = 32768

type Error struct{ Code, Argument string }

func (e *Error) Error() string       { return e.Code + ": " + e.Argument }
func invalid(code, arg string) error { return &Error{code, arg} }

// decimal scans bounds before allocating arbitrary precision coefficients/powers.
func decimal(raw, arg string) (*big.Rat, error) {
	if len(raw) == 0 || len(raw) > maxText {
		return nil, invalid("number-length", arg)
	}
	i := 0
	negative := false
	if raw[0] == '-' {
		negative = true
		i++
	}
	start := i
	if i == len(raw) {
		return nil, invalid("number-syntax", arg)
	}
	if raw[i] == '0' {
		i++
		if i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			return nil, invalid("number-syntax", arg)
		}
	} else {
		if raw[i] < '1' || raw[i] > '9' {
			return nil, invalid("number-syntax", arg)
		}
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
	}
	digits := raw[start:i]
	frac := 0
	if i < len(raw) && raw[i] == '.' {
		i++
		a := i
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
		if a == i {
			return nil, invalid("number-syntax", arg)
		}
		frac = i - a
		digits += raw[a:i]
	}
	if len(digits) > 32768 {
		return nil, invalid("number-digits", arg)
	}
	trimmed := strings.TrimLeft(digits, "0")
	trailing := len(trimmed) - len(strings.TrimRight(trimmed, "0"))
	coefficient := strings.TrimRight(trimmed, "0")
	// Saturate before multiply; the maximum allowed explicit exponent can only
	// differ from the effective bound by the already bounded fractional/trailing counts.
	limit := 4096 + frac + trailing
	exp := 0
	sign := 1
	if i < len(raw) && (raw[i] == 'e' || raw[i] == 'E') {
		i++
		if i < len(raw) && (raw[i] == '+' || raw[i] == '-') {
			if raw[i] == '-' {
				sign = -1
			}
			i++
		}
		a := i
		for i < len(raw) && raw[i] >= '0' && raw[i] <= '9' {
			d := int(raw[i] - '0')
			if exp > (limit-d)/10 {
				exp = limit + 1
			} else {
				exp = exp*10 + d
			}
			i++
		}
		if a == i {
			return nil, invalid("number-syntax", arg)
		}
	}
	if i != len(raw) {
		return nil, invalid("number-syntax", arg)
	}
	if coefficient == "" {
		return new(big.Rat), nil
	}
	effective := sign*exp - frac + trailing
	if exp > limit || effective < -4096 || effective > 4096 {
		return nil, invalid("number-exponent", arg)
	}
	c, ok := new(big.Int).SetString(coefficient, 10)
	if !ok {
		return nil, invalid("number-syntax", arg)
	}
	if negative {
		c.Neg(c)
	}
	magnitude := effective
	if magnitude < 0 {
		magnitude = -magnitude
	}
	power := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(magnitude)), nil)
	r := new(big.Rat).SetInt(c)
	if effective >= 0 {
		r.Mul(r, new(big.Rat).SetInt(power))
	} else {
		r.Quo(r, new(big.Rat).SetInt(power))
	}
	f, _ := r.Float64()
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return nil, invalid("number-finite", arg)
	}
	return r, nil
}
func safe(r *big.Rat) bool {
	return r.IsInt() && r.Num().IsInt64() && r.Num().Int64() >= -SafeInteger && r.Num().Int64() <= SafeInteger
}
func ExactInteger(raw string) (int64, error) {
	r, err := decimal(raw, "value")
	if err != nil {
		return 0, err
	}
	if !safe(r) {
		return 0, invalid("number-integer", "value")
	}
	return r.Num().Int64(), nil
}
func Integer(value float64) (int64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || math.Trunc(value) != value || value < -float64(SafeInteger) || value > float64(SafeInteger) {
		return 0, invalid("number-integer", "value")
	}
	return int64(value), nil
}
