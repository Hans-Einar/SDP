// Package sourceinput loads explicit design inputs. The parser remains I/O-free.
package sourceinput

import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "fmt"
 "io"
 "os"
 "path/filepath"
 "strings"

 "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

const Profile = "sdl-source-set/0.1"
const MaxFiles = 128
const MaxManifestBytes = 64 << 10

type Manifest struct {
 Profile string `json:"profile"`
 Sources []string `json:"sources"`
}
type File struct {
 Path string `json:"path"`
 SHA256 string `json:"sha256"`
}
type Origin struct {
 Path string `json:"path"`
 Span parser.Span `json:"span"`
}
type Input struct {
 // Source is canonical, combined design-core for a set; unchanged for a single file.
 Source string
 Revision string
 Files []File
 Origins map[string][]Origin
 paths []string
}
func IsSet(path string) bool { return strings.HasSuffix(path, ".design-set.json") }
func hash(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func read(path string, limit int) ([]byte,error) {
 f,e:=os.Open(path); if e!=nil {return nil,e}; defer f.Close()
 info,e:=f.Stat(); if e!=nil {return nil,e}; if !info.Mode().IsRegular() {return nil,fmt.Errorf("%s: require regular file",path)}
 b,e:=io.ReadAll(io.LimitReader(f,int64(limit)+1)); if e!=nil{return nil,e}
 if len(b)>limit{return nil,fmt.Errorf("%s: source input limit %d",path,limit)};return b,nil
}
func Load(path string) (*Input,error) {
 abs,e:=filepath.Abs(path); if e!=nil{return nil,e}
 if !IsSet(path) {
  b,e:=read(abs,parser.MaxBytes);if e!=nil{return nil,e}
  return &Input{Source:string(b),Revision:hash(b),Files:[]File{{filepath.Base(abs),hash(b)}},paths:[]string{abs}},nil
 }
 b,e:=read(abs,MaxManifestBytes);if e!=nil{return nil,e}
 var spec Manifest
 d:=json.NewDecoder(bytes.NewReader(b));d.DisallowUnknownFields()
 if e=d.Decode(&spec);e!=nil{return nil,fmt.Errorf("%s: %w",path,e)}
 var extra any
 if d.Decode(&extra)!=io.EOF{return nil,fmt.Errorf("%s: trailing manifest JSON",path)}
 if spec.Profile!=Profile || len(spec.Sources)==0 || len(spec.Sources)>MaxFiles{return nil,fmt.Errorf("%s: require %s and 1..%d sources",path,Profile,MaxFiles)}
 in:=&Input{Origins:map[string][]Origin{},paths:[]string{abs}}
 merged:=&parser.Model{}
 seen:=map[string]bool{}
 var ranges []sourceRange
 revision:=[]File{{"@manifest",hash(b)}}
 total,offset,line:=0,0,0
 for _,name:=range spec.Sources {
  if name=="" || filepath.IsAbs(name) || filepath.Ext(name)!=".design" || strings.ContainsAny(name,"\x00\\") {return nil,fmt.Errorf("%s: invalid relative .design member %q",path,name)}
  p:=filepath.Join(filepath.Dir(abs),filepath.FromSlash(name))
  resolved,e:=filepath.EvalSymlinks(p);if e!=nil{return nil,e}
  if seen[resolved]{return nil,fmt.Errorf("duplicate source member: %s",name)};seen[resolved]=true
  data,e:=read(resolved,parser.MaxBytes-total);if e!=nil{return nil,e};total+=len(data)
  m,e:=parser.Parse(string(data));if e!=nil{return nil,fmt.Errorf("%s: %w",name,e)}
  for _,decl:=range m.Declarations {key:=decl.Kind+" "+decl.Name.Name+".";in.Origins[key]=append(in.Origins[key],Origin{name,decl.Span})}
  for _,s:=range m.Statements {in.Origins[s.Sentence()]=append(in.Origins[s.Sentence()],Origin{name,s.Span})}
  ranges=append(ranges,sourceRange{name,offset,line,len(data)})
  shiftModel(m,offset,line)
  if len(merged.Declarations)==0 && len(merged.Statements)==0 {merged.Header=m.Header}
  merged.Declarations=append(merged.Declarations,m.Declarations...)
  merged.Statements=append(merged.Statements,m.Statements...)
  merged.Span.End=m.Span.End;merged.Span.EndLine=m.Span.EndLine;merged.Span.EndColumn=m.Span.EndColumn
  file:=File{name,hash(data)};in.Files=append(in.Files,file);revision=append(revision,file)
  in.paths=append(in.paths,p,resolved)
  offset+=len(data)+1;line+=strings.Count(string(data),"\n")+1
 }
 canonical,ds:=parser.Canonical(merged)
 if len(ds)>0 {return nil,locateDiagnostic(ds[0],ranges)}
 if len(canonical)>parser.MaxBytes{return nil,fmt.Errorf("combined model limit")}
 in.Source=canonical
 record,_:=json.Marshal(revision);in.Revision=hash(record)
 return in,nil
}

// OutsideOutput prevents replacement of the descriptor or any listed source.
func (in *Input) OutsideOutput(output string) error {
 target,e:=filepath.Abs(output);if e!=nil{return e}
 // Resolve existing ancestors too, so a symlinked output parent cannot alias input.
 target,e=resolveDestination(target);if e!=nil{return e}
 for _,p:=range in.paths {
  p,e=resolveDestination(p);if e!=nil{return e}
  rel,e:=filepath.Rel(target,p);if e!=nil{return e}
  if rel!=".." && !strings.HasPrefix(rel,".."+string(filepath.Separator)){return fmt.Errorf("keep source outside output: %s",p)}
 }
 return nil
}
func resolveDestination(p string) (string,error) {
 resolved,e:=filepath.EvalSymlinks(p);if e==nil{return resolved,nil}
 if !os.IsNotExist(e){return "",e}
 parent:=filepath.Dir(p);if parent==p{return "",e}
 base,e:=resolveDestination(parent);if e!=nil{return "",e};return filepath.Join(base,filepath.Base(p)),nil
}
