package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestLeaderCoordinatorInventoryExcludesInheritedServers(t *testing.T) {
	for _, kind := range []string{"exclusive", "other-server", "incomplete"} {
		t.Run(kind, func(t *testing.T) {
			left, right := net.Pipe()
			c := newClient(left, left)
			defer c.Close()
			defer right.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			checked := make(chan bool, 1)
			go func() {
				scan := bufio.NewScanner(right)
				if !scan.Scan() {
					return
				}
				var f frame
				json.Unmarshal(scan.Bytes(), &f)
				var p map[string]any
				json.Unmarshal(f.Params, &p)
				_, filtered := p["serverName"]
				checked <- !filtered && p["threadId"] == "leader-thread"
				rows := []any{map[string]any{"name": "sdp_governance", "tools": map[string]any{"coordination_read": map[string]any{}, "report_acknowledge": map[string]any{}}}}
				result := map[string]any{"data": rows}
				if kind == "other-server" {
					result["data"] = append(rows, map[string]any{"name": "inherited", "tools": map[string]any{"external_write": map[string]any{}}})
				}
				if kind == "incomplete" {
					result["nextCursor"] = "more"
				}
				json.NewEncoder(right).Encode(map[string]any{"id": f.ID, "result": result})
			}()
			names, _, e := c.CoordinatorTools(ctx, "leader-thread")
			if kind == "exclusive" {
				if e != nil || len(names) != 2 {
					t.Fatal(names, e)
				}
			} else if e == nil {
				t.Fatal("non-exclusive inventory accepted")
			}
			if !<-checked {
				t.Fatal("inventory was filtered before authority check")
			}
		})
	}
}
