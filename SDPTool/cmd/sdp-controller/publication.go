package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"io"
	"os"
)

func publicationCommand(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("publication requires prepare|inspect|apply|abandon")
	}
	action := args[0]
	if action != "prepare" && action != "inspect" && action != "apply" && action != "abandon" {
		return fmt.Errorf("unknown publication action")
	}
	f := flag.NewFlagSet("publication "+action, flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational store")
	input := f.String("input", "", "exact PublicationRequest JSON for prepare")
	id := f.String("publication", "", "publication UUID for inspect/apply/abandon")
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	if *root == "" || *state == "" || f.NArg() != 0 || action == "prepare" && (*input == "" || *id != "") || action != "prepare" && (*id == "" || *input != "") {
		return fmt.Errorf("use --project ROOT --state DIR and --input FILE (prepare) or --publication UUID (inspect/apply/abandon)")
	}
	project, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	s, e := governance.New(project.Root, project.Area, project.Inventory.ProjectID, *state)
	if e != nil {
		return e
	}
	var v governance.PublicationView
	switch action {
	case "prepare":
		file, e := os.Open(*input)
		if e != nil {
			return e
		}
		defer file.Close()
		b, e := io.ReadAll(io.LimitReader(file, governance.MaxRequest+1))
		if e != nil {
			return e
		}
		var q governance.PublicationRequest
		if e = governance.Decode(b, &q); e != nil {
			return e
		}
		v, e = s.PreparePublication(governance.LocalController(), q)
		if e != nil {
			return e
		}
	case "inspect":
		v, e = s.PublicationRead(*id)
	case "apply":
		v, e = s.ApplyPublication(governance.LocalController(), *id)
	case "abandon":
		v, e = s.AbandonPublication(governance.LocalController(), *id)
	}
	if e != nil {
		return e
	}
	return json.NewEncoder(out).Encode(v)
}
