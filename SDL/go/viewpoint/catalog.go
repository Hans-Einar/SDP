package viewpoint

type Spec struct {
	ID, Title, Note string
	Levels          []string
}

var Catalog = []Spec{
	{"VP01", "Use cases and traceability", "pursues, supports and contributes-to; modeled scope without an invented System boundary.", []string{"A0", "A1"}},
	{"VP02", "Architecture and logical decomposition", "Container/Unit and contains; library structure is not deployment allocation.", []string{"A2", "A3"}},
	{"VP03", "Responsibilities and capabilities across the architecture", "owns, realizes and provides. Capability is not Feature.", []string{"A1", "A2", "A3"}},
	{"VP04", "Interfaces and collaboration", "consumes records use; no provider, Channel or executable message flow is inferred.", []string{"A2", "A3"}},
	{"VP05", "Dependencies by mode", "requires in mode; modes have no implicit inheritance.", []string{"A2", "A3"}},
	{"VP06", "Activities and delivery plan", "refines, addresses, delivers and depends-on; status is an explicit source claim.", []string{}},
	{"VP07", "Features across the architecture", "contributes-to, owns and explicit allocated-to per mode. Unspecified allocation is reported as a gap.", []string{"A1", "A2", "A3"}},
	{"VP08", "Channel contracts and sequences", "Explicit scenario steps validated against permits, participation, mode and request/result correlation.", []string{"A1", "A4"}},
	{"VP09", "Dataset, Datagram and persistent Database", "Explicit holders, sources, contracts, variants, fields and projections.", []string{"A1", "A4"}},
	{"VP10", "Datagram encoding and packets", "Only closed contracts with validated Encoding and explicit bit positions.", []string{"A4"}},
	{"VP11", "Properties, traceability and model gaps", "Declarations and all facts with source positions; support boundaries are retained.", []string{"A0", "A1", "A2", "A3", "A4", "A5"}},
}

func SpecFor(id string) (Spec, bool) {
	for _, s := range Catalog {
		if s.ID == id {
			return s, true
		}
	}
	return Spec{}, false
}
