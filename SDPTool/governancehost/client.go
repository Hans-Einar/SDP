// Package governancehost is a version-pinned local app-server adapter. PGI2
// supports negotiation and observed state, not model submission or account login.
package governancehost

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

const HostVersion = "0.161.0"

func UUID() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic(e)
	}
	b[6] = (b[6] & 15) | 64
	b[8] = (b[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type frame struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code int `json:"code"`
	} `json:"error"`
}
type response struct {
	f   frame
	err error
}
type Client struct {
	approvalMu    sync.Mutex
	ask           bool
	approvals     map[string]pendingApproval
	cmd           *exec.Cmd
	in            io.WriteCloser
	out           io.ReadCloser
	mu            sync.Mutex
	writeMu       sync.Mutex
	id            uint64
	responses     chan response
	Notifications chan Event
	done          chan struct{}
	closed        sync.Once
	errMu         sync.Mutex
	err           error
	Epoch         string
	HostIdentity  string
	UserAgent     string
}

// Open starts an owned stdio process with its existing managed-account setup.
// It never edits account configuration, signs in or submits a model turn.
func Open(ctx context.Context, binary, cwd string) (*Client, error) {
	v, e := exec.CommandContext(ctx, binary, "--version").Output()
	if e != nil {
		return nil, fmt.Errorf("host version unavailable: %w", e)
	}
	if strings.TrimSpace(string(v)) != "codex-cli "+HostVersion {
		return nil, fmt.Errorf("incompatible_profile: require codex-cli %s", HostVersion)
	}
	cmd := exec.CommandContext(ctx, binary, "app-server", "--listen", "stdio://")
	cmd.Dir = cwd
	in, e := cmd.StdinPipe()
	if e != nil {
		return nil, e
	}
	out, e := cmd.StdoutPipe()
	if e != nil {
		in.Close()
		return nil, e
	}
	cmd.Stderr = io.Discard // Never persist unfiltered host logs or auth payloads.
	if e = cmd.Start(); e != nil {
		in.Close()
		return nil, e
	}
	c := newClient(in, out)
	c.cmd = cmd
	return c.initialize(ctx)
}

func (c *Client) initialize(ctx context.Context) (*Client, error) {
	var init struct {
		UserAgent      string `json:"userAgent"`
		PlatformOs     string `json:"platformOs"`
		PlatformFamily string `json:"platformFamily"`
		CodexHome      string `json:"codexHome"`
	}
	e := c.call(ctx, "initialize", map[string]any{"clientInfo": map[string]string{"name": "sdp_governance", "title": "SDP ProjectGovernance", "version": "PGI2"}, "capabilities": map[string]any{"experimentalApi": false}}, &init)
	if e == nil && (init.UserAgent == "" || init.PlatformOs == "") {
		e = fmt.Errorf("incompatible_profile: missing initialize capabilities")
	}
	if e == nil && c.HostIdentity != "" && !strings.HasPrefix(init.UserAgent, "sdp_governance/"+HostVersion+" ") {
		e = fmt.Errorf("incompatible_profile: socket server version mismatch")
	}
	if e == nil {
		e = c.send(map[string]any{"method": "initialized", "params": map[string]any{}})
	}
	if e != nil {
		c.Close()
		return nil, e
	}
	c.UserAgent = init.UserAgent
	return c, nil
}
func newClient(in io.WriteCloser, out io.ReadCloser) *Client {
	c := &Client{approvals: map[string]pendingApproval{}, in: in, out: out, responses: make(chan response, 1), Notifications: make(chan Event, 256), done: make(chan struct{}), Epoch: UUID()}
	go func() {
		defer out.Close()
		defer close(c.Notifications)
		s := bufio.NewScanner(out)
		s.Buffer(make([]byte, 4096), governance.MaxRequest+1)
		for s.Scan() {
			b := s.Bytes()
			if e := governance.CheckWire(b); e != nil {
				c.stop(e)
				return
			}
			var f frame
			if e := json.Unmarshal(b, &f); e != nil {
				c.stop(e)
				return
			}
			if f.Method != "" {
				var event Event
				var keep bool
				var e error
				if len(f.ID) > 0 {
					event, keep, e = c.approvalRequest(f)
					if e == nil && !keep {
						e = c.send(map[string]any{"id": f.ID, "error": map[string]any{"code": -32601, "message": "unsupported controller request"}})
					}
				} else if f.Method == "serverRequest/resolved" {
					event, keep, e = c.resolveApproval(f.Params)
				} else {
					event, keep, e = Visible(f.Method, f.Params, c.Epoch, time.Now().UTC())
					if keep && event.Observation.Kind == "turn.completed" {
						c.approvalMu.Lock()
						for key, p := range c.approvals {
							if p.prompt.ThreadID == event.ThreadID {
								delete(c.approvals, key)
							}
						}
						c.approvalMu.Unlock()
					}
				}
				if e != nil {
					c.stop(e)
					return
				}
				if keep {
					select {
					case c.Notifications <- event:
					case <-c.done:
						return
					default:
						c.stop(fmt.Errorf("observation_gap: notification queue overflow"))
						return
					}
				}
				continue
			}
			select {
			case c.responses <- response{f: f}:
			case <-c.done:
				return
			default:
				c.stop(fmt.Errorf("unexpected response"))
				return
			}
		}
		e := s.Err()
		if e == nil {
			e = io.EOF
		}
		c.stop(e)
	}()
	return c
}
func (c *Client) stop(e error) {
	c.closed.Do(func() {
		c.errMu.Lock()
		c.err = e
		c.errMu.Unlock()
		close(c.done)
		c.in.Close()
		c.out.Close()
	})
}
func (c *Client) Err() error { c.errMu.Lock(); defer c.errMu.Unlock(); return c.err }
func (c *Client) Close() error {
	c.stop(io.EOF)
	if c.cmd != nil {
		_ = c.cmd.Process.Kill()
		return c.cmd.Wait()
	}
	return nil
}
func (c *Client) send(v any) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	if len(b) > governance.MaxRequest {
		return fmt.Errorf("request exceeds frame limit")
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	_, e = c.in.Write(append(b, '\n'))
	return e
}
func (c *Client) call(ctx context.Context, method string, params, out any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	stop := context.AfterFunc(ctx, func() { c.stop(ctx.Err()) })
	defer stop()
	select {
	case <-c.done:
		return fmt.Errorf("reconciliation_required: connection closed")
	default:
	}
	c.id++
	id := c.id
	if e := c.send(map[string]any{"id": id, "method": method, "params": params}); e != nil {
		c.stop(e)
		return e
	}
	select {
	case r := <-c.responses:
		var got uint64
		if e := json.Unmarshal(r.f.ID, &got); e != nil || got != id {
			c.stop(fmt.Errorf("response identity mismatch"))
			return fmt.Errorf("reconciliation_required: response identity mismatch")
		}
		if r.f.Error != nil {
			return fmt.Errorf("host_rpc_error: %s code %d", method, r.f.Error.Code)
		}
		if len(r.f.Result) == 0 {
			return fmt.Errorf("invalid host result")
		}
		return json.Unmarshal(r.f.Result, out)
	case <-ctx.Done():
		c.stop(ctx.Err())
		return fmt.Errorf("reconciliation_required: request outcome uncertain: %w", ctx.Err())
	case <-c.done:
		return fmt.Errorf("reconciliation_required: connection ended: %w", c.Err())
	}
}

// AccountMode returns only the mode, excluding email, tokens and plan identifiers.
func (c *Client) AccountMode(ctx context.Context) (string, error) {
	var out struct {
		Account *struct {
			Type string `json:"type"`
		} `json:"account"`
	}
	if e := c.call(ctx, "account/read", map[string]bool{"refreshToken": false}, &out); e != nil {
		return "", e
	}
	if out.Account == nil {
		return "unavailable", nil
	}
	return out.Account.Type, nil
}

type Thread struct {
	Cwd        string `json:"cwd"`
	CLIVersion string `json:"cliVersion"`
	ID         string `json:"id"`
	Status     struct {
		Type string `json:"type"`
	} `json:"status"`
	Turns []Turn `json:"turns"`
}
type Turn struct {
	ItemsView string            `json:"itemsView"`
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Items     []json.RawMessage `json:"items"`
}

// StartProbeThread creates an ephemeral, idle, read-only thread. It does not
// submit a turn. The optional MCP configuration belongs to this thread only.
func (c *Client) StartProbeThread(ctx context.Context, cwd string, mcpConfig map[string]any) (Thread, error) {
	var out struct {
		Thread Thread `json:"thread"`
	}
	params := map[string]any{"cwd": cwd, "ephemeral": true, "sandbox": "read-only", "approvalPolicy": "never"}
	if mcpConfig != nil {
		params["config"] = map[string]any{"mcp_servers": mcpConfig}
	}
	e := c.call(ctx, "thread/start", params, &out)
	return out.Thread, e
}
func (c *Client) ReadThread(ctx context.Context, id string) (Thread, error) {
	var out struct {
		Thread Thread `json:"thread"`
	}
	e := c.call(ctx, "thread/read", map[string]any{"threadId": id, "includeTurns": true}, &out)
	if e == nil && out.Thread.ID != id {
		e = fmt.Errorf("host thread mismatch")
	}
	return out.Thread, e
}

// MCPTools reports only the named probe server's tool names and auth status.
// No account configuration or other server details leave the adapter.
func (c *Client) MCPTools(ctx context.Context, name, threadID string) ([]string, string, error) {
	var out struct {
		Data []struct {
			Name       string                     `json:"name"`
			Tools      map[string]json.RawMessage `json:"tools"`
			AuthStatus string                     `json:"authStatus"`
		} `json:"data"`
		NextCursor *string `json:"nextCursor"`
	}
	e := c.call(ctx, "mcpServerStatus/list", map[string]any{"limit": 64, "serverName": name, "threadId": threadID}, &out)
	if e != nil {
		return nil, "", e
	}
	for _, v := range out.Data {
		if v.Name == name {
			names := []string{}
			for k := range v.Tools {
				names = append(names, k)
			}
			sort.Strings(names)
			return names, v.AuthStatus, nil
		}
	}
	return nil, "", fmt.Errorf("unavailable: selected MCP server not present in bounded status page")
}
