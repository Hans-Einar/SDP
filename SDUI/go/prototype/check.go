// Package prototype defines the standalone local-prototype readiness contract.
package prototype

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
)

type Report struct {
	Schema     string `json:"schema"`
	Operation  string `json:"operation"`
	Status     string `json:"status"`
	Revision   string `json:"revision"`
	Entry      string `json:"entry"`
	Source     string `json:"source"`
	Diagnostic string `json:"diagnostic"`
}

func Check(path, entry, revision string) (reload.Candidate, Report, error) {
	c := reload.Read(path, entry, 1, nil)
	r := Report{Schema: "sdptool/0.2", Operation: "sdui-check", Entry: entry, Source: path, Revision: c.Hash}
	if c.Err != nil {
		return c, r, c.Err
	}
	if revision != "" && c.Hash != revision {
		return c, r, fmt.Errorf("stale: SDUI source changed; regenerate the preview")
	}
	prepared, err := preparation.Prepare(preparation.Request{
		Document: c.Document, Entry: entry, SessionID: "preflight", SourceRevision: c.Hash,
		Mode: preparation.Prototype, Capabilities: admission.Capabilities(),
		ValidateLayout: func(root *parser.Instance) error { return admission.Layout(root, layout.Size{W: 1280, H: 800}) },
	})
	if err != nil {
		return c, r, err
	}
	defer prepared.Close()
	unbound := prepared.Unbound
	r.Status = "prototype"
	r.Diagnostic = "Local Fyne prototype; no SDL runtime"
	if unbound > 0 {
		r.Status = "prototype-unbound"
		r.Diagnostic = fmt.Sprintf("Local Fyne prototype; %d unbound SDL callbacks (no SDL execution)", unbound)
	}
	return c, r, nil
}
