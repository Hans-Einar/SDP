package main
import("fmt";"os";"path/filepath";"strings";"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser")
func main(){root:=os.Args[1];if e:=os.MkdirAll(root,0700);e!=nil{panic(e)};const h="language design-core version 0.6.\n";var includes,facts strings.Builder
for i:=0;i<127;i++{name:=fmt.Sprintf("Unit%03d",i);fmt.Fprintf(&includes,"includes %q.\n",name+".design");fmt.Fprintf(&facts,"Demo contains %s.\n",name);if e:=os.WriteFile(filepath.Join(root,name+".design"),[]byte(h+"unit "+name+".\n"),0600);e!=nil{panic(e)}}
m,e:=parser.Parse(h+"system Demo.\n"+includes.String()+facts.String());if e!=nil{panic(e)};if e=os.WriteFile(filepath.Join(root,"System.design"),[]byte(parser.FormatFile(m)),0600);e!=nil{panic(e)}}
