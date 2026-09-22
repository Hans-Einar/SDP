package parser

import "fmt"

func ValidateClasses(m *ClassModel) error {
	if m == nil {
		return Diagnostic{Code: "CLASS_MODEL", Message: "Missing model"}
	}
	bad := func(code, msg string, s Span) error { return Diagnostic{code, msg, s} }
	if m.Header.Language != "class-core" || m.Header.Version != "0.1" {
		return bad("UNSUPPORTED_VERSION", "Expected class-core 0.1", m.Header.Span)
	}
	kinds := map[string]string{}
	classes, assocs := 0, 0
	for _, d := range m.Declarations {
		if !namePattern.MatchString(d.Name.Name) || kinds[d.Name.Name] != "" {
			return bad("CLASS_DECLARATION", "Invalid or repeated name", d.Span)
		}
		kinds[d.Name.Name] = d.Kind
		switch d.Kind {
		case "class":
			classes++
		case "association":
			assocs++
		default:
			return bad("CLASS_DECLARATION", "Unknown declaration kind", d.Span)
		}
	}
	if classes > 128 || assocs > 256 {
		return bad("CLASS_LIMIT", "128 classes / 256 associations", m.Span)
	}
	seen := map[string]bool{}
	counts := map[string]int{}
	links := map[string]ClassStatement{}
	ownership := map[string]string{}
	roles := map[string]bool{}
	for _, s := range m.Statements {
		subject := s.Subject.Name
		k := kinds[subject]
		if k == "" {
			return bad("UNDECLARED_REFERENCE", subject, s.Subject.Span)
		}
		key := subject + "/" + s.Verb
		switch s.Verb {
		case "attribute", "operation":
			if k != "class" {
				return bad("CLASS_SUBJECT", "Member requires class", s.Span)
			}
			if !namePattern.MatchString(s.Name.Name) {
				return bad("CLASS_MEMBER", "Invalid member name", s.Name.Span)
			}
			key += "/" + s.Name.Name
			if !has([]string{"integer", "text", "boolean"}, s.Type) && kinds[s.Type] != "class" {
				return bad("CLASS_TYPE", s.Type, s.Span)
			}
			counts[subject]++
			if counts[subject] > 64 {
				return bad("CLASS_LIMIT", "64 members per class", s.Span)
			}
		case "links":
			if k != "association" {
				return bad("CLASS_SUBJECT", "links requires association", s.Span)
			}
			for _, end := range []ClassEnd{s.First, s.Second} {
				if kinds[end.Class.Name] != "class" {
					return bad("CLASS_END", "End requires declared class", end.Span)
				}
				if !namePattern.MatchString(end.Role.Name) || roles[end.Class.Name+"/"+end.Role.Name] {
					return bad("CLASS_ROLE", "Invalid or repeated role", end.Role.Span)
				}
				roles[end.Class.Name+"/"+end.Role.Name] = true
				if end.Minimum.Value < 0 || end.Minimum.Value > 65536 || end.Maximum != nil && (end.Maximum.Value < end.Minimum.Value || end.Maximum.Value > 65536) {
					return bad("CLASS_MULTIPLICITY", "Invalid multiplicity interval", end.Span)
				}
			}
			if s.First.Role.Name == s.Second.Role.Name {
				return bad("CLASS_ROLE", "Distinct end roles required", s.Span)
			}
			links[subject] = s
		case "has":
			if k != "association" || !has([]string{"none", "aggregation", "composition"}, s.Ownership) {
				return bad("CLASS_OWNERSHIP", "Explicit none/aggregation/composition required", s.Span)
			}
			ownership[subject] = s.Ownership
		default:
			return bad("UNSUPPORTED_SYNTAX", "Unknown class statement", s.Span)
		}
		if seen[key] {
			return bad("CLASS_CARDINALITY", "Repeated member or association fact", s.Span)
		}
		seen[key] = true
	}
	children := map[string][]string{}
	for _, d := range m.Declarations {
		if d.Kind != "association" {
			continue
		}
		s, ok := links[d.Name.Name]
		if !ok || ownership[d.Name.Name] == "" {
			return bad("CLASS_CARDINALITY", "Association requires one links and one ownership", d.Span)
		}
		if ownership[d.Name.Name] == "composition" {
			if s.First.Maximum == nil || s.First.Maximum.Value != 1 {
				return bad("CLASS_COMPOSITION", "Part must have at most one whole", s.First.Span)
			}
			children[s.First.Class.Name] = append(children[s.First.Class.Name], s.Second.Class.Name)
		}
	}
	var visit func(string) error
	active, done := map[string]bool{}, map[string]bool{}
	visit = func(n string) error {
		if active[n] {
			return bad("CLASS_COMPOSITION_CYCLE", fmt.Sprintf("Composition cycle at %s", n), m.Span)
		}
		if done[n] {
			return nil
		}
		active[n] = true
		for _, c := range children[n] {
			if e := visit(c); e != nil {
				return e
			}
		}
		active[n] = false
		done[n] = true
		return nil
	}
	for _, d := range m.Declarations {
		if d.Kind == "class" {
			if e := visit(d.Name.Name); e != nil {
				return e
			}
		}
	}
	return nil
}
