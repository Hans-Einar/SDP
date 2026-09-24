package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	for _, command := range []string{"ast", "check", "format"} {
		var out, err bytes.Buffer
		if code := execute([]string{command, "-"}, strings.NewReader("language design-core version 0.5.\nunit A.\n"), &out, &err); code != 0 || out.Len() == 0 {
			t.Fatal(code, out.String(), err.String())
		}
	}
	var out, err bytes.Buffer
	if execute([]string{"format", "-"}, strings.NewReader("language design-core version 0.5.\nfunctionality MissingOwner.\n"), &out, &err) != 1 || out.Len() != 0 {
		t.Fatal("invalid model produced formatted output")
	}
}
