package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	model "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/generatedmodel"
	"os"
	"strings"
)

func main() {
	values := flag.String("values", "430,invalid,440", "Ordered simulated drafts")
	flag.Parse()
	app, e := application.New(model.Program(), model.Document(), model.Root())
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
