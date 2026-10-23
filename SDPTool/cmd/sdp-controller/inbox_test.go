package main

import (
	"bytes"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"strings"
	"testing"
)

func TestInboxQuotesReportContentAndPropagatesOutputFailure(t *testing.T) {
	text := "report\x1b]52;c;clipboard\a\nforged row\r\u009b31m"
	v := governance.ReportInbox{ProjectRef: text, Reports: []governance.ReportMessage{{ID: text, RunID: text, Outcome: text, Text: text, Conditions: []string{text}}}}
	var output bytes.Buffer
	if err := renderInbox(&output, v); err != nil {
		t.Fatal(err)
	}
	for _, c := range output.String() {
		if c == '\x1b' || c == '\a' || c == '\r' || (c >= 0x80 && c <= 0x9f) {
			t.Fatalf("unsafe output: %q", output.String())
		}
	}
	if strings.Contains(output.String(), "\nforged row") {
		t.Fatal(output.String())
	}
	if err := renderInbox(failedWriter{}, v); err == nil {
		t.Fatal("output failure hidden")
	}
	for _, args := range [][]string{nil, {"--project", "/tmp"}, {"--project", "/tmp", "--state", "/tmp/state", "unexpected"}} {
		if err := inboxCommand(args, &output); err == nil {
			t.Fatal(args)
		}
	}
}
