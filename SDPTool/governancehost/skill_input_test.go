package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestSkillTurnProtocolAndHistoryCorrelation(t *testing.T) {
	v := governance.OwnerInput{ID: UUID(), Payload: "$sdp $sdp-master\nexact", Skills: []governance.SkillInput{{Name: "sdp", Path: "/project/Skills/sdp/SKILL.md", Digest: strings.Repeat("a", 64)}, {Name: "sdp-master", Path: "/project/Skills/sdp-master/SKILL.md", Digest: strings.Repeat("b", 64)}}}
	left, right := net.Pipe()
	c := newClient(left, left)
	defer c.Close()
	defer right.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	captured := make(chan json.RawMessage, 1)
	go func() {
		s := bufio.NewScanner(right)
		if !s.Scan() {
			return
		}
		var f frame
		json.Unmarshal(s.Bytes(), &f)
		captured <- f.Params
		json.NewEncoder(right).Encode(map[string]any{"id": f.ID, "result": map[string]any{"turn": map[string]string{"id": "turn", "status": "inProgress"}}})
	}()
	if _, e := c.StartSkillTurn(ctx, "thread", v.ID, v.Payload, v.Skills); e != nil {
		t.Fatal(e)
	}
	raw := <-captured
	var q struct {
		Input    json.RawMessage `json:"input"`
		ClientID string          `json:"clientUserMessageId"`
	}
	json.Unmarshal(raw, &q)
	if q.ClientID != v.ID || governance.ObservedInputDigest(q.Input) != v.WireDigest() {
		t.Fatal(string(raw))
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.Marshal([]map[string]any{{"definition": "TurnStartParams", "value": raw}})
		os.WriteFile(filepath.Join(dir, "skill-rpc-cases.json"), b, 0600)
	}
	for _, name := range []string{"exact", "changed-name", "changed-path", "missing", "reordered", "extra", "duplicate", "summary", "other-summary", "text-elements"} {
		t.Run(name, func(t *testing.T) {
			content := governance.TurnInputs(v.Payload, v.Skills)
			switch name {
			case "changed-name":
				content[1]["name"] = "other"
			case "changed-path":
				content[1]["path"] = "/other"
			case "missing":
				content = content[:2]
			case "reordered":
				content[1], content[2] = content[2], content[1]
			case "extra":
				content = append(content, content[1])
			}
			b, _ := json.Marshal(content)
			if name == "text-elements" {
				b = []byte(strings.Replace(string(b), `"type":"text"`, `"type":"text","text_elements":[{}]`, 1))
			} else {
				b = []byte(strings.Replace(string(b), `"type":"text"`, `"type":"text","text_elements":[]`, 1))
			}
			item := map[string]any{"type": "userMessage", "id": "item", "clientId": v.ID, "content": json.RawMessage(b)}
			raw, _ := json.Marshal(item)
			turn := Turn{ID: "turn", Status: "completed", ItemsView: "full", Items: []json.RawMessage{raw}}
			if name == "summary" {
				turn.ItemsView = "summary"
			}
			if name == "duplicate" {
				turn.Items = append(turn.Items, raw)
			}
			thread := Thread{ID: "thread", Turns: []Turn{turn}}
			if name == "other-summary" {
				thread.Turns = append(thread.Turns, Turn{ID: "other", ItemsView: "summary"})
			}
			_, inputs, e := correlatedInputs(thread, "", []governance.OwnerInput{v})
			if name == "exact" {
				if e != nil || len(inputs) != 1 {
					t.Fatal(inputs, e)
				}
				events, _, e := History(thread, "turn", UUID(), time.Now())
				if e != nil || len(events) != 1 || events[0].Observation.InputDigest != v.WireDigest() {
					t.Fatal(events, e)
				}
			} else if e == nil {
				t.Fatal("inexact input accepted")
			}
		})
	}
}
