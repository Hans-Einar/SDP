package viewpoint

import (
	"fmt"
	"strings"
)

// ChannelOverviewContext uses checked declarations/facts; overview indexes
// never create diagrams or infer ordering between separately declared Scenarios.
func (v *Views) ChannelOverviewContext(mode string) string {
	ref := func(name string) string {
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				return v.SpanReference(name, d.Span)
			}
		}
		return name
	}
	refs := func(names []string) string {
		out := []string{}
		for _, n := range unique(names) {
			out = append(out, ref(n))
		}
		return strings.Join(out, ", ")
	}
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nChannel topology is partitioned into one diagram per declared Channel, retaining every participant and binding. Repeated participants refer to the same source declarations. The diagrams show explicit participation in one Mode, not runtime delivery, transport or execution order. Sender arrows point from participant to Channel; receiver arrows point from Channel to participant. Each arrow groups all messages for that participant/channel/role; the register retains every binding and source proof. Modes have no implicit inheritance.\n\n")
	if mode == "" {
		b.WriteString("## Choose a mode\n\nSelect a named mode under VP08 in the navigation tree. No topology or combined scenario timeline is generated for this index.\n\n")
		if len(v.names("mode")) == 0 {
			b.WriteString("No Modes are declared.\n\n")
		} else {
			b.WriteString("| Mode | Participation bindings | Scenarios |\n| --- | --- | --- |\n")
			for _, m := range v.names("mode") {
				fmt.Fprintf(&b, "| %s | %d | %d |\n", ref(m), len(v.channelBindings(m)), len(v.channelScenarios(m)))
			}
			b.WriteString("\n")
		}
	} else {
		fmt.Fprintf(&b, "Selected mode: **%s**. Mode source: %s.\n\n", mode, ref(mode))
		bindings := v.channelBindings(mode)
		if len(bindings) == 0 {
			b.WriteString("No participation is declared for this mode. This does not prove independence or absence of runtime communication.\n\n")
		} else {
			b.WriteString("## Participation register\n\n| Participant | Channel | Role | Message | Source fact |\n| --- | --- | --- | --- | --- |\n")
			for _, f := range bindings {
				fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", ref(f.S("subject")), ref(f.S("channel")), f.S("role"), ref(f.S("message")), v.FactReference(f.S("id")))
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("## Channel register\n\nAll declared Channels are retained, including those without participants in the selected mode. Contracts and permitted messages are global model declarations, not evidence of delivery in any mode.\n\n")
	channels := v.names("channel")
	if len(channels) == 0 {
		b.WriteString("No Channels are declared.\n\n")
	} else {
		b.WriteString("| Channel | Contract | Permitted messages | Source facts |\n| --- | --- | --- | --- |\n")
		for _, ch := range channels {
			governing := eq(v.Relations["upholds"], "subject", ch)
			contracts := values(governing, "object")
			allowed := filter(v.Relations["permits"], func(f Fact) bool { return has(contracts, f.S("subject")) })
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", ref(ch), refs(contracts), refs(values(allowed, "object")), v.FactReferences(append(values(governing, "id"), values(allowed, "id")...)))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Scenario index\n\nOpen an individual Scenario in the navigation tree for its validated sequence. Source links below open declarations; steps are not concatenated across Scenarios. Completeness is a source assertion, not evidence of execution or test coverage. Preconditions and postconditions are not inferred.\n\n")
	scenarios := v.channelScenarios(mode)
	if len(scenarios) == 0 {
		b.WriteString("No Scenarios are declared in this scope.\n\n")
	} else {
		b.WriteString("| Scenario | Mode | Completeness | Steps | Source facts |\n| --- | --- | --- | --- | --- |\n")
		for _, name := range scenarios {
			fs := eq(v.Facts, "subject", name)
			fmt.Fprintf(&b, "| %s | %s | %s | %d | %s |\n", ref(name), refs(values(eq(fs, "verb", "runs-in"), "object")), or(v.prop(name, "completeness"), "Unspecified"), len(eq(fs, "verb", "step")), v.FactReferences(values(fs, "id")))
		}
		b.WriteString("\n")
	}
	return b.String()
}
