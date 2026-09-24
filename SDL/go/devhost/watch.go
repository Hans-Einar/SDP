package devhost

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDUI/go/sourcewatch"
)

// Watch observes only Go sources and module metadata below the registered root.
// SDL/SDUI files remain the live app's model reload responsibility.
func Watch(ctx context.Context, root string, interval time.Duration) <-chan sourcewatch.Change {
	return sourcewatch.Poll(ctx, interval, func() sourcewatch.Change {
		hash := sha256.New()
		count, total := 0, 0
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == ".git" || d.Name() == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}
			name := d.Name()
			if !strings.HasSuffix(name, ".go") && name != "go.mod" && name != "go.sum" {
				return nil
			}
			count++
			if count > 4096 {
				return fmt.Errorf("go-watch: more than 4096 source files")
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			data, err := io.ReadAll(io.LimitReader(f, 4<<20+1))
			closeErr := f.Close()
			if err != nil {
				return err
			}
			if closeErr != nil {
				return closeErr
			}
			total += len(data)
			if len(data) > 4<<20 || total > 32<<20 {
				return fmt.Errorf("go-watch: source size limit")
			}
			relative, _ := filepath.Rel(root, path)
			fmt.Fprintf(hash, "%s\x00%d\x00", relative, len(data))
			hash.Write(data)
			return nil
		})
		return sourcewatch.Change{Hash: fmt.Sprintf("%x", hash.Sum(nil)), Err: err}
	})
}
