package fynehost
import (
 "bytes"
 "fmt"
 "testing"
 "github.com/Hans-Einar/SDP/SDUI/go/markdown"
)
func TestReviewerEmptyOriginCompatibility(t *testing.T) {
 for _, box := range []string{"0 0 400 200", "40 20 400 200", "-40 -20 400 200"} {
  for _, attr := range []string{"", ` transform="scale(.5)"`, ` transform="translate(1 2) scale(2e0)"`} {
   for _, ending := range []string{`/>`, `></svg>`, `><!--empty--><title>A &amp; B</title></svg>`} {
    source := []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="%s"%s%s`,box,attr,ending))
    original:=append([]byte(nil),source...)
    r:=markdown.Resource{SVG:source,Width:400,Height:200}
    if _,err:=markdown.ValidateSVGResource(r);err!=nil {t.Fatalf("source invalid %s: %v",source,err)}
    data,err:=nativePreviewSVG(r)
    if err!=nil {t.Errorf("native empty/metadata compatibility %s: %v",source,err);continue}
    if _,err:=markdown.ValidateSVGResource(markdown.Resource{SVG:data,Width:400,Height:200});err!=nil {t.Errorf("derived invalid %s: %v",data,err)}
    if !bytes.Equal(source,original) {t.Fatal("source mutated")}
   }
  }
 }
}
