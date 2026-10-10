package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("closed observer") }

func TestOverviewTerminalSafetyAndWriteErrors(t *testing.T) {
	malicious := "task\x1b]52;c;clipboard\a\nforged row\r\u009b31m"
	v := governance.Overview{ProjectRef: malicious, Root: malicious, Runs: []governance.WorkView{{RunID: malicious, Outcome: malicious, Selection: governance.Selection{Session: malicious}, Conditions: []string{malicious}}}}
	var out bytes.Buffer
	if err := renderOverview(&out, v); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, c := range text {
		if c == '\x1b' || c == '\a' || c == '\r' || (c >= 0x80 && c <= 0x9f) {
			t.Fatalf("unescaped control in %q", text)
		}
	}
	if strings.Contains(text, "\nforged row") || !strings.Contains(text, `\x1b`) {
		t.Fatal(text)
	}
	if err := renderOverview(failedWriter{}, v); err == nil {
		t.Fatal("output failure hidden")
	}
}

func TestOverviewRequiresExplicitSelection(t *testing.T) {
	for _, args := range [][]string{nil, {"--project", "/tmp"}, {"--project", "/tmp", "--state", "/tmp/state", "unexpected"}, {"--run", "forged"}} {
		if err := overviewCommand(args, &bytes.Buffer{}); err == nil {
			t.Fatal(args)
		}
	}
}
