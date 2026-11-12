package governanceapp

import (
	"context"
	"errors"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestRequestOperationsRejectUnknownFieldsBeforeDomainCall(t *testing.T) {
	for _, op := range Operations() {
		if op.Input != RequestInput {
			continue
		}
		t.Run(op.Name, func(t *testing.T) {
			// A nil service makes accidental dispatch observable; strict decoding must fail first.
			_, err := op.Execute(context.Background(), nil, governance.LocalOwner(), []byte(`{"unexpected":true}`), "")
			if err == nil {
				t.Fatal("unknown request field accepted")
			}
		})
	}
}
func TestCancellationBoundaries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, op := range Operations() {
		_, err := op.Execute(ctx, nil, governance.LocalOwner(), nil, "")
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("%s: %v", op.Name, err)
		}
	}
	ctx, cancel = context.WithCancel(context.Background())
	committed := Operation{execute: func(_ *governance.Service, _ governance.Caller, _ []byte, _ string) (any, error) {
		cancel()
		return "committed", nil
	}}
	result, err := committed.Execute(ctx, nil, governance.LocalOwner(), nil, "")
	if err != nil || result != "committed" {
		t.Fatalf("committed result lost: %v %v", result, err)
	}
}

func TestCLIArgumentContracts(t *testing.T) {
	cases := []struct{ name, input, run, instance, want string }{
		{name: "session-list", input: "request", want: "operation takes no request/run"},
		{name: "bind", run: "run", want: "bind takes no request/run"},
		{name: "read", want: "read requires --run only"},
		{name: "session-host-read", input: "request", instance: "instance", want: "session-read requires --instance only"},
		{name: "session-read", want: "session-read requires --instance only"},
		{name: "session-adopt", run: "run", want: "operation requires --input only"},
		{name: "session-definition", instance: "instance", want: "--instance is only valid for session-read"},
	}
	for _, tc := range cases {
		op, ok := Find(tc.name)
		if !ok {
			t.Fatal(tc.name)
		}
		if err := op.Validate(tc.input, tc.run, tc.instance); err == nil || err.Error() != tc.want {
			t.Fatalf("%s: %v, want %q", tc.name, err, tc.want)
		}
	}
	for _, op := range Operations() {
		var input, run, instance string
		switch op.Input {
		case RequestInput:
			input = "request"
		case RunInput:
			run = "run"
		case InstanceInput:
			instance = "instance"
		}
		if err := op.Validate(input, run, instance); err != nil {
			t.Fatalf("%s rejected valid input kind: %v", op.Name, err)
		}
	}
}
