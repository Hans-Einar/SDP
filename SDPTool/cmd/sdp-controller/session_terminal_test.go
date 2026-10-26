package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestTerminalSessionSelectorAndEscaping(t *testing.T) {
	a := governance.SessionView{SessionID: "SESSION-TEST-0001", InstanceID: "a"}
	b := a
	b.InstanceID = "b"
	list := governance.SessionList{Instances: []governance.SessionView{a, b}}
	if _, e := selectTerminalSession(list, a.SessionID, ""); e == nil {
		t.Fatal("ambiguous canonical Session selected")
	}
	if v, e := selectTerminalSession(list, "", "b"); e != nil || v.InstanceID != "b" {
		t.Fatal(v, e)
	}
	a.SessionID = "session\x1b]52;c;clipboard\a\nforged"
	a.Conditions = []string{a.SessionID}
	a.Steps = []governance.SessionStepView{{ID: a.SessionID, Holds: a.Conditions}}
	var out bytes.Buffer
	if e := emitTerminalSession(nil, a, &out, false); e != nil {
		t.Fatal(e)
	}
	if strings.ContainsAny(out.String(), "\x1b\a\r") || strings.Contains(out.String(), "\nforged") {
		t.Fatal(out.String())
	}
	if e := emitTerminalSession(nil, a, failedWriter{}, false); e == nil {
		t.Fatal("lost output error")
	}
}

func TestTerminalSessionRejectsInvalidArgumentsBeforeDiscovery(t *testing.T) {
	for _, args := range [][]string{
		nil, {"forged"}, {"status", "--project", "/missing", "--state", "/missing"},
		{"list", "--project", "/missing", "--state", "/missing", "--session", "extra"},
		{"select", "--project", "/missing", "--state", "/missing", "--session", "extra", "--reason", "missing step"},
		{"answer", "--project", "/missing", "--state", "/missing", "--session", "extra"},
		{"accept", "--project", "/missing", "--state", "/missing", "--session", "extra", "--reason", "missing exact gate"},
	} {
		if e := sessionCommand(args, &bytes.Buffer{}, strings.NewReader(""), &bytes.Buffer{}); e == nil || strings.Contains(e.Error(), "manifest") {
			t.Fatal(args, e)
		}
	}
}
