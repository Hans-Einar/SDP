// Package sourcewatch provides bounded content-change notifications shared by
// language hosts. It has no parser, runtime or GUI dependencies.
package sourcewatch

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

type Change struct {
	Sequence uint64
	Hash     string
	Bytes    []byte
	Err      error
}

func Read(path string, limit int) Change {
	c := Change{}
	f, err := os.Open(path)
	if err != nil {
		c.Err = err
		return c
	}
	defer f.Close()
	c.Bytes, c.Err = io.ReadAll(io.LimitReader(f, int64(limit+1)))
	if c.Err != nil {
		return c
	}
	sum := sha256.Sum256(c.Bytes)
	c.Hash = fmt.Sprintf("%x", sum)
	return c
}
func Watch(ctx context.Context, path string, limit int, interval time.Duration) <-chan Change {
	return Poll(ctx, interval, func() Change { return Read(path, limit) })
}

// Poll coalesces intermediate snapshots while keeping latest content and errors.
// Probe must be bounded and must not mutate a published model.
func Poll(ctx context.Context, interval time.Duration, probe func() Change) <-chan Change {
	out := make(chan Change, 1)
	if interval < 10*time.Millisecond {
		interval = 10 * time.Millisecond
	}
	go func() {
		defer close(out)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		last := ""
		var sequence uint64
		for {
			sequence++
			c := probe()
			c.Sequence = sequence
			key := c.Hash
			if c.Err != nil {
				key += "|" + c.Err.Error()
			}
			if key != last {
				last = key
				select {
				case <-ctx.Done():
					return
				default:
				}
				select {
				case out <- c:
				default:
					select {
					case <-out:
					default:
					}
					select {
					case out <- c:
					case <-ctx.Done():
						return
					}
				}
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return out
}
