package runtime

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

func program(t *testing.T) *parser.Program {
	t.Helper()
	source, e := os.ReadFile("../examples/echo.sdl")
	if e != nil {
		t.Fatal(e)
	}
	p, e := parser.CompileActions(string(source))
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func registry(fn Handler) Registry {
	return Registry{"GoEcho": {Input: parser.RecordType{"Value": parser.TextType}, Output: parser.RecordType{"Value": parser.TextType}, Call: fn}}
}
func TestTypedCallsAndRegistration(t *testing.T) {
	p := program(t)
	if _, e := New(p, Registry{}); e == nil {
		t.Fatal("missing binding")
	}
	bad := registry(func(context.Context, Record) (Record, error) { return nil, nil })
	bad["GoEcho"].Input["Value"] = parser.IntegerType
	if _, e := New(p, bad); e == nil {
		t.Fatal("mismatched binding")
	}
	calls := 0
	e, err := New(p, registry(func(_ context.Context, r Record) (Record, error) {
		calls++
		r["Value"] = Text(strings.ToUpper(r["Value"].Text))
		return r, nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	input := Record{"Value": Text("hello")}
	out, err := e.Execute(context.Background(), Request{"Echo", 1, 1, input})
	if err != nil || out.Output["Value"].Text != "HELLO" || input["Value"].Text != "hello" || calls != 1 {
		t.Fatal(out, err, calls, input)
	}
	for _, request := range []Request{{"Echo", 1, 1, input}, {"Echo", 0, 2, input}, {"Echo", 1, 2, Record{"Value": Integer(2)}}, {"Echo", 1, 2, Record{"Extra": Text("x")}}} {
		if _, err = e.Execute(context.Background(), request); err == nil {
			t.Fatal("invalid call")
		}
	}
	if calls != 1 {
		t.Fatal("invalid call reached handler")
	}
	p.Records["EchoInput"]["Value"] = parser.IntegerType
	if _, err = e.Execute(context.Background(), Request{"Echo", 1, 2, input}); err != nil {
		t.Fatal("caller mutated compiled runtime", err)
	}
}
func TestInvalidResultAndPanicConsumed(t *testing.T) {
	for _, fn := range []Handler{func(context.Context, Record) (Record, error) { return Record{"Value": Integer(1)}, nil }, func(context.Context, Record) (Record, error) { panic("domain error") }} {
		e, err := New(program(t), registry(fn))
		if err != nil {
			t.Fatal(err)
		}
		r := Request{"Echo", 1, 1, Record{"Value": Text("x")}}
		if _, err = e.Execute(context.Background(), r); err == nil {
			t.Fatal("bad result accepted")
		}
		if _, err = e.Execute(context.Background(), r); err == nil || !strings.Contains(err.Error(), "duplicate-command") {
			t.Fatal("bad result retried")
		}
	}
}
