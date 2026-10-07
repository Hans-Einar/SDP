// sdui-preview is the SDUI-owned static preview/check adapter for native clients.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
	"github.com/Hans-Einar/SDP/SDUI/go/prototype"
	"github.com/Hans-Einar/SDP/SDUI/go/sourcewatch"
	"os"
	"path/filepath"
)

func run() error {
	source := flag.String("source", "", "Original SDUI source")
	entry := flag.String("entry", "", "Root frame")
	revision := flag.String("revision", "", "Required SHA256")
	output := flag.String("output", "", "Caller-owned empty bundle directory")
	check := flag.Bool("check", false, "Return local prototype readiness only")
	flag.Parse()
	if *check {
		_, r, e := prototype.Check(*source, *entry, *revision)
		if e != nil {
			return e
		}
		return json.NewEncoder(os.Stdout).Encode(r)
	}
	c := sourcewatch.Read(*source, parser.MaxBytes)
	if c.Err != nil {
		return c.Err
	}
	if *revision != "" && c.Hash != *revision {
		return fmt.Errorf("stale: SDUI source changed")
	}
	document, roots, e := parser.Compile(string(c.Bytes))
	if e != nil {
		return e
	}
	root := roots[*entry]
	if root == nil || root.Kind != "frame" {
		return fmt.Errorf("entry must be a defined frame")
	}
	text, e := presentation.Combined(root, 160, document)
	if e != nil {
		return e
	}
	if *output == "" {
		return fmt.Errorf("output required")
	}
	out, e := filepath.Abs(*output)
	if e != nil {
		return e
	}
	info, e := os.Lstat(out)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("output must be an existing real empty directory")
	}
	entries, e := os.ReadDir(out)
	if e != nil {
		return e
	}
	if len(entries) != 0 {
		return fmt.Errorf("output must be empty")
	}
	latest := sourcewatch.Read(*source, parser.MaxBytes)
	if latest.Err != nil {
		return latest.Err
	}
	if latest.Hash != c.Hash {
		return fmt.Errorf("stale: source changed during generation")
	}
	files := map[string][]byte{"entry.md": []byte(text), "source.sdui": c.Bytes}
	hashes := map[string]string{}
	for name, data := range files {
		f, e := os.OpenFile(filepath.Join(out, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		_, e = f.Write(data)
		closeErr := f.Close()
		if e != nil {
			return e
		}
		if closeErr != nil {
			return closeErr
		}
		hashes[name] = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	manifest, _ := json.MarshalIndent(map[string]any{"schema": "sdui-combined/1", "revision": c.Hash, "source": *source, "frame": *entry, "spanUnits": "UTF-8 bytes; end exclusive", "outputs": hashes}, "", "  ")
	if e = os.WriteFile(filepath.Join(out, "sdui.json"), manifest, 0600); e != nil {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(map[string]string{"schema": "sdptool/0.2", "operation": "sdui-preview", "entry": filepath.Join(out, "entry.md"), "directory": out, "revision": c.Hash})
}
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
