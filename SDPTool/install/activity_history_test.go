package install

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

type emptyActivityInventory struct{}

func (emptyActivityInventory) InspectExecutions([]activities.ExecutionStore) (activities.ExecutionSnapshot, error) {
	return activities.ExecutionSnapshot{Known: true}, nil
}
func TestActivityHistoryUsesDomainReplayDuringInstallation(t *testing.T) {
	root := t.TempDir()
	area := filepath.Join(root, "SDP")
	if err := os.MkdirAll(filepath.Join(area, "ProjectManagement"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(area, "authority.md"), []byte("Owner authority"), 0600); err != nil {
		t.Fatal(err)
	}
	service := activities.New(area, "TEST")
	caller := activities.Caller{Actor: "fixture", Authority: activities.Authority{Role: "project-owner", Source: "authority.md"}}
	before, err := projecthistory.Read(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create(activities.CreateRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: before.Revision, Title: "Need", Owner: "owner", Mandate: "Need", Reason: "Owner request", Scope: []string{"Source"}}, caller)
	if err != nil {
		t.Fatal(err)
	}
	session, err := service.CreateSession(activities.CreateSessionRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: created.HistoryRevision, ActivityID: created.ActivityID, ExpectedActivityEventID: created.EventID, ExpectedDescriptorRevision: created.DescriptorRevision, Title: "Study", Objective: "Study this need", Phase: "04-Design", Reason: "Owner selects preparation"}, caller)
	if err != nil {
		t.Fatal(err)
	}
	other, err := service.Create(activities.CreateRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: session.HistoryRevision, Title: "Next", Owner: "owner", Mandate: "Next need", Reason: "Owner preparation", Scope: []string{"OtherSource"}}, caller)
	if err != nil {
		t.Fatal(err)
	}
	active, err := service.Transition(activities.TransitionRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: other.HistoryRevision, ActivityID: created.ActivityID, ExpectedActivityEventID: session.EventID, ExpectedDescriptorRevision: created.DescriptorRevision, Action: "activate", Reason: "Owner selects work"}, caller, nil)
	if err != nil {
		t.Fatal(err)
	}
	adopted, err := service.AdoptAdmission(activities.AdoptionRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: active.HistoryRevision, Stores: []activities.ExecutionStore{{StatePath: t.TempDir()}}}, caller, emptyActivityInventory{})
	if err != nil {
		t.Fatal(err)
	}
	switched, err := service.Transition(activities.TransitionRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: adopted.Catalog.HistoryRevision, ActivityID: created.ActivityID, ExpectedActivityEventID: active.EventID, ExpectedDescriptorRevision: created.DescriptorRevision, Action: "switch", Reason: "Owner changes selected need", ResumeCondition: "Owner selects predecessor again", TargetActivityID: other.ActivityID, ExpectedTargetEventID: other.EventID, ExpectedTargetDescriptorRevision: other.DescriptorRevision}, caller, emptyActivityInventory{})
	if err != nil || !strings.HasSuffix(switched.EventID, "-TARGET") {
		t.Fatal(switched, err)
	}
	// A later record must also accept the switch TARGET identity as predecessor.
	if _, err = service.CreateSession(activities.CreateSessionRequest{OperationID: activityOperation(t), ExpectedHistoryRevision: switched.HistoryRevision, ActivityID: other.ActivityID, ExpectedActivityEventID: switched.EventID, ExpectedDescriptorRevision: other.DescriptorRevision, Title: "Next study", Objective: "Study selected successor", Phase: "02-Study", Reason: "Owner prepares next step"}, caller); err != nil {
		t.Fatal(err)
	}
	// Validate actual Go-generated IDs/pins, including alphabetic UUID digits,
	// Session predecessors and the atomic switch TARGET event, through Python.
	for _, name := range []string{"ledger-event.schema.json", "activity-descriptor-0.1.schema.json", "activity-descriptor-0.2.schema.json", "activity-payload-0.1.schema.json", "activity-payload-0.2.schema.json", "activity-admission-payload-0.1.schema.json"} {
		data, err := os.ReadFile(filepath.Join("../../Template/sdp-root/ProjectManagement", name))
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(area, "ProjectManagement", name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if output, err := exec.Command("python3", "../../Toolkit/scripts/activity_history.py", area).CombinedOutput(); err != nil {
		t.Fatalf("actual Go/Python Activity history parity: %v %s", err, output)
	}
	files := map[string][]byte{}
	if err = filepath.WalkDir(area, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		b, err := os.ReadFile(path)
		files[filepath.ToSlash(rel)] = b
		return err
	}); err != nil {
		t.Fatal(err)
	}
	history := files["SDP/ProjectManagement/Ledger.ndjson"]
	if err = ValidateManagementHistory(history, files); err != nil {
		t.Fatal("Activity reinterpreted as management", err)
	}
	if err = validateHistory(history, files); err != nil {
		t.Fatal("installation rejected supported Activity history", err)
	}
	for path := range files {
		if filepath.Base(path) == "activity.yaml" {
			files[path] = []byte("schema: unsupported\n")
			break
		}
	}
	if err = ValidateManagementHistory(history, files); err == nil {
		t.Fatal("invalid Activity silently ignored")
	}
}

func activityOperation(t *testing.T) string {
	t.Helper()
	for {
		var data [16]byte
		if _, err := rand.Read(data[:]); err != nil {
			t.Fatal(err)
		}
		data[6] = data[6]&0x0f | 0x40
		data[8] = data[8]&0x3f | 0x80
		id := fmt.Sprintf("%x-%x-%x-%x-%x", data[:4], data[4:6], data[6:8], data[8:10], data[10:])
		if strings.ContainsAny(id, "abcdef") {
			return id
		}
	}
}
