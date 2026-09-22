package main

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type options struct {
	source, output, format, entry string
	columns                       int
	width, height                 float64
	syntax                        bool
}

func arguments(args []string) (options, error) {
	o := options{format: "ast", columns: 160, width: 1920, height: 1200}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "--syntax-only":
			o.syntax = true
		case "-o", "--format", "--entry", "--columns", "--width", "--height":
			i++
			if i >= len(args) {
				return o, fmt.Errorf("Missing value for %s", a)
			}
			switch a {
			case "-o":
				o.output = args[i]
			case "--format":
				o.format = args[i]
			case "--entry":
				o.entry = args[i]
			case "--width", "--height":
				v, e := strconv.ParseFloat(args[i], 64)
				if e != nil {
					return o, e
				}
				if a == "--width" {
					o.width = v
				} else {
					o.height = v
				}
			case "--columns":
				n, e := strconv.Atoi(args[i])
				if e != nil {
					return o, e
				}
				o.columns = n
			}
		default:
			if strings.HasPrefix(a, "-") && a != "-" {
				return o, fmt.Errorf("Unknown option %s", a)
			}
			if o.source != "" {
				return o, fmt.Errorf("Only one source allowed")
			}
			o.source = a
		}
	}
	if o.source == "" {
		return o, fmt.Errorf("Usage: sdui source|- [--format ast|dump|markdown|svg|prototype-svg|prototype-html] [--entry name] [-o file]")
	}
	return o, nil
}
func report(w io.Writer, e error, code int) int {
	if _, ok := e.(*parser.Diagnostic); !ok {
		e = &parser.Diagnostic{Code: "io", Message: e.Error()}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"error": e})
	return code
}
func execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	o, e := arguments(args)
	if e != nil {
		return report(stderr, e, 2)
	}
	r := stdin
	if o.source != "-" {
		if o.output != "" {
			a, _ := filepath.Abs(o.source)
			b, _ := filepath.Abs(o.output)
			sa, ea := os.Stat(a)
			sb, eb := os.Stat(b)
			if a == b || ea == nil && eb == nil && os.SameFile(sa, sb) {
				return report(stderr, fmt.Errorf("Output cannot overwrite source"), 3)
			}
		}
		f, e := os.Open(o.source)
		if e != nil {
			return report(stderr, e, 3)
		}
		defer f.Close()
		r = f
	}
	source, e := io.ReadAll(io.LimitReader(r, parser.MaxBytes+1))
	if e != nil {
		return report(stderr, e, 3)
	}
	d, e := parser.Parse(string(source))
	if e != nil {
		return report(stderr, e, 2)
	}
	validation := "syntax-only"
	var roots map[string]*parser.Instance
	if !o.syntax {
		roots, e = parser.Normalize(d)
		if e != nil {
			return report(stderr, e, 2)
		}
		validation = "local-profile"
	}
	var out []byte
	if o.format == "ast" {
		out, e = json.MarshalIndent(map[string]any{"astFormat": "sdui-ast/0.2", "validation": validation, "document": parser.Data(d)}, "", "  ")
		out = append(out, '\n')
	} else {
		if o.syntax {
			return report(stderr, fmt.Errorf("Export requires validation"), 2)
		}
		if o.entry == "" {
			for name, n := range roots {
				if n.Kind == "frame" {
					if o.entry != "" {
						return report(stderr, fmt.Errorf("Select --entry"), 2)
					}
					o.entry = name
				}
			}
		}
		root := roots[o.entry]
		if root == nil || root.Kind != "frame" {
			return report(stderr, fmt.Errorf("Entry must be a defined frame"), 2)
		}
		var text string
		switch o.format {
		case "svg":
			var box *layout.Box
			box, e = (&layout.Engine{}).Layout(root, layout.Size{W: o.width, H: o.height})
			if e == nil {
				text, e = svg.Render(box, svg.Options{Width: o.width, Height: o.height})
			}
		case "dump":
			text, e = presentation.Dump(root, o.columns)
		case "markdown":
			text, e = presentation.Markdown(root, o.columns)
		case "prototype-svg":
			text, e = presentation.PrototypeSVG(root)
		case "prototype-html":
			text, e = presentation.PrototypeHTML(root)
		default:
			e = fmt.Errorf("Unknown format %s", o.format)
		}
		if e != nil {
			return report(stderr, e, 2)
		}
		out = []byte(text)
	}
	if e != nil {
		return report(stderr, e, 3)
	}
	if o.output == "" {
		_, e = stdout.Write(out)
	} else {
		e = writeAtomic(o.output, out)
	}
	if e != nil {
		return report(stderr, e, 3)
	}
	return 0
}
func writeAtomic(path string, data []byte) error {
	f, e := os.CreateTemp(filepath.Dir(path), ".sdui-*")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	if _, e = f.Write(data); e != nil {
		f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	return os.Rename(f.Name(), path)
}
func main() { os.Exit(execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }
