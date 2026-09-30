package install

import (
	"context"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/bootstrap"
	"github.com/Hans-Einar/SDP/SDPTool/presentation"
	"io"
	"os"
)

func Run(ctx context.Context, root, op string, args []string, out, errs io.Writer) int {
	fs := flag.NewFlagSet(op, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	o := Options{Root: root, Operation: op}
	fs.StringVar(&o.Artifact, "artifact", "", "explicit local development descriptor")
	fs.StringVar(&o.PreviousArtifact, "previous-artifact", "", "original development descriptor")
	fs.StringVar(&o.Manifest, "manifest", "", "adoption JSON/YAML")
	fs.BoolVar(&o.AllowUnreleased, "allow-unreleased", false, "allow development artifacts")
	dest := fs.String("plan-output", "", "new plan file outside project")
	jsonMode := fs.Bool("json", false, "machine-readable result")
	apply := fs.String("apply", "", "apply saved plan")
	resume := fs.String("resume", "", "resume operation")
	release := fs.String("release", "", "exact signed descriptor path or HTTPS URL")
	key := fs.String("test-key", os.Getenv("SDP_TEST_KEY"), "explicit nonproduction signer key")
	offline := fs.Bool("offline", os.Getenv("SDP_OFFLINE") == "true", "verified cache only")
	e := fs.Parse(args)
	if value := os.Getenv("SDP_OFFLINE"); e == nil && value != "" && value != "true" && value != "false" {
		e = fmt.Errorf("SDP_OFFLINE must be true or false")
	}
	if e == nil && fs.NArg() != 0 {
		e = fmt.Errorf("unexpected positional arguments")
	}
	if e == nil && (*apply != "" || *resume != "") {
		if *apply != "" && *resume != "" || o.Artifact != "" || o.PreviousArtifact != "" || o.Manifest != "" || o.AllowUnreleased || *dest != "" || *release != "" {
			return output(out, errs, *jsonMode, root, op, nil, fail("arguments", 2, "preview/apply/resume flags cannot be combined"))
		}
		canonical, err := Root(root)
		if err != nil {
			return output(out, errs, *jsonMode, root, op, nil, err)
		}
		var r Result
		if *apply != "" {
			var p Plan
			p, err = LoadPlan(*apply)
			if err == nil {
				if p.Operation != op {
					err = fail("arguments", 2, "plan operation mismatch")
				} else {
					r, err = (Executor{}).Apply(ctx, canonical, p)
				}
			}
		} else {
			r, err = (Executor{}).Resume(ctx, canonical, *resume)
		}
		return output(out, errs, *jsonMode, canonical, op, r, err)
	}
	if e == nil && o.Artifact != "" && *release != "" {
		e = fmt.Errorf("--artifact and --release are exclusive")
	}
	if e == nil && o.Artifact == "" && *release == "" {
		*release = os.Getenv("SDP_RELEASE")
		if *release == "" {
			*release = bootstrap.DefaultRelease
		}
	}
	if e != nil {
		return output(out, errs, *jsonMode, root, op, nil, fail("arguments", 2, "%v", e))
	}
	select {
	case <-ctx.Done():
		return output(out, errs, *jsonMode, root, op, nil, fail("canceled", 2, "%v", ctx.Err()))
	default:
	}
	var p Plan
	if o.Artifact != "" {
		p, e = Preview(o)
	} else {
		var v bootstrap.Verified
		v, e = bootstrap.Resolve(ctx, bootstrap.Config{Descriptor: *release, TestKey: *key, Offline: *offline, CacheDir: os.Getenv("SDP_CACHE_DIR")})
		if e == nil {
			var root string
			root, e = Root(o.Root)
			if e == nil {
				p, e = previewInput(o, root, Input{Path: v.Path, SHA256: v.SHA256, Bytes: v.Bytes, Signature: v.Signature, KeyID: v.KeyID, Provenance: v.Provenance})
			}
		} else {
			e = fail("verification", 4, "%v", e)
		}
	}
	if e == nil && *dest != "" {
		e = SavePlan(p, *dest)
	}
	if e == nil && !p.CanApply {
		e = fail("conflict", 3, "preview contains %d conflicts", len(p.Conflicts))
	}
	return output(out, errs, *jsonMode, p.ProjectRoot, op, &p, e)
}
func output(out, errs io.Writer, jsonMode bool, root, op string, result any, e error) int {
	status := "ok"
	exit := 0
	var problem *Error
	if e != nil {
		status = "error"
		var ok bool
		problem, ok = e.(*Error)
		if !ok {
			problem = &Error{"io", e.Error(), 4}
		}
		exit = problem.Exit
	}
	envelope := map[string]any{"schemaVersion": Protocol, "operation": op, "projectRoot": root, "status": status}
	if result != nil {
		envelope["result"] = result
	}
	if problem != nil {
		envelope["error"] = problem
	}
	if err := presentation.Default().Write(out, envelope, jsonMode); err != nil {
		fmt.Fprintln(errs, err)
		return 4
	}
	if !jsonMode && problem != nil {
		fmt.Fprintln(errs, problem.Error())
	}
	return exit
}

// ArgumentError preserves the installation command's error envelope and exit
// contract when the outer CLI rejects a shared output-selection option.
func ArgumentError(root, op string, jsonMode bool, out, errs io.Writer, err error) int {
	return output(out, errs, jsonMode, root, op, nil, fail("arguments", 2, "%v", err))
}
