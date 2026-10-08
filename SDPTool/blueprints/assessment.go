package blueprints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

const EvidenceSchema = "sdp-blueprint-evidence/1"
const AssessmentSchema = "sdp-blueprint-assessment/1"

type CodeSnapshot struct {
	Revision string            `json:"revision"`
	Root     string            `json:"root"`
	Digest   string            `json:"digest"`
	Files    map[string]string `json:"files"`
}
type Mapping struct {
	Side       string `json:"side"`
	Element    string `json:"element"`
	Role       string `json:"role"`
	Path       string `json:"path"`
	Symbol     string `json:"symbol"`
	FileDigest string `json:"fileDigest"`
}
type CheckRequirement struct {
	ID         string `json:"id"`
	Side       string `json:"side"`
	Obligation string `json:"obligation"`
}
type CheckReceipt struct {
	ID             string            `json:"id"`
	Side           string            `json:"side"`
	Obligation     string            `json:"obligation"`
	SourceDigest   string            `json:"sourceDigest"`
	TaskDigest     string            `json:"taskDigest"`
	CodeDigest     string            `json:"codeDigest"`
	Command        string            `json:"command"`
	Environment    string            `json:"environment"`
	InputsDigest   string            `json:"inputsDigest"`
	Inputs         map[string]string `json:"inputs"`
	Result         string            `json:"result"`
	Artifact       string            `json:"artifact"`
	ArtifactDigest string            `json:"artifactDigest"`
	Traceability   string            `json:"traceability"`
}
type Disposition struct {
	Unknown      string `json:"unknown"`
	Actor        string `json:"actor"`
	Role         string `json:"role"`
	Scope        string `json:"scope"`
	Rationale    string `json:"rationale"`
	Traceability string `json:"traceability"`
}
type Evidence struct {
	Schema            string                  `json:"schema"`
	BlueprintRevision string                  `json:"blueprintRevision"`
	RetainedRevision  string                  `json:"retainedRevision"`
	TaskDigest        string                  `json:"taskDigest"`
	Scope             string                  `json:"scope"`
	Code              map[string]CodeSnapshot `json:"code"`
	Mappings          []Mapping               `json:"mappings"`
	RequiredChecks    []CheckRequirement      `json:"requiredChecks"`
	Receipts          []CheckReceipt          `json:"receipts"`
	Dispositions      []Disposition           `json:"dispositions"`
}
type UnknownAssessment struct {
	ID          string       `json:"id"`
	Detail      string       `json:"detail"`
	Disposition *Disposition `json:"disposition,omitempty"`
}
type CheckAssessment struct {
	ID        string `json:"id"`
	Side      string `json:"side"`
	Reported  string `json:"reported"`
	Effective string `json:"effective"`
}
type Assessment struct {
	Schema            string              `json:"schema"`
	Operation         string              `json:"operation"`
	Status            string              `json:"status"`
	Identity          string              `json:"identity"`
	BlueprintRevision string              `json:"blueprintRevision"`
	RetainedRevision  string              `json:"retainedRevision"`
	EvidenceDigest    string              `json:"evidenceDigest"`
	Scope             string              `json:"scope"`
	SourceFreshness   string              `json:"sourceFreshness"`
	Authority         string              `json:"authority"`
	Blockers          []string            `json:"blockers"`
	Unknowns          []UnknownAssessment `json:"unknowns"`
	Checks            []CheckAssessment   `json:"checks"`
	Observed          map[string]string   `json:"observed"`
}

// Assess reads bounded, pinned evidence as data. It never executes commands,
// mutates the retained bundle, or grants assignment/implementation authority.
func Assess(bundlePath, evidencePath string) (Assessment, error) {
	return assess(bundlePath, evidencePath, nil)
}

func assess(bundlePath, evidencePath string, beforeFinalize func()) (Assessment, error) {
	var out Assessment
	d, bundle, err := verify(bundlePath, &budget{})
	if err != nil {
		return out, err
	}
	retained := model.Digest(model.Files(bundle.Files))
	if filepath.Base(bundlePath) != retained || filepath.Base(filepath.Dir(bundlePath)) != key(d) {
		return out, fmt.Errorf("assessment requires an intact retained revision")
	}
	b := &budget{}
	raw, err := b.read(evidencePath, 1<<20)
	if err != nil {
		return out, err
	}
	var in Evidence
	if err = decodeRecord(raw, &in); err != nil {
		return out, err
	}
	if err = exactFields(raw, reflect.TypeOf(in)); err != nil {
		return out, err
	}
	if in.Schema != EvidenceSchema || in.BlueprintRevision != d.Revision || in.RetainedRevision != retained || in.TaskDigest != d.TaskDigest || strings.TrimSpace(in.Scope) == "" {
		return out, fmt.Errorf("unsupported evidence or stale blueprint/task binding")
	}
	if len(in.Mappings) > 4000 || len(in.RequiredChecks) > 1000 || len(in.Receipts) > 1000 || len(in.Dispositions) > 10000 {
		return out, fmt.Errorf("evidence record limit")
	}
	out = Assessment{Schema: AssessmentSchema, Operation: "assess-blueprint", Status: "blocked", BlueprintRevision: d.Revision, RetainedRevision: retained, EvidenceDigest: documents.Hash(raw), Scope: in.Scope, SourceFreshness: "not-evaluated", Authority: "attributed-only; no execution or lifecycle permission", Blockers: []string{}, Unknowns: []UnknownAssessment{}, Checks: []CheckAssessment{}, Observed: map[string]string{}}
	unknown := map[string]string{}
	addUnknown := func(id, detail string) { unknown[id] = detail }
	for _, u := range d.Analysis.Unknowns {
		v, _ := json.Marshal(u)
		addUnknown("analysis/"+documents.Hash(v), u.Side+" "+u.Code+" "+u.Subject+": "+u.Detail)
	}
	addUnknown("coverage/code-conformance", "Code mappings and reported tests do not prove complete implementation conformance")
	addUnknown("coverage/runtime-dependencies", "Runtime dependency coverage is not established by this assessment")
	addUnknown("coverage/sdui-semantics", "Semantic SDUI conformance is outside this analyzer")
	if !d.Analysis.ConstraintsPass {
		out.Blockers = append(out.Blockers, "structural constraints failed")
	}
	obligations := map[string]bool{}
	for _, ob := range d.Analysis.Obligations {
		obligations[ob.ID] = true
		if ob.Result != "pass" {
			out.Blockers = append(out.Blockers, "obligation "+ob.ID+": "+ob.Result)
		}
	}
	root := filepath.Dir(evidencePath)
	// Missing files are explicit unknowns. Unsafe files and exhausted budgets fail closed.
	capturedReads := map[string][]byte{}
	read := func(rel string) ([]byte, error) {
		if !safeRelative(rel) {
			return nil, fmt.Errorf("unsafe evidence path %q", rel)
		}
		if value, ok := capturedReads[rel]; ok {
			return value, nil
		}
		data, e := b.read(filepath.Join(root, filepath.FromSlash(rel)), MaxCatalogueBytes)
		if os.IsNotExist(e) {
			out.Observed[rel] = "missing"
			capturedReads[rel] = nil
			return nil, nil
		}
		if e != nil {
			return nil, e
		}
		capturedReads[rel] = data
		out.Observed[rel] = documents.Hash(data)
		return data, nil
	}
	sources := map[string]string{"NOW": d.Now.SourceDigest, "TARGET": d.Target.SourceDigest}
	actual := map[string]map[string]string{}
	codeOK := map[string]bool{}
	for side := range in.Code {
		if _, ok := sources[side]; !ok {
			return out, fmt.Errorf("invalid code side")
		}
	}
	for _, side := range []string{"NOW", "TARGET"} {
		snap, ok := in.Code[side]
		actual[side] = map[string]string{}
		if !ok {
			addUnknown("code/"+side, "No scoped code snapshot supplied")
			continue
		}
		if snap.Revision == "" || !safeRelative(snap.Root) || !hashPattern.MatchString(snap.Digest) || len(snap.Files) == 0 || len(snap.Files) > 4000 {
			return out, fmt.Errorf("invalid code snapshot")
		}
		captured := model.Files{}
		paths := []string{}
		for p := range snap.Files {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		complete := true
		for _, p := range paths {
			if !safeRelative(p) || !hashPattern.MatchString(snap.Files[p]) {
				return out, fmt.Errorf("invalid code file")
			}
			rel := snap.Root + "/" + p
			data, e := read(rel)
			if e != nil {
				return out, e
			}
			if out.Observed[rel] == "missing" {
				complete = false
				continue
			}
			captured[p] = data
			actual[side][p] = documents.Hash(data)
			if actual[side][p] != snap.Files[p] {
				complete = false
			}
		}
		codeOK[side] = complete && model.Digest(captured) == snap.Digest
		if !codeOK[side] {
			addUnknown("code/"+side, "Scoped code capture is missing or stale")
		}
	}
	selected := map[string]map[string]bool{"NOW": {}, "TARGET": {}}
	for _, n := range d.Analysis.Nodes {
		if n.Change != "added" {
			selected["NOW"][n.Key] = true
		}
		if n.Change != "removed" {
			selected["TARGET"][n.Key] = true
		}
	}
	mapped := map[string]bool{}
	seenMappings := map[string]bool{}
	for _, m := range in.Mappings {
		if !selected[m.Side][m.Element] || m.Role == "" || m.Symbol == "" || !safeRelative(m.Path) || !hashPattern.MatchString(m.FileDigest) {
			return out, fmt.Errorf("invalid mapping")
		}
		mk := m.Side + "/" + m.Element
		roleKey := mk + "/" + m.Role
		if seenMappings[roleKey] {
			addUnknown("mapping/"+mk, "Ambiguous mapping role")
		}
		seenMappings[roleKey] = true
		mapped[mk] = true
		if !codeOK[m.Side] || actual[m.Side][m.Path] != m.FileDigest {
			addUnknown("mapping/"+mk, "Stale or missing code locator")
		}
	}
	for _, side := range []string{"NOW", "TARGET"} {
		for element := range selected[side] {
			mk := side + "/" + element
			if !mapped[mk] {
				addUnknown("mapping/"+mk, "No code locator for selected element")
			}
		}
	}
	required := map[string]CheckRequirement{}
	for _, req := range in.RequiredChecks {
		if req.ID == "" || sources[req.Side] == "" || (!obligations[req.Obligation] && !strings.HasPrefix(req.Obligation, "behavior:")) || req.Obligation == "behavior:" {
			return out, fmt.Errorf("invalid required check")
		}
		if _, exists := required[req.ID]; exists {
			return out, fmt.Errorf("duplicate required check")
		}
		required[req.ID] = req
	}
	receipts := map[string]CheckReceipt{}
	for _, r := range in.Receipts {
		req, ok := required[r.ID]
		if !ok || r.Side != req.Side || r.Obligation != req.Obligation {
			return out, fmt.Errorf("receipt outside declared check scope")
		}
		if _, exists := receipts[r.ID]; exists {
			return out, fmt.Errorf("ambiguous duplicate receipt")
		}
		if r.Command == "" || r.Environment == "" || r.Traceability == "" || !hashPattern.MatchString(r.InputsDigest) {
			return out, fmt.Errorf("incomplete check provenance")
		}
		switch r.Result {
		case "pass", "fail", "not-run", "unavailable", "stale":
		default:
			return out, fmt.Errorf("invalid receipt result")
		}
		receipts[r.ID] = r
	}
	ids := []string{}
	for id := range required {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		req := required[id]
		r, ok := receipts[id]
		c := CheckAssessment{ID: id, Side: req.Side, Reported: "not-run", Effective: "not-run"}
		if ok {
			c.Reported = r.Result
			c.Effective = r.Result
			if r.SourceDigest != sources[r.Side] || r.TaskDigest != d.TaskDigest || !codeOK[r.Side] || r.CodeDigest != in.Code[r.Side].Digest {
				c.Effective = "stale"
			}
			inputs := model.Files{}
			names := []string{}
			for name := range r.Inputs {
				names = append(names, name)
			}
			sort.Strings(names)
			if len(names) > 1000 {
				return out, fmt.Errorf("check input limit")
			}
			for _, name := range names {
				if !hashPattern.MatchString(r.Inputs[name]) {
					return out, fmt.Errorf("invalid check input hash")
				}
				data, e := read(name)
				if e != nil {
					return out, e
				}
				if out.Observed[name] != r.Inputs[name] {
					c.Effective = "stale"
				}
				if out.Observed[name] != "missing" {
					inputs[name] = data
				}
			}
			if model.Digest(inputs) != r.InputsDigest {
				c.Effective = "stale"
			}
			if (r.Result == "pass" || r.Result == "fail") && len(names) == 0 {
				return out, fmt.Errorf("observed result requires pinned fixture/config inputs")
			}
			if r.Artifact != "" {
				if !hashPattern.MatchString(r.ArtifactDigest) {
					return out, fmt.Errorf("invalid evidence hash")
				}
				_, e := read(r.Artifact)
				if e != nil {
					return out, e
				}
				if out.Observed[r.Artifact] != r.ArtifactDigest {
					c.Effective = "stale"
				}
			} else if r.Result == "pass" || r.Result == "fail" {
				return out, fmt.Errorf("observed result requires evidence artifact")
			}
			if r.Result == "fail" {
				out.Blockers = append(out.Blockers, "failed check "+id)
			}
		}
		if c.Effective != "pass" && c.Effective != "fail" {
			addUnknown("check/"+id, "Required check is "+c.Effective)
		}
		out.Checks = append(out.Checks, c)
	}
	dispositions := map[string]Disposition{}
	for _, disp := range in.Dispositions {
		if _, exists := unknown[disp.Unknown]; !exists {
			return out, fmt.Errorf("disposition names unknown outside this assessment: %s", disp.Unknown)
		}
		if _, exists := dispositions[disp.Unknown]; exists {
			return out, fmt.Errorf("duplicate disposition")
		}
		if strings.TrimSpace(disp.Actor) == "" || (disp.Role != "owner" && disp.Role != "reviewer") || disp.Scope != in.Scope || strings.TrimSpace(disp.Rationale) == "" || disp.Traceability == "" {
			return out, fmt.Errorf("incomplete or out-of-scope disposition")
		}
		dispositions[disp.Unknown] = disp
	}
	ids = nil
	for id := range unknown {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		u := UnknownAssessment{ID: id, Detail: unknown[id]}
		if disp, ok := dispositions[id]; ok {
			u.Disposition = &disp
		} else {
			out.Blockers = append(out.Blockers, "undisposed unknown "+id)
		}
		out.Unknowns = append(out.Unknowns, u)
	}
	sort.Strings(out.Blockers)
	if len(out.Blockers) == 0 {
		out.Status = "ready"
	}
	if beforeFinalize != nil {
		beforeFinalize()
	}
	current, e := b.read(evidencePath, 1<<20)
	if e != nil {
		return out, e
	}
	if !bytes.Equal(current, raw) {
		return out, fmt.Errorf("evidence changed during assessment")
	}
	observedPaths := []string{}
	for p := range out.Observed {
		observedPaths = append(observedPaths, p)
	}
	sort.Strings(observedPaths)
	for _, p := range observedPaths {
		current, e := b.read(filepath.Join(root, filepath.FromSlash(p)), MaxCatalogueBytes)
		if os.IsNotExist(e) && out.Observed[p] == "missing" {
			continue
		}
		if e != nil {
			return out, e
		}
		if documents.Hash(current) != out.Observed[p] {
			return out, fmt.Errorf("input changed during assessment: %s", p)
		}
	}
	encoded, _ := json.Marshal(out)
	out.Identity = documents.Hash(encoded)
	return out, nil
}

// Struct fields require their exact wire spelling. Dynamic map keys (paths,
// sides, identifiers) retain their case; Go's case-insensitive struct matching
// must not allow a conflicting Result/result pair to override a failed receipt.
func exactFields(raw json.RawMessage, typ reflect.Type) error {
	switch typ.Kind() {
	case reflect.Struct:
		var object map[string]json.RawMessage
		if err := json.Unmarshal(raw, &object); err != nil {
			return err
		}
		fields := map[string]reflect.Type{}
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			name := strings.Split(f.Tag.Get("json"), ",")[0]
			if name == "" {
				name = f.Name
			}
			fields[name] = f.Type
		}
		for name, child := range object {
			t, ok := fields[name]
			if !ok {
				return fmt.Errorf("unknown or noncanonical JSON field %q", name)
			}
			if err := exactFields(child, t); err != nil {
				return err
			}
		}
	case reflect.Map:
		var object map[string]json.RawMessage
		if err := json.Unmarshal(raw, &object); err != nil {
			return err
		}
		for _, child := range object {
			if err := exactFields(child, typ.Elem()); err != nil {
				return err
			}
		}
	case reflect.Slice:
		var children []json.RawMessage
		if err := json.Unmarshal(raw, &children); err != nil {
			return err
		}
		for _, child := range children {
			if err := exactFields(child, typ.Elem()); err != nil {
				return err
			}
		}
	}
	return nil
}
