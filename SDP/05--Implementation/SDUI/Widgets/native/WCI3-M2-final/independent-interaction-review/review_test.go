package bridge_test
import (
 "context"
 "strings"
 "testing"
 ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
 "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
 fixture "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/commands"
 sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
 sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)
func TestReviewerExtendedInteractionReceiver(t *testing.T) {
 t.Run("command",func(t *testing.T){
 source:=strings.Replace(fixture.Source,`preview=input("Run result",value="No command")`,`preview=input("Run result",value="No command",readOnly=true)`,1)
 d,s,f:=commandSession(t,source);bindCommands(t,d,s,f)
 before,_:=s.Widget(fixture.PreviewPath)
 result,err:=invoke(t,s,"page/view/header/toolbar/runButton","button")
 after,_:=s.Widget(fixture.PreviewPath)
 t.Logf("result=%+v err=%v calls=%v before=%q after=%q draftrev=%d",result,err,f.Calls(),before.Value,after.Value,before.DraftRevision)
 if err!=nil || result.Domain!=ui.DomainSucceeded || after.Value==before.Value || f.Calls()["Run"]!=1 { t.Fatal("admitted command failed to publish extended readonly receiver") }
 })
 t.Run("page",func(t *testing.T){
 d,s:=paneSession(t,strings.Replace(paneSource,`preview=input("Result",value="initial")`,`preview=input("Result",value="initial",readOnly=true)`,1))
 calls:=0
 e:=paneEngine(t,paneActions,sp.TextType,func(_ context.Context,r sdl.Record)(sdl.Record,error){calls++;return sdl.Record{"Preview":sdl.Text("loaded")},nil})
 if _,err:=bridge.Bind(context.Background(),s,d,map[string]*sdl.Engine{"api":e},panePlans(),nil);err!=nil{t.Fatal(err)}
 result,err:=s.DispatchInteraction(pageEvent(t,s,"second"));w,_:=s.Widget("page/preview");h,_:=s.Pane("page/tabs");tabs,_:=s.Tabs(h)
 t.Logf("result=%+v err=%v calls=%d value=%q selected=%s draftrev=%d",result,err,calls,w.Value,tabs.Selected,w.DraftRevision)
 if err!=nil || w.Value!="loaded" || calls!=1 || tabs.Selected!="second" {t.Fatal("admitted page failed to publish extended readonly receiver")}
 })
}

func TestReviewerCommandConflictReplay(t *testing.T) {
 source:=strings.Replace(fixture.Source,`preview=input("Run result",value="No command")`,`preview=input("Run result",value="No command",multiline=false)`,1)
 d,s,f:=commandSession(t,source);bindCommands(t,d,s,f)
 if err:=f.NextAction("Run","draft-conflict");err!=nil{t.Fatal(err)}
 e,err:=s.CaptureCommand(origin(t,s,"page/view/header/toolbar/runButton"),"button",nil);if err!=nil{t.Fatal(err)}
 r,err:=s.DispatchInteraction(e);if err==nil || r.Domain!=ui.DomainSucceeded || f.Calls()["Run"]!=1{t.Fatal(r,err,f.Calls())}
 r,err=s.DispatchInteraction(e);if err==nil || r.Domain!=ui.DomainNotCalled || f.Calls()["Run"]!=1{t.Fatal("replayed domain",r,err,f.Calls())}
}
