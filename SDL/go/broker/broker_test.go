package broker

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/reader"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type recording struct {
	mu         sync.Mutex
	deliveries []reader.Delivery
}

func (r *recording) Open(_ context.Context, d reader.Delivery) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deliveries = append(r.deliveries, d)
	return nil
}
func setup(t *testing.T) (*Broker, string, *recording) {
	t.Helper()
	root := t.TempDir()
	source := filepath.Join(root, "source.design")
	os.WriteFile(source, []byte("language design-core version 0.5.\nunit Host.\n"), 0600)
	s, e := NewStore(filepath.Join(root, "store"), 1<<20, 2)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(s.Close)
	r := &recording{}
	return New(map[string]Project{"demo": {Source: source}}, reader.Registry{"xfmd": r}, s), source, r
}
func req(seq uint64) Request {
	return Request{URI: "sdl-view://demo/VP02?diagram=VP02-roots", Window: "one", Client: "test", Sequence: seq, Open: true}
}
func TestCacheInvalidationLeaseRestartAndQuota(t *testing.T) {
	b, source, r := setup(t)
	ctx := context.Background()
	a, e := b.Select(ctx, req(1))
	if e != nil {
		t.Fatal(e)
	}
	second, e := b.Select(ctx, req(2))
	if e != nil || !second.Cached || second.Entry != a.Entry {
		t.Fatal(second, e)
	}
	if a.Lease == second.Lease {
		t.Fatal("shared lease token")
	}
	os.WriteFile(source, []byte("language design-core version 0.5.\nunit Next.\n"), 0600)
	next, e := b.Select(ctx, req(3))
	if e != nil || next.Cached || next.Revision == a.Revision {
		t.Fatal(next, e)
	}
	if len(r.deliveries) != 3 {
		t.Fatal("missing delivery")
	}
	b.Store.Close()
	b.Store, e = NewStore(b.Store.root, 1<<20, 2)
	if e != nil {
		t.Fatal(e)
	}
	defer b.Store.Close()
	if e = b.Store.Sweep(); e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(a.Entry); e != nil {
		t.Fatal("restart removed leased images", e)
	}
	os.WriteFile(source, []byte("language design-core version 0.5.\nunit Third.\n"), 0600)
	if _, e = b.Select(ctx, req(4)); e == nil {
		t.Fatal("leased quota exceeded")
	}
	b.Store.Release(a.Lease)
	if _, e = b.Select(ctx, req(5)); e == nil {
		t.Fatal("other reader lease ignored")
	}
	b.Store.Release(second.Lease)
	third, e := b.Select(ctx, req(6))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(a.Entry); !os.IsNotExist(e) {
		t.Fatal("released bundle not evicted")
	}
	os.WriteFile(source, []byte("invalid"), 0600)
	if _, e = b.Select(ctx, req(7)); e == nil {
		t.Fatal("invalid source")
	}
	if _, e = os.Stat(third.Entry); e != nil {
		t.Fatal("failed reload lost document")
	}
	if _, e = b.Select(ctx, req(6)); e == nil {
		t.Fatal("stale request accepted")
	}
}

type blockingRenderer struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func (r *blockingRenderer) Identity() string { return "blocking" }
func (r *blockingRenderer) Render(ctx context.Context, _ string) ([]byte, error) {
	r.once.Do(func() { close(r.entered) })
	select {
	case <-r.release:
		return []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
func TestSupersededAndSourceChangedDuringRender(t *testing.T) {
	b, source, r := setup(t)
	renderer := &blockingRenderer{entered: make(chan struct{}), release: make(chan struct{})}
	b.Projects["demo"] = Project{source, renderer}
	done := make(chan error, 1)
	go func() { _, e := b.Select(context.Background(), req(1)); done <- e }()
	<-renderer.entered
	bad := req(2)
	bad.URI = "sdl-view://demo/VP99"
	if _, e := b.Select(context.Background(), bad); e == nil {
		t.Fatal("bad query")
	}
	close(renderer.release)
	if e := <-done; e == nil {
		t.Fatal("older result delivered after newer click")
	}
	if len(r.deliveries) != 0 {
		t.Fatal("stale open")
	}
	renderer = &blockingRenderer{entered: make(chan struct{}), release: make(chan struct{})}
	b.Projects["demo"] = Project{source, renderer}
	os.WriteFile(source, []byte("language design-core version 0.5.\nunit Second.\n"), 0600)
	go func() { _, e := b.Select(context.Background(), req(3)); done <- e }()
	<-renderer.entered
	os.WriteFile(source, []byte("language design-core version 0.5.\nunit Changed.\n"), 0600)
	close(renderer.release)
	if e := <-done; e == nil {
		t.Fatal("changed revision delivered")
	}
}
func TestStoreByteQuota(t *testing.T) {
	s, e := NewStore(filepath.Join(t.TempDir(), "private"), 10, 1)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	_, e = s.Put(documents.Hash([]byte("x")), &documents.Bundle{Files: map[string][]byte{"entry.md": make([]byte, 11)}})
	if e == nil {
		t.Fatal("quota")
	}
}
func TestIPCSelectReleaseAndRestart(t *testing.T) {
	b, _, _ := setup(t)
	root := t.TempDir()
	os.Chmod(root, 0700)
	socket := filepath.Join(root, "broker.sock")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- b.Serve(ctx, socket) }()
	for i := 0; i < 100; i++ {
		if _, e := os.Stat(socket); e == nil {
			break
		}
		time.Sleep(time.Millisecond * 5)
	}
	r, e := Send(context.Background(), socket, Call{Operation: "select", Request: req(1)})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = Send(context.Background(), socket, Call{Operation: "release", Lease: r.Result.Lease}); e != nil {
		t.Fatal(e)
	}
	if _, e = Send(context.Background(), socket, Call{Operation: "sweep"}); e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(r.Result.Entry); !os.IsNotExist(e) {
		t.Fatal("released artifact retained after sweep")
	}
	cancel()
	if e = <-done; e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(socket); !os.IsNotExist(e) {
		t.Fatal("socket left behind")
	}
}
