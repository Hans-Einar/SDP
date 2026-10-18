//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly

package sdptool

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProgramCancellationStopsOwnedDescendants(t *testing.T) {
	root, d := programFixture(t)
	if err := os.WriteFile(filepath.Join(root, "runner"), []byte("#!/bin/sh\n(sleep 1; echo leaked > descendant) &\necho started > started\nwait\n"), 0700); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	finished := make(chan int, 1)
	go func() {
		var out, errs bytes.Buffer
		finished <- Run(ctx, []string{root, "run", "--program", d.ID}, &out, &errs)
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(root, "started")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("runner did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case code := <-finished:
		if code == 0 {
			t.Fatal("cancellation returned success")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancellation did not finish")
	}
	time.Sleep(1100 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(root, "descendant")); !os.IsNotExist(err) {
		t.Fatal("descendant survived cancellation")
	}
}
