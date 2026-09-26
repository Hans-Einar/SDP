package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/devhost"
	"os"
	"os/signal"
	"time"
)

func main() {
	root := flag.String("root", ".", "Go module directory")
	pkg := flag.String("package", "./cmd/sdl-demo", "Explicit Go package")
	compiler := flag.String("go", "", "Go executable (defaults to this toolchain)")
	tags := flag.String("tags", "desktop", "Go build tags")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	manager, err := devhost.New(devhost.Config{Root: *root, Package: *pkg, Go: *compiler, Tags: *tags, Args: flag.Args(), Log: os.Stdout})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer manager.Close()
	for change := range devhost.Watch(ctx, *root, 250*time.Millisecond) {
		if change.Err != nil {
			fmt.Fprintln(os.Stderr, change.Err)
			continue
		}
		if err = manager.Rebuild(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Fprintln(os.Stderr, "Go app rebuilt and restarted; pid", manager.PID())
		}
	}
}
