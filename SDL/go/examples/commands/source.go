// Package commands provides bounded, real SDL command and dialog acceptance
// fixtures. Its control channel never invokes or simulates a user's command.
package commands

import "strings"

const (
	WindowTitle     = "SDUI WCI2 Commands"
	DialogTitle     = "SDUI WCI2 Settings"
	ItemsPath       = "page/view/body/items"
	PreviewPath     = "page/view/footer/preview"
	FlagPreviewPath = "page/view/footer/flagPreview"
	ItemPreviewPath = "page/view/footer/itemPreview"
	SettingsPath    = "page/settings"
	NamePath        = "page/settings/form/body/name"
	NotePath        = "page/settings/form/body/note"
)

const Source = `sdui 0.3;
ref: actions "commands.sdl";
page=[
 run=command("Run",icon="fixture-run",callback=actions.Run.@invoke,key="Primary+R",tooltip="Run the shared action");
 flag=command("Flag",toggle=true,callback=actions.Toggle.@invoke,key="Primary+K",tooltip="Toggle the shared flag");
 alpha=command("Alpha",toggle=true,checked=true,exclusive="choice");
 beta=command("Beta",toggle=true,exclusive="choice");
 inspect=command("Inspect",context="item",target="view/body/items",callback=actions.Inspect.@invoke);
 openSettings=command("Settings",effect="open",target="settings",key="Primary+D");
 itemMenu=menu("Item actions",mode="context",target="view/body/items")[inspectItem=item(command="inspect")];
 view=[
  header=[
   actions=menu("Actions")[
    runItem=item(command="run");
    separator();
    options=menuGroup("Options")[flagItem=item(command="flag")];
    more=menu("More",mode="submenu")[alphaItem=item(command="alpha"); betaItem=item(command="beta")]
   ];
   toolbar=<runButton=button(command="run"),flagButton=button(command="flag"),openButton=button(command="openSettings"),inspectButton=button(command="inspect"),alphaButton=button(command="alpha"),betaButton=button(command="beta")>
  ];
  body=<items=tree("Items") {x=1fr,y=fill,overflow-x=scroll,overflow-y=scroll},mainDraft=input("Main draft",value="Main accepted") {x=1fr}> {x=fill,y=fill};
  footer=<preview=input("Run result",value="No command") {x=1fr},flagPreview=input("Flag result",value="No toggle") {x=1fr},itemPreview=input("Item result",value="No item") {x=1fr}> {x=fill}
 ] {x=fill,y=fill};
 settings=dialog("SDUI WCI2 Settings",modal=true,callback=actions.Save.@invoke)[
  mark=command("Mark",toggle=true);
  openChild=command("Open child",effect="open",target="x");
  editMenu=menu("Text actions",mode="context",target="settings/form/body/name")[markItem=item(command="settings/mark");childItem=item(command="settings/openChild")];
  form=[
   header="Edit settings";
   body=<name=input("Name",value="Initial name") {x=1fr},note=input("Note",value="Initial note") {x=1fr}> {x=fill,y=fill};
   footer=<acceptButton=button("Accept",effect="accept"),cancelButton=button("Cancel",effect="cancel"),closeButton=button("Close",effect="close")>
  ] {x=fill,y=fill}
 ] {scale-x=0.65,scale-y=0.65};
 x=dialog("SDUI WCI2 Child",modal=true)[
  childInput=input("Child input",value="Child accepted");
  childClose=button("Close",effect="close")
 ] {scale-x=0.4,scale-y=0.4}
] {x=fill,y=fill};
actions.Run.setHandle(page.view.footer.preview);
actions.Toggle.setHandle(page.view.footer.flagPreview);
actions.Inspect.setHandle(page.view.footer.itemPreview);
`

func FixtureSource(nonmodal bool) string {
	if nonmodal {
		return strings.Replace(Source, "modal=true", "modal=false", 1)
	}
	return Source
}

const Actions = `language action-core version 0.1.
action Inspect.
record InspectInput.
action Run.
record RunInput.
action Save.
record SaveInput.
record SaveOutput.
record TextOutput.
action Toggle.
record ToggleInput.
Inspect invokes GoInspect.
Inspect returns TextOutput.
Inspect takes InspectInput.
InspectInput field ItemId as text.
Run invokes GoRun.
Run returns TextOutput.
Run takes RunInput.
RunInput field Token as text.
Save invokes GoSave.
Save returns SaveOutput.
Save takes SaveInput.
SaveInput field Name as text.
SaveInput field Note as text.
SaveOutput field Accepted as boolean.
SaveOutput field Message as text.
TextOutput field Preview as text.
Toggle invokes GoToggle.
Toggle returns TextOutput.
Toggle takes ToggleInput.
ToggleInput field Checked as boolean.
`
