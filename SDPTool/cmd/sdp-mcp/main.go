// sdp-mcp serves only the project and optional assignment explicitly selected by
// its local launcher. stdout is reserved for the MCP protocol.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancemcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	root := flag.String("project", "", "explicit project root")
	state := flag.String("state", "", "private operational state directory")
	run := flag.String("run", "", "controller-selected run; omitted means read-only inspection")
	attempt := flag.String("attempt", "", "controller-selected observed attempt")
	epoch := flag.String("epoch", "", "controller connection epoch")
	flag.Parse()
	if *root == "" || *state == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "require --project ROOT --state DIR [--run UUID [--attempt UUID --epoch UUID]]")
		os.Exit(2)
	}
	check := func(e error) {
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
	}
	p, e := sdptool.Discover(*root)
	check(e)
	service, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	check(e)
	server, e := governancemcp.New(service, *run, *attempt, *epoch)
	check(e)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	check(server.Run(ctx, &mcp.StdioTransport{MaxLineLength: governance.MaxRequest + 65536}))
}
