// Package text provides the bounded WCI3-M2 actual SDL/native editing fixture.
package text

import (
	"fmt"
	"strconv"
	"strings"
)

const WindowTitle = "SDUI WCI3 Text Editing"
const FormPath = "page/textForm"
const WorkspacePath = "page/view/body/workspace"
const SplitPath = "page/view/body"
const OuterPath = WorkspacePath + "/editor/outer"

var Targets = map[string]string{
	"Single": WorkspacePath + "/editor/single", "Required": WorkspacePath + "/editor/required",
	"Multi": OuterPath + "/multi", "Other": WorkspacePath + "/other/otherText",
	"Preview": "page/view/body/side/preview", "FormName": FormPath + "/formName",
	"FormMulti": FormPath + "/formMulti", "Child": FormPath + "/child",
}
var LongText = longText()

func longText() string {
	var b strings.Builder
	for i := 1; i <= 36; i++ {
		fmt.Fprintf(&b, "%02d Blåbær 日本語 🙂 é — native Unicode editing, selection and word wrap across a deliberately long line.\n", i)
	}
	return b.String()
}

const sourceTemplate = `sdui 0.3;
ref: text "text.sdl";
page=[
 openText=command("Text form",effect="open",target="textForm");
 view=[
  header=<loadButton=button("Load",callback=text.Load.@invoke),dialogButton=button(command="openText")>;
  body=split(axis="horizontal",proportion=0.7,minFirst=0.25,minSecond=0.15)[
   workspace=tabs("Editors",selected="editor")[
    editor=page("Editor")[
     single=input("Single",multiline=false,value="Initial single",callback=text.SaveSingle.@invoke) {x=fill};
     required=input("Required",required=true,value="Required text",callback=text.SaveRequired.@invoke) {x=fill};
     outer=[
      multi=input("Multiline",multiline=true,value=LONG_TEXT,callback=text.SaveMulti.@invoke) {x=fill};
      FILLER_ROWS
     ] {x=fill,y=fill,overflow-y=scroll}
    ];
    other=page("Other")[otherText=input("Other",placeholder="Type here",value="Other text",callback=text.SaveOther.@invoke) {x=fill}]
   ];
   side=[preview=input("",readOnly=true,multiline=true,placeholder="Read-only preview",value="Not loaded") {x=fill}]
  ] {x=fill,y=fill}
 ] {x=fill,y=fill};
 textForm=dialog("SDUI WCI3 Text Form",modal=true,callback=text.TextSave.@invoke)[
  formName=input("Name",required=true,value="Initial name") {x=fill};
  formMulti=input("Body",multiline=true,value="Initial body") {x=fill};
  child=input("Child Save",multiline=false,value="Initial child",callback=text.SaveChild.@invoke) {x=fill};
  formActions=<formAccept=button("Accept",effect="accept"),formCancel=button("Cancel",effect="cancel"),formClose=button("Close",effect="close")>
 ] {scale-x=0.7,scale-y=0.8}
] {x=fill,y=fill};
text.Load.setHandle(page.view.body.side.preview);
text.SaveSingle.setHandle(page.view.body.workspace.editor.single);
text.SaveRequired.setHandle(page.view.body.workspace.editor.required);
text.SaveMulti.setHandle(page.view.body.workspace.editor.outer.multi);
text.SaveOther.setHandle(page.view.body.workspace.other.otherText);
text.SaveChild.setHandle(page.textForm.child);
`

func FixtureSource(nonmodal bool) string {
	s := strings.Replace(sourceTemplate, "LONG_TEXT", strconv.Quote(LongText), 1)
	rows := []string{}
	for i := 1; i <= 28; i++ {
		rows = append(rows, fmt.Sprintf("filler%d=\"Outer viewport row %02d: scroll the background independently\"", i, i))
	}
	s = strings.Replace(s, "FILLER_ROWS", strings.Join(rows, ";\n"), 1)
	if nonmodal {
		s = strings.Replace(s, "modal=true", "modal=false", 1)
	}
	return s
}

const Actions = `language action-core version 0.1.
record Decision.
record FormValue.
action Load.
action SaveChild.
action SaveMulti.
action SaveOther.
action SaveRequired.
action SaveSingle.
action TextSave.
record TextValue.
Decision field Accepted as boolean.
Decision field Message as text.
FormValue field Body as text.
FormValue field Name as text.
Load invokes GoLoad.
Load returns TextValue.
Load takes TextValue.
SaveChild invokes GoSaveChild.
SaveChild returns TextValue.
SaveChild takes TextValue.
SaveMulti invokes GoSaveMulti.
SaveMulti returns TextValue.
SaveMulti takes TextValue.
SaveOther invokes GoSaveOther.
SaveOther returns TextValue.
SaveOther takes TextValue.
SaveRequired invokes GoSaveRequired.
SaveRequired returns TextValue.
SaveRequired takes TextValue.
SaveSingle invokes GoSaveSingle.
SaveSingle returns TextValue.
SaveSingle takes TextValue.
TextSave invokes GoTextSave.
TextSave returns Decision.
TextSave takes FormValue.
TextValue field Value as text.
`

func Slot(name string) (string, error) {
	switch name {
	case "empty":
		return "", nil
	case "short":
		return "Programmatic text", nil
	case "unicode":
		return "Blåbær 日本語 🙂 é", nil
	case "lines":
		return "First line\nSecond line", nil
	case "crlf":
		return "First line\r\nSecond line", nil
	case "spaces":
		return "   ", nil
	case "long":
		return LongText, nil
	}
	return "", fmt.Errorf("unknown text slot %q", name)
}

// RequiredEmptySource changes only the required main field's initial accepted text.
func RequiredEmptySource(nonmodal bool) string {
	return strings.Replace(FixtureSource(nonmodal), `required=true,value="Required text"`, `required=true,value=""`, 1)
}
