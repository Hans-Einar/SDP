package governance

import (
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
)

func sessionDesignFixture(t *testing.T) (*Service, SessionDesignRequest) {
	return sessionDesignConfigured(t, false)
}
func sessionDesignConfigured(t *testing.T, required bool) (*Service, SessionDesignRequest) {
	t.Helper()
	q := SessionDesignRequest{Schema: SessionSchema, ModelArea: "Models", WorkRef: "work:Target", Entry: "System.design", TaskPath: "evidence/design-task.json"}
	s, a := assessmentConfigured(t, "", "", func(t *testing.T, s *Service, b *SessionBinding) {
		multiCardImplementation(t, s, b)
		area := filepath.Join(s.area, q.ModelArea)
		if err := os.MkdirAll(area, 0700); err != nil {
			t.Fatal(err)
		}
		w, err := model.CreateWork(area, "Target", "", true)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(w.Path, q.Entry), []byte("language design-core version 0.6.\nsystem Example.\nunit Worker.\nExample contains Worker.\n"), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err = model.Commit(area, q.WorkRef, "Pin fixture target", false); err != nil {
			t.Fatal(err)
		}
		task := blueprint.Task{ID: "session-design-fixture", Intent: "Inspect the Session-owned target before actual owner disposition", AllowedCodePaths: b.Steps[0].WritablePaths}
		raw, _ := json.Marshal(task)
		evidenceFile(t, s, q.TaskPath, string(raw))
		if !required {
			return
		}
		// Blueprint assignment creation requires an already active canonical plan.
		originalLedger := mustFile(t, filepath.Join(s.root, managementLedger))
		events, err := publicationEvents(mustFile(t, filepath.Join(s.root, managementLedger)))
		if err != nil {
			t.Fatal(err)
		}
		lines := []string{}
		for _, ev := range events {
			if ev["subjectId"] == b.Plan.ID {
				ev["payload"].(map[string]any)["to"] = "active"
			}
			raw, _ := json.Marshal(ev)
			lines = append(lines, string(raw))
		}
		activeLedger := strings.Join(lines, "\n") + "\n"
		evidenceFile(t, s, managementLedger, activeLedger)

		bundle, err := blueprints.GenerateRetained(context.Background(), blueprints.Options{Area: area, From: q.WorkRef, To: q.WorkRef, Entry: q.Entry, Task: filepath.Join(s.root, q.TaskPath)}, filepath.Join(s.area, "Blueprints"))
		if err != nil {
			t.Fatal(err)
		}
		var doc blueprints.Document
		if err = json.Unmarshal(mustFile(t, filepath.Join(bundle.Path, "blueprint.json")), &doc); err != nil {
			t.Fatal(err)
		}
		if !doc.Analysis.ConstraintsPass {
			t.Fatal("fixture constraints failed", doc.Analysis)
		}
		rel, _ := filepath.Rel(s.area, bundle.Path)
		req := blueprintstate.Request{Schema: blueprintstate.Schema, ID: "BPA-session-design", EventID: "EVT-BPA-SESSION-DESIGN-1", OccurredAt: "2026-10-18T00:00:00Z", Action: "create", Reason: "Retain exact Session design fixture", Binding: &blueprintstate.Binding{Bundle: filepath.ToSlash(rel), Task: doc.TaskID, System: doc.Analysis.System, Revision: bundle.Revision, Retained: bundle.RetainedRevision, Plan: b.Plan.ID, Milestone: b.Steps[0].Milestone, ModelArea: q.ModelArea, From: q.WorkRef, To: q.WorkRef}}
		if _, err = blueprints.ApplyAssignment(s.area, blueprintstate.Principal{Actor: "fixture-controller", Authority: "controller"}, req); err != nil {
			t.Fatal(err)
		}
		// Restore the pre-selection synthetic plan state, retaining the provider event.
		afterLedger := mustFile(t, filepath.Join(s.root, managementLedger))
		evidenceFile(t, s, managementLedger, string(originalLedger)+string(afterLedger[len(activeLedger):]))
		q.BlueprintAssignment = req.ID
		q.BlueprintEvent = req.EventID
	})
	a.Criteria[0].Impact = "changed"
	if required {
		a.Blueprint.Requirement = "required"
	}
	r := commitAssessment(t, s, a)
	q.InstanceID = r.InstanceID
	q.ExpectedRevision = r.Revision
	q.AssessmentID = a.OperationID
	return s, q
}

func TestSessionDesignContextReadOnlyAndStable(t *testing.T) {
	s, q := sessionDesignFixture(t)
	inventory := inventoryDesign(t, filepath.Dir(s.root))
	_, head, err := s.load()
	if err != nil {
		t.Fatal(err)
	}
	before, err := model.Snapshot(filepath.Join(s.area, q.ModelArea), q.WorkRef)
	if err != nil {
		t.Fatal(err)
	}
	v, err := s.InspectSessionDesign(q)
	if err != nil || v.Status != "current" || v.Basis == nil || v.ReviewRevision == "" || v.DesignAcceptance != "not-recorded" || v.Dispatch != "unavailable" {
		t.Fatal(v, err)
	}
	again, err := s.InspectSessionDesign(q)
	if err != nil || !reflect.DeepEqual(v, again) {
		t.Fatal(again, err)
	}
	_, afterHead, err := s.load()
	if err != nil || head != afterHead {
		t.Fatal("inspection mutated state", err)
	}
	after, err := model.Snapshot(filepath.Join(s.area, q.ModelArea), q.WorkRef)
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("inspection mutated WORK", err)
	}
	if after := inventoryDesign(t, filepath.Dir(s.root)); !reflect.DeepEqual(inventory, after) {
		t.Fatal("inspection wrote files")
	}
	view, err := s.ReadSession(q.InstanceID)
	if err != nil || view.Steps[0].Stage != "design" || view.Steps[0].RunID != "" {
		t.Fatal("inspection bypassed design gate", view, err)
	}
}

func TestSessionDesignContextRejectsChangedOrIncompleteBasis(t *testing.T) {
	for _, kind := range []string{"dirty-work", "invalid-sdl", "wrong-entry", "task-scope", "stale-session", "wrong-assessment", "blueprint-required", "changed-during-inspection", "symlink", "traversal"} {
		t.Run(kind, func(t *testing.T) {
			s, q := sessionDesignFixture(t)
			var during func()
			source := filepath.Join(s.area, "Models", "WORK--Target", "System.design")
			switch kind {
			case "dirty-work":
				f, e := os.OpenFile(source, os.O_APPEND|os.O_WRONLY, 0600)
				if e != nil {
					t.Fatal(e)
				}
				f.WriteString("// Uncommitted source change.\n")
				f.Close()
			case "invalid-sdl":
				os.WriteFile(source, []byte("language design-core version 0.6.\nsystem Example.\nExample contains Missing.\n"), 0600)
				if _, e := model.Commit(filepath.Dir(filepath.Dir(source)), q.WorkRef, "Invalid target fixture", false); e != nil {
					t.Fatal(e)
				}
			case "wrong-entry":
				q.Entry = "Missing.design"
			case "task-scope":
				b, _ := json.Marshal(blueprint.Task{ID: "scope-drift", Intent: "Try a broader scope", AllowedCodePaths: []string{"extra.txt"}})
				evidenceFile(t, s, q.TaskPath, string(b))
			case "stale-session":
				q.ExpectedRevision++
			case "wrong-assessment":
				q.AssessmentID = newTestID()
			case "blueprint-required":
				st, _, _ := s.load()
				a := st.Sessions[q.InstanceID].Assessments[0].Request
				a.ExpectedRevision = q.ExpectedRevision
				a.OperationID = newTestID()
				a.Blueprint.Requirement = "required"
				r := commitAssessment(t, s, a)
				q.ExpectedRevision = r.Revision
				q.AssessmentID = a.OperationID
			case "changed-during-inspection":
				during = func() { evidenceFile(t, s, q.TaskPath, string(mustFile(t, filepath.Join(s.root, q.TaskPath)))+"\n") }
			case "symlink":
				if e := os.Symlink(filepath.Join(s.area, "Models"), filepath.Join(s.area, "Link")); e != nil {
					t.Fatal(e)
				}
				q.ModelArea = "Link"
			case "traversal":
				q.ModelArea = "../Models"
			}
			v, e := s.inspectSessionDesign(q, during)
			if e == nil && (v.Status != "blocked" || len(v.Holds) == 0 || v.Basis != nil || v.ReviewRevision != "") {
				t.Fatal("unsafe design basis accepted", v)
			}
		})
	}
}

func TestSessionDesignContextRequiredBlueprint(t *testing.T) {
	s, q := sessionDesignConfigured(t, true)
	inventory := inventoryDesign(t, filepath.Dir(s.root))
	v, err := s.InspectSessionDesign(q)
	if err != nil || v.Status != "current" || v.Basis == nil || v.Basis.Blueprint == nil {
		t.Fatal(v, err)
	}
	if !reflect.DeepEqual(inventory, inventoryDesign(t, filepath.Dir(s.root))) {
		t.Fatal("Blueprint inspection wrote files")
	}
	for _, kind := range []string{"wrong-event", "task-drift", "work-drift"} {
		t.Run(kind, func(t *testing.T) {
			bad := q
			switch kind {
			case "wrong-event":
				bad.BlueprintEvent = "EVT-WRONG"
			case "task-drift":
				path := filepath.Join(s.root, q.TaskPath)
				before := mustFile(t, path)
				defer os.WriteFile(path, before, 0600)
				if err := os.WriteFile(path, append(before, '\n'), 0600); err != nil {
					t.Fatal(err)
				}
			case "work-drift":
				path := filepath.Join(s.area, q.ModelArea, "WORK--Target", "System.design")
				if err := os.WriteFile(path, append(mustFile(t, path), []byte("// New committed basis.\n")...), 0600); err != nil {
					t.Fatal(err)
				}
				if _, err := model.Commit(filepath.Dir(filepath.Dir(path)), q.WorkRef, "Change fixture basis", false); err != nil {
					t.Fatal(err)
				}
			}
			v, err := s.InspectSessionDesign(bad)
			if err != nil || v.Status != "blocked" || len(v.Holds) == 0 || v.Basis != nil {
				t.Fatal("stale Blueprint accepted", v, err)
			}
		})
	}
}
