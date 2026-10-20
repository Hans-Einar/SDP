// sdp-controller supplies bounded owner control and protocol inspection.
// Host completion never accepts SDP work or publishes project records.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "attention" {
		if e := attentionCommand(os.Args[2:], os.Stdout, os.Stdin, os.Stderr); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "inbox" {
		if e := inboxCommand(os.Args[2:], os.Stdout); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "assignment" {
		if e := assignmentCommand(os.Args[2:], os.Stdout); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "overview" {
		if e := overviewCommand(os.Args[2:], os.Stdout); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
		return
	}
	if len(os.Args) > 1 && (os.Args[1] == "run" || os.Args[1] == "status" || os.Args[1] == "recover" || os.Args[1] == "attach" || os.Args[1] == "provenance") {
		if e := ownerCommand(os.Args[1:]); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
		return
	}
	if len(os.Args) < 2 || os.Args[1] != "probe" {
		fmt.Fprintln(os.Stderr, "usage: sdp-controller probe|run|status|recover|attach|provenance|overview|inbox|attention|assignment --project ROOT [options]; each command accepts --help")
		os.Exit(2)
	}
	f := flag.NewFlagSet("probe", flag.ExitOnError)
	root := f.String("project", "", "explicit cwd")
	binary := f.String("codex", "codex", "pinned Codex executable")
	mcpBinary := f.String("mcp-binary", "", "optional sdp-mcp binary for ephemeral-thread negotiation")
	state := f.String("state", "", "MCP state directory (read-only inspection)")
	f.Parse(os.Args[2:])
	check := func(e error) {
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
	}
	if *root == "" || f.NArg() != 0 || (*mcpBinary != "" && *state == "") {
		check(fmt.Errorf("explicit --project and paired --mcp-binary/--state required"))
	}
	cwd, e := filepath.Abs(*root)
	check(e)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	c, e := governancehost.Open(ctx, *binary, cwd)
	check(e)
	mode, e := c.AccountMode(ctx)
	if e != nil {
		c.Close()
		check(e)
	}
	out := map[string]any{"schema": "sdp-host-probe/1", "hostVersion": governancehost.HostVersion, "initialized": true, "accountMode": mode, "managedAccount": mode == "chatgpt", "modelSubmitted": false}
	if *mcpBinary != "" {
		mcpPath, e := filepath.Abs(*mcpBinary)
		if e != nil {
			c.Close()
			check(e)
		}
		cfg := map[string]any{"sdp_pgi2": map[string]any{"command": mcpPath, "args": []string{"--project", cwd, "--state", *state}, "startup_timeout_sec": 20}}
		thread, e := c.StartProbeThread(ctx, cwd, cfg)
		if e != nil {
			c.Close()
			check(e)
		}
		tools, auth, e := c.MCPTools(ctx, "sdp_pgi2", thread.ID)
		if e != nil {
			c.Close()
			check(e)
		}
		if len(tools) != 3 {
			c.Close()
			check(fmt.Errorf("incompatible_profile: expected three read-only inspection tools, got %d", len(tools)))
		}
		out["mcpTools"] = tools
		out["mcpAuthStatus"] = auth
		out["ephemeralThread"] = true
		out["threadObserved"] = thread.ID != ""
		out["mcpProtocol"] = "2025-11-25"
		out["mcpSDK"] = "v1.8.0"
	}
	c.Close()
	check(json.NewEncoder(os.Stdout).Encode(out))
}
