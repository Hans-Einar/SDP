package parser

type index struct {
	out, in           map[string][]string
	props             map[string]string
	placements, steps map[string][]Statement
}

func modelIndex(m *Model) *index {
	x := &index{out: map[string][]string{}, in: map[string][]string{}, props: map[string]string{}, placements: map[string][]Statement{}, steps: map[string][]Statement{}}
	for _, s := range m.Statements {
		switch s.Kind {
		case "Relation":
			a, b := key(s.Subject.Name, s.Verb), key(s.Object.Name, s.Verb)
			x.out[a] = append(x.out[a], s.Object.Name)
			x.in[b] = append(x.in[b], s.Subject.Name)
		case "PropertyAssignment":
			x.props[key(s.Subject.Name, s.Property)] = s.Value
		case "Placement":
			x.placements[s.Subject.Name] = append(x.placements[s.Subject.Name], s)
		case "Step":
			x.steps[s.Subject.Name] = append(x.steps[s.Subject.Name], s)
		}
	}
	return x
}
func key(a, b string) string                   { return a + "/" + b }
func (x *index) outgoing(a, b string) []string { return x.out[key(a, b)] }
func (x *index) incoming(a, b string) []string { return x.in[key(a, b)] }
func (x *index) property(a, b string) string   { return x.props[key(a, b)] }
