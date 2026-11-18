package sdptool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

func TestPreviewBundleCapacity(t *testing.T) {
	// Reuse one backing slice: the budget counts published bytes per file.
	data := make([]byte, 1<<20)
	for _, tc := range []struct {
		name         string
		files, bytes int
		reject       bool
	}{
		{"file-boundary", 512, 0, false},
		{"file-overflow", 513, 0, true},
		{"byte-boundary", 128, 1 << 20, false},
		{"byte-overflow", 129, 1 << 20, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string][]byte{}
			for n := 0; n < tc.files; n++ {
				files[fmt.Sprintf("file-%03d", n)] = data[:tc.bytes]
			}
			err := checkPreviewBundle(files, "bundle")
			if (err != nil) != tc.reject {
				t.Fatalf("reject=%v: %v", tc.reject, err)
			}
			if err != nil && (!strings.Contains(err.Error(), "512 files or 128 MiB") || !strings.Contains(err.Error(), fmt.Sprintf("actual: %d files", tc.files))) {
				t.Fatal(err)
			}
		})
	}
}

func TestPreviewAndSourceReadingBeyondOldBundleLimit(t *testing.T) {
	root, inventory := projectFixture(t)
	inventory.Models = []Model{{"model", "Demo", "System.design", "design-core/0.6"}}
	saveInventory(t, root, inventory)
	const header = "language design-core version 0.6.\n"
	text := header + "system Demo.\n"
	facts := ""
	for n := 0; n < 127; n++ {
		name := fmt.Sprintf("Unit%03d", n)
		text += fmt.Sprintf("includes %q.\n", name+".design")
		facts += fmt.Sprintf("Demo contains %s.\n", name)
		if err := os.WriteFile(filepath.Join(root, name+".design"), []byte(header+"unit "+name+".\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	syntax, err := parser.Parse(text + facts)
	if err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(root, "System.design")
	if err = os.WriteFile(entry, []byte(parser.FormatFile(syntax)), 0600); err != nil {
		t.Fatal(err)
	}
	p, err := DiscoverWithOptions(root, DiscoverOptions{Navigation: "lazy-sdl"})
	if err != nil {
		t.Fatal(err)
	}
	p.Inventory = inventory
	preview, err := Preview(context.Background(), PreviewOptions{Source: entry, Output: filepath.Join(t.TempDir(), "preview")})
	if err != nil {
		t.Fatal(err)
	}
	source, err := SelectProject(context.Background(), p, "model", PreviewOptions{URI: "sdl-view://" + inventory.ProjectID + "/VP11?projection=source", Revision: preview.Revision, Output: filepath.Join(t.TempDir(), "source")})
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range []Result{preview, source} {
		count := 0
		err = filepath.WalkDir(result.Directory, func(_ string, d os.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if !d.IsDir() {
				count++
			}
			return nil
		})
		if err != nil || count <= 128 || count > 512 || len(result.SourceLinks) != 128 {
			t.Fatalf("files=%d sources=%d error=%v", count, len(result.SourceLinks), err)
		}
	}
}
