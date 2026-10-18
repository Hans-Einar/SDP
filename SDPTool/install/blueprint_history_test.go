package install

import (
	"bytes"
	"os"
	"testing"
)

func TestInstallerPreservesBlueprintHistory(t *testing.T) {
	raw, e := os.ReadFile("../blueprintstate/testdata/lifecycle.ndjson")
	if e != nil {
		t.Fatal(e)
	}
	files := map[string][]byte{"SDP/05--Implementation/Test.md": []byte("| Field | Value |\n| --- | --- |\n| id | PLAN-SDP-0001 |\n| state | active |\n| PlanType | ImplementationPlan |\n| BranchPolicy | current |\n| CommitPolicy | milestone |\n")}
	if e = validateHistory(raw, files); e != nil {
		t.Fatal(e)
	}
	bad := bytes.Replace(raw, []byte(`"authority":"reviewer"`), []byte(`"authority":"assignee"`), 1)
	if e = validateHistory(bad, files); e == nil {
		t.Fatal("invalid blueprint history accepted")
	}
}
