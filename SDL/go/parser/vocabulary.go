package parser

import "strings"

type signature struct{ subject, object []string }

func sig(a, b string) signature { return signature{strings.Fields(a), strings.Fields(b)} }

var kinds = strings.Fields("unit container functionality capability interface activity mode actor usecase feature dataset database datagram contract variant field encoding channel message scenario")
var signatures = map[string]signature{
	"contains": sig("unit", "unit"), "owns": sig("unit", "functionality database"), "realizes": sig("functionality", "capability"), "provides": sig("unit", "capability"), "consumes": sig("unit", "interface"), "refines": sig("activity", "activity"), "pursues": sig("actor", "usecase"), "supports": sig("feature", "usecase"), "contributes-to": sig("functionality", "feature usecase"), "addresses": sig("activity", "functionality"), "delivers": sig("activity", "feature"), "depends-on": sig("activity", "activity"), "illustrates": sig("scenario", "activity"),
	"upholds": sig("dataset datagram message channel", "contract"), "from": sig("datagram", "dataset"), "holds": sig("unit database", "dataset"), "defines": sig("contract", "variant"), "has-field": sig("contract variant", "field"), "encodes": sig("encoding", "variant"), "permits": sig("contract", "message datagram"), "replies-to": sig("message", "message"), "runs-in": sig("scenario", "mode"), "exercises": sig("scenario", "usecase"),
}
var properties = map[string][]string{
	"state-retention": {"stateful", "stateless"}, "repeatability": {"deterministic", "nondeterministic"}, "implementation-status": {"planned", "implemented", "verified"}, "completeness": {"open", "closed"}, "value-type": {"unsigned", "signed", "boolean", "text", "bytes", "decimal"}, "presence": {"required", "optional"}, "byte-order": {"big-endian", "little-endian"}, "bit-order": {"most-significant-first", "least-significant-first"}, "message-kind": {"request", "result", "event"},
}
var propertyKinds = map[string]string{"completeness": "contract", "value-type": "field", "presence": "field", "byte-order": "encoding", "bit-order": "encoding"}

func has(items []string, v string) bool {
	for _, x := range items {
		if x == v {
			return true
		}
	}
	return false
}
