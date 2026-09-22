// Package broker coordinates bounded on-demand projections and immutable reader leases.
package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/reader"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"io"
	"os"
	"sync"
)

type Project struct {
	Source   string
	Renderer documents.Renderer
}
type Request struct {
	URI, Window, Client string
	Sequence            uint64
	Open                bool
}
type Result struct {
	Entry, Revision, Lease string
	Sequence               uint64
	Cached                 bool
}
type Broker struct {
	mu       sync.Mutex
	Projects map[string]Project
	Readers  reader.Registry
	Store    *Store
	Endpoint string
	jobs     chan struct{}
}

func New(projects map[string]Project, readers reader.Registry, store *Store) *Broker {
	p := map[string]Project{}
	for k, v := range projects {
		p[k] = v
	}
	return &Broker{Projects: p, Readers: readers, Store: store, jobs: make(chan struct{}, 4)}
}
func read(path string) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
	if len(b) > parser.MaxBytes {
		return nil, fmt.Errorf("source limit")
	}
	return b, e
}
func (b *Broker) Select(ctx context.Context, req Request) (Result, error) {
	s, e := viewpoint.ParseURI(req.URI)
	if e != nil {
		return Result{}, e
	}
	project, ok := b.Projects[s.Project]
	if !ok {
		return Result{}, fmt.Errorf("unregistered project")
	}
	d := reader.Delivery{Window: req.Window, Pane: s.Target, Client: req.Client, Sequence: req.Sequence, Entry: "/validation"}
	if e = d.Validate(); e != nil {
		return Result{}, e
	}
	route := req.Client + "/" + req.Window + "/" + s.Target
	if e = b.Store.Begin(route, req.Sequence); e != nil {
		return Result{}, e
	}
	select {
	case b.jobs <- struct{}{}:
		defer func() { <-b.jobs }()
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
		return Result{}, fmt.Errorf("projection busy")
	}
	source, e := read(project.Source)
	if e != nil {
		return Result{}, e
	}
	v, e := viewpoint.New(string(source))
	if e != nil {
		return Result{}, e
	}
	if e = s.Query.Validate(v); e != nil {
		return Result{}, e
	}
	rid := "none"
	if project.Renderer != nil {
		rid = project.Renderer.Identity()
	}
	query, _ := json.Marshal(s.Query.Canonical())
	key := documents.Hash([]byte(v.Revision + viewpoint.Version + rid + string(query)))
	entry, hit := b.Store.Lookup(key)
	if !hit {
		bundle, e := documents.Selected(ctx, v, s.Query, project.Renderer)
		if e != nil {
			return Result{}, e
		}
		entry, e = b.Store.Put(key, bundle)
		if e != nil {
			return Result{}, e
		}
	}
	// Source may have changed while a renderer was working; never open a stale revision.
	current, e := read(project.Source)
	if e != nil {
		return Result{}, e
	}
	if documents.Hash(current) != v.Revision {
		return Result{}, fmt.Errorf("source changed during projection")
	}
	if e = ctx.Err(); e != nil {
		return Result{}, e
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.Store.Latest(route, req.Sequence) {
		return Result{}, fmt.Errorf("superseded request")
	}
	lease, e := b.Store.Acquire(key)
	if e != nil {
		return Result{}, e
	}
	d.Entry = entry
	d.Revision = v.Revision
	d.Lease = lease
	d.Broker = b.Endpoint
	if req.Open {
		if e = b.Readers.Open(ctx, s.Consumer, d); e != nil {
			// An unacknowledged open may still be visible. Retain the lease conservatively.
			return Result{Entry: entry, Revision: v.Revision, Lease: lease, Sequence: req.Sequence, Cached: hit}, e
		}
	}
	return Result{entry, v.Revision, lease, req.Sequence, hit}, nil
}
