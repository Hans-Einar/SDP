// sdui-preview is a thin development adapter; distribution uses sdptool sdui.
package main

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/preview"
	"io"
	"os"
)

func execute(args []string, stdout io.Writer) error { return preview.Execute(args, stdout) }
func main() {
	if err := execute(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
