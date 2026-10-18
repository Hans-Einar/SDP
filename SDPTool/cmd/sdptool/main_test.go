package main

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

func TestContextInputCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	result := make(chan error, 1)
	go func() { _, e := (contextInput{ctx, r}).Read(make([]byte, 16)); result <- e }()
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked request did not cancel")
	}
}
