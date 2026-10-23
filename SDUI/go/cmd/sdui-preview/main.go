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
	"io"
	"os"
	"path/filepath"
)

func run() error {
	return execute(os.Args[1:], os.Stdout)
}

func execute(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("sdui-preview", flag.ContinueOnError)
	source := flags.String("source", "", "Original SDUI source")
	entry := flags.String("entry", "", "Root frame")
	revision := flags.String("revision", "", "Required SHA256")
	output := flags.String("output", "", "Caller-owned empty bundle directory")
	check := flags.Bool("check", false, "Return local prototype readiness only")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *check {
		_, r, e := prototype.Check(*source, *entry, *revision)
		if e != nil {
			return e
		}
		return json.NewEncoder(stdout).Encode(r)
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
	metadata := map[string]any{"schema": "sdui-combined/1", "revision": c.Hash, "source": *source, "frame": *entry, "spanUnits": "UTF-8 bytes; end exclusive", "outputs": hashes}
	result := map[string]string{"schema": "sdptool/0.2", "operation": "sdui-preview", "entry": filepath.Join(out, "entry.md"), "directory": out, "revision": c.Hash}
	// Keep the legacy 0.2 protocol bytes; identify the development profile when
	// serving 0.3. These are static artifacts, never connected readiness claims.
	if document.Profile == "sdui/0.3" {
		metadata["profile"] = document.Profile
		result["profile"] = document.Profile
	}
	manifest, _ := json.MarshalIndent(metadata, "", "  ")
	if e = os.WriteFile(filepath.Join(out, "sdui.json"), manifest, 0600); e != nil {
		return e
	}
	return json.NewEncoder(stdout).Encode(result)
}
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
