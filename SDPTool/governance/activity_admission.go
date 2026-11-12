package governance

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

func (s *Service) activityService() *activities.Service { return activities.New(s.area, "") }

// lockActivityState serializes opt-in adoption with admissions, acquiring canonical
// history before operational state even when policy is not yet enabled.
func (s *Service) lockActivityState() (*projecthistory.Transaction, func(), error) {
	tx, err := s.canonicalTransaction()
	if err != nil {
		if errors.Is(err, projecthistory.ErrBusy) {
			return nil, nil, fail("busy", err.Error())
		}
		return nil, nil, err
	}
	if err = tx.Check(""); err != nil {
		tx.Close()
		return nil, nil, err
	}
	release, err := lock(s.state)
	if err != nil {
		tx.Close()
		return nil, nil, err
	}
	return tx, func() { release(); tx.Close() }, nil
}

func (s *Service) checkActivity(tx *projecthistory.Transaction, owner string, sel Selection, a Assignment, settling bool) (*activities.AdmissionBinding, error) {
	if tx == nil {
		return nil, nil
	}
	service := s.activityService()
	view, err := service.ReadAdmission(tx, owner)
	if err != nil {
		return nil, err
	}
	if !view.Enabled {
		return nil, nil
	}
	// Only settlement of an exact already-sending input may cross a pending gate.
	// The caller proves that claim in its locked PG snapshot.
	if settling {
		view.Pending = nil
	}
	binding, err := activities.CheckAdmission(view, activities.WorkMembership{SessionPath: sel.Session, PlanPath: sel.Plan, CardPaths: []string{sel.Card}, WritablePaths: a.WritablePaths})
	if err != nil {
		return nil, fail("activity_admission", err.Error())
	}
	if !settling {
		if err = service.RegisterExecutionStore(tx, owner, activities.ExecutionStore{StatePath: s.state}); err != nil {
			return nil, err
		}
	}
	return &binding, nil
}

// CheckActivityDispatch revalidates a retained exact sending claim immediately
// before the host RPC. It grants no new send and does not change PG input state.
// A pending hold counts this claim as in-flight until existing host recovery settles it.
func (s *Service) CheckActivityDispatch(c Caller, runID, inputID string) error {
	if !c.controller {
		return fail("unauthorized", "trusted controller required")
	}
	tx, release, err := s.lockActivityState()
	if err != nil {
		return err
	}
	defer release()
	if tx == nil {
		return s.CheckRunReady(runID)
	}
	st, _, err := s.load()
	if err != nil {
		return err
	}
	r, ok := st.Runs[runID]
	if !ok {
		return fail("missing_coverage", "run absent")
	}
	claimed := false
	for _, in := range r.Inputs {
		if in.ID == inputID && in.State == "sending" && (in.Kind == "launch" || in.Kind == "steer" || in.Kind == "approval") {
			claimed = true
		}
	}
	if !claimed {
		return fail("reconciliation_required", "exact sending input required")
	}
	binding, err := s.checkActivity(tx, "", r.Selection, r.Assignment, true)
	if err != nil {
		return err
	}
	if binding == nil {
		return s.CheckRunReady(runID)
	}
	if r.Activity == nil || *r.Activity != *binding {
		return fail("activity_admission", "run lacks matching Activity admission binding")
	}
	return s.CheckRunReady(runID)
}

type activityExecutionInspector struct{ root, area, project string }

// ActivityExecutionInspector reads authoritative registered PG stores. Its caller
// holds canonical history first; this method acquires only each store's state lock.
func ActivityExecutionInspector(root, area, project string) activities.ExecutionInspector {
	return activityExecutionInspector{root, area, project}
}
func (i activityExecutionInspector) InspectExecutions(stores []activities.ExecutionStore) (activities.ExecutionSnapshot, error) {
	out := activities.ExecutionSnapshot{Known: true, InFlight: []activities.ExecutionRef{}}
	if len(stores) == 0 || len(stores) > 128 {
		return activities.ExecutionSnapshot{}, fmt.Errorf("explicit bounded execution inventory required")
	}
	for _, store := range stores {
		if _, err := os.Stat(filepath.Join(store.StatePath, "HEAD")); err != nil {
			return activities.ExecutionSnapshot{}, fmt.Errorf("execution store unavailable: %w", err)
		}
		s, err := New(i.root, i.area, i.project, store.StatePath)
		if err != nil {
			return activities.ExecutionSnapshot{}, err
		}
		release, err := lock(s.state)
		if err != nil {
			return activities.ExecutionSnapshot{}, err
		}
		st, _, err := s.load()
		release()
		if err != nil {
			return activities.ExecutionSnapshot{}, err
		}
		for _, r := range st.Runs {
			reason := ""
			if r.State == "running" && len(r.Attempts) == 0 {
				reason = "running routine has no observed host disposition"
			}
			for _, in := range r.Inputs {
				if in.State == "sending" || in.State == "uncertain" {
					reason = "outbound host input " + in.State
				}
			}
			for _, a := range r.Attempts {
				if a.Status == "running" || a.Status == "unknown" {
					reason = "host attempt active or unresolved"
				}
			}
			if reason != "" {
				activityID := ""
				if r.Activity != nil {
					activityID = r.Activity.ActivityID
				}
				out.InFlight = append(out.InFlight, activities.ExecutionRef{ActivityID: activityID, StorePath: store.StatePath, RunID: r.ID, Reason: reason})
			}
		}
		for _, cmd := range st.SessionHostCommands {
			if cmd.State == "claimed" || cmd.State == "uncertain" {
				out.InFlight = append(out.InFlight, activities.ExecutionRef{StorePath: store.StatePath, RunID: cmd.Request.RunID, Reason: "host command " + cmd.State})
			}
		}
	}
	return out, nil
}

// validateActivitySessionReference extends document ownership only. Execution
// bindings, claims and operational instances retain their existing canonical home.
func (s *Service) validateActivitySessionReference(ref SessionRecordRef) error {
	if !strings.HasPrefix(ref.Path, "SDP/Activities/") {
		return nil
	}
	namespace, err := activities.ResolveProjectID(s.area, "")
	if err != nil {
		return err
	}
	raw, err := readFile(filepath.Join(s.root, managementLedger), projecthistory.MaxBytes)
	if err != nil {
		return err
	}
	catalog, err := activities.Replay(s.area, namespace, projecthistory.Snapshot{Bytes: raw, Revision: hash(raw)})
	if err != nil {
		return err
	}
	members := map[string]bool{}
	for _, a := range catalog.Activities {
		for _, session := range a.Sessions {
			if session.ID == ref.ID && "SDP/"+session.Path == ref.Path {
				members[a.Descriptor.ID] = true
			}
		}
		for _, linked := range a.Descriptor.References {
			if linked.Role == "session" && linked.Availability != "unavailable" && (linked.ProjectID == "" || linked.ProjectID == namespace) && linked.ID == ref.ID && "SDP/"+linked.Path == ref.Path {
				members[a.Descriptor.ID] = true
			}
		}
	}
	if len(members) != 1 {
		return fail("activity_admission", "Activity-local Session must have exactly one catalog membership")
	}
	return nil
}
