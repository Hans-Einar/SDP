//go:build linux

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/broker"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/reader"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	source := flag.String("source", "", "Registered SDL source")
	project := flag.String("project", "local", "Project ID")
	root := flag.String("runtime", "", "Private runtime directory")
	renderer := flag.String("renderer", "", "Registered mmdr executable")
	xfmd := flag.String("xfmd", "", "Registered XFMD executable")
	flag.Parse()
	must := func(e error) {
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
	}
	if *source == "" {
		must(fmt.Errorf("--source required"))
	}
	var e error
	if *root == "" {
		*root, e = broker.RuntimeDirectory()
		must(e)
	}
	p := broker.Project{Source: *source}
	if *renderer != "" {
		p.Renderer, e = documents.NewMmdr(*renderer)
		must(e)
	}
	store, e := broker.NewStore(filepath.Join(*root, "views"), 128<<20, 64)
	must(e)
	defer store.Close()
	registry := reader.Registry{}
	if *xfmd != "" {
		registry["xfmd"] = reader.XFMD{Program: *xfmd}
	}
	b := broker.New(map[string]broker.Project{*project: p}, registry, store)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	fmt.Println("socket=" + filepath.Join(*root, "views.sock"))
	must(b.Serve(ctx, filepath.Join(*root, "views.sock")))
}
