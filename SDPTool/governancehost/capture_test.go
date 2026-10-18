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
)

func TestVisibleCaptureAndHistory(t *testing.T) {
	epoch := UUID()
	now := time.Now().UTC()
	for _, kind := range []string{"reasoning", "functionCallOutput", "unknownFutureItem"} {
		b := []byte(`{"threadId":"t","turnId":"u","item":{"id":"i","type":"` + kind + `","text":"secret must not persist"}}`)
		_, keep, e := Visible("item/completed", b, epoch, now)
		if e != nil || keep {
			t.Fatal(kind, e)
		}
	}
	for _, phase := range []string{"commentary", "final_answer", "unknown"} {
		b := []byte(`{"threadId":"t","turnId":"u","item":{"id":"i","type":"agentMessage","phase":"` + phase + `","text":"visible","untrustedRole":"owner"}}`)
		v, keep, e := Visible("item/completed", b, epoch, now)
		if e != nil || !keep || v.Observation.Phase != phase || v.Observation.Text != "visible" {
			t.Fatal(v, e)
		}
	}
	_, keep, e := Visible("item/agentMessage/delta", []byte(`{"delta":"provisional"}`), epoch, now)
	if e != nil || keep {
		t.Fatal(e)
	}
	b := []byte(`{"threadId":"t","turnId":"u","item":{"id":"i","type":"commandExecution","command":"secret","aggregatedOutput":"secret","status":"completed"}}`)
	v, keep, e := Visible("item/completed", b, epoch, now)
	out, _ := json.Marshal(v)
	if e != nil || !keep || strings.Contains(string(out), "secret") {
		t.Fatal(string(out), e)
	}
	item := json.RawMessage(`{"id":"i","type":"agentMessage","text":"final","phase":"final_answer"}`)
	thread := Thread{ID: "t", Turns: []Turn{{ID: "u", Status: "interrupted", Items: []json.RawMessage{item}}}}
	events, status, e := History(thread, "u", epoch, now)
	if e != nil || len(events) != 1 || status != "interrupted" {
		t.Fatal(events, status, e)
	}
	_, status, e = History(thread, "absent", epoch, now)
	if e == nil || status != "unknown" {
		t.Fatal(status, e)
	}
	_, _, e = Visible("turn/completed", []byte(`{"threadId":"t","threadId":"other","turn":{"id":"u","status":"completed"}}`), epoch, now)
	if e == nil {
		t.Fatal("duplicate fields accepted")
	}
	b = []byte(`{"threadId":"t","turnId":"u","item":{"id":"user","type":"userMessage","content":[{"type":"text","text":"inspect attached"},{"type":"localImage","path":"/selected/image.png"},{"type":"audio","url":"private-token-url"}]}}`)
	v, keep, e = Visible("item/completed", b, epoch, now)
	if e != nil || !keep || len(v.Observation.Attachments) != 2 || v.Observation.Attachments[0].Reference != "/selected/image.png" || v.Observation.Attachments[1].Availability != "unavailable" {
		t.Fatal(v, e)
	}
	out, _ = json.Marshal(v)
	if strings.Contains(string(out), "private-token-url") {
		t.Fatal("remote attachment payload persisted")
	}
	thread.Turns[0].Status = "inProgress"
	events, status, e = History(thread, "u", epoch, now)
	if e != nil || len(events) != 0 || status != "running" {
		t.Fatal("partial history finalized", events, status, e)
	}
	thread.Turns[0].ItemsView = "summary"
	if _, _, err := History(thread, "u", epoch, now); err == nil {
		t.Fatal("summary treated as full turn")
	}
	thread.Turns[0].ItemsView = "full"
	thread.Turns[0].Status = "interrupted"
	thread.Turns[0].Items = nil
	events, status, e = History(thread, "u", epoch, now)
	if e != nil || len(events) != 0 || status != "interrupted" {
		t.Fatal("fabricated final response", events, status, e)
	}
}

func TestUnsupportedHostRefusesBeforeConnection(t *testing.T) {
	p := filepath.Join(t.TempDir(), "codex")
	if e := os.WriteFile(p, []byte("#!/bin/sh\nprintf 'codex-cli 0.0.0\\n'\n"), 0700); e != nil {
		t.Fatal(e)
	}
	_, e := Open(context.Background(), p, t.TempDir())
	if e == nil || !strings.Contains(e.Error(), "incompatible_profile") {
		t.Fatal(e)
	}
}
func TestConnectionInterleavingAndLoss(t *testing.T) {
	left, right := net.Pipe()
	c := newClient(left, left)
	defer c.Close()
	defer right.Close()
	go func() {
		s := bufio.NewScanner(right)
		if !s.Scan() {
			return
		}
		var request frame
		json.Unmarshal(s.Bytes(), &request)
		right.Write([]byte(`{"method":"item/completed","params":{"threadId":"t","turnId":"u","item":{"id":"i","type":"agentMessage","text":"hello","phase":"commentary"}}}` + "\n"))
		b, _ := json.Marshal(map[string]any{"id": request.ID, "result": map[string]any{"account": map[string]any{"type": "chatgpt", "email": "private", "accessToken": "never record"}}})
		right.Write(append(b, '\n'))
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mode, e := c.AccountMode(ctx)
	if e != nil || mode != "chatgpt" {
		t.Fatal(mode, e)
	}
	v := <-c.Notifications
	if v.Observation.Text != "hello" {
		t.Fatal(v)
	}
	right.Close()
	select {
	case <-c.done:
	case <-ctx.Done():
		t.Fatal("disconnect not observed")
	}
	if _, e = c.AccountMode(ctx); e == nil {
		t.Fatal("closed connection reused")
	}
}
func TestOversizeAndCancellation(t *testing.T) {
	left, right := net.Pipe()
	c := newClient(left, left)
	defer c.Close()
	defer right.Close()
	go func() { right.Write([]byte(strings.Repeat("x", (1<<20)+2))) }()
	select {
	case <-c.done:
	case <-time.After(2 * time.Second):
		t.Fatal("oversize not bounded")
	}
	left2, right2 := net.Pipe()
	c2 := newClient(left2, left2)
	defer c2.Close()
	defer right2.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, e := c2.AccountMode(ctx); e == nil {
		t.Fatal("unread request did not time out")
	}
}
