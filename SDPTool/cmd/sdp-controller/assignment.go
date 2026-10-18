package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func assignmentCommand(args []string, output io.Writer) error {
	return assignmentCommandIO(args, output, os.Stdin, os.Stderr)
}

func assignmentCommandIO(args []string, output io.Writer, ownerInput io.Reader, prompt io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("assignment requires context|prepare|revise|retire|inspect|design-check|design-review|design-decide")
	}
	action := args[0]
	if action != "context" && action != "prepare" && action != "revise" && action != "retire" && action != "inspect" && action != "design-check" && action != "design-review" && action != "design-decide" {
		return fmt.Errorf("unsupported assignment operation")
	}
	f := flag.NewFlagSet("assignment "+action, flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	input := f.String("input", "", "bounded JSON request file (all operations except inspect)")
	owner := f.String("owner", "", "owner label attributed by local controller (design-decide only)")
	id := f.String("assignment", "", "optional prepared UUID (inspect)")
	if err := f.Parse(args[1:]); err != nil {
		return err
	}
	if *root == "" || *state == "" || f.NArg() != 0 {
		return fmt.Errorf("explicit --project and --state required")
	}
	if action == "inspect" && *input != "" || action != "inspect" && (*input == "" || *id != "") {
		return fmt.Errorf("only inspect permits --assignment; other operations require --input")
	}
	if (action == "design-decide" && *owner == "") || (action != "design-decide" && *owner != "") {
		return fmt.Errorf("--owner is required only for design-decide")
	}
	p, err := sdptool.Discover(*root)
	if err != nil {
		return err
	}
	service, err := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if err != nil {
		return err
	}
	var request []byte
	if *input != "" {
		file, err := os.Open(*input)
		if err != nil {
			return err
		}
		defer file.Close()
		request, err = io.ReadAll(io.LimitReader(file, governance.MaxRequest+1))
		if err != nil {
			return err
		}
	}
	var result any
	switch action {
	case "context":
		var q governance.PreparationContextRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = service.PreparationContext(q)
		}
	case "prepare":
		var q governance.PrepareRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = service.PrepareAssignment(governance.LocalController(), q)
		}
	case "inspect":
		result, err = service.InspectAssignments(*id)
	case "revise":
		var q governance.RevisePreparationRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = service.ReviseAssignment(governance.LocalController(), q)
		}
	case "retire":
		var q governance.RetirePreparationRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = service.RetireAssignment(governance.LocalController(), q)
		}
	case "design-review":
		var q governance.DesignCheckRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = service.ReviewAssignmentDesign(q)
		}
	case "design-decide":
		var q governance.DesignDecisionRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = captureDesignDecision(service, q, *owner, ownerInput, prompt)
		}
	case "design-check":
		var q governance.DesignCheckRequest
		if err = governance.Decode(request, &q); err == nil {
			result, err = service.CheckAssignmentDesign(q)
		}
	}
	if err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(result)
}
