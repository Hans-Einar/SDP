package main

import (
	"context"
	tool "github.com/Hans-Einar/SDP/SDPTool"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(tool.RunWithInput(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
