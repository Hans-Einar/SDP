package documents

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"strings"
	"testing"
)

func TestStructureAttachmentsPreserveModeAndActualParticipant(t *testing.T) {
	v := &viewpoint.Views{ReadingProject: "trial", Kinds: map[string]string{"Runtime": "container", "Worker": "unit", "External": "unit", "Bus": "channel", "Request": "message", "Stored": "datagram"}, Relations: map[string][]viewpoint.Fact{"contains": {{"subject": "Runtime", "object": "Worker"}}}, Facts: []viewpoint.Fact{
		{"verb": "uses", "subject": "Worker", "channel": "Bus", "role": "sender", "message": "Request", "mode": "Online"},
		{"verb": "uses", "subject": "External", "channel": "Bus", "role": "receiver", "message": "Request", "mode": "Online"},
		{"verb": "uses", "subject": "Worker", "channel": "Bus", "role": "receiver", "message": "Stored", "mode": "Retained"},
	}}
	text, assets, e := ReadingStructure(v, viewpoint.Query{Viewpoint: "VP02", Projection: "reading", Focus: "Runtime", Mode: "Online"})
	if e != nil {
		t.Fatal(e)
	}
	svg := string(assets["structure-0-0.svg"])
	for _, want := range []string{"Worker", "External", "Bus", "marker-end", "marker-start", "Request", "Online"} {
		if !strings.Contains(svg, want) {
			t.Fatal("missing topology", want, svg)
		}
	}
	if strings.Contains(svg, "Retained") || strings.Contains(svg, "Stored") {
		t.Fatal("mode leaked")
	}
	if !strings.Contains(text, "Not applicable") || strings.Contains(text, "inferred sender") {
		t.Fatal(text)
	}
	if strings.Count(svg, ">Runtime</text>") != 1 {
		t.Fatal("container copied as child")
	}
	if e = ValidateSVG(assets["structure-0-0.svg"]); e != nil {
		t.Fatal(e)
	}
}
func TestStructureCollapsedContainerDoesNotInventParticipants(t *testing.T) {
	v := &viewpoint.Views{Kinds: map[string]string{"Runtime": "container", "Worker": "unit", "Bus": "channel"}, Relations: map[string][]viewpoint.Fact{"contains": {{"subject": "Runtime", "object": "Worker"}}}, Facts: []viewpoint.Fact{{"verb": "uses", "subject": "Runtime", "channel": "Bus", "role": "sender", "message": "Read", "mode": "Online"}}}
	_, assets, e := ReadingStructure(v, viewpoint.Query{Viewpoint: "VP08", Projection: "reading", Focus: "Bus"})
	if e != nil {
		t.Fatal(e)
	}
	svg := string(assets["structure-0-0.svg"])
	if strings.Contains(svg, ">Worker</text>") {
		t.Fatal("invented internal unit participant")
	}
	if !strings.Contains(svg, ">Runtime</text>") {
		t.Fatal("missing declared participant")
	}
}

func TestStructureNestedUnitsAreScopedAndCollapsed(t *testing.T) {
	v := &viewpoint.Views{ReadingProject: "trial", Kinds: map[string]string{"Runtime": "container", "Kernel": "unit", "Registry": "unit", "Bus": "channel"}, Relations: map[string][]viewpoint.Fact{"contains": {{"subject": "Runtime", "object": "Kernel"}, {"subject": "Kernel", "object": "Registry"}}}, Facts: []viewpoint.Fact{{"verb": "uses", "subject": "Registry", "channel": "Bus", "role": "sender", "message": "Read", "mode": "Online"}}}
	text, assets, e := ReadingStructure(v, viewpoint.Query{Viewpoint: "VP02", Projection: "reading"})
	if e != nil {
		t.Fatal(e)
	}
	if len(assets) != 3 || !strings.Contains(text, "Kernel — participants") {
		t.Fatal("nested unit scope missing", text)
	}
	svg := string(assets["structure-0-0.svg"])
	if strings.Contains(svg, ">Registry</text>") || strings.Contains(svg, ">Kernel</text>") {
		t.Fatal("nested units falsely shown as roots")
	}
	for _, b := range assets {
		if e = ValidateSVG(b); e != nil {
			t.Fatal(e)
		}
	}
}

func TestStructureScopeParticipantUsesOwnBoundary(t *testing.T) {
	v := &viewpoint.Views{Kinds: map[string]string{"Runtime": "container", "Kernel": "unit", "Worker": "unit", "Bus": "channel"}, Relations: map[string][]viewpoint.Fact{"contains": {{"subject": "Runtime", "object": "Kernel"}, {"subject": "Kernel", "object": "Worker"}}}, Facts: []viewpoint.Fact{{"verb": "uses", "subject": "Kernel", "channel": "Bus", "role": "sender", "message": "Read", "mode": "Online"}}}
	_, assets, e := ReadingStructure(v, viewpoint.Query{Viewpoint: "VP02", Projection: "reading", Focus: "Kernel"})
	if e != nil {
		t.Fatal(e)
	}
	svg := string(assets["structure-0-0.svg"])
	if !strings.Contains(svg, "M262,45") || strings.Contains(svg, ">Runtime</text>") {
		t.Fatal("scope participant escaped its own boundary", svg)
	}
}

func TestStructurePartitionRetainsEveryBinding(t *testing.T) {
	v := &viewpoint.Views{Kinds: map[string]string{"Runtime": "container", "Bus": "channel"}, Relations: map[string][]viewpoint.Fact{}}
	parents := map[string]string{}
	for i := 0; i < 11; i++ {
		n := fmt.Sprintf("Worker%d", i)
		v.Kinds[n] = "unit"
		parents[n] = "Runtime"
		for _, role := range []string{"sender", "receiver"} {
			for _, mode := range []string{"Online", "Offline"} {
				id := n + role + mode
				v.Facts = append(v.Facts, viewpoint.Fact{"id": id, "verb": "uses", "subject": n, "channel": "Bus", "role": role, "mode": mode, "message": "Request"})
			}
		}
	}
	panels := structurePanels(v, "Runtime", "", parents)
	seen := map[string]int{}
	for _, p := range panels {
		for _, f := range p.Facts {
			seen[f.S("id")]++
		}
		if e := ValidateSVG(p.SVG); e != nil {
			t.Fatal(e)
		}
	}
	if len(panels) < 8 {
		t.Fatal("dense figures not partitioned", len(panels))
	}
	for _, f := range v.Facts {
		if seen[f.S("id")] != 1 {
			t.Fatal("lost or duplicated binding", f, seen)
		}
	}
}
