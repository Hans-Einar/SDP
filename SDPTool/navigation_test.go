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
	p, e := Discover("..")
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

func TestKanBanReferenceScope(t *testing.T) {
	for _, tc := range []struct {
		name, primary, target, wantLocal, wantExternal, wantError string
		namespaces                                                []string
		brokenHistory                                             bool
	}{
		{name: "foreign XFMD reference", primary: "KB-SDP-014", wantExternal: "KB-SDP-014"},
		{name: "invalid board namespace", namespaces: []string{"../bad"}, primary: "KB-SDP-014", wantError: "invalid KanBan namespace"},
		{name: "unknown foreign project", primary: "KB-UNKNOWN-999", wantExternal: "KB-UNKNOWN-999"},
		{name: "local reference", primary: "KB-XFMD-001", target: "KB-XFMD-001", wantLocal: "kanban/card/KB-XFMD-001"},
		{name: "missing project reference", primary: "KB-XFMD-999", wantError: "unresolved Ref"},
		{name: "shared namespace reference", namespaces: []string{"XFMD", "SDL"}, primary: "KB-SDL-001", target: "KB-SDL-001", wantLocal: "kanban/card/KB-SDL-001"},
		{name: "missing shared reference", namespaces: []string{"XFMD", "SDL"}, primary: "KB-SDL-999", wantError: "unresolved Ref"},
		{name: "malformed foreign ID", primary: "KB-SDP-nope", wantError: "invalid Ref ID"},
		{name: "path is not ID", primary: "../SDP/KB-SDP-014", wantError: "invalid Ref ID"},
		{name: "URI is not ID", primary: "https://example.com/KB-SDP-014", wantError: "invalid Ref ID"},
		{name: "foreign reference cannot hide broken history", primary: "KB-SDP-014", brokenHistory: true, wantError: "broken card chain"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, r := projectFixture(t)
			r.KanBan = "SDP/KanBan"
			saveRegistration(t, root, r)
			board := filepath.Join(root, r.KanBan)
			write := func(name string, data []byte) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(board, name), data, 0600); err != nil {
					t.Fatal(err)
				}
			}
			for _, folder := range []string{"backlog", "active", "onHold", "completed", "canceled", "superseded", "irrelevant"} {
				if err := os.MkdirAll(filepath.Join(board, folder), 0700); err != nil {
					t.Fatal(err)
				}
			}
			descriptor, err := json.Marshal(map[string]any{"schemaVersion": "0.2", "projectId": "XFMD", "namespaces": tc.namespaces, "ledger": "Ledger.ndjson", "profile": "sdp-project-management/0.2"})
			if err != nil {
				t.Fatal(err)
			}
			write("board.json", descriptor)
			var history []byte
			addCard := func(id, primary string) {
				name := "backlog/#" + id + ".md"
				card := "| Field | Value |\n| --- | --- |\n| id | " + id + " |\n| CardState | backlog |\n"
				if primary != "" {
					card += "| primary | " + primary + " |\n"
				}
				write(name, []byte(card))
				var previous any
				if tc.brokenHistory {
					previous = "missing-event"
				}
				ev, err := json.Marshal(map[string]any{"eventId": "created-" + id, "eventType": "x-kanban:created", "subjectId": id, "payload": map[string]any{"schemaVersion": "0.2", "previousEventId": previous, "toPath": name}})
				if err != nil {
					t.Fatal(err)
				}
				history = append(history, append(ev, '\n')...)
			}
			addCard("KB-XFMD-012", tc.primary)
			if tc.target != "" {
				addCard(tc.target, "")
			}
			write("Ledger.ndjson", history)
			p, err := Discover(root)
			if err != nil {
				t.Fatal(err)
			}
			nav, err := Navigation(p, "")
			if err != nil {
				t.Fatal(err)
			}
			foundTab, foundCard := false, false
			for _, n := range nav.Nodes {
				if n.ID == "kanban" {
					foundTab = true
					if tc.wantError != "" {
						if n.State != "unavailable" || !strings.Contains(n.Diagnostic, tc.wantError) {
							t.Fatalf("bad diagnostic: %+v", n)
						}
					} else if n.State != "validated" {
						t.Fatalf("board unavailable: %+v", n)
					}
				}
				if n.ID == "kanban/card/KB-XFMD-012" {
					foundCard = true
					if n.Reference != tc.wantLocal || n.ExternalReference != tc.wantExternal || n.State != "available" || n.Target == nil || n.Target.Operation != "open" || n.Target.Revision == "" {
						t.Fatalf("bad card: %+v", n)
					}
					if n.Target.Path != filepath.Join(board, "backlog/#KB-XFMD-012.md") {
						t.Fatal("target must open the local card")
					}
					b, err := json.Marshal(n)
					if err != nil {
						t.Fatal(err)
					}
					if tc.wantExternal != "" && (!strings.Contains(string(b), `"externalReference":"`+tc.wantExternal+`"`) || strings.Contains(string(b), `"reference":`)) {
						t.Fatalf("foreign reference must not be a local node link: %s", b)
					}
				}
			}
			if !foundTab || (tc.wantError == "" && !foundCard) {
				t.Fatal("missing navigation nodes")
			}
		})
	}
}
