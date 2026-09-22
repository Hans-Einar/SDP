// Package reload watches SDUI sources and prepares candidates off the UI thread.
// Publishing and UI state ownership remain with the consuming host.
package reload

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	uiruntime "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type Candidate struct {
	Sequence uint64
	Hash     string
	Document *parser.Document
	Root     *parser.Instance
	Err      error
}

// Prepare is an optional pure validation/measurement gate. It must not mutate
// the published session or create native widgets on the watcher goroutine.
type Prepare func(*parser.Instance) error

func Read(path, entry string, sequence uint64, prepare Prepare) Candidate {
	c, source := readSource(path, sequence)
	return compile(c, source, entry, prepare)
}
func readSource(path string, sequence uint64) (Candidate, []byte) {
	c := Candidate{Sequence: sequence}
	f, err := os.Open(path)
	if err != nil {
		c.Err = err
		return c, nil
	}
	defer f.Close()
	source, err := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
	if err != nil {
		c.Err = err
		return c, nil
	}
	hash := sha256.Sum256(source)
	c.Hash = fmt.Sprintf("%x", hash)
	return c, source
}
func compile(c Candidate, source []byte, entry string, prepare Prepare) Candidate {
	if c.Err != nil {
		return c
	}
	d, roots, err := parser.Compile(string(source))
	if err != nil {
		c.Err = err
		return c
	}
	root := roots[entry]
	if root == nil || root.Kind != "frame" {
		c.Err = fmt.Errorf("reload-entry: %s is not a frame", entry)
		return c
	}
	if prepare != nil {
		if err = prepare(root); err != nil {
			c.Err = err
			return c
		}
	}
	c.Document = d
	c.Root = root
	return c
}

// Watch checks content hashes, so editor atomic-renames and same-size rewrites
// are handled. The bounded latest-candidate channel coalesces intermediate saves.
func Watch(ctx context.Context, path, entry string, interval time.Duration, prepare Prepare) <-chan Candidate {
	out := make(chan Candidate, 1)
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
			c, source := readSource(path, sequence)
			key := c.Hash
			if c.Err != nil {
				key += "|" + c.Err.Error()
			}
			if key != last {
				last = key
				c = compile(c, source, entry, prepare)
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

type Controller struct {
	Session    *uiruntime.Session
	Sequence   uint64
	SourceHash string
	LastError  error
}

func (c *Controller) Adopt(candidate Candidate) error {
	if candidate.Sequence <= c.Sequence {
		return fmt.Errorf("stale-candidate: result superseded")
	}
	c.Sequence = candidate.Sequence
	if candidate.Err != nil {
		c.LastError = candidate.Err
		return candidate.Err
	}
	if candidate.Root == nil {
		c.LastError = fmt.Errorf("reload-root: missing candidate root")
		return c.LastError
	}
	if err := c.Session.Reload(candidate.Root); err != nil {
		c.LastError = err
		return err
	}
	c.SourceHash = candidate.Hash
	c.LastError = nil
	return nil
}
