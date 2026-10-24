package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestCoordinationInputAndTerminalBoundaries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	input := readCoordinationInput(ctx, strings.NewReader("{bad}\n{\"action\":\"refresh\"}\n"))
	first := <-input
	if first.Err == nil {
		t.Fatal(first)
	}
	second := <-input
	if second.Err != nil || second.Command.Action != "refresh" {
		t.Fatal(second)
	}
	if _, ok := <-input; ok {
		t.Fatal("input did not close")
	}
	text := "attack\x1b]52;c;x\a\r\nforged row\u009b31m"
	v := governance.CoordinationView{ProjectRef: text, Items: []governance.CoordinationItem{{ID: text, RunID: text, Attention: &governance.AttentionView{Status: text, Question: governance.AttentionQuestion{Prompt: text}}}}, NextCursor: text}
	var out bytes.Buffer
	if e := renderCoordination(&out, v); e != nil {
		t.Fatal(e)
	}
	for _, c := range out.String() {
		if c == '\x1b' || c == '\a' || c == '\r' || (c >= 0x80 && c <= 0x9f) {
			t.Fatal("unsafe terminal output")
		}
	}
	if strings.Contains(out.String(), "\nforged row") {
		t.Fatal("injected row")
	}
	if e := renderCoordination(failedWriter{}, v); e == nil {
		t.Fatal("write failure hidden")
	}
}
func TestCoordinationCommandsDoNotInferAuthority(t *testing.T) {
	q := governance.CoordinationRequest{RunID: "selected"}
	for _, v := range []coordinationCommand{
		{Action: "start"}, {Action: "quit", Reply: &coordinationReply{}}, {Action: "acknowledge"},
		{Action: "acknowledge", Acknowledgement: &governance.ReportAcknowledgementRequest{Coordinator: "forged", RunID: "selected"}},
		{Action: "acknowledge", Acknowledgement: &governance.ReportAcknowledgementRequest{Coordinator: "leader", RunID: "other"}},
		{Action: "reply", Disposition: &governance.AttentionDisposition{RunID: "other"}, Reply: &coordinationReply{Text: "answer"}},
	} {
		if _, e := applyCoordinationCommand(nil, q, "leader", "owner", v); e == nil {
			t.Fatal(v)
		}
	}
	for _, v := range []coordinationCommand{{Action: "refresh"}, {Action: "quit"}} {
		if r, e := applyCoordinationCommand(nil, q, "", "", v); r != nil || e != nil {
			t.Fatal(r, e)
		}
	}
	var out bytes.Buffer
	for _, args := range [][]string{nil, {"--project", "/tmp", "--state", "/tmp/no-state", "--coordinator", "leader"}, {"--project", "/tmp", "--state", "/tmp/no-state", "--watch", "--cursor", "forged"}} {
		if e := coordinationCommandMain(args, strings.NewReader(""), &out); e == nil {
			t.Fatal(args)
		}
	}
}

func TestCoordinationWatchCoalescesTicksAndStopsOnOutputFailure(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "project")
	if e := os.Mkdir(root, 0700); e != nil {
		t.Fatal(e)
	}
	s, e := governance.New(root, filepath.Join(root, "SDP"), "fixture", filepath.Join(base, "state"))
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	commands := make(chan coordinationInput)
	ticks := make(chan time.Time)
	done := make(chan error, 1)
	var out bytes.Buffer
	go func() {
		done <- watchCoordination(ctx, s, governance.CoordinationRequest{}, "", "", commands, ticks, &out, true)
	}()
	for i := 0; i < 3; i++ {
		select {
		case ticks <- time.Now():
		case <-ctx.Done():
			t.Fatal("watch blocked")
		}
	}
	for _, action := range []string{"refresh", "quit"} {
		select {
		case commands <- coordinationInput{Command: coordinationCommand{Action: action}}:
		case <-ctx.Done():
			t.Fatal("input blocked")
		}
	}
	if e := <-done; e != nil {
		t.Fatal(e)
	}
	decoder := json.NewDecoder(&out)
	count := 0
	for decoder.More() {
		var event coordinationEvent
		if e := decoder.Decode(&event); e != nil {
			t.Fatal(e)
		}
		if event.Kind != "snapshot" {
			t.Fatal(event)
		}
		count++
	}
	if count != 2 {
		t.Fatal("unchanged ticks emitted snapshots", count)
	}
	if e := watchCoordination(ctx, s, governance.CoordinationRequest{}, "", "", nil, nil, failedWriter{}, true); e == nil {
		t.Fatal("output failure hidden")
	}
	if _, e := os.Stat(filepath.Join(base, "state")); !os.IsNotExist(e) {
		t.Fatal("observer created operational state", e)
	}
}
