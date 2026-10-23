package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestAttentionTerminalQuotesProjectText(t *testing.T) {
	malicious := "\x1b[31m\nforged status"
	v := governance.AttentionInbox{ProjectRef: malicious, Items: []governance.AttentionView{{RunID: "run", Status: "answered", Applicability: "stale", Coverage: "stale", Question: governance.AttentionQuestion{ID: "q", Prompt: malicious, Options: []governance.AttentionOption{{ID: malicious, Label: malicious}}, Reply: &governance.AttentionReply{Actor: malicious, Text: malicious}}}}}
	var out bytes.Buffer
	if e := renderAttention(&out, v); e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out.String(), malicious) || strings.Contains(out.String(), "\x1b") || !strings.Contains(out.String(), `\x1b[31m\nforged status`) || !strings.Contains(out.String(), `applicability: "stale"`) {
		t.Fatal(out.String())
	}
}
