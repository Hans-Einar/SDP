package parser

import (
	"fmt"
	"sort"
	"strings"
)

type edge struct {
	from, to string
	span     Span
}

func Validate(m *Model) []Diagnostic {
	diagnostics := []Diagnostic{}
	symbols := Symbols(m)
	declared := map[string]bool{}
	names := []string{}
	for _, d := range m.Declarations {
		name := d.Name.Name
		if declared[name] {
			diagnostics = append(diagnostics, Diagnostic{"DUPLICATE_DECLARATION", name + " is already declared", d.Name.Span})
		} else {
			names = append(names, name)
		}
		declared[name] = true
	}
	add := func(code, message string, span Span) {
		diagnostics = append(diagnostics, Diagnostic{code, message, span})
	}
	check := func(id Identifier, expected []string, role string) bool {
		d, ok := symbols[id.Name]
		if !ok {
			add("UNDECLARED_NAME", id.Name+" is not declared", id.Span)
			return false
		}
		if !has(expected, d.Kind) && !(d.Kind == "container" && has(expected, "unit")) {
			add(role+"_TYPE_MISMATCH", fmt.Sprintf("%s %s expects %s, received %s", strings.ToLower(role), id.Name, strings.Join(expected, " or "), d.Kind), id.Span)
			return false
		}
		return true
	}
	facts := map[string]bool{}
	propertiesSeen := map[string]string{}
	owners, parents := map[string]map[string]bool{}, map[string]map[string]bool{}
	allocations := map[string]string{}
	edges := map[string][]edge{"contains": {}, "refines": {}, "depends-on": {}}
	for _, s := range m.Statements {
		fact := s.Sentence()
		if facts[fact] {
			add("DUPLICATE_FACT", "Repeated fact: "+fact, s.Span)
		}
		facts[fact] = true
		switch s.Kind {
		case "PropertyAssignment":
			expected := propertyKinds[s.Property]
			if expected == "" {
				expected = "functionality"
			}
			switch s.Property {
			case "completeness":
				expected = "contract scenario"
			case "message-kind":
				expected = "message"
			case "implementation-status":
				expected = "activity"
			}
			check(s.Subject, strings.Fields(expected), "PROPERTY")
			if !has(properties[s.Property], s.Value) {
				add("PROPERTY_TYPE_MISMATCH", s.Property+" does not accept "+s.Value, s.Span)
			}
			key := s.Subject.Name + "/" + s.Property
			if old, ok := propertiesSeen[key]; ok && old != s.Value {
				add("PROPERTY_CONFLICT", s.Subject.Name+" has conflicting "+s.Property+" values", s.Span)
			} else if !ok {
				propertiesSeen[key] = s.Value
			}
		case "Participation":
			check(s.Subject, []string{"unit"}, "SUBJECT")
			check(s.Channel, []string{"channel"}, "OBJECT")
			check(s.Message, []string{"message", "datagram"}, "OBJECT")
			check(s.Mode, []string{"mode"}, "QUALIFIER")
		case "Step":
			check(s.Subject, []string{"scenario"}, "SUBJECT")
			check(s.Message, []string{"message", "datagram"}, "OBJECT")
			check(s.Sender, []string{"unit"}, "OBJECT")
			check(s.Receiver, []string{"unit"}, "OBJECT")
			check(s.Channel, []string{"channel"}, "OBJECT")
			if s.Variant != nil {
				check(*s.Variant, []string{"variant"}, "QUALIFIER")
			}
		case "Projection":
			check(s.Subject, []string{"functionality"}, "SUBJECT")
			check(s.Dataset, []string{"dataset"}, "OBJECT")
			check(s.Datagram, []string{"datagram"}, "QUALIFIER")
		case "Placement":
			check(s.Subject, []string{"encoding"}, "SUBJECT")
			check(s.Field, []string{"field"}, "OBJECT")
		case "Dependency":
			check(s.Subject, []string{"capability"}, "SUBJECT")
			check(s.Interface, []string{"interface"}, "OBJECT")
			check(s.Mode, []string{"mode"}, "QUALIFIER")
		case "Allocation":
			a := check(s.Subject, []string{"functionality"}, "SUBJECT")
			b := check(s.Container, []string{"container"}, "OBJECT")
			c := check(s.Mode, []string{"mode"}, "QUALIFIER")
			if a && b && c {
				key := s.Subject.Name + "/" + s.Mode.Name
				if old, ok := allocations[key]; ok && old != s.Container.Name {
					add("ALLOCATION_CARDINALITY", s.Subject.Name+" has multiple Containers in mode "+s.Mode.Name, s.Span)
				} else if !ok {
					allocations[key] = s.Container.Name
				}
			}
		case "Relation":
			sig := signatures[s.Verb]
			a := check(s.Subject, sig.subject, "SUBJECT")
			b := check(s.Object, sig.object, "OBJECT")
			if !a || !b {
				continue
			}
			if _, ok := edges[s.Verb]; ok {
				edges[s.Verb] = append(edges[s.Verb], edge{s.Subject.Name, s.Object.Name, s.Span})
			}
			if s.Verb == "owns" || s.Verb == "contains" {
				mapping, code, role := owners, "OWNERSHIP_CARDINALITY", "owners"
				if s.Verb == "contains" {
					mapping = parents
					code = "CONTAINMENT_CARDINALITY"
					role = "parents"
				}
				if old := mapping[s.Object.Name]; len(old) > 0 && !old[s.Subject.Name] {
					add(code, s.Object.Name+" has multiple immediate "+role, s.Span)
				}
				if mapping[s.Object.Name] == nil {
					mapping[s.Object.Name] = map[string]bool{}
				}
				mapping[s.Object.Name][s.Subject.Name] = true
			}
		}
	}
	for _, name := range names {
		d := symbols[name]
		if (d.Kind == "functionality" || d.Kind == "database") && len(owners[name]) == 0 {
			add("OWNERSHIP_CARDINALITY", name+" has no valid immediate owner", d.Name.Span)
		}
	}
	for _, relation := range []string{"contains", "refines", "depends-on"} {
		diagnostics = append(diagnostics, cycles(edges[relation], relation)...)
	}
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, validateData(m, symbols, names)...)
	}
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, validateChannels(m, symbols, names)...)
	}
	return diagnostics
}
func cycles(edges []edge, relation string) []Diagnostic {
	graph := map[string][]edge{}
	for _, e := range edges {
		graph[e.from] = append(graph[e.from], e)
	}
	keys := []string{}
	for k := range graph {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	colors := map[string]int{}
	diagnostics := []Diagnostic{}
	type frame struct {
		name string
		next int
	}
	for _, root := range keys {
		if colors[root] != 0 {
			continue
		}
		colors[root] = 1
		stack := []frame{{root, 0}}
		for len(stack) > 0 {
			last := &stack[len(stack)-1]
			if last.next >= len(graph[last.name]) {
				colors[last.name] = 2
				stack = stack[:len(stack)-1]
				continue
			}
			e := graph[last.name][last.next]
			last.next++
			if colors[e.to] == 1 {
				diagnostics = append(diagnostics, Diagnostic{"STRUCTURE_CYCLE", fmt.Sprintf("%s edge %s -> %s closes a cycle", relation, e.from, e.to), e.span})
			} else if colors[e.to] == 0 {
				colors[e.to] = 1
				stack = append(stack, frame{e.to, 0})
			}
		}
	}
	return diagnostics
}
