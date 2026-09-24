package parser

import (
	"strings"
	"testing"
)

const EchoSource = "language action-core version 0.1.\naction Echo.\nrecord EchoInput.\nrecord EchoOutput.\nEcho invokes GoEcho.\nEcho returns EchoOutput.\nEcho takes EchoInput.\nEchoInput field Value as text.\nEchoOutput field Value as text.\n"

func TestActionProfile(t *testing.T) {
	p, e := CompileActions(EchoSource)
	if e != nil {
		t.Fatal(e)
	}
	if p.Actions["Echo"].GoSymbol != "GoEcho" || p.Records["EchoInput"]["Value"] != TextType {
		t.Fatal(p)
	}
	for _, s := range []string{strings.Replace(EchoSource, "version 0.1", "version 0.2", 1), strings.Replace(EchoSource, "takes EchoInput", "takes Missing", 1), strings.Replace(EchoSource, "Echo takes EchoInput.\n", "", 1), strings.Replace(EchoSource, "as text", "as object", 1), strings.Replace(EchoSource, "record EchoInput.", "record EchoInput.\nrecord EchoInput.", 1), EchoSource + "EchoInput field Value as integer.\n", strings.Replace(EchoSource, "Echo takes EchoInput.", "EchoInput takes EchoInput.", 1), strings.Replace(EchoSource, "EchoOutput field Value as text.\n", "", 1)} {
		if _, e := CompileActions(s); e == nil {
			t.Fatalf("accepted invalid action source %s", s)
		}
	}
	for i := range EchoSource {
		_, _ = CompileActions(EchoSource[:i])
	}
}
