// Package presentation renders command results without owning domain operations.
// Adapters are compiled Go functions, selected by schema and operation.
package presentation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Key struct{ Schema, Operation string }
type Renderer func(io.Writer, json.RawMessage) error
type Registry map[Key]Renderer

func Default() Registry {
	r := Registry{}
	for _, op := range []string{"create", "status", "history", "commit", "restore", "merge", "recover", "snapshot"} {
		r[Key{"sdp-model/0.1", op}] = fields
	}
	for _, op := range []string{"discover", "select", "preview", "sdui-preview", "view"} {
		r[Key{"sdptool/0.2", op}] = fields
	}
	r[Key{"sdptool/0.2", "discover"}] = discovery
	r[Key{"sdptool/0.2", "tree"}] = tree
	r[Key{"sdptool/0.2", "version"}] = version
	r[Key{"sdptool/0.2", ""}] = fields // diagnostic envelope
	r[Key{"sdptool/0.2", "release-log"}] = releaseLog
	for _, op := range []string{"install", "upgrade"} {
		r[Key{"sdp-install-command/1", op}] = installation
	}
	return r
}

func (r Registry) Write(w io.Writer, value any, jsonMode bool) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.WriteJSON(w, b, jsonMode)
}

// WriteJSON preserves unknown results byte-for-byte (plus a final newline).
// Render in memory first so a failed adapter cannot leave partial human output
// followed by a JSON fallback. A known renderer failure is an error, not a retry.
func (r Registry) WriteJSON(w io.Writer, b json.RawMessage, jsonMode bool) error {
	if !json.Valid(b) {
		return fmt.Errorf("invalid presentation JSON")
	}
	var h struct {
		Schema        string
		SchemaVersion string
		Operation     string
	}
	// Non-object JSON and unfamiliar header types have no registered adapter.
	headerErr := json.Unmarshal(b, &h)
	schema := h.Schema
	if schema == "" {
		schema = h.SchemaVersion
	}
	render := r[Key{schema, h.Operation}]
	if jsonMode || headerErr != nil || render == nil {
		if _, err := w.Write(b); err != nil {
			return err
		}
		if !bytes.HasSuffix(b, []byte("\n")) {
			_, err := io.WriteString(w, "\n")
			return err
		}
		return nil
	}
	var output bytes.Buffer
	if err := render(&output, b); err != nil {
		return err
	}
	_, err := output.WriteTo(w)
	return err
}

func safe(s string) string {
	var b strings.Builder
	for _, c := range s {
		if unicode.IsControl(c) {
			q := strconv.QuoteRune(c)
			b.WriteString(q[1 : len(q)-1])
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func fields(w io.Writer, b json.RawMessage) error {
	var value any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if err := d.Decode(&value); err != nil {
		return err
	}
	var print func(string, string, any)
	print = func(indent, label string, v any) {
		switch x := v.(type) {
		case map[string]any:
			if label != "" {
				fmt.Fprintf(w, "%s%s:\n", indent, safe(label))
				indent += "  "
			}
			keys := make([]string, 0, len(x))
			for k := range x {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				print(indent, k, x[k])
			}
		case []any:
			fmt.Fprintf(w, "%s%s: (%d)\n", indent, safe(label), len(x))
			for i, item := range x {
				print(indent+"  ", strconv.Itoa(i+1), item)
			}
		default:
			fmt.Fprintf(w, "%s%s: %s\n", indent, safe(label), safe(fmt.Sprint(x)))
		}
	}
	print("", "", value)
	return nil
}

func version(w io.Writer, b json.RawMessage) error {
	var v struct{ Version, Revision, InstallationProtocol string }
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "SDPTool %s\nRevision: %s\nInstallation protocol: %s\n", safe(v.Version), safe(v.Revision), safe(v.InstallationProtocol))
	return err
}

func releaseLog(w io.Writer, b json.RawMessage) error {
	var v struct {
		Markdown string
		Count    int
		Status   string
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Markdown != "" {
		_, err := io.WriteString(w, v.Markdown)
		return err
	}
	_, err := fmt.Fprintf(w, "%d release log(s) %s\n", v.Count, safe(v.Status))
	return err
}

func installation(w io.Writer, b json.RawMessage) error {
	var v struct {
		Operation string
		Result    *struct {
			Actions                                 []json.RawMessage
			Preserved                               []json.RawMessage
			Conflicts, Warnings                     []string
			PlanDigest, Status, OperationID, Report string
		}
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	p := v.Result
	if p == nil {
		return nil
	}
	if p.PlanDigest != "" {
		fmt.Fprintf(w, "%s preview: %d changes, %d preserved, %d conflicts; plan %s\n", safe(v.Operation), len(p.Actions), len(p.Preserved), len(p.Conflicts), safe(p.PlanDigest))
		for _, c := range p.Conflicts {
			fmt.Fprintln(w, "Conflict:", safe(c))
		}
		for _, warning := range p.Warnings {
			fmt.Fprintln(w, "Warning:", safe(warning))
		}
	} else if p.Status != "" {
		fmt.Fprintf(w, "%s: operation %s; report %s\n", safe(p.Status), safe(p.OperationID), safe(p.Report))
	}
	return nil
}

func discovery(w io.Writer, b json.RawMessage) error {
	var v struct {
		Root, Area, Status string
		Inventory          struct{ ProjectID string }
		Capabilities       map[string]string
		Sources            []struct{ Source, Profile, State, Diagnostic string }
		Programs           []struct{ ID, Label, Source, Entry, State, Diagnostic string }
		ProgramDiscovery   struct{ State, Diagnostic string }
		Navigation         struct{ InventoryRevision string }
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	fmt.Fprintf(w, "%s [%s]\nProject: %s\nSDP area: %s\n", safe(v.Inventory.ProjectID), safe(v.Status), safe(v.Root), safe(v.Area))
	for _, name := range []string{"sdl", "sdui", "kanban", "sessions", "implementation-plan"} {
		fmt.Fprintf(w, "%s: %s\n", name, safe(v.Capabilities[name]))
	}
	fmt.Fprintf(w, "Sources: %d\n", len(v.Sources))
	for _, s := range v.Sources {
		fmt.Fprintf(w, "  %s [%s; %s]\n", safe(s.Source), safe(s.Profile), safe(s.State))
		if s.Diagnostic != "" {
			fmt.Fprintf(w, "    %s\n", safe(s.Diagnostic))
		}
	}
	if v.Navigation.InventoryRevision != "" {
		fmt.Fprintf(w, "Snapshot: %s\n", safe(v.Navigation.InventoryRevision))
	}
	fmt.Fprintf(w, "SDUI programs: %d [%s]\n", len(v.Programs), safe(v.ProgramDiscovery.State))
	if v.ProgramDiscovery.Diagnostic != "" {
		fmt.Fprintf(w, "  %s\n", safe(v.ProgramDiscovery.Diagnostic))
	}
	for _, p := range v.Programs {
		fmt.Fprintf(w, "  %s — %s [%s]\n    %s :: %s\n", safe(p.ID), safe(p.Label), safe(p.State), safe(p.Source), safe(p.Entry))
		if p.Diagnostic != "" {
			fmt.Fprintf(w, "    %s\n", safe(p.Diagnostic))
		} else {
			fmt.Fprintf(w, "    run --program %s\n", safe(p.ID))
		}
	}
	return nil
}
