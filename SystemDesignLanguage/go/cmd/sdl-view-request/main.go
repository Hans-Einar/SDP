//go:build linux

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/broker"
	"os"
)

func main() {
	socket := flag.String("socket", "", "Broker socket")
	uri := flag.String("uri", "", "sdl-view URI")
	window := flag.String("window", "", "Captured window ID")
	client := flag.String("client", "cli", "Client/session ID")
	seq := flag.Uint64("sequence", 1, "Monotonic sequence per client/window/pane")
	open := flag.Bool("open", false, "Open through registered reader")
	release := flag.String("release", "", "Release reader lease")
	sweep := flag.Bool("sweep", false, "Remove only released bundles")
	flag.Parse()
	call := broker.Call{Operation: "select", Request: broker.Request{URI: *uri, Window: *window, Client: *client, Sequence: *seq, Open: *open}}
	if *release != "" {
		call.Operation = "release"
		call.Lease = *release
	}
	if *sweep {
		call.Operation = "sweep"
	}
	r, e := broker.Send(context.Background(), *socket, call)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	json.NewEncoder(os.Stdout).Encode(r)
}
