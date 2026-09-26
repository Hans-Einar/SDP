package sdptool

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalBoardInventory(t *testing.T) {
	p, e := Discover("../..")
	if e != nil {
		t.Fatal(e)
	}
	nodes, hash, e := BoardNodes(p)
	if e != nil {
		t.Fatal(e)
	}
	count, grouped, refs := 0, 0, 0
	for _, n := range nodes {
		if n.Kind == "card" {
			count++
			if n.Sprint != "" {
				grouped++
			}
			if n.Reference != "" {
				refs++
			}
			if n.Target == nil || n.Target.Revision == "" {
				t.Fatal("missing revision")
			}
		}
	}
	if count < 35 || grouped < 6 || refs == 0 || hash == "" {
		t.Fatalf("cards=%d sprint=%d refs=%d", count, grouped, refs)
	}
}
func TestMovedCardRefreshAndMalformedBoard(t *testing.T) {
	root, r := projectFixture(t)
	r.KanBan = "SDP/KanBan"
	saveRegistration(t, root, r)
	board := filepath.Join(root, r.KanBan)
	for _, folder := range []string{"backlog", "active", "onHold", "completed", "canceled", "superseded", "irrelevant"} {
		os.MkdirAll(filepath.Join(board, folder), 0700)
	}
	os.WriteFile(filepath.Join(board, "board.json"), []byte(`{"schemaVersion":"0.2","projectId":"SDP","namespaces":["SDP"],"ledger":"Ledger.ndjson","profile":"sdp-project-management/0.1"}`), 0600)
	name := "#001--Change--Test.md"
	card := "| Field | Value |\n| --- | --- |\n| id | KB-SDP-001 |\n| CardState | backlog |\n"
	os.WriteFile(filepath.Join(board, "backlog", name), []byte(card), 0600)
	ev := `{"eventId":"one","eventType":"x-kanban:created","subjectId":"KB-SDP-001","payload":{"schemaVersion":"0.2","previousEventId":null,"toPath":"backlog/` + name + `"}}` + "\n"
	ledger := filepath.Join(board, "Ledger.ndjson")
	os.WriteFile(ledger, []byte(ev), 0600)
	p, _ := Discover(root)
	_, before, e := BoardNodes(p)
	if e != nil {
		t.Fatal(e)
	}
	os.Rename(filepath.Join(board, "backlog", name), filepath.Join(board, "active", name))
	os.WriteFile(filepath.Join(board, "active", name), []byte(strings.Replace(card, "backlog", "ready", 1)), 0600)
	if _, _, e = BoardNodes(p); e == nil {
		t.Fatal("stale placement accepted")
	}
	ev += `{"eventId":"two","eventType":"x-kanban:moved","subjectId":"KB-SDP-001","payload":{"schemaVersion":"0.2","previousEventId":"one","toPath":"active/` + name + `"}}` + "\n"
	os.WriteFile(ledger, []byte(ev), 0600)
	_, after, e := BoardNodes(p)
	if e != nil || before == after {
		t.Fatalf("refresh %v", e)
	}
	os.WriteFile(filepath.Join(board, "board.json"), []byte(`{"schemaVersion":"99"}`), 0600)
	nav, e := Navigation(p, "")
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, n := range nav.Nodes {
		if n.ID == "kanban" && n.State == "unavailable" && n.Diagnostic != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("no unsupported diagnostic")
	}
}
func TestSDUIServiceAndOptionalTabs(t *testing.T) {
	root, r := projectFixture(t)
	p, _ := Discover(root)
	tree, e := Navigation(p, "")
	if e != nil || len(tree.Roots) != 3 {
		t.Fatalf("optional tabs %v", e)
	}
	r.SDUI = []Model{{"ui", "SDUI", "page.sdui", "sdui/0.2"}}
	saveRegistration(t, root, r)
	source := filepath.Join(root, "page.sdui")
	os.WriteFile(source, []byte("sdui 0.2;\npage = [<\"# Hello\", button(\"OK\")>];\n"), 0600)
	p, _ = Discover(root)
	nodes, revision, e := UINodes(p)
	if e != nil {
		t.Fatal(e)
	}
	if nodes[0].Target == nil {
		t.Fatal("no preview target")
	}
	out := filepath.Join(t.TempDir(), "ui")
	res, e := UIPreview(context.Background(), p, "ui", "page", out, nodes[0].Target.Revision)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(res.Entry)
	if !strings.Contains(string(b), "OK") {
		t.Fatal(string(b))
	}
	tree, e = Navigation(p, "")
	if e != nil || tree.InventoryRevision == "" || revision == "" {
		t.Fatal(e)
	}
	j, _ := json.Marshal(tree)
	if !strings.Contains(string(j), "sdui-preview") {
		t.Fatal(string(j))
	}
	os.WriteFile(source, []byte("sdui 0.2;\npage = [];\n"), 0600)
	if _, e = UIPreview(context.Background(), p, "ui", "page", out, res.Revision); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatalf("%v", e)
	}
}
