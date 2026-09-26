package main

import (
	"context"
	tool "github.com/Hans-Einar/SDP/Toolkit/SDPTool"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	os.Exit(tool.Run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
