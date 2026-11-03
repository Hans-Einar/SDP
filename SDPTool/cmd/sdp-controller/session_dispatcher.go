package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

type sessionProcess struct {
	done chan error
	log  string
}
type sessionDispatcher struct {
	service                          *governance.Service
	sessions                         []string
	root, state, codex, mcp, sandbox string
	attempted                        map[string]bool
	active                           map[string]sessionProcess
}

// The launcher owns these controllers, not the model. Each child independently
// acquires the run lease and durable command claim before any host effect.
func (d *sessionDispatcher) tick(ctx context.Context, emit *json.Encoder) error {
	for id, p := range d.active {
		select {
		case err := <-p.done:
			result := map[string]any{"requestId": id, "log": p.log, "controller": "finished"}
			if err != nil {
				result["error"] = err.Error()
			}
			if e := emit.Encode(map[string]any{"sessionHostExecution": result}); e != nil {
				return e
			}
			delete(d.active, id)
		default:
		}
	}
	for _, id := range d.sessions {
		pending, e := d.service.ReadSessionHostCommands(id)
		if e != nil {
			return e
		}
		for _, command := range pending.Commands {
			q := command.Request
			if command.State != "queued" || d.attempted[q.OperationID] || (q.Action != "launch" && q.Action != "continue") || len(d.active) >= 2 {
				continue
			}
			executable, e := os.Executable()
			if e != nil {
				return e
			}
			logDir := filepath.Join(d.state, "session-controller-logs")
			if e = os.MkdirAll(logDir, 0700); e != nil {
				return e
			}
			logPath := filepath.Join(logDir, q.OperationID+".jsonl")
			log, e := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if e != nil {
				return e
			}
			cmd := exec.CommandContext(ctx, executable, "execute", "--project", d.root, "--state", d.state, "--host-request", q.OperationID, "--codex", d.codex, "--mcp-binary", d.mcp, "--sandbox", d.sandbox, "--reporting", "automatic")
			cmd.Stdout, cmd.Stderr = log, log
			if e = cmd.Start(); e != nil {
				log.Close()
				return fmt.Errorf("start Session controller: %w", e)
			}
			d.attempted[q.OperationID] = true
			done := make(chan error, 1)
			d.active[q.OperationID] = sessionProcess{done: done, log: logPath}
			go func() { err := cmd.Wait(); log.Close(); done <- err }()
			if e = emit.Encode(map[string]any{"sessionHostExecution": map[string]any{"requestId": q.OperationID, "runId": q.RunID, "controller": "started", "log": logPath}}); e != nil {
				return e
			}
		}
	}
	return nil
}
