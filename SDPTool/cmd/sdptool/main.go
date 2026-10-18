package main

import (
	"context"
	tool "github.com/Hans-Einar/SDP/SDPTool"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(tool.RunWithInput(ctx, os.Args[1:], contextInput{ctx, os.Stdin}, os.Stdout, os.Stderr))
}

// Only the executable wraps stdin: a blocked OS read can outlive cancellation,
// but the process exits after writing its response. Never reuse the caller buffer
// in that pending read, and do not impose goroutines on the library reader API.
type contextInput struct {
	ctx context.Context
	r   io.Reader
}

func (in contextInput) Read(p []byte) (int, error) {
	if err := in.ctx.Err(); err != nil {
		return 0, err
	}
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 1)
	go func() { b := make([]byte, len(p)); n, err := in.r.Read(b); ch <- result{b[:n], err} }()
	select {
	case <-in.ctx.Done():
		return 0, in.ctx.Err()
	case r := <-ch:
		return copy(p, r.data), r.err
	}
}
