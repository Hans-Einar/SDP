package blueprint

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/sourcegraph"
)

func compile(t *testing.T, sources map[string]string) Input {
	t.Helper()
	c := parser.SyntaxCache{}
	files := []*parser.File{}
	for _, name := range keys(sources) {
		f, e := c.Parse(name, sources[name])
		if e != nil {
			t.Fatal(e)
		}
		files = append(files, f)
	}
	s, issues := sourcegraph.Compile("System.design", files, false)
	if len(issues) > 0 {
		t.Fatal(issues)
	}
	return Input{Snapshot: s}
}
func pilot(t *testing.T, side string) map[string]string {
	t.Helper()
	root := filepath.Join("../../../experiments/blueprint_mvp1", side)
	out := map[string]string{}
	e := filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			name, _ := filepath.Rel(root, p)
			out[filepath.ToSlash(name)] = string(b)
		}
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	return out
}
func task() Task { return Task{ID: "calibration", Intent: "Extract calibration ownership"} }
func run(t *testing.T, a, b Input, task Task) *Analysis {
	t.Helper()
	r, e := Analyze(context.Background(), a, b, task, DefaultLimits())
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func code(t *testing.T, e error, want string) {
	t.Helper()
	v, ok := e.(*Error)
	if !ok || v.Code != want {
		t.Fatalf("wanted %s, got %v", want, e)
	}
}
func TestPilotAndPreserve(t *testing.T) {
	a, b := compile(t, pilot(t, "NOW")), compile(t, pilot(t, "TARGET"))
	task := task()
	task.AllowedChanges = []Permission{
		{Element: "MVP1InspectionPilot/unit/CalibrationProcessor", Change: "added"},
		{Fact: "MachineService owns CalibrateMeasurement.", Change: "removed"},
		{Fact: "CalibrationProcessor owns CalibrateMeasurement.", Change: "added"},
		{Fact: "MachineService contains CalibrationProcessor.", Change: "added"},
	}
	task.Rules = []Rule{{ID: "owner", Source: "owner-task", Fact: "MachineService owns ReduceMachineState.", Now: true, Target: true}}
	r := run(t, a, b, task)
	if !r.Coverage.ModeledClosureComplete || !r.ConstraintsPass || len(r.Unknowns) == 0 || len(r.Facts) == 0 {
		t.Fatalf("%+v", r)
	}
	added, removed := false, false
	for _, f := range r.Facts {
		if strings.Contains(f.Statement, "owns CalibrateMeasurement") {
			added = added || f.Change == "added"
			removed = removed || f.Change == "removed"
		}
	}
	if !added || !removed {
		t.Fatal("lost ownership delta")
	}
	encoded, _ := json.Marshal(r)
	again, _ := json.Marshal(run(t, a, b, task))
	if string(encoded) != string(again) {
		t.Fatal("nondeterministic")
	}
	bad := pilot(t, "TARGET")
	bad["Containers/MachineService.design"] = strings.ReplaceAll(bad["Containers/MachineService.design"], "MachineService owns ReduceMachineState.", "BuckingWeb owns ReduceMachineState.")
	if run(t, a, compile(t, bad), task).ConstraintsPass {
		t.Fatal("accepted ownership drift")
	}
}
func TestFailures(t *testing.T) {
	a, b := compile(t, pilot(t, "NOW")), compile(t, pilot(t, "TARGET"))
	tests := []struct {
		name   string
		mutate func(*Task, *Input, *Limits)
		want   string
	}{
		{"stale", func(_ *Task, a *Input, _ *Limits) { a.ExpectedRevision = "old" }, "STALE"},
		{"exclude", func(x *Task, _ *Input, _ *Limits) {
			x.Exclude = []string{"MVP1InspectionPilot/unit/CalibrationProcessor"}
		}, "SCOPE_CONFLICT"},
		{"unknown context", func(x *Task, _ *Input, _ *Limits) { x.Context = []string{"missing"} }, "TASK_REFERENCE"},
		{"nodes", func(_ *Task, _ *Input, l *Limits) { l.Nodes = 1 }, "LIMIT"},
		{"facts", func(_ *Task, _ *Input, l *Limits) { l.Facts = 1 }, "LIMIT"},
		{"bytes", func(_ *Task, _ *Input, l *Limits) { l.Bytes = 1 }, "LIMIT"},
		{"conflict", func(x *Task, _ *Input, _ *Limits) {
			x.Rules = []Rule{{"a", "task", "MachineService owns ReduceMachineState.", true, true}, {"b", "task", "MachineService owns ReduceMachineState.", true, false}}
		}, "RULE_CONFLICT"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			x := task()
			n := a
			l := DefaultLimits()
			tc.mutate(&x, &n, &l)
			_, e := Analyze(context.Background(), n, b, x, l)
			code(t, e, tc.want)
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := Analyze(ctx, a, b, task(), DefaultLimits()); e != context.Canceled {
		t.Fatal(e)
	}
	_, e := Analyze(context.Background(), Input{}, b, task(), DefaultLimits())
	code(t, e, "INPUT")
}
func TestMissingPeerAndFormatting(t *testing.T) {
	src := pilot(t, "NOW")
	a := compile(t, src)
	formatted := pilot(t, "NOW")
	formatted["System.design"] = "\n" + formatted["System.design"]
	// actual whitespace, not a literal backslash
	formatted["System.design"] = strings.Replace(formatted["System.design"], "\\n", "\n", 1)
	b := compile(t, formatted)
	r := run(t, a, b, task())
	if len(r.Nodes) != 0 || a.Snapshot.Revision() == b.Snapshot.Revision() {
		t.Fatal("formatting must change identity, not semantics")
	}
	for p, s := range src {
		lines := strings.Split(s, "\n")
		kept := []string{}
		for _, line := range lines {
			if strings.Contains(line, "as receiver") && strings.Contains(line, "ReadBaseline") {
				continue
			}
			kept = append(kept, line)
		}
		src[p] = strings.Join(kept, "\n")
	}
	r = run(t, a, compile(t, src), task())
	found := false
	for _, u := range r.Unknowns {
		found = found || u.Code == "CHANNEL_PEER" && u.Side == "TARGET"
	}
	if !found {
		t.Fatal("missing peer hidden by union")
	}
}
func TestSystemAndProfile(t *testing.T) {
	mk := func(name, version string) Input {
		return compile(t, map[string]string{"System.design": "language design-core version " + version + ".\nsystem " + name + ".\n"})
	}
	_, e := Analyze(context.Background(), mk("A", "0.6"), mk("B", "0.6"), task(), DefaultLimits())
	code(t, e, "SYSTEM")
	a := compile(t, map[string]string{"System.design": "language design-core version 0.5.\nunit X.\n"})
	_, e = Analyze(context.Background(), a, mk("A", "0.6"), task(), DefaultLimits())
	code(t, e, "PROFILE")
}
func TestCycleKindRenameAndBoundaries(t *testing.T) {
	text := "language design-core version 0.6.\nsystem S.\nactivity A.\nactivity B.\nactivity C.\nunit U.\nunit V.\nS contains U.\nS contains V.\nA depends-on B.\nA depends-on C.\nB depends-on C.\n"
	a := compile(t, map[string]string{"System.design": text})
	x := task()
	x.Context = []string{"S/activity/A"}
	r := run(t, a, a, x)
	if len(r.Nodes) != 3 || len(r.Facts) != 3 {
		t.Fatalf("cycle closure: %+v", r)
	}
	x.Context = []string{"S/unit/U"}
	r = run(t, a, a, x)
	if len(r.Nodes) != 2 || len(r.Frontier) != 1 {
		t.Fatal("System expanded siblings")
	}
	b := compile(t, map[string]string{"System.design": strings.ReplaceAll(text, "unit V.", "container V.")})
	r = run(t, a, b, task())
	removed, added := false, false
	for _, n := range r.Nodes {
		if n.Name == "V" {
			removed = removed || n.Change == "removed"
			added = added || n.Change == "added"
		}
	}
	if !removed || !added {
		t.Fatal("kind change lost")
	}
	b = compile(t, map[string]string{"System.design": strings.ReplaceAll(text, " V", " W")})
	r = run(t, a, b, task())
	if len(r.Nodes) == 0 {
		t.Fatal("rename missing")
	}
}
func TestTypedOperands(t *testing.T) {
	// Closed AST policy unit checks complement parser-backed workflow tests.
	cases := []parser.Statement{
		{Kind: "Dependency", Subject: parser.Identifier{Name: "a"}, Interface: parser.Identifier{Name: "b"}, Mode: parser.Identifier{Name: "c"}},
		{Kind: "Projection", Subject: parser.Identifier{Name: "a"}, Dataset: parser.Identifier{Name: "b"}, Datagram: parser.Identifier{Name: "c"}},
		{Kind: "Placement", Subject: parser.Identifier{Name: "a"}, Field: parser.Identifier{Name: "b"}, Width: parser.Integer{Value: 8}},
		{Kind: "Step", Subject: parser.Identifier{Name: "a"}, Message: parser.Identifier{Name: "b"}, Sender: parser.Identifier{Name: "c"}, Receiver: parser.Identifier{Name: "d"}, Channel: parser.Identifier{Name: "e"}, Variant: &parser.Identifier{Name: "f"}, ReplyTo: &parser.Integer{Value: 1}},
	}
	for _, s := range cases {
		refs, e := endpoints(s)
		if e != nil || len(refs) < 2 {
			t.Fatal(refs, e)
		}
		old := s.Sentence()
		s.Subject.Name = "different"
		if old == s.Sentence() {
			t.Fatal("semantic operand ignored")
		}
	}
	_, e := endpoints(parser.Statement{Kind: "Future"})
	code(t, e, "UNSUPPORTED_STATEMENT")
	_, e = endpoints(parser.Statement{Kind: "Relation", Verb: "future"})
	code(t, e, "UNSUPPORTED_RELATION")
}

func TestNoParticipantsAndPermissions(t *testing.T) {
	src := pilot(t, "TARGET")
	a := compile(t, pilot(t, "NOW"))
	for p, s := range src {
		kept := []string{}
		for _, line := range strings.Split(s, "\n") {
			if strings.Contains(line, " uses ") && strings.Contains(line, " of ReadBaseline ") {
				continue
			}
			kept = append(kept, line)
		}
		src[p] = strings.Join(kept, "\n")
	}
	r := run(t, a, compile(t, src), task())
	found := false
	for _, u := range r.Unknowns {
		found = found || u.Code == "CHANNEL_PEER" && u.Side == "TARGET" && strings.Contains(u.Subject, "ReadBaseline")
	}
	if !found {
		t.Fatal("zero participants not diagnosed")
	}
	if r.ConstraintsPass {
		t.Fatal("unapproved semantic changes passed")
	}
	x := task()
	x.AllowedCodePaths = []string{"../outside.go"}
	_, e := Analyze(context.Background(), a, a, x, DefaultLimits())
	code(t, e, "TASK_PATH")
	x = task()
	x.AllowedChanges = []Permission{{Fact: "MachineService owns ReduceMachineState.", Change: "removed"}}
	x.Rules = []Rule{{"keep", "task", "MachineService owns ReduceMachineState.", true, true}}
	_, e = Analyze(context.Background(), a, a, x, DefaultLimits())
	code(t, e, "RULE_CONFLICT")
}
func TestScopedProtectionAndTypedFacts(t *testing.T) {
	a, b := compile(t, pilot(t, "NOW")), compile(t, pilot(t, "TARGET"))
	x := task()
	x.Protect = []Protection{{"machine", "task", "MVP1InspectionPilot/container/MachineService", "subtree"}}
	r := run(t, a, b, x)
	failed := false
	typed := false
	for _, o := range r.Obligations {
		if strings.HasPrefix(o.ID, "protect/") && o.Result == "fail" {
			failed = true
		}
	}
	for _, f := range r.Facts {
		if f.Semantic.Kind == "Participation" && f.Semantic.Role != "" && f.Semantic.Channel.Name != "" {
			typed = true
		}
	}
	if !failed || !typed {
		t.Fatal("protection/typed facts missing")
	}
	x.Protect[0].Scope = "boundary"
	r = run(t, a, b, x)
	for _, o := range r.Obligations {
		if strings.HasPrefix(o.ID, "protect/") && o.Result == "fail" {
			t.Fatalf("internal calibration extraction violated boundary: %+v", o)
		}
	}
}
func TestSemanticOperandsAndMoves(t *testing.T) {
	s := parser.Statement{Kind: "Placement", Subject: parser.Identifier{Name: "E"}, Field: parser.Identifier{Name: "F"}, Width: parser.Integer{Value: 8}}
	n := s
	n.Width.Value = 16
	if digest(semantic(s)) == digest(semantic(n)) {
		t.Fatal("width ignored")
	}
	step := parser.Statement{Kind: "Step", Ordinal: parser.Integer{Value: 1}, ReplyTo: &parser.Integer{Value: 1}}
	other := step
	other.ReplyTo = &parser.Integer{Value: 2}
	if digest(semantic(step)) == digest(semantic(other)) {
		t.Fatal("reply target ignored")
	}
	src := pilot(t, "NOW")
	a := compile(t, src)
	src["Moved/BuckingUI.design"] = src["Containers/BuckingUI.design"]
	delete(src, "Containers/BuckingUI.design")
	src["System.design"] = strings.ReplaceAll(src["System.design"], "Containers/BuckingUI.", "Moved/BuckingUI.")
	b := compile(t, src)
	r := run(t, a, b, task())
	if len(r.Nodes) != 0 || len(r.Facts) != 0 || a.Snapshot.Revision() == b.Snapshot.Revision() {
		t.Fatal("move comparison wrong")
	}
}
