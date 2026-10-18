package main

import (
	"bytes"
	"testing"
)

func TestAssignmentCommandDoesNotImplicitlySelectOrLaunch(t *testing.T) {
	for _, args := range [][]string{nil, {"launch"}, {"prepare"}, {"revise"}, {"retire"}, {"retire", "--project", "/tmp", "--state", "/tmp/a", "--assignment", "x"}, {"context", "--project", "/tmp", "--state", "/tmp/a", "--input", "a", "--assignment", "x"}, {"inspect", "--project", "/tmp", "--state", "/tmp/a", "--input", "a"}} {
		if err := assignmentCommand(args, &bytes.Buffer{}); err == nil {
			t.Fatal("ambiguous operation accepted", args)
		}
	}
}
