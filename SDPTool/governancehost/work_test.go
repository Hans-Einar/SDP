package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWorkProtocolRequests(t *testing.T) {
	left, right := net.Pipe()
	c := newClient(left, left)
	defer c.Close()
	defer right.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	captured := make(chan []map[string]any, 1)
	go func() {
		var cases []map[string]any
		s := bufio.NewScanner(right)
		definitions := []string{"ThreadStartParams", "TurnStartParams", "TurnSteerParams", "TurnInterruptParams"}
		results := []string{`{"thread":{"id":"thread","status":{"type":"idle"}}}`, `{"turn":{"id":"turn","status":"inProgress"}}`, `{"turnId":"turn"}`, `{}`}
		for i, definition := range definitions {
			if !s.Scan() {
				return
			}
			var f frame
			if json.Unmarshal(s.Bytes(), &f) != nil {
				return
			}
			cases = append(cases, map[string]any{"definition": definition, "value": f.Params})
			response := map[string]any{"id": f.ID, "result": json.RawMessage(results[i])}
			b, _ := json.Marshal(response)
			if _, e := right.Write(append(b, '\n')); e != nil {
				return
			}
		}
		captured <- cases
	}()
	if _, e := c.StartWorkThread(ctx, "/fixture", "read-only", "bounded Master", "never", map[string]any{}); e != nil {
		t.Fatal(e)
	}
	id := UUID()
	if _, e := c.StartTurn(ctx, "thread", id, "exact input"); e != nil {
		t.Fatal(e)
	}
	if e := c.SteerTurn(ctx, "thread", "turn", UUID(), "steering"); e != nil {
		t.Fatal(e)
	}
	if e := c.InterruptTurn(ctx, "thread", "turn"); e != nil {
		t.Fatal(e)
	}
	cases := <-captured
	var start, steer map[string]any
	b, _ := json.Marshal(cases[1]["value"])
	json.Unmarshal(b, &start)
	b, _ = json.Marshal(cases[2]["value"])
	json.Unmarshal(b, &steer)
	if start["clientUserMessageId"] != id || steer["expectedTurnId"] != "turn" {
		t.Fatal(start, steer)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.MarshalIndent(cases, "", "  ")
		if e := os.WriteFile(filepath.Join(dir, "work-rpc-cases.json"), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}

func TestCorrelationRequiresExactUniqueRecordedInput(t *testing.T) {
	id := UUID()
	input := governance.OwnerInput{ID: id, Payload: "exact"}
	item := func(client, text string) json.RawMessage {
		b, _ := json.Marshal(map[string]any{"type": "userMessage", "id": "item", "clientId": client, "content": []any{map[string]string{"type": "text", "text": text}}})
		return b
	}
	for _, tc := range []struct {
		name   string
		items  []json.RawMessage
		mapped string
		good   bool
	}{
		{"match", []json.RawMessage{item(id, "exact")}, "", true},
		{"text-only", []json.RawMessage{item("", "exact")}, "", false},
		{"changed", []json.RawMessage{item(id, "different")}, "", false},
		{"duplicate", []json.RawMessage{item(id, "exact"), item(id, "exact")}, "", false},
		{"other-turn", []json.RawMessage{item(id, "exact")}, "other", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			thread := Thread{ID: "thread", Turns: []Turn{{ID: "turn", Status: "completed", Items: tc.items}}}
			turn, inputs, e := correlatedInputs(thread, tc.mapped, []governance.OwnerInput{input})
			if tc.good {
				if e != nil || turn != "turn" || len(inputs) != 1 || inputs[0].ThreadID != "thread" {
					t.Fatal(turn, inputs, e)
				}
			} else if e == nil {
				t.Fatal("ambiguous correlation accepted")
			}
		})
	}
}

func TestExplicitApprovalModeUsesStableOwnerPolicy(t *testing.T) {
	left, right := net.Pipe()
	c := newClient(left, left)
	defer c.Close()
	defer right.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	captured := make(chan frame, 1)
	go func() {
		scanner := bufio.NewScanner(right)
		if !scanner.Scan() {
			return
		}
		var f frame
		json.Unmarshal(scanner.Bytes(), &f)
		json.NewEncoder(right).Encode(map[string]any{"id": f.ID, "result": map[string]any{"thread": map[string]any{"id": "thread", "status": map[string]string{"type": "idle"}}}})
		captured <- f
	}()
	if _, e := c.StartWorkThread(ctx, "/fixture", "read-only", "bounded", "ask", map[string]any{}); e != nil {
		t.Fatal(e)
	}
	f := <-captured
	var params map[string]any
	json.Unmarshal(f.Params, &params)
	if params["approvalPolicy"] != "on-request" || params["approvalsReviewer"] != "user" || params["sandbox"] != "read-only" {
		t.Fatal(params)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.MarshalIndent([]map[string]any{{"definition": "ThreadStartParams", "value": f.Params}}, "", "  ")
		if e := os.WriteFile(filepath.Join(dir, "approval-start-cases.json"), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
