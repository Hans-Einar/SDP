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
	run := flag.String("run", "", "controller-selected run; omitted means read-only inspection unless --coordinator is selected")
	attempt := flag.String("attempt", "", "controller-selected observed attempt")
	coordinator := flag.String("coordinator", "", "explicit coordinator identity; read all registered work and acknowledge reports (excludes run/attempt/epoch)")
	session := flag.String("session-controller", "", "explicit registered Session instance; privileged Session operations without owner decisions or host dispatch")
	leaderSkill := flag.String("leader-skill", "", "project-local Project Leader SKILL.md; mandatory for session-controller")
	epoch := flag.String("epoch", "", "controller connection epoch")
	flag.Parse()
	if *root == "" || *state == "" || flag.NArg() != 0 || (*coordinator != "" && (*run != "" || *attempt != "" || *epoch != "")) || (*session != "" && (*coordinator != "" || *run != "" || *attempt != "" || *epoch != "")) || ((*session == "") != (*leaderSkill == "")) {
		fmt.Fprintln(os.Stderr, "require --project ROOT --state DIR [--run UUID [--attempt UUID --epoch UUID] | --coordinator ID | --session-controller UUID --leader-skill PATH]")
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
	var server *mcp.Server
	if *session != "" {
		skill, err := service.SessionLeaderSkill(*session, *leaderSkill)
		check(err)
		server, e = governancemcp.NewSessionController(service, *session, skill)
	} else if *coordinator != "" {
		server, e = governancemcp.NewCoordinator(service, *coordinator)
	} else {
		server, e = governancemcp.New(service, *run, *attempt, *epoch)
	}
	check(e)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	check(server.Run(ctx, &mcp.StdioTransport{MaxLineLength: governance.MaxRequest + 65536}))
}
