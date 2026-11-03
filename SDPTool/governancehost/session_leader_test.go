package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestSessionLeaderInventoryRequiresExactGrants(t *testing.T) {
	for _, kind := range []string{"exact", "extra-server", "missing-server", "extra-tool", "missing-tool", "duplicate-server", "page"} {
		t.Run(kind, func(t *testing.T) {
			left, right := net.Pipe()
			c := newClient(left, left)
			defer c.Close()
			defer right.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			checked := make(chan bool, 1)
			go func() {
				scanner := bufio.NewScanner(right)
				if !scanner.Scan() {
					return
				}
				var f frame
				json.Unmarshal(scanner.Bytes(), &f)
				var params map[string]any
				json.Unmarshal(f.Params, &params)
				_, filtered := params["serverName"]
				checked <- !filtered && params["threadId"] == "leader"
				a := map[string]any{"name": "sdp_governance", "tools": map[string]any{"coordination_read": map[string]any{}}}
				tools := map[string]any{"session_read": map[string]any{}, "session_host_request": map[string]any{}}
				b := map[string]any{"name": "sdp_session_1", "tools": tools}
				rows := []any{a, b}
				result := map[string]any{}
				switch kind {
				case "extra-server":
					rows = append(rows, map[string]any{"name": "inherited", "tools": map[string]any{}})
				case "missing-server":
					rows = rows[:1]
				case "extra-tool":
					tools["owner_accept"] = map[string]any{}
				case "missing-tool":
					delete(tools, "session_host_request")
				case "duplicate-server":
					rows = []any{a, a}
				case "page":
					result["nextCursor"] = "next"
				}
				result["data"] = rows
				json.NewEncoder(right).Encode(map[string]any{"id": f.ID, "result": result})
			}()
			err := c.SessionLeaderTools(ctx, "leader", map[string][]string{"sdp_governance": {"coordination_read"}, "sdp_session_1": {"session_host_request", "session_read"}})
			if (err == nil) != (kind == "exact") {
				t.Fatal(kind, err)
			}
			if !<-checked {
				t.Fatal("inventory filtered before grant check")
			}
		})
	}
}
