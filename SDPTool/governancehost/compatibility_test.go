package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestInitializeObservedVersionAndContract(t *testing.T) {
	for _, tc := range []struct {
		name, executable, server, os string
		pass                         bool
	}{
		{"baseline", "0.164.0", "0.164.0", "linux", true},
		{"future", "0.190.3", "0.190.3", "linux", true},
		{"future socket", "", "0.190.3", "linux", true},
		{"different executable", "0.164.0", "0.190.3", "linux", false},
		{"old socket", "", "0.161.0", "linux", false},
		{"preview", "", "0.164.0-dev", "linux", false},
		{"missing capability", "0.164.0", "0.164.0", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			left, right := net.Pipe()
			defer right.Close()
			c := newClient(left, left)
			c.version = tc.executable
			defer c.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			go func() {
				s := bufio.NewScanner(right)
				if !s.Scan() {
					return
				}
				var f frame
				if json.Unmarshal(s.Bytes(), &f) != nil {
					return
				}
				json.NewEncoder(right).Encode(map[string]any{"id": f.ID, "result": map[string]any{"userAgent": "sdp_governance/" + tc.server + " fixture", "platformOs": tc.os, "futureOptionalField": true}})
				// Drain initialized; the producer must not submit a work turn.
				s.Scan()
			}()
			_, err := c.initialize(ctx)
			if (err == nil) != tc.pass {
				t.Fatalf("pass=%v: %v", tc.pass, err)
			}
			if tc.pass && c.Version() != tc.server {
				t.Fatalf("actual version lost: %q", c.Version())
			}
		})
	}
}

func TestNewerVersionStillRequiresTheWireContract(t *testing.T) {
	for _, result := range []string{
		`{"error":{"code":-32601}}`,
		`{"error":{"code":-32602}}`,
		`{"result":null}`,
		`{"result":{"turn":{"id":"turn","status":"unknownNewStatus"}}}`,
	} {
		t.Run(result, func(t *testing.T) {
			left, right := net.Pipe()
			defer right.Close()
			c := newClient(left, left)
			c.version = "0.190.3"
			defer c.Close()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			go func() {
				s := bufio.NewScanner(right)
				if !s.Scan() {
					return
				}
				var request frame
				if json.Unmarshal(s.Bytes(), &request) != nil {
					return
				}
				var response map[string]any
				json.Unmarshal([]byte(result), &response)
				response["id"] = request.ID
				json.NewEncoder(right).Encode(response)
			}()
			if _, err := c.StartTurn(ctx, "thread", "message", "bounded fixture"); err == nil {
				t.Fatal("newer version bypassed required turn contract")
			} else if strings.Contains(result, "error") && !strings.Contains(err.Error(), "incompatible_profile") {
				t.Fatal("missing wire capability was not identified", err)
			}
		})
	}
}

func TestReplayRequestRetainsRecordedHostVersion(t *testing.T) {
	for _, version := range []string{"0.160.1", "0.161.0", "0.164.0", "0.190.3"} {
		b := Binding{AttemptID: "fixture", Result: governance.Result{Run: governance.Run{Attempts: []governance.Attempt{{ID: "fixture", HostVersion: version}}}}}
		if got := b.Request("reconcile", governance.Observation{}).HostVersion; got != version {
			t.Fatalf("%s relabeled as %s", version, got)
		}
	}
	if got := (Binding{AttemptID: "unmapped"}).Request("observe", governance.Observation{}).HostVersion; got != "" {
		t.Fatalf("unmapped attempt acquired invented host version %q", got)
	}
}
