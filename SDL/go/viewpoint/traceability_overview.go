package viewpoint

// traceabilityOverview retains the complete checked facts and reported gaps.
// It is an inspection index; a union of unrelated diagrams would imply no
// additional semantics and obscure the source registry.
func (v *Views) traceabilityOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	return &out
}
