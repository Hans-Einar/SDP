// Package preview owns the source-based static preview/check service.
package preview

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
	"github.com/Hans-Einar/SDP/SDUI/go/prototype"
	"github.com/Hans-Einar/SDP/SDUI/go/sourcewatch"
	"os"
	"path/filepath"
)

type Options struct {
	Source, Entry, Revision, Output string
	Check                           bool
}

// Generate is the shared typed service for source preview/readiness. Callers own
// command parsing and output presentation; this service never parses CLI tokens.
func Generate(options Options) (any, error) {
	if options.Check {
		_, r, e := prototype.Check(options.Source, options.Entry, options.Revision)
		if e != nil {
			return nil, e
		}
		return r, nil
	}
	c := sourcewatch.Read(options.Source, parser.MaxBytes)
	if c.Err != nil {
		return nil, c.Err
	}
	if options.Revision != "" && c.Hash != options.Revision {
		return nil, fmt.Errorf("stale: SDUI source changed")
	}
	document, roots, e := parser.Compile(string(c.Bytes))
	if e != nil {
		return nil, e
	}
	root := roots[options.Entry]
	if root == nil || root.Kind != "frame" {
		return nil, fmt.Errorf("entry must be a defined frame")
	}
	text, e := presentation.Combined(root, 160, document)
	if e != nil {
		return nil, e
	}
	if options.Output == "" {
		return nil, fmt.Errorf("output required")
	}
	out, e := filepath.Abs(options.Output)
	if e != nil {
		return nil, e
	}
	info, e := os.Lstat(out)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("output must be an existing real empty directory")
	}
	entries, e := os.ReadDir(out)
	if e != nil {
		return nil, e
	}
	if len(entries) != 0 {
		return nil, fmt.Errorf("output must be empty")
	}
	latest := sourcewatch.Read(options.Source, parser.MaxBytes)
	if latest.Err != nil {
		return nil, latest.Err
	}
	if latest.Hash != c.Hash {
		return nil, fmt.Errorf("stale: source changed during generation")
	}
	files := map[string][]byte{"entry.md": []byte(text), "source.sdui": c.Bytes}
	hashes := map[string]string{}
	for name, data := range files {
		f, e := os.OpenFile(filepath.Join(out, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return nil, e
		}
		_, e = f.Write(data)
		closeErr := f.Close()
		if e != nil {
			return nil, e
		}
		if closeErr != nil {
			return nil, closeErr
		}
		hashes[name] = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	metadata := map[string]any{"schema": "sdui-combined/1", "revision": c.Hash, "source": options.Source, "frame": options.Entry, "spanUnits": "UTF-8 bytes; end exclusive", "outputs": hashes}
	result := map[string]string{"schema": "sdptool/0.2", "operation": "sdui-preview", "entry": filepath.Join(out, "entry.md"), "directory": out, "revision": c.Hash}
	// Keep the legacy 0.2 protocol bytes; identify the development profile when
	// serving 0.3. These are static artifacts, never connected readiness claims.
	if document.Profile == "sdui/0.3" {
		metadata["profile"] = document.Profile
		result["profile"] = document.Profile
	}
	manifest, _ := json.MarshalIndent(metadata, "", "  ")
	if e = os.WriteFile(filepath.Join(out, "sdui.json"), manifest, 0600); e != nil {
		return nil, e
	}
	return result, nil
}
