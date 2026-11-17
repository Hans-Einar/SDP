//go:build desktop

package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"github.com/Hans-Einar/SDP/SDUI/go/prototype"
)

const documentStatus = "Local DocumentHost prototype · values in memory · actions unbound · no SDL execution"

// All disk parsing/preparation happens before app startup or on the watcher.
// Native request guards read only watcher-observed generation identity.
func runDocument(source, entry, contextPath string, watch bool, initial *prototype.DocumentInput) error {
	a := app.NewWithID("no.sdp.sdui.prototype")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("SDUI · " + entry + " · native document prototype")
	w.SetPadded(false)
	host := fynehost.NewDocumentHost(w.Canvas(), initial.Size)
	status := widget.NewLabel(documentStatus)
	reportError := func(err error) {
		if err != nil {
			status.SetText(err.Error())
			fmt.Fprintln(os.Stderr, err)
		}
	}
	host.OnStatus = reportError
	var observed atomic.Uint64
	observed.Store(1)
	var closed atomic.Bool
	request := func(input *prototype.DocumentInput, generation uint64) fynehost.DocumentRequest {
		return fynehost.DocumentRequest{Document: input.Candidate.Document, Entry: entry, SessionID: "document-prototype", Sequence: generation, SourceRevision: input.Candidate.Hash, Mode: preparation.Prototype, Choices: input.Choices, SVGResources: input.SVGResources, Guard: func() error {
			if closed.Load() || observed.Load() != generation {
				return fmt.Errorf("stale: observed source/context/resource generation changed")
			}
			return nil
		}}
	}
	w.SetContent(container.NewBorder(nil, status, nil, nil, host.Container))
	// The context viewport describes native document content; the status row
	// occupies additional outer window space. BorderLayout adds one gap.
	w.Resize(fyne.NewSize(float32(initial.Size.W), float32(initial.Size.H)+status.MinSize().Height+theme.Padding()))
	// Establish the actual content viewport before preparing a native bundle.
	// SetContent may first lay out a tiny not-yet-sized window; an empty host can
	// accept that transient size without rejecting an already published document.
	if err := host.Adopt(request(initial, 1)); err != nil {
		host.Close()
		return err
	}
	fmt.Fprintf(os.Stderr, "Document prototype ready: source=%s context=%s\n", initial.Candidate.Hash, initial.ContextRevision)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.SetOnClosed(func() { closed.Store(true); cancel(); host.Close() })
	if watch {
		go watchDocuments(ctx, source, entry, contextPath, initial, &observed, func(input *prototype.DocumentInput, generation uint64, err error) {
			fyne.Do(func() {
				if closed.Load() || observed.Load() != generation {
					return
				}
				if err != nil {
					reportError(err)
					return
				}
				if err = host.Adopt(request(input, generation)); err != nil {
					reportError(err)
					return
				}
				status.SetText(documentStatus)
				fmt.Fprintf(os.Stderr, "Document prototype reloaded: generation=%d source=%s context=%s\n", generation, input.Candidate.Hash, input.ContextRevision)
			})
		})
	}
	w.ShowAndRun()
	return nil
}

// watchDocuments treats failures as observed generations too, invalidating queued
// older candidates without touching the last valid published bundle. LoadDocument
// rechecks disk snapshots before this callback; native guards perform no I/O.
func watchDocuments(ctx context.Context, source, entry, contextPath string, initial *prototype.DocumentInput, observed *atomic.Uint64, post func(*prototype.DocumentInput, uint64, error)) {
	last := initial.Candidate.Hash + ":" + initial.ContextRevision
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		input, err := prototype.LoadDocument(source, entry, "", contextPath)
		key := ""
		if err != nil {
			key = "error:" + err.Error()
		} else {
			key = input.Candidate.Hash + ":" + input.ContextRevision
		}
		if key == last {
			continue
		}
		last = key
		generation := observed.Add(1)
		post(input, generation, err)
	}
}
