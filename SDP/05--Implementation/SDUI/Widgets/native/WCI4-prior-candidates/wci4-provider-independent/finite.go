package main
import("fmt";"math";"bytes";"github.com/Hans-Einar/SDP/SDUI/go/markdown";"github.com/fyne-io/oksvg")
func main(){
 cases:=[]string{`<g transform="matrix(1e308 0 -1e308 1 0 0)"><line x1="1" y1="-1" x2="0" y2="0"/></g>`,`<rect x="1e308" y="0" width="1e308" height="1"/>`,`<path d="M0 0 Q-1e308 0 1e308 0 T0 0"/>`}
 for _,s:=range cases{b:=[]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 20">`+s+`</svg>`); _,e:=markdown.ValidateSVGResource(markdown.Resource{SVG:b,Width:40,Height:20});_,backend:=oksvg.ReadIconStream(bytes.NewReader(b),oksvg.StrictErrorMode);fmt.Printf("source=%s\nprovider=%v backend=%v\n",s,e,backend)}
 a:=1e308;fmt.Printf("actual transformed line endpoint x=%v infinite=%v\n",a*1+(-a)*(-1),math.IsInf(a*1+(-a)*(-1),0))
}
