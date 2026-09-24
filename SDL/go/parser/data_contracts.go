package parser

import "sort"

func validateData(m *Model, symbols map[string]Declaration, names []string) []Diagnostic {
	x := modelIndex(m)
	errors := []Diagnostic{}
	add := func(code, name, message string, span *Span) {
		s := symbols[name].Span
		if span != nil {
			s = *span
		}
		errors = append(errors, Diagnostic{code, name + ": " + message, s})
	}
	one := func(name, verb string, reverse bool) string {
		values := x.outgoing(name, verb)
		if reverse {
			values = x.incoming(name, verb)
		}
		if len(values) != 1 {
			add("DATA_CARDINALITY", name, "exactly one "+verb+" relation required", nil)
			return ""
		}
		return values[0]
	}
	for _, name := range names {
		kind := symbols[name].Kind
		if kind == "dataset" || kind == "datagram" {
			contract := one(name, "upholds", false)
			if kind == "dataset" {
				one(name, "holds", true)
			} else {
				one(name, "from", false)
			}
			if contract != "" && x.property(contract, "completeness") == "closed" {
				variants := x.outgoing(contract, "defines")
				if kind == "dataset" && (len(variants) > 0 || len(x.outgoing(contract, "has-field")) == 0) {
					add("CONTRACT_SHAPE", name, "closed Dataset requires record fields and no variants", nil)
				}
				if kind == "datagram" && len(variants) == 0 {
					add("CONTRACT_SHAPE", name, "closed Datagram requires declared variants", nil)
				}
			}
		}
		if kind == "variant" {
			parent := one(name, "defines", true)
			if parent != "" && x.property(parent, "completeness") == "closed" && len(x.outgoing(name, "has-field")) == 0 {
				add("CONTRACT_SHAPE", name, "closed variant requires payload fields", nil)
			}
		}
		if kind == "field" {
			one(name, "has-field", true)
		}
		for _, prop := range []string{"completeness", "value-type", "presence", "byte-order", "bit-order"} {
			if kind == propertyKinds[prop] && x.property(name, prop) == "" {
				add("MISSING_CONTRACT_PROPERTY", name, "missing "+prop, nil)
			}
		}
	}
	for _, s := range m.Statements {
		if s.Kind == "Projection" {
			sources := x.outgoing(s.Datagram.Name, "from")
			if len(sources) != 1 || sources[0] != s.Dataset.Name {
				add("PROJECTION_SOURCE_MISMATCH", s.Subject.Name, "projection disagrees with Datagram source", &s.Span)
			}
		}
	}
	for _, name := range names {
		if symbols[name].Kind != "encoding" {
			continue
		}
		variant := one(name, "encodes", false)
		parents := x.incoming(variant, "defines")
		if len(parents) != 1 {
			continue
		}
		contract := parents[0]
		if x.property(contract, "completeness") != "closed" {
			add("OPEN_ENCODING", name, "packet requires a closed contract", nil)
		}
		datagram := false
		for _, n := range x.incoming(contract, "upholds") {
			datagram = datagram || symbols[n].Kind == "datagram"
		}
		if !datagram {
			add("ENCODING_FAMILY", name, "variant must belong to a Datagram contract", nil)
		}
		expected := map[string]bool{}
		for _, f := range append(append([]string{}, x.outgoing(contract, "has-field")...), x.outgoing(variant, "has-field")...) {
			expected[f] = true
		}
		entries := append([]Statement{}, x.placements[name]...)
		sort.SliceStable(entries, func(i, j int) bool {
			if entries[i].Offset.Value != entries[j].Offset.Value {
				return entries[i].Offset.Value < entries[j].Offset.Value
			}
			return entries[i].Field.Name < entries[j].Field.Name
		})
		seen := map[string]bool{}
		mismatch := false
		for _, s := range entries {
			if seen[s.Field.Name] || !expected[s.Field.Name] {
				mismatch = true
			}
			seen[s.Field.Name] = true
		}
		if mismatch || len(seen) != len(expected) || len(entries) == 0 {
			add("ENCODING_FIELDS", name, "place every header/variant field exactly once", nil)
		}
		cursor := 0
		for _, s := range entries {
			field, width := s.Field.Name, s.Width.Value
			if s.Offset.Value != cursor || width < 1 {
				add("ENCODING_RANGE", name, "positive contiguous fields required; no overlap/gap", &s.Span)
			}
			cursor = s.Offset.Value + width
			if cursor > 65536 {
				add("LAYOUT_LIMIT", name, "packet exceeds 65536 bits", &s.Span)
			}
			typ := x.property(field, "value-type")
			valid := (typ == "unsigned" || typ == "signed") && width >= 1 && width <= 64 || typ == "boolean" && width == 1 || typ == "bytes" && width > 0 && width%8 == 0
			if !valid || x.property(field, "presence") != "required" {
				add("ENCODING_TYPE", name, "field needs a supported fixed-width required type", &s.Span)
			}
		}
	}
	return errors
}
