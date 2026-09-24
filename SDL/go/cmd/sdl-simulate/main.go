package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	"os"
	"strings"
)

func main() {
	actions := flag.String("actions", "examples/edit-apt-cell.sdl", "Explicit SDL source")
	source := flag.String("ui", "examples/edit-apt-cell.sdui", "SDUI source")
	values := flag.String("values", "430,invalid,440", "Ordered simulated drafts")
	flag.Parse()
	a, e := os.ReadFile(*actions)
	must(e)
	s, e := os.ReadFile(*source)
	must(e)
	app, e := application.Load(string(a), string(s))
	must(e)
	defer app.Close()
	report, e := application.Exercise(app, strings.Split(*values, ","))
	must(e)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	must(enc.Encode(report))
}
func must(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
}
