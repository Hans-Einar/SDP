package viewpoint

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var wordBreak = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func label(s string) string { return wordBreak.ReplaceAllString(s, "$1 $2") }
func (v *Views) channels() {
	r := v.Relations
	parts := eq(v.Facts, "verb", "uses")
	contracts := map[string]string{}
	for _, f := range r["upholds"] {
		contracts[f.S("subject")] = f.S("object")
	}
	modes := map[string]string{}
	for _, f := range r["runs-in"] {
		modes[f.S("subject")] = f.S("object")
	}
	for _, ch := range v.names("channel") {
		contract := contracts[ch]
		allowed := eq(r["permits"], "subject", contract)
		bindings := eq(parts, "channel", ch)
		for _, mode := range values(bindings, "mode") {
			for _, p := range allowed {
				msg := p.S("object")
				members := filter(bindings, func(f Fact) bool { return f.S("mode") == mode && f.S("message") == msg })
				senders := values(eq(members, "role", "sender"), "subject")
				receivers := values(eq(members, "role", "receiver"), "subject")
				proof := values(filter(r["upholds"], func(f Fact) bool { return f.S("subject") == ch || f.S("subject") == msg }), "id")
				proof = unique(append(append(proof, p.S("id")), values(members, "id")...))
				v.MessageSets = append(v.MessageSets, Fact{"channel": ch, "mode": mode, "message": msg, "contract": contract, "payload_contract": contracts[msg], "variants": values(eq(r["defines"], "subject", contracts[msg]), "object"), "senders": senders, "receivers": receivers, "source_facts": proof})
				if len(senders) == 0 || len(receivers) == 0 {
					v.Gaps = append(v.Gaps, Fact{"viewpoint": "VP08", "model_id": ch, "mode": mode, "code": "INCOMPLETE_PARTICIPATION", "message": "Ufullstendig deltakelse for " + msg + "."})
				}
			}
		}
		if len(bindings) == 0 {
			v.Gaps = append(v.Gaps, Fact{"viewpoint": "VP08", "model_id": ch, "code": "NO_PARTICIPANTS", "message": "No explicit participation."})
		}
	}
	for _, scenario := range v.names("scenario") {
		mode := modes[scenario]
		steps := filter(v.Facts, func(f Fact) bool { return f.S("verb") == "step" && f.S("subject") == scenario })
		sort.Slice(steps, func(i, j int) bool { return steps[i].N("ordinal") < steps[j].N("ordinal") })
		order := []string{}
		for _, s := range steps {
			for _, k := range []string{"sender", "receiver"} {
				if n := s.S(k); !has(order, n) {
					order = append(order, n)
				}
			}
		}
		lines := []string{"sequenceDiagram"}
		for _, n := range order {
			lines = append(lines, "    participant n_"+n+" as "+label(n))
		}
		source := values(eq(v.Facts, "subject", scenario), "id")
		elements := []Fact{}
		for _, s := range steps {
			ch, msg := s.S("channel"), s.S("message")
			bindings := filter(parts, func(f Fact) bool {
				return f.S("channel") == ch && f.S("message") == msg && f.S("mode") == mode && ((f.S("subject") == s.S("sender") && f.S("role") == "sender") || (f.S("subject") == s.S("receiver") && f.S("role") == "receiver"))
			})
			permits := filter(r["permits"], func(f Fact) bool { return f.S("subject") == contracts[ch] && f.S("object") == msg })
			governing := filter(r["upholds"], func(f Fact) bool { return f.S("subject") == ch || f.S("subject") == msg })
			types := filter(v.Facts, func(f Fact) bool {
				return f.S("subject") == msg && (f.S("verb") == "has" || f.S("verb") == "replies-to")
			})
			proof := unique(append(values(join(bindings, permits, governing, types), "id"), s.S("id")))
			source = append(source, proof...)
			arrow := "->>"
			if v.prop(msg, "message-kind") == "result" {
				arrow = "-->>"
			}
			caption := fmt.Sprintf("%d: %s", s.N("ordinal"), label(msg))
			if s.S("variant") != "" {
				caption += " / " + label(s.S("variant"))
			}
			caption += " (" + label(ch) + ")"
			if s.N("reply_to") != 0 {
				caption += fmt.Sprintf(" reply-to %d", s.N("reply_to"))
			}
			lines = append(lines, "    n_"+s.S("sender")+arrow+"n_"+s.S("receiver")+": "+caption)
			elements = append(elements, Fact{"ordinal": s.N("ordinal"), "sender": s.S("sender"), "receiver": s.S("receiver"), "message": msg, "variant": s["variant"], "channel": ch, "reply_to": s["reply_to"], "fact": s.S("id"), "proof": proof})
		}
		v.Diagrams = append(v.Diagrams, Diagram{"VP08-" + scenario, "Scenario: " + scenario + " — mode " + mode, "sequence", v.nodes(order), []Edge{}, unique(source), elements, strings.Join(lines, "\n") + "\n"})
	}
	if len(modes) == 0 {
		v.Gaps = append(v.Gaps, Fact{"viewpoint": "VP08", "model_id": "—", "code": "NO_SCENARIOS", "message": "No declared scenario steps; ordering is not inferred."})
	}
}
