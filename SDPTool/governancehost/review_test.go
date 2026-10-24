package governancehost

import (
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"testing"
)

func TestReviewObservationExactIsolatedFinalAnswer(t *testing.T) {
	d := governance.Review{ID: UUID(), ThreadID: "review", Spec: governance.ReviewSpec{HostVersion: "0.164.0"}, Payload: "pinned candidate"}
	input, _ := json.Marshal(map[string]any{"type": "userMessage", "id": "input", "clientId": d.ID, "content": governance.TurnInputs(d.Payload, nil)})
	good := Thread{ID: d.ThreadID, Cwd: "/project", CLIVersion: d.Spec.HostVersion, Turns: []Turn{{ID: "turn", Status: "completed", ItemsView: "full", Items: []json.RawMessage{input, json.RawMessage(`{"id":"private","type":"reasoning","text":"private"}`), json.RawMessage(`{"id":"commentary","type":"agentMessage","phase":"commentary","text":"not the verdict"}`), json.RawMessage(`{"id":"answer","type":"agentMessage","phase":"final_answer","text":"final verdict"}`)}}}}
	id, status, answer, e := ReviewObservation(good, d, "/project", "0.164.0")
	if e != nil || id != "turn" || status != "completed" || answer != "final verdict" {
		t.Fatal(id, status, answer, e)
	}
	for _, mode := range []string{"parent", "extra-turn", "version", "thread", "cwd", "partial", "input"} {
		b, _ := json.Marshal(good)
		var v Thread
		json.Unmarshal(b, &v)
		switch mode {
		case "parent":
			v.ParentThreadID = "implementation"
		case "extra-turn":
			v.Turns = append(v.Turns, v.Turns[0])
		case "version":
			v.CLIVersion = "0.165.0"
		case "thread":
			v.ID = "other"
		case "cwd":
			v.Cwd = "/other"
		case "partial":
			v.Turns[0].ItemsView = "summary"
		case "input":
			v.Turns[0].Items[0] = json.RawMessage(`{"id":"other","type":"userMessage","content":[]}`)
		}
		if _, _, _, e = ReviewObservation(v, d, "/project", "0.164.0"); e == nil {
			t.Fatal("accepted", mode)
		}
	}
}
