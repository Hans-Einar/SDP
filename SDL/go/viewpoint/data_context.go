package viewpoint

import (
	"fmt"
	"strings"
)

func (v *Views) DataOverviewContext() string {
	ref := func(name string) string {
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				return v.SpanReference(name, d.Span)
			}
		}
		return name
	}
	refs := func(names []string) string {
		labels := []string{}
		for _, n := range unique(names) {
			labels = append(labels, ref(n))
		}
		if len(labels) == 0 {
			return "Unspecified"
		}
		return strings.Join(labels, ", ")
	}
	proof := func(fs []Fact) string {
		if len(fs) == 0 {
			return "No facts declared"
		}
		return v.FactReferences(values(fs, "id"))
	}
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nDataset neighborhoods retain explicit from, holds, owns and upholds relations. Database means a persistent data source, not SQL tables, keys, transactions, a storage engine or deployment. Unit-held data has no declared persistence. These are definitions, not runtime instances. Contract completeness is a source assertion, not execution evidence. Optional fields do not imply null, empty, false or zero. Projections do not specify executable transformations. Encoding and bit positions belong to VP10.\n\n")
	b.WriteString("## Data definitions\n\n")
	for _, kind := range []string{"dataset", "datagram"} {
		names := v.names(kind)
		if len(names) == 0 {
			fmt.Fprintf(&b, "No %s declarations.\n\n", kind)
		}
		for _, n := range names {
			contracts := eq(v.Relations["upholds"], "subject", n)
			if kind == "dataset" {
				holders := eq(v.Relations["holds"], "object", n)
				persistence := "Unspecified"
				if len(holders) > 0 {
					persistence = "Unit-held; persistence unspecified"
					if v.Kinds[holders[0].S("subject")] == "database" {
						persistence = "Database-backed persistence"
					}
				}
				fmt.Fprintf(&b, "- %s (Dataset): holder %s; contract %s; %s. Source facts: %s.\n", ref(n), refs(values(holders, "subject")), refs(values(contracts, "object")), persistence, proof(join(holders, contracts)))
			} else {
				origin := eq(v.Relations["from"], "subject", n)
				fmt.Fprintf(&b, "- %s (Datagram): logical source Dataset %s; contract %s. This does not inherit the source's holder or persistence. Source facts: %s.\n", ref(n), refs(values(origin, "object")), refs(values(contracts, "object")), proof(join(origin, contracts)))
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("## Database register\n\n")
	if len(v.names("database")) == 0 {
		b.WriteString("No Databases are declared; no persistent store is inferred.\n\n")
	} else {
		b.WriteString("| Database | Responsible owner | Held Datasets | Source facts |\n| --- | --- | --- | --- |\n")
		for _, n := range v.names("database") {
			owners := eq(v.Relations["owns"], "object", n)
			held := eq(v.Relations["holds"], "subject", n)
			heldText := refs(values(held, "object"))
			if len(held) == 0 {
				heldText = "No Datasets declared"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", ref(n), refs(values(owners, "subject")), heldText, proof(join(owners, held)))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Contract register\n\nAll Contracts are included, even unused/open definitions. Open a Contract in the navigation tree for its structure diagram. Counts describe explicit shape membership (fields / variants / permitted messages), not runtime coverage.\n\n")
	if len(v.names("contract")) == 0 {
		b.WriteString("No Contracts are declared.\n\n")
	} else {
		b.WriteString("| Contract | Completeness | Fields / variants / messages | Source fact |\n| --- | --- | --- | --- |\n")
		for _, n := range v.names("contract") {
			props := filter(v.Facts, func(f Fact) bool { return f.S("subject") == n && f.S("property") == "completeness" })
			fmt.Fprintf(&b, "| %s | %s | %d / %d / %d | %s |\n", ref(n), or(v.prop(n, "completeness"), "Unspecified"), len(eq(v.Relations["has-field"], "subject", n)), len(eq(v.Relations["defines"], "subject", n)), len(eq(v.Relations["permits"], "subject", n)), proof(props))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Variant register\n\n")
	if len(v.names("variant")) == 0 {
		b.WriteString("No Variants are declared.\n\n")
	} else {
		b.WriteString("| Variant | Contract | Source fact |\n| --- | --- | --- |\n")
		for _, n := range v.names("variant") {
			fs := eq(v.Relations["defines"], "object", n)
			fmt.Fprintf(&b, "| %s | %s | %s |\n", ref(n), refs(values(fs, "subject")), proof(fs))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Field register\n\n")
	if len(v.names("field")) == 0 {
		b.WriteString("No Fields are declared.\n\n")
	} else {
		b.WriteString("| Field | Shape owner | Value type | Presence | Source facts |\n| --- | --- | --- | --- | --- |\n")
		for _, n := range v.names("field") {
			owners := eq(v.Relations["has-field"], "object", n)
			props := filter(v.Facts, func(f Fact) bool { return f.S("subject") == n && f.S("verb") == "has" })
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", ref(n), refs(values(owners, "subject")), or(v.prop(n, "value-type"), "Unspecified"), or(v.prop(n, "presence"), "Unspecified"), proof(join(owners, props)))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Shape membership\n\n")
	members := filter(v.Facts, func(f Fact) bool { return has([]string{"defines", "has-field", "permits"}, f.S("verb")) })
	if len(members) == 0 {
		b.WriteString("No shape membership facts are declared.\n\n")
	} else {
		b.WriteString("| Shape | Relation | Member | Source fact |\n| --- | --- | --- | --- |\n")
		for _, f := range members {
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", ref(f.S("subject")), f.S("verb"), ref(f.S("object")), v.FactReference(f.S("id")))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Projection responsibilities\n\nProjects relates Functionality, Dataset and Datagram together; it does not imply a lossless conversion or a runtime call.\n\n")
	projections := eq(v.Facts, "verb", "projects")
	if len(projections) == 0 {
		b.WriteString("No projections are declared.\n\n")
	} else {
		b.WriteString("| Functionality | Dataset | Datagram | Responsible owner | Source facts |\n| --- | --- | --- | --- | --- |\n")
		for _, f := range projections {
			owners := eq(v.Relations["owns"], "object", f.S("subject"))
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", ref(f.S("subject")), ref(f.S("dataset")), ref(f.S("datagram")), refs(values(owners, "subject")), proof(join([]Fact{f}, owners)))
		}
		b.WriteString("\n")
	}
	return b.String()
}
