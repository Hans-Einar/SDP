# Exact bounded numeric grids

This standard-library-only package is shared by the SDUI frontend, runtime and SDL
bridge. It imports none of those packages or any GUI library. Callers retain source
spans; `Error.Code` and `Error.Argument` identify the failing numeric rule/argument.
Grid is immutable and exposes no mutable big.Int/big.Rat operands.

- `NewGrid(min,max,step)` accepts original decimal lexemes and validates constraints.
- `Parse(raw)` requires exact rational range/grid membership; it never snaps text.
- `Tick(number)` validates a typed binary64 by equality to the correctly rounded
  exact reconstructed point at its nearest legal tick.
- `At(tick)` reconstructs that exact point and rounds once; `Text(tick)` returns its
  exact finite decimal for editable drafts, avoiding rounded-float formatting drift.
- `LastTick()` returns floor((max-min)/step) as uint64.
- `Nearest(number)` explicitly clamps/snaps native gesture input. It is not an
  admission path for source, text Commit or programmatic Number values.
- `ValidateSDL()` checks exact safe53 integral constraints. `ExactInteger(raw)`
  checks decimal provenance; `Integer(number)` checks a typed finite/integral value.
  Both enforce ±9007199254740991 before conversion.

Decimal admission scans grammar and bounds before arbitrary-precision conversion:
at most 32768 total bytes and mantissa digits; effective exponent bounded to ±4096
for nonzero coefficients after trailing-zero cancellation. Explicit exponent scanning
saturates before multiplication. An all-zero coefficient becomes exact zero without
constructing an exponent power. NaN/infinity/overflow and non-source number syntax
reject. Exact rational operations preserve distinctions lost by binary64 rounding.

Only exact integral min/max/step within safe53 receive the no-tick-cap exception;
uint64 covers the full signed-safe53 interval. Every other Go grid has at most 2^26
intervals and step strictly greater than the largest adjacent binary64 spacing in
both directions at both rounded endpoints. Missing/nonfinite neighbors reject.
Equality is insufficient because ties can collapse different ticks. Reconstructed
points use bounded rationals directly, not expanded strings reparsed as source.

Tests cover hostile exponents/zero coefficients, exact .3/.1 membership, decimals
rounding onto a legal float, wide integer intervals, nonrepresentable constraints,
clamped native movement, exact display and concurrent immutable use. There is no
expression evaluator, locale parser, persistence or numeric-string SDL workaround.
