package devhost

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRealBuildFailureAndRestart(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module reloadprobe\n\ngo 1.26.0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	write := func(value string) {
		s := fmt.Sprintf("package main\nimport(\"os\";\"time\")\nfunc main(){os.WriteFile(\"marker\",[]byte(%q),0600);for{time.Sleep(time.Second)}}\n", value)
		if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(s), 0600); err != nil {
			t.Fatal(err)
		}
	}
	wait := func(value string) {
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			b, _ := os.ReadFile(filepath.Join(root, "marker"))
			if string(b) == value {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatalf("child did not publish %s", value)
	}
	write("one")
	m, err := New(Config{Root: root, Package: "."})
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = m.Rebuild(ctx); err != nil {
		t.Fatal(err)
	}
	wait("one")
	first := m.PID()
	if first == 0 {
		t.Fatal("no child")
	}
	os.WriteFile(filepath.Join(root, "main.go"), []byte("not valid Go"), 0600)
	if err = m.Rebuild(ctx); err == nil || m.PID() != first {
		t.Fatal("compile error killed last good process", err)
	}
	write("two")
	if err = m.Rebuild(ctx); err != nil {
		t.Fatal(err)
	}
	wait("two")
	if m.PID() == first {
		t.Fatal("valid Go change did not restart")
	}
	m.Close()
	if m.PID() != 0 {
		t.Fatal("child not closed")
	}
}
