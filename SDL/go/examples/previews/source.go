// Package previews is the connected WCI4 provider-preview acceptance application.
package previews

import (
	"fmt"
	"strconv"
	"strings"
)

const WindowTitle = "SDUI WCI4 Previews"
const DialogTitle = "SDUI WCI4 Preview Detail"
const WorkspacePath = "page/view/body/workspace"
const SplitPath = "page/view/body"
const DetailPath = "page/detail"

var Targets = map[string]string{
	"Hero":       WorkspacePath + "/art/viewport/hero",
	"Prose":      WorkspacePath + "/art/viewport/prose",
	"Tail":       WorkspacePath + "/art/viewport/tail",
	"ReportA":    WorkspacePath + "/reports/reportViewport/reportA",
	"ReportB":    WorkspacePath + "/reports/reportViewport/reportB",
	"ReportSVG":  WorkspacePath + "/reports/reportViewport/reportSVG",
	"Note":       "page/view/body/aside/note",
	"LoadStatus": "page/view/body/aside/loadStatus",
	"MarkStatus": "page/view/body/aside/markStatus",
	"DetailSVG":  DetailPath + "/detailViewport/detailSVG",
	"DetailDoc":  DetailPath + "/detailViewport/detailDoc",
	"DetailNote": DetailPath + "/detailNote",
}

const DiagramText = "# Diagram report\n\nProse before the diagram.\n\n```mermaid\nflowchart LR\n A-->B\n```\n\nProse after the diagram."
const ProseText = "# Preview report\n\nPrepared application content.\n\n- Red region: first site\n- Blue region: second site"

const sourceTemplate = `sdui 0.3;
ref: previews "previews.sdl";
ref: art "application-supplied SVG";
page=[
 mark=command("Mark",toggle=true,callback=previews.Mark.@invoke,key="Primary+M");
 openDetail=command("Open Detail",effect="open",target="detail");
 view=[
  header=[
   actions=menu("Actions")[markItem=item(command="mark")];
   toolbar=<loadButton=button("Load",callback=previews.Load.@invoke),markButton=button(command="mark"),openButton=button(command="openDetail")>
  ];
  body=split(axis="horizontal",proportion=0.7,minFirst=0.25,minSecond=0.15)[
   workspace=tabs("Previews",selected="art")[
    art=page("Art")[viewport=[
     hero=svg(art.Chart.@resource,label="Weekly chart",description="Output by site",fallback=POLICY_Hero) {x=fill,scale-y=0.55};
     prose=markdown(DOC_Prose,description="Preview explanation",fallback=POLICY_Prose) {x=fill};
     tail=markdown(DOC_Tail,description="Long preview report",fallback=POLICY_Tail) {x=fill}
    ] {x=fill,y=fill,overflow-y=scroll}];
    reports=page("Reports")[reportViewport=[
     reportA=markdown(DOC_ReportA,description="First diagram report",fallback=POLICY_ReportA) {x=fill};
     reportB=markdown(DOC_ReportB,description="Second diagram report",fallback=POLICY_ReportB) {x=fill};
     reportSVG=svg(art.Chart.@resource,label="Inactive chart",description="Report artwork",fallback=POLICY_ReportSVG) {x=fill,scale-y=0.4}
    ] {x=fill,y=fill,overflow-y=scroll}]
   ];
   aside=[
    note=input("Saved note",multiline=false,value="Initial note",callback=previews.SaveNote.@invoke) {x=fill};
    loadStatus=input("Loaded metadata",readOnly=true,value="Not loaded") {x=fill};
    markStatus=input("Command result",readOnly=true,value="No mark") {x=fill}
   ]
  ] {x=fill,y=fill}
 ] {x=fill,y=fill};
 detail=dialog("SDUI WCI4 Preview Detail",modal=MODAL,callback=previews.AcceptDetail.@invoke)[
  detailViewport=[
   detailSVG=svg(art.Chart.@resource,label="Detail chart",description="Detail artwork",fallback=POLICY_DetailSVG) {x=fill,scale-y=0.5};
   detailDoc=markdown(DOC_DetailDoc,description="Detail explanation",fallback=POLICY_DetailDoc) {x=fill}
  ] {x=fill,y=fill,overflow-y=scroll};
  detailNote=input("Detail note",multiline=false,value="Initial detail",callback=previews.SaveDetail.@invoke) {x=fill};
  detailActions=<acceptButton=button("Accept",effect="accept"),cancelButton=button("Cancel",effect="cancel"),closeButton=button("Close",effect="close")>
 ] {scale-x=0.7,scale-y=0.8}
] {x=fill,y=fill};
previews.Load.setHandle(page.view.body.aside.loadStatus);
previews.Mark.setHandle(page.view.body.aside.markStatus);
previews.SaveNote.setHandle(page.view.body.aside.note);
previews.SaveDetail.setHandle(page.detail.detailNote);
`

// FixtureSource uses only selected .3 source forms. Every declaration is bound,
// including the inactive page and closed dialog.
func FixtureSource(nonmodal bool) string { return sourceFor(nonmodal, nil, nil) }
func sourceFor(nonmodal bool, policies, documents map[string]string) string {
	s := strings.ReplaceAll(sourceTemplate, "MODAL", strconv.FormatBool(!nonmodal))
	for _, alias := range []string{"Hero", "Prose", "Tail", "ReportA", "ReportB", "ReportSVG", "DetailSVG", "DetailDoc"} {
		policy := "reject"
		if alias == "ReportA" || alias == "ReportB" || alias == "ReportSVG" {
			policy = "label"
		}
		if v, ok := policies[alias]; ok {
			policy = v
		}
		s = strings.ReplaceAll(s, "POLICY_"+alias, strconv.Quote(policy))
	}
	for _, alias := range []string{"Prose", "Tail", "ReportA", "ReportB", "DetailDoc"} {
		doc := ProseText
		switch alias {
		case "ReportA", "ReportB":
			doc = DiagramText
		case "Tail":
			var b strings.Builder
			for i := 1; i <= 18; i++ {
				fmt.Fprintf(&b, "## Report row %02d\n\nScroll this prepared report independently of the artwork.\n\n", i)
			}
			doc = b.String()
		}
		if v, ok := documents[alias]; ok {
			doc = v
		}
		s = strings.ReplaceAll(s, "DOC_"+alias, strconv.Quote(doc))
	}
	return s
}

// Canonical action-core ordering, compiled by the real SDL parser/runtime.
const Actions = `language action-core version 0.1.
action AcceptDetail.
record Decision.
action Load.
action Mark.
record MarkInput.
record NoteInput.
action SaveDetail.
action SaveNote.
record TextValue.
AcceptDetail invokes GoAcceptDetail.
AcceptDetail returns Decision.
AcceptDetail takes NoteInput.
Decision field Accepted as boolean.
Decision field Message as text.
Load invokes GoLoad.
Load returns TextValue.
Load takes TextValue.
Mark invokes GoMark.
Mark returns TextValue.
Mark takes MarkInput.
MarkInput field Checked as boolean.
NoteInput field Note as text.
SaveDetail invokes GoSaveDetail.
SaveDetail returns TextValue.
SaveDetail takes TextValue.
SaveNote invokes GoSaveNote.
SaveNote returns TextValue.
SaveNote takes TextValue.
TextValue field Value as text.
`
