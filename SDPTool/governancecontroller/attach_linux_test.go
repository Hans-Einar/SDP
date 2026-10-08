package governancecontroller

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
	"github.com/coder/websocket"
)

// This is a protocol fixture, not a model. It asserts that reattachment subscribes
// with exactly threadId, never a turn submission or configuration override.
func TestAttachmentSubscribesWithoutDispatch(t *testing.T) {
	for _, scenario := range []string{"running", "unloaded", "other-turn", "summary"} {
		t.Run(scenario, func(t *testing.T) {
			c, host, id := setup(t)
			dir, err := os.MkdirTemp("", "pgi4-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			path := filepath.Join(dir, "host.sock")
			l, err := net.Listen("unix", path)
			if err != nil {
				t.Fatal(err)
			}
			defer l.Close()
			os.Chmod(path, 0600)
			methods := make(chan string, 4)
			server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ws, err := websocket.Accept(w, r, nil)
				if err != nil {
					return
				}
				defer ws.CloseNow()
				for {
					_, b, err := ws.Read(r.Context())
					if err != nil {
						return
					}
					var q struct {
						ID     json.RawMessage `json:"id"`
						Method string          `json:"method"`
						Params json.RawMessage `json:"params"`
					}
					if json.Unmarshal(b, &q) != nil {
						return
					}
					if q.Method == "initialized" {
						continue
					}
					var result any
					switch q.Method {
					case "initialize":
						result = map[string]string{"userAgent": "sdp_governance/" + governancehost.HostVersion + " fixture", "platformOs": "linux"}
					case "thread/read", "thread/resume":
						methods <- q.Method
						var params map[string]any
						json.Unmarshal(q.Params, &params)
						if params["threadId"] != "thread" || (q.Method == "thread/resume" && len(params) != 1) {
							t.Error("unexpected resume overrides", params)
							return
						}
						if q.Method == "thread/resume" && os.Getenv("PGI_SCHEMA_EVIDENCE") != "" {
							cases, _ := json.Marshal([]map[string]any{{"definition": "ThreadResumeParams", "value": params}})
							if err = os.WriteFile(filepath.Join(os.Getenv("PGI_SCHEMA_EVIDENCE"), "attach-rpc-cases.json"), cases, 0600); err != nil {
								t.Error(err)
							}
						}
						thread := governancehost.Thread{ID: "thread", Cwd: c.Root, CLIVersion: governancehost.HostVersion, Turns: []governancehost.Turn{{ID: "turn", Status: "inProgress", ItemsView: "full"}}}
						thread.Status.Type = "active"
						if scenario == "unloaded" {
							thread.Status.Type = "notLoaded"
						}
						if scenario == "other-turn" {
							thread.Turns[0].ID = "other"
						}
						if scenario == "summary" {
							thread.Turns[0].ItemsView = "summary"
						}
						result = map[string]any{"thread": thread}
					default:
						t.Error("unexpected RPC", q.Method)
						return
					}
					out, _ := json.Marshal(map[string]any{"id": q.ID, "result": result})
					if ws.Write(r.Context(), websocket.MessageText, out) != nil {
						return
					}
				}
			})}
			go server.Serve(l)
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client, err := governancehost.OpenSocket(ctx, path)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()
			c.HostIdentity = client.HostIdentity
			r, err := c.Launch(ctx, id, "Fix the selected fixture")
			if err != nil {
				t.Fatal(err)
			}
			recovery, err := client.PrepareAttachment(ctx, c.Service, id, c.Root, "never")
			if scenario != "running" {
				if err == nil {
					t.Fatal("uncertain history accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			r, err = c.Service.RecordHost(governance.LocalController(), recovery.ReattachRequest(client.HostIdentity))
			if err != nil {
				t.Fatal(err)
			}
			if host.starts != 1 || r.Run.Attempts[0].Epoch != client.Epoch || r.Run.Attempts[0].CapabilityEpoch() != c.Epoch {
				t.Fatal("dispatch or capability changed", r)
			}
			if <-methods != "thread/read" || <-methods != "thread/resume" {
				t.Fatal("subscription sequence changed")
			}
		})
	}
}
