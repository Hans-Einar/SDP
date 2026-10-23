package governance

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
)

type designFixture struct {
	s       *Service
	q       DesignCheckRequest
	prepare PrepareRequest
	work    string
	request blueprintstate.Request
}

func providerDesignFixture(t *testing.T) designFixture {
	t.Helper()
	s, _, q := preparationFixture(t)
	put := func(path string, b []byte) {
		t.Helper()
		if err := os.WriteFile(path, b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	// Add the provider's active typed-plan metadata to this synthetic ledger.
	path := filepath.Join(s.area, "ProjectManagement/Ledger.ndjson")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	for i, line := range lines {
		var e map[string]any
		if err = json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatal(err)
		}
		e["actor"], e["occurredAt"], e["commit"] = "fixture-controller", "2026-10-18T00:00:00Z", nil
		e["eventId"] = "EVT-FIXTURE-" + e["eventId"].(string)
		if e["subjectId"] == "PLAN-TEST-0002" {
			p := e["payload"].(map[string]any)
			p["kind"], p["planType"] = "Plan", "ImplementationPlan"
		}
		b, _ := json.Marshal(e)
		lines[i] = string(b)
	}
	put(path, []byte(strings.Join(lines, "\n")+"\n"))
	area := filepath.Join(s.area, "Models")
	if err = os.MkdirAll(area, 0700); err != nil {
		t.Fatal(err)
	}
	works := []model.Result{}
	for _, name := range []string{"Before", "After"} {
		w, err := model.CreateWork(area, name, "", true)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.CopyFS(w.Path, os.DirFS("../../experiments/blueprint_mvp1/NOW")); err != nil {
			t.Fatal(err)
		}
		works = append(works, w)
	}
	taskPath := filepath.Join(s.root, "provider-task.json")
	task, _ := json.Marshal(blueprint.Task{ID: "design-fixture", Intent: "Review unchanged bounded model", AllowedCodePaths: q.Assignment.WritablePaths})
	put(taskPath, task)
	bundle, err := blueprints.GenerateRetained(context.Background(), blueprints.Options{Area: area, From: "work:Before", To: "work:After", Entry: "System.design", Task: taskPath}, filepath.Join(s.area, "Blueprints"))
	if err != nil {
		t.Fatal(err)
	}
	blueprintPath := filepath.Join(bundle.Path, "blueprint.json")
	raw, err = os.ReadFile(blueprintPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc blueprints.Document
	if err = json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	rel, _ := filepath.Rel(s.area, bundle.Path)
	request := blueprintstate.Request{Schema: blueprintstate.Schema, EventID: "EVT-BPA-DESIGN-1", OccurredAt: "2026-10-18T00:00:00Z", ID: "BPA-design-fixture", Action: "create", Reason: "Synthetic provider fixture", Binding: &blueprintstate.Binding{Bundle: filepath.ToSlash(rel), Task: doc.TaskID, System: doc.Analysis.System, Revision: bundle.Revision, Retained: bundle.RetainedRevision, Plan: "PLAN-TEST-0002", Milestone: "M1", ModelArea: "Models", From: "work:Before", To: "work:After"}}
	if _, err = blueprints.ApplyAssignment(s.area, blueprintstate.Principal{Actor: "fixture-controller", Authority: "controller"}, request); err != nil {
		t.Fatal(err)
	}
	paths := []string{filepath.Join(works[0].Path, filepath.Base(works[0].Path)+".yaml"), filepath.Join(works[1].Path, filepath.Base(works[1].Path)+".yaml"), taskPath, blueprintPath, filepath.Join(s.root, "fixture-design-disposition.json")}
	q.Artifacts = nil
	for i, kind := range []string{"now", "target", "task", "blueprint", "design-disposition"} {
		data, err := os.ReadFile(paths[i])
		if err != nil {
			t.Fatal(err)
		}
		rel, _ := filepath.Rel(s.root, paths[i])
		q.Artifacts = append(q.Artifacts, ArtifactPin{Kind: kind, Path: filepath.ToSlash(rel), SHA256: hash(data)})
	}
	if _, err = s.Bind(LocalController()); err != nil {
		t.Fatal(err)
	}
	q = refreshPreparation(t, s, q)
	p, err := s.PrepareAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	_, head, err := s.load()
	if err != nil {
		t.Fatal(err)
	}
	return designFixture{s: s, prepare: q, work: works[1].Path, request: request, q: DesignCheckRequest{Schema: DesignCheckSchema, AssignmentID: p.Assignment.ID, ExpectedRevision: 1, ExpectedGeneration: head, BlueprintAssignment: request.ID, BlueprintEvent: request.EventID}}
}

func inventoryDesign(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			out[path] = "directory"
			return nil
		}
		b, err := os.ReadFile(path)
		if err == nil {
			out[path] = hash(b)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestDesignCheckUsesLiveProvidersWithoutWritesOrApproval(t *testing.T) {
	f := providerDesignFixture(t)
	before := inventoryDesign(t, filepath.Dir(f.s.root))
	a, err := f.s.CheckAssignmentDesign(f.q)
	if err != nil || a.Status != "current" || !digestPattern.MatchString(a.ReviewRevision) || a.Provider == nil {
		t.Fatal(a, err)
	}
	if a.DesignAcceptance != "not-recorded" || a.Dispatch != "unavailable" || !a.Provider.Target.Preliminary || a.Provider.Readiness != "not-adopted" {
		t.Fatal("provider facts promoted to approval/readiness", a)
	}
	b, err := f.s.CheckAssignmentDesign(f.q)
	if err != nil || !reflect.DeepEqual(a, b) || !reflect.DeepEqual(before, inventoryDesign(t, filepath.Dir(f.s.root))) {
		t.Fatal("read check changed result or filesystem", err)
	}
	// Source bytes drift without any metadata edit: file-pin inspection alone misses it.
	path := filepath.Join(f.work, "System.design")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, append(raw, []byte("\n// uncommitted TARGET drift\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	inspected, err := f.s.InspectAssignments(f.q.AssignmentID)
	if err != nil || inspected.Assignments[0].Readiness != "prepared" {
		t.Fatal("metadata pin unexpectedly changed", inspected, err)
	}
	before = inventoryDesign(t, filepath.Dir(f.s.root))
	b, err = f.s.CheckAssignmentDesign(f.q)
	if err != nil || b.Status != "blocked" || b.ReviewRevision != "" || len(b.Conditions) == 0 {
		t.Fatal(b, err)
	}
	if !reflect.DeepEqual(before, inventoryDesign(t, filepath.Dir(f.s.root))) {
		t.Fatal("blocked check wrote files")
	}
}

func TestDesignCheckRefusesWrongPinsScopeAndStaleEvent(t *testing.T) {
	for _, change := range []string{"metadata", "blueprint-copy", "task", "scope", "event", "retired", "hold", "retained-file"} {
		t.Run(change, func(t *testing.T) {
			f := providerDesignFixture(t)
			state, head, err := f.s.load()
			if err != nil {
				t.Fatal(err)
			}
			p := state.Prepared[f.q.AssignmentID]
			switch change {
			case "event":
				f.q.BlueprintEvent = "other-event"
			case "retained-file":
				path := filepath.Join(f.s.area, f.request.Binding.Bundle, "changes.md")
				if err = os.WriteFile(path, []byte("changed retained document"), 0600); err != nil {
					t.Fatal(err)
				}
			case "retired":
				if _, err = f.s.RetireAssignment(LocalController(), RetirePreparationRequest(preparationChange(t, f.s, p))); err != nil {
					t.Fatal(err)
				}
				f.q.ExpectedRevision++
			case "hold":
				r := f.request
				r.Action, r.Expected, r.EventID, r.Binding = "hold", r.EventID, "EVT-BPA-HOLD", nil
				if _, err = blueprints.ApplyAssignment(f.s.area, blueprintstate.Principal{Actor: "fixture-controller", Authority: "controller"}, r); err != nil {
					t.Fatal(err)
				}
				f.q.BlueprintEvent = r.EventID
				// Update only the preparation's coverage via the real revision operation.
				q := revisionRequest(t, f.s, p)
				q.InputRevision = refreshPreparation(t, f.s, f.prepare).InputRevision
				if _, err = f.s.ReviseAssignment(LocalController(), q); err != nil {
					t.Fatal(err)
				}
				f.q.ExpectedRevision++
			default:
				q := revisionRequest(t, f.s, p)
				if change == "scope" {
					q.Assignment.WritablePaths = []string{"work"}
				} else {
					index := map[string]int{"metadata": 0, "blueprint-copy": 3, "task": 2}[change]
					pin := &q.Artifacts[index]
					b, err := os.ReadFile(filepath.Join(f.s.root, pin.Path))
					if err != nil {
						t.Fatal(err)
					}
					if change == "task" {
						b = append(b, '\n')
					}
					pin.Path, pin.SHA256 = "other-pin.json", hash(b)
					if err = os.WriteFile(filepath.Join(f.s.root, pin.Path), b, 0600); err != nil {
						t.Fatal(err)
					}
				}
				if _, err = f.s.ReviseAssignment(LocalController(), q); err != nil {
					t.Fatal(err)
				}
				f.q.ExpectedRevision++
			}
			_, head, err = f.s.load()
			if err != nil {
				t.Fatal(err)
			}
			f.q.ExpectedGeneration = head
			v, err := f.s.CheckAssignmentDesign(f.q)
			if err != nil || v.Status != "blocked" || v.ReviewRevision != "" {
				t.Fatal(v, err)
			}
		})
	}
}

func TestDesignCheckRejectsDriftDuringRead(t *testing.T) {
	for _, what := range []string{"model", "preparation", "coverage"} {
		t.Run(what, func(t *testing.T) {
			f := providerDesignFixture(t)
			v, err := f.s.checkAssignmentDesign(f.q, func() {
				if what == "preparation" {
					state, _, err := f.s.load()
					if err != nil {
						t.Fatal(err)
					}
					q := revisionRequest(t, f.s, state.Prepared[f.q.AssignmentID])
					if _, err = f.s.ReviseAssignment(LocalController(), q); err != nil {
						t.Fatal(err)
					}
					return
				}
				path := filepath.Join(f.work, "System.design")
				if what == "coverage" {
					path = filepath.Join(f.s.root, f.prepare.Selection.Authority)
				}
				file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					t.Fatal(err)
				}
				_, err = file.WriteString("\n// changed during read\n")
				file.Close()
				if err != nil {
					t.Fatal(err)
				}
			})
			if err != nil || v.Status != "blocked" || v.ReviewRevision != "" {
				t.Fatal(v, err)
			}
		})
	}
}

func TestDesignReviewDigestBindsPreparationAndPreservesProviderLimits(t *testing.T) {
	f := providerDesignFixture(t)
	first, err := f.s.CheckAssignmentDesign(f.q)
	if err != nil || first.Status != "current" {
		t.Fatal(first, err)
	}
	q := revisionRequest(t, f.s, first.Assignment)
	q.Assignment.Outcome = "Changed owner scope description"
	revised, err := f.s.ReviseAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.s.CheckAssignmentDesign(f.q)
	code(t, err, "stale_revision")
	f.q.ExpectedRevision = revised.Assignment.Revision
	_, f.q.ExpectedGeneration, err = f.s.load()
	if err != nil {
		t.Fatal(err)
	}
	second, err := f.s.CheckAssignmentDesign(f.q)
	if err != nil || second.Status != "current" || second.ReviewRevision == first.ReviewRevision {
		t.Fatal(second, err)
	}
	// The public provider API checks at most two distinct live bindings. Never
	// manufacture current status when this selected assignment is beyond that budget.
	for _, id := range []string{"BPA-aaa", "BPA-aab"} {
		r := f.request
		b := *r.Binding
		b.Milestone = id
		r.ID, r.EventID, r.Binding = id, "EVT-"+strings.ToUpper(id), &b
		if _, err = blueprints.ApplyAssignment(f.s.area, blueprintstate.Principal{Actor: "fixture-controller", Authority: "controller"}, r); err != nil {
			t.Fatal(err)
		}
	}
	q = revisionRequest(t, f.s, revised.Assignment)
	q.InputRevision = refreshPreparation(t, f.s, f.prepare).InputRevision
	revised, err = f.s.ReviseAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	f.q.ExpectedRevision = revised.Assignment.Revision
	_, f.q.ExpectedGeneration, err = f.s.load()
	if err != nil {
		t.Fatal(err)
	}
	blocked, err := f.s.CheckAssignmentDesign(f.q)
	if err != nil || blocked.Status != "blocked" || !strings.Contains(strings.Join(blocked.Conditions, " "), "budget exhausted") {
		t.Fatal(blocked, err)
	}
}
