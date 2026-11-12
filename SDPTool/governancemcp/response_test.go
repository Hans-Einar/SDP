package governancemcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestResponsePoliciesPreserveWireErrors(t *testing.T) {
	for _, policy := range []responsePolicy{{governance.Schema, "use bounded observation pages"}, {governance.SessionSchema, "bounded Session response exceeded"}} {
		t.Run(policy.schema, func(t *testing.T) {
			for _, tc := range []struct {
				value         any
				err           error
				code, message string
			}{
				{nil, errors.New("bad request"), "invalid_request", "bad request"},
				{nil, &governance.Error{Code: "stale_revision", Message: "changed"}, "stale_revision", "changed"},
				{strings.Repeat("x", 256<<10), nil, "limit", policy.limitMessage},
			} {
				calls := 0
				r, err := toolResponse(context.Background(), nil, policy, func([]byte) (any, error) { calls++; return tc.value, tc.err })
				if err != nil || !r.IsError || calls != 1 {
					t.Fatalf("result %v, calls %d", err, calls)
				}
				var got struct {
					Schema string
					Error  governance.Error
				}
				if err = json.Unmarshal(r.StructuredContent.(json.RawMessage), &got); err != nil || got.Schema != policy.schema || got.Error.Code != tc.code || got.Error.Message != tc.message {
					t.Fatalf("wire error: %+v %v", got, err)
				}
			}
		})
	}
}

func TestResponseRetryRetainsArgumentsAndCommitOutcome(t *testing.T) {
	policy := responsePolicy{governance.Schema, "use bounded observation pages"}
	args := []byte(`{"operationId":"fixed"}`)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	r, err := toolResponse(ctx, args, policy, func(got []byte) (any, error) {
		calls++
		if !bytes.Equal(got, args) {
			t.Fatal("retry changed request")
		}
		if calls == 1 {
			return nil, &governance.Error{Code: "busy", Message: "locked"}
		}
		cancel() // A completed mutation must not be reported undone.
		return map[string]any{"committed": true}, nil
	})
	if err != nil || r.IsError || calls != 2 {
		t.Fatalf("lost committed result: %+v %v %d", r, err, calls)
	}
	if !bytes.Contains(r.StructuredContent.(json.RawMessage), []byte(`"committed":true`)) {
		t.Fatal("lost result")
	}
	ctx, cancel = context.WithCancel(context.Background())
	calls = 0
	_, err = toolResponse(ctx, args, policy, func([]byte) (any, error) {
		calls++
		cancel()
		return nil, &governance.Error{Code: "busy", Message: "locked"}
	})
	if !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatalf("busy cancellation: %v %d", err, calls)
	}
	_, err = toolResponse(ctx, args, policy, func([]byte) (any, error) { t.Fatal("pre-canceled call executed"); return nil, nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
