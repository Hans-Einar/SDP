//go:build desktop

package main

import (
	"flag"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func main() {
	entry := flag.String("entry", "page", "Root frame")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage: sdui-fyne [-entry page] source.sdui")
		os.Exit(2)
	}
	source, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	_, roots, err := parser.Compile(string(source))
	if err != nil {
		panic(err)
	}
	root := roots[*entry]
	if root == nil {
		panic("Unknown entry")
	}
	a := app.NewWithID("no.sdp.sdui.prototype")
	w := a.NewWindow("SDUI · native prototype")
	status := widget.NewLabel("Lokal prototype · ingen SDL-runtime")
	v := fynehost.New(root)
	v.OnStatus = func(err error) {
		if err != nil {
			status.SetText(err.Error())
		} else {
			status.SetText("OK")
		}
	}
	// The demo binding is explicit Go application code, not inferred from SDL names.
	if _, ok := v.Controls[*entry+"/ok"]; ok {
		v.Actions[*entry+"/ok"] = func(path, value string) error { fmt.Println("local callback:", path, value); return nil }
	}
	w.SetContent(container.NewBorder(nil, status, nil, nil, v.Container))
	w.Resize(fyne.NewSize(1280, 800))
	w.SetOnClosed(v.Close)
	w.ShowAndRun()
}
