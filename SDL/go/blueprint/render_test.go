package blueprint

import (
	"strings"
	"testing"
)

func TestRenderProtectionIDsAndEscaping(t *testing.T) {
	a := compile(t, pilot(t, "NOW"))
	task := task()
	task.ID = "[](<unsafe>)"
	task.Intent = "<script>\n~~~"
	task.Protect = []Protection{{"owner/boundary", "task", "MVP1InspectionPilot/container/MachineService", "subtree"}}
	r := run(t, a, a, task)
	docs := Documents(r, task)
	if !strings.Contains(string(docs["context.mmd"]), "PRESERVE") {
		t.Fatal("protection with slash lost")
	}
	if strings.Contains(string(docs["index.md"]), "<script>") || strings.Contains(string(docs["index.md"]), "[](") {
		t.Fatal("authored text not escaped")
	}
}
