package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
)

func activityIntegrationFixture(t *testing.T) (string, *activities.Service) {
	t.Helper()
	root, _ := projectFixture(t)
	for _, dir := range []string{"KanBan", "ProjectManagement"} {
		if e := os.MkdirAll(filepath.Join(root, "SDP", dir), 0700); e != nil {
			t.Fatal(e)
		}
	}
	if e := os.WriteFile(filepath.Join(root, "SDP/KanBan/board.json"), []byte(`{"projectId":"TEST"}`), 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(root, "SDP/owner.md"), []byte("# Owner decision\n\nCreate this disposable need.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	return root, activities.New(filepath.Join(root, "SDP"), "TEST")
}
func activityIntegrationID(n int) string { return fmt.Sprintf("%08d-1234-4123-8123-%012d", n, n) }
func activityIntegrationRequest(t *testing.T, s *activities.Service, n int) activities.CreateRequest {
	t.Helper()
	c, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	return activities.CreateRequest{OperationID: activityIntegrationID(n), ExpectedHistoryRevision: c.HistoryRevision, Title: "Disposable need", Owner: "project-owner", Mandate: "Study and deliver the selected disposable need.", Reason: "Owner-selected integration fixture", Scope: []string{"Source"}}
}
func activityIntegrationJSON(t *testing.T, root, name string, value any) string {
	t.Helper()
	b, e := json.Marshal(value)
	if e != nil {
		t.Fatal(e)
	}
	p := filepath.Join(root, name)
	if e = os.WriteFile(p, b, 0600); e != nil {
		t.Fatal(e)
	}
	return p
}
func activityIntegrationCLI(t *testing.T, root string, args ...string) (int, []byte, string) {
	t.Helper()
	var out, errs bytes.Buffer
	code := Run(context.Background(), append(append([]string{root, "activity"}, args...), "--json"), &out, &errs)
	return code, out.Bytes(), errs.String()
}
func activityIntegrationResult(t *testing.T, data []byte) activities.Result {
	t.Helper()
	var envelope struct {
		Schema string            `json:"schema"`
		Data   activities.Result `json:"data"`
	}
	if e := json.Unmarshal(data, &envelope); e != nil || envelope.Schema != activities.ApplicationSchema {
		t.Fatalf("%s %v", data, e)
	}
	return envelope.Data
}
func activityIntegrationCreateCLI(t *testing.T, root string, s *activities.Service, n int) activities.Result {
	t.Helper()
	p := activityIntegrationJSON(t, root, fmt.Sprintf("create-%d.json", n), activityIntegrationRequest(t, s, n))
	code, data, errs := activityIntegrationCLI(t, root, "create", "--request", p, "--authority-source", "owner.md")
	if code != 0 {
		t.Fatal(code, errs)
	}
	return activityIntegrationResult(t, data)
}

func TestActivityIntegrationCLICreateSessionRetryAndStale(t *testing.T) {
	root, s := activityIntegrationFixture(t)
	q := activityIntegrationRequest(t, s, 1)
	p := activityIntegrationJSON(t, root, "create.json", q)
	code, data, errs := activityIntegrationCLI(t, root, "create", "--request", p, "--authority-source", "owner.md")
	if code != 0 {
		t.Fatal(code, errs)
	}
	first := activityIntegrationResult(t, data)
	code, retry, errs := activityIntegrationCLI(t, root, "create", "--request", p, "--authority-source", "owner.md")
	if code != 0 || activityIntegrationResult(t, retry) != first {
		t.Fatal("retry changed", code, errs)
	}
	q.OperationID = activityIntegrationID(2)
	p = activityIntegrationJSON(t, root, "stale.json", q)
	if code, _, errs = activityIntegrationCLI(t, root, "create", "--request", p, "--authority-source", "owner.md"); code == 0 || !strings.Contains(errs, "revision") {
		t.Fatal("stale accepted", code, errs)
	}
	cat, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	a := cat.Activities[0]
	session := activities.CreateSessionRequest{OperationID: activityIntegrationID(3), ExpectedHistoryRevision: cat.HistoryRevision, ActivityID: a.Descriptor.ID, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Title: "Study discussion", Objective: "Investigate before choosing a plan", Phase: "02-Study", Reason: "Owner requests study"}
	p = activityIntegrationJSON(t, root, "session.json", session)
	code, data, errs = activityIntegrationCLI(t, root, "session", "create", "--request", p, "--authority-source", "owner.md")
	if code != 0 {
		t.Fatal(code, errs)
	}
	created := activityIntegrationResult(t, data)
	if created.SessionID != "SESSION-TEST-0001" {
		t.Fatal(created)
	}
	b, e := os.ReadFile(filepath.Join(s.Area, filepath.FromSlash(created.Path)))
	if e != nil || !strings.Contains(string(b), "| Status | planned |") {
		t.Fatal(string(b), e)
	}
	if _, e = os.Stat(filepath.Join(s.Area, "Sessions/execution")); !os.IsNotExist(e) {
		t.Fatal("record creation initialized runtime binding", e)
	}
	code, data, errs = activityIntegrationCLI(t, root, "session-create", "--request", p, "--authority-source", "owner.md")
	if code != 0 || activityIntegrationResult(t, data) != created {
		t.Fatal(code, errs)
	}
	if code, _, _ = activityIntegrationCLI(t, root, "create", "--request", p); code == 0 {
		t.Fatal("missing authority source accepted")
	}
}

func TestActivityIntegrationInvokeAuthorityStrictRequestsAndParity(t *testing.T) {
	root, s := activityIntegrationFixture(t)
	q := activityIntegrationRequest(t, s, 1)
	activityIntegrationJSON(t, root, "create.json", q)
	params := ActivityRequestInput{RequestFile: "create.json", AuthoritySource: "owner.md"}
	raw := requestBytes(t, root, "activity.create", params)
	for _, caller := range []ActionCaller{{}, {Principal: blueprintstate.Principal{Actor: "agent", Authority: "assignee"}}, {Principal: blueprintstate.Principal{Actor: "controller", Authority: "controller"}}} {
		r := Invoke(context.Background(), bytes.NewReader(raw), caller)
		if r.Error == nil || r.Error.Code != "unauthorized" {
			t.Fatalf("untrusted mutation: %+v", r)
		}
	}
	owner := ActionCaller{Principal: blueprintstate.Principal{Actor: "owner-host", Authority: "owner"}}
	r := Invoke(context.Background(), bytes.NewReader(raw), owner)
	if r.ExitCode != 0 || r.Error != nil {
		t.Fatalf("trusted invocation failed: %+v", r)
	}
	retry := Invoke(context.Background(), bytes.NewReader(raw), owner)
	if !reflect.DeepEqual(r.Result, retry.Result) {
		t.Fatal("typed action retry differs")
	}
	activityIntegrationJSON(t, root, "bad.json", map[string]any{"operationId": activityIntegrationID(2), "unexpected": true, "actor": "forged-owner"})
	bad := Invoke(context.Background(), bytes.NewReader(requestBytes(t, root, "activity.create", ActivityRequestInput{RequestFile: "bad.json", AuthoritySource: "owner.md"})), owner)
	if bad.Error == nil || !strings.Contains(bad.Error.Message, "unknown") {
		t.Fatalf("unknown request accepted: %+v", bad)
	}
	catalog, e := Actions(context.Background(), root)
	if e != nil {
		t.Fatal(e)
	}
	registered := map[string]ActionDescriptor{}
	for _, d := range catalog.Actions {
		registered[d.ID] = d
	}
	for _, op := range activities.Operations() {
		d, ok := registered[op.ActionID]
		if !ok || d.Group != "Activities" || d.Label != op.Label {
			t.Fatalf("operation not registered: %+v", op)
		}
		if _, ok := activities.FindOperation(op.Name); !ok {
			t.Fatal(op.Name)
		}
	}
	if _, ok := cliCommands()["activity"]; !ok {
		t.Fatal("CLI Activity dispatcher absent")
	}
	for _, op := range []string{"list", "read", "history", "validate"} {
		params := any(map[string]any{})
		args := []string{op}
		if op == "read" || op == "history" {
			params = ActivityReadInput{Activity: "ACTIVITY-TEST-0001"}
			args = append(args, "--activity", "ACTIVITY-TEST-0001")
		}
		code, out, errs := activityIntegrationCLI(t, root, args...)
		actual := Invoke(context.Background(), bytes.NewReader(requestBytes(t, root, "activity."+op, params)), ActionCaller{})
		if code != 0 || actual.ExitCode != 0 {
			t.Fatal(op, code, errs, actual.Error)
		}
		var cliValue, actionValue any
		json.Unmarshal(out, &cliValue)
		b, _ := json.Marshal(actual.Result)
		json.Unmarshal(b, &actionValue)
		if !reflect.DeepEqual(cliValue, actionValue) {
			t.Fatalf("%s CLI/action parity diverged", op)
		}
	}
}

func TestActivityIntegrationPagedNavigationSelectorsAndBundles(t *testing.T) {
	root, s := activityIntegrationFixture(t)
	created := activityIntegrationCreateCLI(t, root, s, 1)
	phase := filepath.Join(s.Area, filepath.FromSlash(created.Path), "02-Study")
	for i := 0; i < 205; i++ {
		if e := os.WriteFile(filepath.Join(phase, fmt.Sprintf("Study-%03d.md", i)), []byte(fmt.Sprintf("# Study %d\n", i)), 0600); e != nil {
			t.Fatal(e)
		}
	}
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	hasActivity := false
	queue := []*Target{}
	for _, n := range p.Navigation.Nodes {
		if n.ID == "sessions" {
			t.Fatal("duplicate standalone Sessions tab")
		}
		if n.ID == "activities" {
			hasActivity = true
		}
		if n.ExpansionTarget != nil {
			queue = append(queue, n.ExpansionTarget)
		}
	}
	if !hasActivity {
		t.Fatal("Activities tab absent")
	}
	seenTargets := map[string]bool{}
	documents := map[string]bool{}
	pages := 0
	var phaseTarget *Target
	for len(queue) > 0 {
		target := queue[0]
		queue = queue[1:]
		if seenTargets[target.URI] {
			t.Fatal("expansion cycle", target.URI)
		}
		seenTargets[target.URI] = true
		tree, e := ExpandActivity(p, target.URI, target.Revision)
		if e != nil {
			t.Fatal(e)
		}
		if pages == 0 {
			code, data, errs := activityIntegrationCLI(t, root, "tree", "--uri", target.URI, "--revision", target.Revision)
			var cliTree Tree
			if code != 0 || json.Unmarshal(data, &cliTree) != nil || !reflect.DeepEqual(cliTree, tree) {
				t.Fatal("tree CLI adapter", code, errs)
			}
			invoked := Invoke(context.Background(), bytes.NewReader(requestBytes(t, root, "activity.tree", ActivityTreeInput{URI: target.URI, Revision: target.Revision})), ActionCaller{})
			if invoked.ExitCode != 0 {
				t.Fatal("tree action adapter", invoked.Error)
			}
		}
		pages++
		if len(tree.Nodes) > 102 {
			t.Fatal("unbounded page", len(tree.Nodes))
		}
		for _, n := range tree.Nodes[1:] {
			if n.Label == "Study" {
				phaseTarget = n.Target
			}
			if n.ExpansionTarget != nil {
				queue = append(queue, n.ExpansionTarget)
			}
			if n.Target != nil && n.Target.Operation == "open" {
				if documents[n.Target.Path] {
					t.Fatal("duplicate document", n.Target.Path)
				}
				documents[n.Target.Path] = true
				if _, e := os.Stat(n.Target.Path); e != nil {
					t.Fatal(e)
				}
			}
		}
	}
	if pages < 5 || phaseTarget == nil {
		t.Fatal("did not exercise nested pagination", pages)
	}
	for i := 0; i < 205; i++ {
		if !documents[filepath.Join(phase, fmt.Sprintf("Study-%03d.md", i))] {
			t.Fatal("pagination lost document", i)
		}
	}
	projection, e := activityCatalog(p)
	if e != nil {
		t.Fatal(e)
	}
	uri := "activity-view://" + p.Inventory.ProjectID + "/activity/" + created.ActivityID + "/study"
	out := filepath.Join(t.TempDir(), "view")
	code, bundleData, errs := activityIntegrationCLI(t, root, "view", "--uri", uri, "--revision", projection.Revision, "--output", out)
	var result Result
	if code != 0 || json.Unmarshal(bundleData, &result) != nil {
		t.Fatal("view CLI adapter", code, errs)
	}
	b, e := os.ReadFile(result.Entry)
	if e != nil {
		t.Fatal(e)
	}
	links := regexp.MustCompile(`\]\((file:[^)]+)\)`).FindAllSubmatch(b, -1)
	if len(links) < 205 {
		t.Fatal("source links missing", len(links))
	}
	for _, match := range links {
		u, e := url.Parse(string(match[1]))
		if e != nil {
			t.Fatal(e)
		}
		if _, e = os.Stat(u.Path); e != nil {
			t.Fatal("broken source link", u, e)
		}
	}
	if _, e = ActivityView(context.Background(), p, uri, projection.Revision, filepath.Join(s.Area, "generated-view")); e == nil {
		t.Fatal("view wrote into SDP")
	}
	if _, e = os.Stat(filepath.Join(s.Area, "generated-view")); !os.IsNotExist(e) {
		t.Fatal("refused view created output", e)
	}
	for _, bad := range []string{"activity-tree://other/children?node=a&parent=a&offset=0", "activity-tree://trial/children?node=activities&parent=activities&offset=1", queueURIForTest(p, "absent")} {
		if _, e := ExpandActivity(p, bad, projection.Revision); e == nil {
			t.Fatal("bad selector accepted", bad)
		}
	}
	if _, e := ActivityView(context.Background(), p, "activity-view://other/overview", projection.Revision, filepath.Join(t.TempDir(), "bad")); e == nil {
		t.Fatal("cross-project view accepted")
	}
	if _, e := ActivityView(context.Background(), p, uri, "bad-revision", filepath.Join(t.TempDir(), "stale")); e == nil {
		t.Fatal("wrong bundle revision accepted")
	}
	if e = os.WriteFile(filepath.Join(phase, "Added.md"), []byte("# Added\n"), 0600); e != nil {
		t.Fatal(e)
	}
	for _, targetURI := range keysActivityTargets(seenTargets) {
		if _, e := ExpandActivity(p, targetURI, projection.Revision); e == nil {
			t.Fatal("stale page accepted")
		}
		break
	}
	if _, e := ActivityView(context.Background(), p, uri, projection.Revision, filepath.Join(t.TempDir(), "changed")); e == nil {
		t.Fatal("stale bundle accepted")
	}
}
func queueURIForTest(p Project, node string) string {
	return activityExpansionTarget(p, node, node, 0, "unused").URI
}
func keysActivityTargets(m map[string]bool) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestActivityIntegrationRealManualCatalogReadOnly(t *testing.T) {
	root, e := filepath.Abs("..")
	if e != nil {
		t.Fatal(e)
	}
	ledger := filepath.Join(root, "SDP/ProjectManagement/Ledger.ndjson")
	before, e := os.ReadFile(ledger)
	if e != nil {
		t.Fatal(e)
	}
	code, data, errs := activityIntegrationCLI(t, root, "list")
	if code != 0 {
		t.Fatal(code, errs)
	}
	var response struct {
		Data activities.Catalog `json:"data"`
	}
	if e = json.Unmarshal(data, &response); e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{"ACTIVITY-SDP-0001", "ACTIVITY-SDP-0002"} {
		a, ok := response.Data.Find(id)
		if !ok || len(a.Events) == 0 || a.Events[0].Payload.SchemaVersion != "0.1" {
			t.Fatal("manual compatibility changed", id, string(data))
		}
	}
	after, e := os.ReadFile(ledger)
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("read changed canonical history", e)
	}
}

func TestActivityIntegrationAdoptionUsesBoundGovernanceProjectIdentity(t *testing.T) {
	root, s := activityIntegrationFixture(t)
	created := activityIntegrationCreateCLI(t, root, s, 1)
	state := filepath.Join(t.TempDir(), "operational")
	var out, errs bytes.Buffer
	if code := Run(context.Background(), []string{root, "governance", "bind", "--state", state, "--json"}, &out, &errs); code != 0 {
		t.Fatal("initialize explicit governance store", code, errs.String())
	}
	cat, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	q := activities.AdoptionRequest{OperationID: activityIntegrationID(2), ExpectedHistoryRevision: cat.HistoryRevision, Stores: []activities.ExecutionStore{{StatePath: state}}}
	request := activityIntegrationJSON(t, root, "adopt.json", q)
	if code, _, errs := activityIntegrationCLI(t, root, "adopt", "--request", request, "--authority-source", "owner.md"); code != 0 {
		t.Fatal("adopt existing initialized store", code, errs)
	}
	cat, e = s.Read()
	if e != nil {
		t.Fatal(e)
	}
	a := cat.Activities[0]
	activate := activities.TransitionRequest{OperationID: activityIntegrationID(3), ExpectedHistoryRevision: cat.HistoryRevision, ActivityID: created.ActivityID, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Reason: "Owner selects need"}
	request = activityIntegrationJSON(t, root, "activate.json", activate)
	if code, _, errs := activityIntegrationCLI(t, root, "activate", "--request", request, "--authority-source", "owner.md"); code != 0 {
		t.Fatal(code, errs)
	}
	cat, e = s.Read()
	if e != nil {
		t.Fatal(e)
	}
	a = cat.Activities[0]
	hold := activities.TransitionRequest{OperationID: activityIntegrationID(4), ExpectedHistoryRevision: cat.HistoryRevision, ActivityID: a.Descriptor.ID, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Reason: "Owner pauses need", ResumeCondition: "Owner selects this same need again"}
	request = activityIntegrationJSON(t, root, "hold.json", hold)
	code, data, stderr := activityIntegrationCLI(t, root, "hold", "--request", request, "--authority-source", "owner.md")
	if code != 0 || activityIntegrationResult(t, data).Pending {
		t.Fatal("quiescent actual store failed hold", code, stderr)
	}
}

func TestActivityIntegrationUnavailableExternalReferencesRemainDistinct(t *testing.T) {
	root, s := activityIntegrationFixture(t)
	q := activityIntegrationRequest(t, s, 1)
	for _, project := range []string{"ALPHA", "BETA"} {
		q.References = append(q.References, activities.Reference{ProjectID: project, Role: "study", ID: "STUDY-" + project + "-0001", Path: "Study.md", Availability: "unavailable", UnavailableReason: "External project offline"})
	}
	for _, pair := range [][2]string{{"context/a", "b"}, {"context", "a/b"}} {
		q.References = append(q.References, activities.Reference{ProjectID: "OTHER", Role: pair[0], ID: pair[1], Path: "Study.md", Availability: "unavailable", UnavailableReason: "External project offline"})
	}
	if _, e := s.Create(q, activities.Caller{Actor: "test", Authority: activities.Authority{Role: "project-owner", Source: "owner.md"}}); e != nil {
		t.Fatal(e)
	}
	p, e := discover(root, false)
	if e != nil {
		t.Fatal(e)
	}
	seen := map[string]bool{}
	external := 0
	for _, n := range activityNodes(p) {
		if seen[n.ID] {
			t.Fatal("duplicate node", n.ID)
		}
		seen[n.ID] = true
		if n.Diagnostic == "External project offline" {
			external++
			if n.Target != nil {
				t.Fatal("external source opened locally")
			}
		}
	}
	if external != 4 {
		t.Fatalf("missing external references: %d", external)
	}
}

func TestActivityIntegrationManualPilotAdoptionPreservesSourceHistory(t *testing.T) {
	root, s := activityIntegrationFixture(t)
	ledger, e := os.ReadFile("../SDP/ProjectManagement/Ledger.ndjson")
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson"), ledger, 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(s.Area, "KanBan/board.json"), []byte(`{"projectId":"SDP"}`), 0600); e != nil {
		t.Fatal(e)
	}
	originals := map[string][]byte{}
	for _, home := range []string{"Activity-0001--SDPTool-Restructuring", "Activity-0002--Activity-Management"} {
		rel := "Activities/" + home + "/activity.yaml"
		b, e := os.ReadFile(filepath.Join("../SDP", rel))
		if e != nil {
			t.Fatal(e)
		}
		originals[rel] = b
		target := filepath.Join(s.Area, rel)
		if e = os.MkdirAll(filepath.Dir(target), 0700); e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(target, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	s = activities.New(s.Area, "SDP")
	catalog, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	state := filepath.Join(t.TempDir(), "state")
	var out, errs bytes.Buffer
	if code := Run(context.Background(), []string{root, "governance", "bind", "--state", state, "--json"}, &out, &errs); code != 0 {
		t.Fatal(code, errs.String())
	}
	q := activities.AdoptionRequest{OperationID: activityIntegrationID(90), ExpectedHistoryRevision: catalog.HistoryRevision, Stores: []activities.ExecutionStore{{StatePath: state}}}
	req := activityIntegrationJSON(t, root, "adopt.json", q)
	code, _, diagnostic := activityIntegrationCLI(t, root, "adopt", "--request", req, "--authority-source", "owner.md")
	if code != 0 {
		t.Fatal(code, diagnostic)
	}
	after, e := s.Read()
	if e != nil || len(after.Activities) != 2 || after.Activities[1].State != "active" {
		t.Fatal(after, e)
	}
	stream, _ := os.ReadFile(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson"))
	if !bytes.HasPrefix(stream, ledger) {
		t.Fatal("original history rewritten")
	}
	for rel, b := range originals {
		after, _ := os.ReadFile(filepath.Join(s.Area, rel))
		if !bytes.Equal(b, after) {
			t.Fatal("original descriptor rewritten", rel)
		}
	}
}
