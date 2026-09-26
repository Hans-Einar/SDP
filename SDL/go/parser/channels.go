package parser

import (
	"sort"
	"strings"
)

func validateChannels(m *Model, symbols map[string]Declaration, names []string) []Diagnostic {
	x := modelIndex(m)
	errors := []Diagnostic{}
	add := func(code, name, message string, span *Span) {
		s := symbols[name].Span
		if span != nil {
			s = *span
		}
		errors = append(errors, Diagnostic{code, name + ": " + message, s})
	}
	one := func(name, verb string) string {
		values := x.outgoing(name, verb)
		if len(values) != 1 {
			add("CHANNEL_CARDINALITY", name, "exactly one "+verb+" relation required", nil)
			return ""
		}
		return values[0]
	}
	participation := map[string]bool{}
	for _, s := range m.Statements {
		if s.Kind == "Participation" {
			participation[participationKey(s.Subject.Name, s.Channel.Name, s.Role, s.Message.Name, s.Mode.Name)] = true
		}
	}
	contracts, modes := map[string]string{}, map[string]string{}
	for _, name := range names {
		kind := symbols[name].Kind
		switch kind {
		case "message", "channel":
			contract := one(name, "upholds")
			if contract != "" {
				contracts[name] = contract
				fields, variants, permits := x.outgoing(contract, "has-field"), x.outgoing(contract, "defines"), x.outgoing(contract, "permits")
				closed := x.property(contract, "completeness") == "closed"
				if kind == "channel" && (len(fields) > 0 || len(variants) > 0 || closed && len(permits) == 0) {
					add("CHANNEL_CONTRACT_SHAPE", name, "Channel contract requires permits, without payload fields/variants", nil)
				}
				if kind == "message" && (len(variants) > 0 || len(permits) > 0 || closed && len(fields) == 0) {
					add("MESSAGE_CONTRACT_SHAPE", name, "Message contract requires record fields, without permits/variants", nil)
				}
			}
			if kind == "message" {
				typ := x.property(name, "message-kind")
				if typ == "" {
					add("MISSING_CONTRACT_PROPERTY", name, "missing message-kind", nil)
				}
				replies := x.outgoing(name, "replies-to")
				if typ == "result" {
					if len(replies) != 1 || x.property(replies[0], "message-kind") != "request" {
						add("REPLY_TYPE", name, "result requires exactly one request message type", nil)
					}
				} else if len(replies) > 0 {
					add("REPLY_TYPE", name, "only result types can declare replies-to", nil)
				}
			}
		case "scenario":
			modes[name] = one(name, "runs-in")
			if x.property(name, "completeness") == "" {
				add("MISSING_CONTRACT_PROPERTY", name, "missing completeness", nil)
			}
			if len(x.steps[name]) == 0 {
				add("EMPTY_SCENARIO", name, "scenario requires explicit steps", nil)
			}
		}
	}
	checked := map[string]bool{}
	for _, s := range m.Statements {
		if s.Kind == "Relation" && s.Verb == "permits" && !checked[s.Subject.Name] {
			name := s.Subject.Name
			checked[name] = true
			if len(x.outgoing(name, "has-field")) > 0 || len(x.outgoing(name, "defines")) > 0 {
				add("CHANNEL_CONTRACT_SHAPE", name, "permissions cannot be mixed with payload structure", nil)
			}
		}
		if s.Kind == "Participation" {
			contract := contracts[s.Channel.Name]
			if contract != "" && !has(x.outgoing(contract, "permits"), s.Message.Name) {
				add("MESSAGE_NOT_PERMITTED", s.Subject.Name, "Channel does not permit message", &s.Span)
			}
		}
	}
	// Scenario order follows first step occurrence, matching source-ordered diagnostics.
	scenarios := []string{}
	seen := map[string]bool{}
	for _, s := range m.Statements {
		if s.Kind == "Step" && !seen[s.Subject.Name] {
			scenarios = append(scenarios, s.Subject.Name)
			seen[s.Subject.Name] = true
		}
	}
	for _, scenario := range scenarios {
		ordered := append([]Statement{}, x.steps[scenario]...)
		sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Ordinal.Value < ordered[j].Ordinal.Value })
		orderOK := true
		for i, s := range ordered {
			orderOK = orderOK && s.Ordinal.Value == i+1
		}
		if !orderOK {
			add("STEP_ORDER", scenario, "step indexes must be unique and contiguous from 1", nil)
		}
		prior := map[int]Statement{}
		answered := map[int]bool{}
		for _, s := range ordered {
			message, channel := s.Message.Name, s.Channel.Name
			channelContract := contracts[channel]
			payload := x.outgoing(message, "upholds")
			if channelContract != "" && !has(x.outgoing(channelContract, "permits"), message) {
				add("MESSAGE_NOT_PERMITTED", scenario, "step is outside Channel contract", &s.Span)
			}
			allContracts := append([]string{}, payload...)
			if channelContract != "" {
				allContracts = append([]string{channelContract}, allContracts...)
			}
			for _, contract := range allContracts {
				if x.property(contract, "completeness") != "closed" {
					add("OPEN_SCENARIO_CONTRACT", scenario, "steps require closed contracts", &s.Span)
				}
			}
			for _, p := range []struct{ name, role string }{{s.Sender.Name, "sender"}, {s.Receiver.Name, "receiver"}} {
				if !participation[participationKey(p.name, channel, p.role, message, modes[scenario])] {
					add("PARTICIPATION_MISMATCH", scenario, "missing "+p.role+" in scenario Mode", &s.Span)
				}
			}
			if symbols[message].Kind == "datagram" {
				variants := []string{}
				for _, contract := range payload {
					variants = append(variants, x.outgoing(contract, "defines")...)
				}
				if s.Variant == nil || !has(variants, s.Variant.Name) {
					add("VARIANT_MISMATCH", scenario, "Datagram needs its own contracted variant", &s.Span)
				}
			} else if s.Variant != nil {
				add("VARIANT_MISMATCH", scenario, "record Message cannot select variant", &s.Span)
			}
			kind := x.property(message, "message-kind")
			if kind == "result" {
				target := -1
				if s.ReplyTo != nil {
					target = s.ReplyTo.Value
				}
				request, ok := prior[target]
				replies := x.outgoing(message, "replies-to")
				valid := ok && !answered[target] && x.property(request.Message.Name, "message-kind") == "request" && len(replies) == 1 && replies[0] == request.Message.Name && request.Channel.Name == channel && request.Sender.Name == s.Receiver.Name && request.Receiver.Name == s.Sender.Name
				if !valid {
					add("CORRELATION_MISMATCH", scenario, "result must reverse one earlier matching request", &s.Span)
				} else {
					answered[target] = true
				}
			} else if s.ReplyTo != nil {
				add("CORRELATION_MISMATCH", scenario, "only results can use reply-to", &s.Span)
			}
			prior[s.Ordinal.Value] = s
		}
		if x.property(scenario, "completeness") == "closed" {
			pending := false
			for ordinal, s := range prior {
				pending = pending || x.property(s.Message.Name, "message-kind") == "request" && !answered[ordinal]
			}
			if pending {
				add("INCOMPLETE_SCENARIO", scenario, "closed scenario has unanswered requests", nil)
			}
		}
	}
	return errors
}
func participationKey(unit, channel, role, message, mode string) string {
	return strings.Join([]string{unit, channel, role, message, mode}, "/")
}
