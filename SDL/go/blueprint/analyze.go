// Package blueprint compares checked SDL source graphs. It performs no I/O,
// executes no task, and makes no claim about runtime or code conformance.
package blueprint

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/sourcegraph"
)

const Version = "sdl-blueprint-analysis/1"
const PolicyVersion = "structural-closure/1"

type Input struct {
	Snapshot         *sourcegraph.Snapshot
	ExpectedRevision string
}
type Rule struct {
	ID     string
	Source string
	// Fact is the canonical SDL sentence. NOW/TARGET specify required presence.
	Fact   string
	Now    bool
	Target bool
}

// Permission authorizes exactly one added/removed declaration or canonical fact.
// No wildcard or implicit permission follows from context inclusion.
type Protection struct{ ID, Source, Element, Scope string }

type Permission struct{ Element, Fact, Change string }
type Coverage struct {
	Profile, Policy                                     string
	ModeledClosureComplete                              bool
	CodeConformance, RuntimeDependencies, SDUISemantics string
}
type Task struct {
	ID                string
	Intent            string
	Context           []string
	Exclude           []string
	Rules             []Rule
	AllowedChanges    []Permission
	AllowedModelPaths []string
	AllowedCodePaths  []string
	Protect           []Protection
}
type Limits struct{ Nodes, Facts, Bytes int }

func DefaultLimits() Limits { return Limits{2000, 10000, 8 << 20} }

type Origin struct{ Now, Target []parser.Span }
type Node struct {
	Key     string
	Name    string
	Kind    string
	Change  string
	Origins Origin
	Reasons []string
}
type Fact struct {
	Key       string
	Statement string
	Endpoints []string
	Change    string
	Origins   Origin
	Semantic  parser.Statement
}
type Boundary struct{ Key, Reason string }
type Unknown struct{ Side, Code, Subject, Detail string }
type Obligation struct {
	ID, Result                string
	NowPresent, TargetPresent bool
}
type Analysis struct {
	Schema, Policy, Identity, System string
	NowRevision, TargetRevision      string
	Nodes                            []Node
	Facts                            []Fact
	Frontier                         []Boundary
	Excluded                         []Boundary
	Unknowns                         []Unknown
	Obligations                      []Obligation
	// Structural constraints are distinct from task authorization and code evidence.
	ConstraintsPass bool
	Coverage        Coverage
}
type Error struct{ Code, Detail string }

func (e *Error) Error() string       { return e.Code + ": " + e.Detail }
func fail(code, detail string) error { return &Error{code, detail} }
func digest(v any) string            { b, _ := json.Marshal(v); return fmt.Sprintf("%x", sha256.Sum256(b)) }
func keys[V any](m map[string]V) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return s
}
func unique(s []string) []string {
	m := map[string]bool{}
	for _, v := range s {
		m[v] = true
	}
	return keys(m)
}

// Every currently accepted statement family and relation has an explicit policy.
// All propagate conservatively in both directions, except expansion through
// system/mode context labels. Unknown future language constructs fail closed.
var relations = set("contains", "owns", "realizes", "provides", "consumes", "refines", "pursues", "supports", "contributes-to", "addresses", "delivers", "depends-on", "illustrates", "upholds", "from", "holds", "defines", "has-field", "encodes", "permits", "replies-to", "runs-in", "exercises")

func set(s ...string) map[string]bool {
	m := map[string]bool{}
	for _, v := range s {
		m[v] = true
	}
	return m
}
func endpoints(s parser.Statement) ([]string, error) {
	var ids []parser.Identifier
	switch s.Kind {
	case "Relation":
		if !relations[s.Verb] {
			return nil, fail("UNSUPPORTED_RELATION", s.Verb)
		}
		ids = []parser.Identifier{s.Subject, s.Object}
	case "Dependency":
		ids = []parser.Identifier{s.Subject, s.Interface, s.Mode}
	case "Allocation":
		ids = []parser.Identifier{s.Subject, s.Container, s.Mode}
	case "PropertyAssignment":
		ids = []parser.Identifier{s.Subject}
	case "Projection":
		ids = []parser.Identifier{s.Subject, s.Dataset, s.Datagram}
	case "Placement":
		ids = []parser.Identifier{s.Subject, s.Field}
	case "Participation":
		ids = []parser.Identifier{s.Subject, s.Channel, s.Message, s.Mode}
	case "Step":
		ids = []parser.Identifier{s.Subject, s.Message, s.Sender, s.Receiver, s.Channel}
		if s.Variant != nil {
			ids = append(ids, *s.Variant)
		}
	default:
		return nil, fail("UNSUPPORTED_STATEMENT", s.Kind)
	}
	out := []string{}
	for _, id := range ids {
		if id.Name != "" {
			out = append(out, id.Name)
		}
	}
	return unique(out), nil
}

func semantic(s parser.Statement) parser.Statement {
	s.Span = parser.Span{}
	s.Path = ""
	s.PathSpan = parser.Span{}
	ids := []*parser.Identifier{&s.Subject, &s.Object, &s.Interface, &s.Mode, &s.Container, &s.Dataset, &s.Datagram, &s.Field, &s.Channel, &s.Message, &s.Sender, &s.Receiver}
	for _, id := range ids {
		id.Span = parser.Span{}
	}
	s.Offset.Span = parser.Span{}
	s.Width.Span = parser.Span{}
	s.Ordinal.Span = parser.Span{}
	if s.Variant != nil {
		v := *s.Variant
		v.Span = parser.Span{}
		s.Variant = &v
	}
	if s.ReplyTo != nil {
		v := *s.ReplyTo
		v.Span = parser.Span{}
		s.ReplyTo = &v
	}
	return s
}

type graph struct {
	nodes      map[string]Node
	facts      map[string]Fact
	names      map[string]string
	statements []parser.Statement
}

func graphOf(in Input, limits Limits) (graph, error) {
	g := graph{nodes: map[string]Node{}, facts: map[string]Fact{}, names: map[string]string{}}
	if in.Snapshot == nil {
		return g, fail("INPUT", "missing checked snapshot")
	}
	m := in.Snapshot.Model()
	if m.Header.Version != "0.6" || in.Snapshot.System() == "" {
		return g, fail("PROFILE", "requires checked design-core/0.6 System")
	}
	if in.ExpectedRevision != "" && in.ExpectedRevision != in.Snapshot.Revision() {
		return g, fail("STALE", "source revision differs")
	}
	if len(m.Declarations) > limits.Nodes || len(m.Statements) > limits.Facts {
		return g, fail("LIMIT", "input graph")
	}
	b, _ := json.Marshal(m)
	if len(b) > limits.Bytes {
		return g, fail("LIMIT", "input bytes")
	}
	for _, d := range m.Declarations {
		if !set("system", "unit", "container", "functionality", "capability", "interface", "activity", "mode", "actor", "usecase", "feature", "dataset", "database", "datagram", "contract", "variant", "field", "encoding", "channel", "message", "scenario")[d.Kind] {
			return g, fail("UNSUPPORTED_DECLARATION", d.Kind)
		}
		k := in.Snapshot.System() + "/" + d.Kind + "/" + d.Name.Name
		g.names[d.Name.Name] = k
		g.nodes[k] = Node{Key: k, Name: d.Name.Name, Kind: d.Kind, Origins: Origin{Now: []parser.Span{d.Span}}}
	}
	for _, s := range m.Statements {
		refs, e := endpoints(s)
		if e != nil {
			return g, e
		}
		for i, n := range refs {
			k, ok := g.names[n]
			if !ok {
				return g, fail("REFERENCE", n)
			}
			refs[i] = k
		}
		// Sentence is canonical for the parser's closed typed statement sum, including
		// integer/optional operands. Resolved composition paths are already removed.
		sentence := s.Sentence()
		k := digest(struct {
			System, Sentence string
			Endpoints        []string
		}{in.Snapshot.System(), sentence, unique(refs)})
		f := g.facts[k]
		f.Key = k
		f.Statement = sentence
		f.Semantic = semantic(s)
		f.Endpoints = unique(refs)
		f.Origins.Now = append(f.Origins.Now, s.Span)
		g.facts[k] = f
	}
	g.statements = m.Statements
	return g, nil
}

// Analyze returns a deterministic model-only diagnostic view, never an executable
// assignment. Rules may fail in the returned preview; conflicting rules and
// excluded mandatory closure are errors. A successful result always states the
// absence of runtime/code dependency evidence.
func Analyze(ctx context.Context, now, target Input, task Task, limits Limits) (*Analysis, error) {
	if limits.Nodes <= 0 || limits.Facts <= 0 || limits.Bytes <= 0 {
		return nil, fail("LIMIT", "limits must be positive")
	}
	if task.ID == "" || task.Intent == "" {
		return nil, fail("TASK", "ID and intent required")
	}
	for _, p := range append(append([]string{}, task.AllowedModelPaths...), task.AllowedCodePaths...) {
		if p == "" || path.IsAbs(p) || path.Clean(p) != p || p == ".." || strings.HasPrefix(p, "../") || strings.ContainsAny(p, "\\:*?[]") || strings.IndexFunc(p, unicode.IsControl) >= 0 {
			return nil, fail("TASK_PATH", p)
		}
	}
	for _, p := range task.AllowedChanges {
		if (p.Element == "") == (p.Fact == "") || (p.Change != "added" && p.Change != "removed") {
			return nil, fail("PERMISSION", "exact element or fact and added/removed required")
		}
	}
	b, _ := json.Marshal(task)
	if len(b) > limits.Bytes {
		return nil, fail("LIMIT", "task bytes")
	}
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	a, e := graphOf(now, limits)
	if e != nil {
		return nil, e
	}
	bgraph, e := graphOf(target, limits)
	if e != nil {
		return nil, e
	}
	if now.Snapshot.System() != target.Snapshot.System() {
		return nil, fail("SYSTEM", "different System identities")
	}
	nodes := map[string]Node{}
	facts := map[string]Fact{}
	selected := map[string][]string{}
	for k, n := range a.nodes {
		n.Change = "removed"
		nodes[k] = n
	}
	for k, n := range bgraph.nodes {
		n.Origins.Target = n.Origins.Now
		n.Origins.Now = nil
		n.Change = "added"
		if old, ok := nodes[k]; ok {
			n.Origins.Now = old.Origins.Now
			n.Change = "unchanged"
		}
		nodes[k] = n
	}
	for k, f := range a.facts {
		f.Change = "removed"
		facts[k] = f
	}
	for k, f := range bgraph.facts {
		f.Origins.Target = f.Origins.Now
		f.Origins.Now = nil
		f.Change = "added"
		if old, ok := facts[k]; ok {
			f.Origins.Now = old.Origins.Now
			f.Change = "unchanged"
		}
		facts[k] = f
	}
	if len(nodes) > limits.Nodes || len(facts) > limits.Facts {
		return nil, fail("LIMIT", "union graph")
	}
	for k, n := range nodes {
		if n.Change != "unchanged" {
			selected[k] = []string{"changed declaration"}
		}
	}
	for _, k := range keys(facts) {
		f := facts[k]
		if f.Change != "unchanged" {
			for _, ref := range f.Endpoints {
				selected[ref] = append(selected[ref], "changed fact "+k)
			}
		}
	}
	// Authored keys are typed System/kind/name IDs, not ambiguous display names.
	for _, k := range task.Context {
		if _, ok := nodes[k]; !ok {
			return nil, fail("TASK_REFERENCE", k)
		}
		selected[k] = append(selected[k], "authored context")
	}
	excluded := set(task.Exclude...)
	for k := range excluded {
		if _, ok := nodes[k]; !ok {
			return nil, fail("TASK_REFERENCE", k)
		}
	}
	r := &Analysis{Schema: Version, Policy: PolicyVersion, System: now.Snapshot.System(), NowRevision: now.Snapshot.Revision(), TargetRevision: target.Snapshot.Revision(), Nodes: []Node{}, Facts: []Fact{}, Frontier: []Boundary{}, Excluded: []Boundary{}, Unknowns: []Unknown{}, Obligations: []Obligation{}, ConstraintsPass: true}
	r.Coverage = Coverage{"design-core/0.6", PolicyVersion, true, "unknown", "unknown", "unsupported"}
	protected := map[string]string{}
	protectionIDs := map[string]bool{}
	for _, p := range task.Protect {
		if p.ID == "" || p.Source == "" || protectionIDs[p.ID] || (p.Scope != "boundary" && p.Scope != "subtree") {
			return nil, fail("PROTECTION", "unique ID/provenance and boundary/subtree scope required")
		}
		protectionIDs[p.ID] = true
		if _, ok := nodes[p.Element]; !ok {
			return nil, fail("TASK_REFERENCE", p.Element)
		}
		members := set(p.Element)
		for {
			changed := false
			for _, f := range facts {
				st := f.Semantic
				if st.Kind != "Relation" || (st.Verb != "contains" && st.Verb != "owns") {
					continue
				}
				for _, g := range []graph{a, bgraph} {
					parent, child := g.names[st.Subject.Name], g.names[st.Object.Name]
					if members[parent] && child != "" && !members[child] {
						members[child] = true
						changed = true
					}
				}
			}
			if !changed {
				break
			}
		}
		for k := range members {
			if p.Scope == "subtree" || k == p.Element {
				protected[k] = p.ID
				selected[k] = append(selected[k], "protection "+p.ID)
			}
		}
		for k, f := range facts {
			inside, outside := false, false
			for _, ref := range f.Endpoints {
				if members[ref] {
					inside = true
				} else {
					outside = true
				}
			}
			if inside && (p.Scope == "subtree" || outside) {
				protected[k] = p.ID
				for _, ref := range f.Endpoints {
					selected[ref] = append(selected[ref], "protection "+p.ID)
				}
			}
		}
	}
	for _, p := range task.AllowedChanges {
		for k, f := range facts {
			if p.Fact == f.Statement && protected[k] != "" {
				return nil, fail("RULE_CONFLICT", p.Fact)
			}
		}
		if protected[p.Element] != "" {
			return nil, fail("RULE_CONFLICT", p.Element)
		}
	}
	for _, k := range keys(protected) {
		change := ""
		np, tp := false, false
		if n, ok := nodes[k]; ok {
			change = n.Change
			np = len(n.Origins.Now) > 0
			tp = len(n.Origins.Target) > 0
		} else {
			f := facts[k]
			change = f.Change
			np = len(f.Origins.Now) > 0
			tp = len(f.Origins.Target) > 0
		}
		result := "pass"
		if change != "unchanged" {
			result = "fail"
			r.ConstraintsPass = false
		}
		r.Obligations = append(r.Obligations, Obligation{"protect/" + protected[k] + "/" + k, result, np, tp})
	}
	checkPermission := func(key, sentence, change string, origins Origin) {
		if change == "unchanged" {
			return
		}
		allowed := false
		for _, p := range task.AllowedChanges {
			if p.Change == change && ((sentence == "" && p.Element == key) || (sentence != "" && p.Fact == sentence)) {
				allowed = true
			}
		}
		if len(task.AllowedModelPaths) > 0 {
			paths := set(task.AllowedModelPaths...)
			spans := origins.Now
			if change == "added" {
				spans = origins.Target
			}
			for _, sp := range spans {
				if !paths[sp.Source] {
					allowed = false
				}
			}
		}
		result := "pass"
		if !allowed {
			result = "fail"
			r.ConstraintsPass = false
		}
		r.Obligations = append(r.Obligations, Obligation{"change/" + change + "/" + key, result, len(origins.Now) > 0, len(origins.Target) > 0})
	}
	for _, key := range keys(nodes) {
		n := nodes[key]
		checkPermission(key, "", n.Change, n.Origins)
	}
	for _, key := range keys(facts) {
		f := facts[key]
		checkPermission(key, f.Statement, f.Change, f.Origins)
	}
	// Rules also seed context; spelling must refer to an actual fact on either side.
	rules := map[string]Rule{}
	ruleIDs := map[string]bool{}
	for _, rule := range task.Rules {
		if rule.ID == "" || rule.Source == "" || ruleIDs[rule.ID] {
			return nil, fail("RULE", "unique ID and provenance required")
		}
		ruleIDs[rule.ID] = true
		if old, ok := rules[rule.Fact]; ok && (old.Now != rule.Now || old.Target != rule.Target) {
			return nil, fail("RULE_CONFLICT", rule.Fact)
		}
		for _, p := range task.AllowedChanges {
			if p.Fact == rule.Fact && rule.Now == rule.Target {
				return nil, fail("RULE_CONFLICT", rule.Fact)
			}
		}
		rules[rule.Fact] = rule
		found := false
		np, tp := false, false
		for _, k := range keys(facts) {
			f := facts[k]
			if f.Statement == rule.Fact {
				found = true
				np = np || len(f.Origins.Now) > 0
				tp = tp || len(f.Origins.Target) > 0
				for _, ref := range f.Endpoints {
					selected[ref] = append(selected[ref], "constraint "+rule.ID)
				}
			}
		}
		if !found {
			return nil, fail("RULE_REFERENCE", rule.Fact)
		}
		result := "pass"
		if np != rule.Now || tp != rule.Target {
			result = "fail"
			r.ConstraintsPass = false
		}
		r.Obligations = append(r.Obligations, Obligation{rule.ID, result, np, tp})
	}
	sort.Slice(r.Obligations, func(i, j int) bool { return r.Obligations[i].ID < r.Obligations[j].ID })
	adjacency := map[string][]string{}
	for _, k := range keys(facts) {
		for _, ref := range facts[k].Endpoints {
			adjacency[ref] = append(adjacency[ref], k)
		}
	}
	visited, chosen := map[string]bool{}, map[string]bool{}
	for {
		pending := []string{}
		for _, k := range keys(selected) {
			if !visited[k] {
				pending = append(pending, k)
			}
		}
		if len(pending) == 0 {
			break
		}
		for _, k := range pending {
			if e := ctx.Err(); e != nil {
				return nil, e
			}
			if excluded[k] {
				return nil, fail("SCOPE_CONFLICT", k)
			}
			visited[k] = true
			n := nodes[k]
			if n.Kind == "system" || n.Kind == "mode" {
				continue
			}
			for _, fk := range adjacency[k] {
				chosen[fk] = true
				for _, ref := range facts[fk].Endpoints {
					if _, ok := selected[ref]; !ok {
						selected[ref] = []string{"relation closure " + fk}
					}
				}
			}
		}
	}
	// A fact changed solely on a System/mode still belongs to the delta.
	for k, f := range facts {
		if f.Change != "unchanged" {
			chosen[k] = true
		}
	}
	for _, k := range keys(nodes) {
		n := nodes[k]
		if reasons, ok := selected[k]; ok {
			n.Reasons = unique(reasons)
			r.Nodes = append(r.Nodes, n)
		} else {
			reason := "outside modeled closure"
			if excluded[k] {
				reason = "authored exclusion"
			}
			r.Excluded = append(r.Excluded, Boundary{k, reason})
		}
	}
	for _, k := range keys(facts) {
		f := facts[k]
		if chosen[k] {
			r.Facts = append(r.Facts, f)
		} else {
			for _, ref := range f.Endpoints {
				if _, ok := selected[ref]; ok {
					r.Frontier = append(r.Frontier, Boundary{k, "System/mode propagation boundary"})
					break
				}
			}
		}
	}
	for _, side := range []struct {
		name string
		g    graph
	}{{"NOW", a}, {"TARGET", bgraph}} {
		roles := map[string]map[string]bool{}
		for _, s := range side.g.statements {
			if s.Kind == "Participation" {
				if _, ok := selected[side.g.names[s.Channel.Name]]; !ok {
					continue
				}
				key := s.Channel.Name + "/" + s.Message.Name + "/" + s.Mode.Name
				if roles[key] == nil {
					roles[key] = map[string]bool{}
				}
				roles[key][s.Role] = true
			}
		}
		for _, g := range []graph{a, bgraph} {
			for _, st := range g.statements {
				if st.Kind != "Participation" {
					continue
				}
				if _, ok := selected[side.g.names[st.Channel.Name]]; !ok {
					continue
				}
				permitted := false
				for _, up := range side.g.statements {
					if up.Kind != "Relation" || up.Verb != "upholds" || up.Subject.Name != st.Channel.Name {
						continue
					}
					for _, pe := range side.g.statements {
						if pe.Kind == "Relation" && pe.Verb == "permits" && pe.Subject.Name == up.Object.Name && pe.Object.Name == st.Message.Name {
							permitted = true
						}
					}
				}
				key := st.Channel.Name + "/" + st.Message.Name + "/" + st.Mode.Name
				if permitted && roles[key] == nil {
					roles[key] = map[string]bool{}
				}
			}
		}
		// Contract-permitted messages with no role tuple cannot disappear from diagnostics.
		for _, s := range side.g.statements {
			if s.Kind != "Relation" || s.Verb != "upholds" {
				continue
			}
			n := side.g.nodes[side.g.names[s.Subject.Name]]
			if n.Kind != "channel" {
				continue
			}
			if _, ok := selected[n.Key]; !ok {
				continue
			}
			any := false
			for _, permit := range side.g.statements {
				if permit.Kind != "Relation" || permit.Verb != "permits" || permit.Subject.Name != s.Object.Name {
					continue
				}
				if side.g.nodes[side.g.names[permit.Object.Name]].Kind != "message" {
					continue
				}
				any = true
				found := false
				for key := range roles {
					if strings.HasPrefix(key, s.Subject.Name+"/"+permit.Object.Name+"/") {
						found = true
					}
				}
				if !found {
					r.Unknowns = append(r.Unknowns, Unknown{side.name, "CHANNEL_PEER", s.Subject.Name + "/" + permit.Object.Name, "no modeled participants or mode for permitted message"})
				}
			}
			if !any {
				r.Unknowns = append(r.Unknowns, Unknown{side.name, "CHANNEL_PROTOCOL", s.Subject.Name, "no modeled permitted messages"})
			}
		}
		for _, k := range keys(roles) {
			if !roles[k]["sender"] || !roles[k]["receiver"] {
				r.Unknowns = append(r.Unknowns, Unknown{side.name, "CHANNEL_PEER", k, "modeled sender/receiver pair missing"})
			}
		}
		for _, s := range side.g.statements {
			if s.Kind != "Dependency" {
				continue
			}
			if _, ok := selected[side.g.names[s.Subject.Name]]; !ok {
				continue
			}
			r.Unknowns = append(r.Unknowns, Unknown{side.name, "RUNTIME_PROVIDER", s.Interface.Name, "structural dependency does not establish an available runtime provider"})
		}
	}
	r.Unknowns = append(r.Unknowns, Unknown{"BOTH", "CODE_MAPPING", "", "No code mappings, complete runtime dependencies or behavioral evidence supplied"})
	sort.Slice(r.Unknowns, func(i, j int) bool {
		a, _ := json.Marshal(r.Unknowns[i])
		b, _ := json.Marshal(r.Unknowns[j])
		return string(a) < string(b)
	})
	r.Identity = digest(struct {
		Version, Policy, Now, Target string
		Task                         Task
		Limits                       Limits
	}{Version, PolicyVersion, r.NowRevision, r.TargetRevision, task, limits})
	output, _ := json.Marshal(r)
	if len(output) > limits.Bytes {
		return nil, fail("LIMIT", "analysis bytes")
	}
	return r, nil
}
