// sdui-preview is a thin development adapter; distribution uses sdptool sdui.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/preview"
	"io"
	"os"
)

func execute(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("sdui-preview", flag.ContinueOnError)
	var options preview.Options
	flags.StringVar(&options.Source, "source", "", "Original SDUI source")
	flags.StringVar(&options.Entry, "entry", "", "Root frame")
	flags.StringVar(&options.Revision, "revision", "", "Required SHA256")
	flags.StringVar(&options.Output, "output", "", "Caller-owned empty bundle directory")
	flags.BoolVar(&options.Check, "check", false, "Return local prototype readiness only")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	result, err := preview.Generate(options)
	if err != nil {
		return err
	}
	return json.NewEncoder(stdout).Encode(result)
}
func main() {
	if err := execute(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
