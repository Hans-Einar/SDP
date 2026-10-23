package sourceinput

import (
 "fmt"
 "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

type sourceRange struct {name string; offset,line,size int}
func locateDiagnostic(d parser.Diagnostic, ranges []sourceRange) error {
 for _,r:=range ranges {
  if d.Span.Start>=r.offset && d.Span.Start<=r.offset+r.size {
   d.Span.Start-=r.offset;d.Span.End-=r.offset;d.Span.Line-=r.line;d.Span.EndLine-=r.line
   return fmt.Errorf("%s: %w",r.name,d)
  }
 }
 return d
}
func shiftModel(m *parser.Model, offset,line int) {
 shift:=func(s *parser.Span){if s.Line!=0{s.Start+=offset;s.End+=offset;s.Line+=line;s.EndLine+=line}}
 shift(&m.Span);shift(&m.Header.Span)
 for i:=range m.Declarations {d:=&m.Declarations[i];shift(&d.Span);shift(&d.Name.Span)}
 for i:=range m.Statements {
  s:=&m.Statements[i];shift(&s.Span)
  for _,id:=range []*parser.Identifier{&s.Subject,&s.Object,&s.Interface,&s.Mode,&s.Container,&s.Dataset,&s.Datagram,&s.Field,&s.Channel,&s.Message,&s.Sender,&s.Receiver,s.Variant} {if id!=nil{shift(&id.Span)}}
  for _,n:=range []*parser.Integer{&s.Offset,&s.Width,&s.Ordinal,s.ReplyTo}{if n!=nil{shift(&n.Span)}}
 }
}
