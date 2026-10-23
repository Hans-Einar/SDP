package governancecontroller

import (
	"context"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type waitingHost struct {
	*fakeHost
	entered chan struct{}
	release chan struct{}
}

type namedHost struct {
	*fakeHost
	threadID string
}

func (h *namedHost) StartWorkThread(ctx context.Context, a, b, c, d string, config map[string]any) (governancehost.Thread, error) {
	thread, err := h.fakeHost.StartWorkThread(ctx, a, b, c, d, config)
	thread.ID = h.threadID
	return thread, err
}

func (h *waitingHost) StartTurn(ctx context.Context, thread, input, payload string) (governancehost.Turn, error) {
	close(h.entered)
	select {
	case <-h.release:
	case <-ctx.Done():
		return governancehost.Turn{}, ctx.Err()
	}
	return h.fakeHost.StartTurn(ctx, thread, input, payload)
}

func TestDisjointControllerLaunchProgressesWhileOtherHostWaits(t *testing.T) {
	a, fa, id := setup(t)
	r, e := a.Service.Read(id)
	if e != nil {
		t.Fatal(e)
	}
	q := governance.StartRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RequestID: governancehost.UUID(), InputRevision: r.Run.InputRevision, RoutineDigest: r.Run.RoutineDigest, Selection: r.Run.Selection, Assignment: governance.Assignment{Outcome: "second fixture", WritablePaths: []string{"second.txt"}, ExcludedPaths: []string{"unrelated.txt"}}}
	second, e := a.Service.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	b := *a
	b.Epoch = governancehost.UUID()
	fb := &fakeHost{t: t, c: &b, runID: second.Run.ID, mode: "chatgpt"}
	b.Host = &namedHost{fakeHost: fb, threadID: "thread-b"}
	h := &waitingHost{fakeHost: fa, entered: make(chan struct{}), release: make(chan struct{})}
	a.Host = h
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := a.Launch(ctx, id, "Fix the selected fixture"); done <- err }()
	select {
	case <-h.entered:
	case <-ctx.Done():
		t.Fatal("first host did not start")
	}
	if _, err := a.Launch(ctx, id, "Fix the selected fixture"); err == nil {
		t.Error("duplicate same-run controller accepted")
	}
	result, e := b.Launch(ctx, second.Run.ID, "Fix the selected fixture")
	close(h.release)
	firstErr := <-done
	if e != nil || firstErr != nil {
		t.Fatal(e, firstErr)
	}
	if result.Run.State != "running" || fb.starts != 1 || fa.starts != 1 {
		t.Fatal(result, fb.starts, fa.starts)
	}
}
