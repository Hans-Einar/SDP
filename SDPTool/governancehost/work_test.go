package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
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
	if _, e := c.StartWorkThread(ctx, "/fixture", "read-only", "bounded Master", map[string]any{}); e != nil {
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
