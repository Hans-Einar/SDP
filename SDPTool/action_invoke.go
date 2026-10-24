package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

func validCorrelation(s string) bool {
	if len(s) < 1 || len(s) > 128 {
		return false
	}
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}
func checkInvokeDepth(data []byte) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("invalid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	depth := 0
	for {
		t, err := d.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if v, ok := t.(json.Delim); ok {
			switch v {
			case '{', '[':
				depth++
				if depth > MaxInvokeDepth {
					return fmt.Errorf("request nesting exceeds %d", MaxInvokeDepth)
				}
			case '}', ']':
				depth--
			}
		}
	}
}
func invokeFailure(r InvokeResponse, code int, name string, err error) InvokeResponse {
	r.ExitCode = code
	r.Status = "error"
	if code == 130 {
		r.Status = "canceled"
	}
	r.Error = &InvokeError{Code: name, Message: err.Error()}
	return r
}

// Invoke reads a single bounded request. It does not authenticate a caller.
// Hosted callers must supply independently established ActionCaller authority.
func Invoke(ctx context.Context, in io.Reader, caller ActionCaller) InvokeResponse {
	r := InvokeResponse{Schema: InvokeSchema, Status: "ok"}
	if err := ctx.Err(); err != nil {
		return invokeFailure(r, 130, "canceled", err)
	}
	data, err := io.ReadAll(io.LimitReader(in, MaxInvokeBytes+1))
	if err != nil {
		if canceled := ctx.Err(); canceled != nil {
			return invokeFailure(r, 130, "canceled", canceled)
		}
		return invokeFailure(r, 2, "invalid-request", err)
	}
	if canceled := ctx.Err(); canceled != nil {
		return invokeFailure(r, 130, "canceled", canceled)
	}
	if len(data) > MaxInvokeBytes {
		return invokeFailure(r, 2, "invalid-request", fmt.Errorf("request exceeds %d bytes", MaxInvokeBytes))
	}
	if err = checkInvokeDepth(data); err != nil {
		return invokeFailure(r, 2, "invalid-request", err)
	}
	var req InvokeRequest
	if err = blueprintstate.Strict(data, &req); err != nil {
		return invokeFailure(r, 2, "invalid-request", err)
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(data, &fields)
	if _, present := fields["requestId"]; present {
		if !validCorrelation(req.RequestID) {
			return invokeFailure(r, 2, "invalid-request", fmt.Errorf("requestId must be 1..128 printable ASCII characters"))
		}
		r.RequestID = &req.RequestID
	}
	if !validCorrelation(req.Action) {
		return invokeFailure(r, 2, "invalid-request", fmt.Errorf("invalid action identity"))
	}
	r.Action = &req.Action
	if req.Schema != InvokeSchema {
		return invokeFailure(r, 2, "unsupported-schema", fmt.Errorf("unsupported invocation schema"))
	}
	spec, err := findAction(req.Action)
	if err != nil {
		return invokeFailure(r, 2, "unknown-action", err)
	}
	params, err := spec.decode(req.Parameters)
	if err != nil {
		return invokeFailure(r, 2, "invalid-parameters", err)
	}
	if req.Action != "blueprint.assignment.apply" && req.Action != "governance.session.adopt" && req.Action != "governance.session.resume" && (caller.Principal.Actor != "" || caller.Principal.Authority != "") {
		return invokeFailure(r, 2, "invalid-parameters", fmt.Errorf("caller attribution applies only to assignment or Session mutation"))
	}
	c, err := resolveActionContext(req.Context.Project, caller)
	if err != nil {
		return invokeFailure(r, 1, "invalid-context", err)
	}
	if err = ctx.Err(); err != nil {
		return invokeFailure(r, 130, "canceled", err)
	}
	result, code, err := spec.run(ctx, c, params)
	if err != nil {
		name := "execution-failed"
		code = 1
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			name = "canceled"
			code = 130
		}
		if errors.Is(err, projecthistory.ErrConflict) || errors.Is(err, projecthistory.ErrStale) {
			name = "stale-revision"
		}
		var typed *actionFailure
		if errors.As(err, &typed) {
			name = typed.code
			code = typed.exit
		}
		return invokeFailure(r, code, name, err)
	}
	r.Result = result
	r.ExitCode = code
	if code == 3 {
		r.Status = "blocked"
	}
	return r
}

type actionFailure struct {
	code    string
	exit    int
	message string
}

func (e *actionFailure) Error() string { return e.message }
func writeInvocation(out, errs io.Writer, r InvokeResponse) int {
	b, err := json.Marshal(r)
	if err == nil {
		b = append(b, '\n')
		var n int
		n, err = out.Write(b)
		if err == nil && n != len(b) {
			err = io.ErrShortWrite
		}
	}
	if err != nil {
		fmt.Fprintln(errs, "invocation output failed:", err)
		return 1
	}
	return r.ExitCode
}
