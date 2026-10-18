package main
import("encoding/json";"os";"fyne.io/fyne/v2";"fyne.io/fyne/v2/app";"fyne.io/fyne/v2/widget")
func log(k string,v any){json.NewEncoder(os.Stdout).Encode(map[string]any{"event":k,"data":v})}
type entry struct{widget.Entry}
func(e *entry)TypedKey(k *fyne.KeyEvent){log("typed-key",k.Name);e.Entry.TypedKey(k)}
func main(){a:=app.NewWithID("sdui.ime.probe");w:=a.NewWindow("SDUI IME Probe"); e:=&entry{};e.ExtendBaseWidget(e);e.OnChanged=func(s string){log("change",s)};e.OnSubmitted=func(s string){log("submit",s)};w.SetContent(e);w.Resize(fyne.NewSize(600,150));w.Show();w.Canvas().Focus(e);log("ready",true);a.Run()}
