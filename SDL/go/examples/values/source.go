// Package values supplies the bounded WCI3-M1 actual SDL/native fixture.
package values

import "strings"

const WindowTitle = "SDUI WCI3 Values"
const MixedPath = "page/mixed"
const TextPath = "page/textForm"
const TextNamePath = TextPath + "/textName"
const TextEditPath = "page/view/body/textEdit"
const PreviewPath = "page/view/footer/loadPreview"

var Targets = map[string]string{
	"Flag": "page/view/body/flag", "Level": "page/view/body/level", "Mode": "page/view/body/mode", "Count": "page/view/body/count",
	"FormFlag": "page/mixed/form/body/formFlag", "FormLevel": "page/mixed/form/body/formLevel", "FormMode": "page/mixed/form/body/formMode", "FormCount": "page/mixed/form/body/formCount", "ChildCount": "page/mixed/form/body/childCount",
}

const FormNotePath = "page/mixed/form/body/formNote"
const Source = `sdui 0.3;
ref: values "values.sdl";
page=[
 openMixed=command("Mixed form",effect="open",target="mixed");
 openText=command("Text form",effect="open",target="textForm");
 view=[
  header=<loadButton=button("Load",callback=values.Load.@invoke),mixedButton=button(command="openMixed"),textButton=button(command="openText")>;
  body=[
   flag=checkbox("Flag",value=false,callback=values.AcceptFlag.@invoke) {x=fill};
   level=slider("Level",min=0,max=100,step=5,value=25,callback=values.AcceptLevel.@invoke) {x=fill};
   mode=select("Mode",value="alpha",callback=values.AcceptMode.@invoke) {x=fill};
   count=number("Count",min=-10,max=100,step=1,value=1,callback=values.AcceptCount.@invoke) {x=fill};
   textEdit=input("Text Save",value="Initial text",callback=values.SaveText.@invoke) {x=fill}
  ] {x=fill,y=fill};
  footer=<loadPreview=input("Loaded text",value="Not loaded") {x=fill}>
 ] {x=fill,y=fill};
 mixed=dialog("SDUI WCI3 Mixed",modal=true)[
  form=[
   header="Mixed values";
   body=[
    formFlag=checkbox("Form flag",value=false) {x=fill};
    formLevel=slider("Form level",min=0,max=100,step=5,value=25) {x=fill};
    formMode=select("Form mode",value="alpha") {x=fill};
    formCount=number("Form count",min=-10,max=100,step=1,value=1) {x=fill};
    formNote=input("Form note",value="Initial note") {x=fill};
    childCount=number("Child count",min=0,max=100,step=1,value=2,callback=values.ChildCount.@invoke) {x=fill}
   ] {x=fill,y=fill};
   footer=<mixedAccept=button("Accept",effect="accept"),mixedCancel=button("Cancel",effect="cancel"),mixedClose=button("Close",effect="close")>
  ] {x=fill,y=fill}
 ] {scale-x=0.8,scale-y=0.85};
 textForm=dialog("SDUI WCI3 Text",modal=true,callback=values.TextSave.@invoke)[
  textName=input("Text name",value="Initial name") {x=fill};
  textActions=<textAccept=button("Accept",effect="accept"),textCancel=button("Cancel",effect="cancel"),textClose=button("Close",effect="close")>
 ] {scale-x=0.65,scale-y=0.45}
] {x=fill,y=fill};
values.AcceptFlag.setHandle(page.view.body.flag);
values.AcceptLevel.setHandle(page.view.body.level);
values.AcceptMode.setHandle(page.view.body.mode);
values.AcceptCount.setHandle(page.view.body.count);
values.ChildCount.setHandle(page.mixed.form.body.childCount);
values.Load.setHandle(page.view.footer.loadPreview);
values.SaveText.setHandle(page.view.body.textEdit);
`

func FixtureSource(nonmodal bool) string {
	if nonmodal {
		return strings.ReplaceAll(Source, "modal=true", "modal=false")
	}
	return Source
}

// RequiredEmptySource changes only the main Mode startup condition. The mixed
// form and supplied option sets retain the default fixture contract.
func RequiredEmptySource(nonmodal bool) string {
	return strings.Replace(FixtureSource(nonmodal),
		`mode=select("Mode",value="alpha",callback=values.AcceptMode.@invoke)`,
		`mode=select("Mode",value="",required=true,callback=values.AcceptMode.@invoke)`, 1)
}

const Actions = `language action-core version 0.1.
action AcceptCount.
action AcceptFlag.
action AcceptLevel.
action AcceptMode.
record BoolValue.
action ChildCount.
record Decision.
record IntValue.
action Load.
action SaveText.
action TextSave.
record TextValue.
AcceptCount invokes GoAcceptCount.
AcceptCount returns IntValue.
AcceptCount takes IntValue.
AcceptFlag invokes GoAcceptFlag.
AcceptFlag returns BoolValue.
AcceptFlag takes BoolValue.
AcceptLevel invokes GoAcceptLevel.
AcceptLevel returns IntValue.
AcceptLevel takes IntValue.
AcceptMode invokes GoAcceptMode.
AcceptMode returns TextValue.
AcceptMode takes TextValue.
BoolValue field Value as boolean.
ChildCount invokes GoChildCount.
ChildCount returns IntValue.
ChildCount takes IntValue.
Decision field Accepted as boolean.
Decision field Message as text.
IntValue field Value as integer.
Load invokes GoLoad.
Load returns TextValue.
Load takes TextValue.
SaveText invokes GoSaveText.
SaveText returns TextValue.
SaveText takes TextValue.
TextSave invokes GoTextSave.
TextSave returns Decision.
TextSave takes TextValue.
TextValue field Value as text.
`
