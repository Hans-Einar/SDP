// Package reload prepares SDL action models on content changes. Hosts publish
// and rebind at their owner-thread boundary after complete validation.
package reload

import (
	"context"
	"github.com/Hans-Einar/SDP/SDUI/go/sourcewatch"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"time"
)

type Candidate struct {
	Sequence uint64
	Hash     string
	Program  *parser.Program
	Err      error
}

func Watch(ctx context.Context, path string, interval time.Duration) <-chan Candidate {
	out := make(chan Candidate, 1)
	go func() {
		defer close(out)
		for c := range sourcewatch.Watch(ctx, path, parser.MaxBytes, interval) {
			candidate := Candidate{Sequence: c.Sequence, Hash: c.Hash, Err: c.Err}
			if c.Err == nil {
				candidate.Program, candidate.Err = parser.CompileActions(string(c.Bytes))
			}
			select {
			case out <- candidate:
			default:
				select {
				case <-out:
				default:
				}
				select {
				case out <- candidate:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
