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

func evidenceCommand(args []string, output io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("evidence requires capture|inspect|attach|rework")
	}
	action := args[0]
	if action != "capture" && action != "inspect" && action != "attach" && action != "rework" {
		return fmt.Errorf("unknown evidence action")
	}
	f := flag.NewFlagSet("evidence "+action, flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	run := f.String("run", "", "run UUID for capture/inspect")
	input := f.String("input", "", "exact JSON request for attach/rework")
	jsonOutput := f.Bool("json", false, "JSON output (capture/attach/rework always emit JSON)")
	dependencies := []string{}
	f.Func("dependency", "additional literal candidate dependency; repeat up to 64 times", func(v string) error { dependencies = append(dependencies, v); return nil })
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	read := action == "capture" || action == "inspect"
	if *root == "" || *state == "" || f.NArg() != 0 || (read && (*run == "" || *input != "")) || (!read && (*input == "" || *run != "")) || (action != "capture" && len(dependencies) > 0) {
		return fmt.Errorf("evidence requires --project ROOT --state DIR and --run UUID (capture/inspect) or --input FILE (attach/rework)")
	}
	p, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	s, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if e != nil {
		return e
	}
	if action == "capture" {
		v, err := s.CaptureCandidate(*run, dependencies)
		if err != nil {
			return err
		}
		return json.NewEncoder(output).Encode(v)
	}
	if action == "inspect" {
		v, err := s.Evidence(*run)
		if err != nil {
			return err
		}
		if *jsonOutput {
			return json.NewEncoder(output).Encode(v)
		}
		if _, err = fmt.Fprintf(output, "PG evidence: %d record(s)\nRun: %q / %q / revision %d\n", len(v.Items), v.RunID, v.State, v.Revision); err != nil {
			return err
		}
		for _, item := range v.Items {
			q := item.Record.Request
			if _, err = fmt.Fprintf(output, "\n%s: %q %q by %q\n  Candidate: %s\n  Applicable: %t; artifact current: %t\n  Note: %q\n", q.OperationID, q.Kind, q.Verdict, q.Issuer, q.Candidate.Digest, item.Applicable, item.ArtifactCurrent, q.Note); err != nil {
				return err
			}
			for _, c := range item.Conditions {
				if _, err = fmt.Fprintf(output, "  Condition: %q\n", c); err != nil {
					return err
				}
			}
		}
		for _, l := range v.Limitations {
			if _, err = fmt.Fprintf(output, "Limit: %s\n", l); err != nil {
				return err
			}
		}
		return nil
	}
	file, e := os.Open(*input)
	if e != nil {
		return e
	}
	defer file.Close()
	b, e := io.ReadAll(io.LimitReader(file, governance.MaxRequest+1))
	if e != nil {
		return e
	}
	var result governance.Result
	if action == "attach" {
		var q governance.EvidenceRequest
		if e = governance.Decode(b, &q); e != nil {
			return e
		}
		result, e = s.AttachEvidence(governance.LocalController(), q)
	} else {
		var q governance.EvidenceReworkRequest
		if e = governance.Decode(b, &q); e != nil {
			return e
		}
		result, e = s.ReworkEvidence(governance.LocalController(), q)
	}
	if e != nil {
		return e
	}
	return json.NewEncoder(output).Encode(result)
}
