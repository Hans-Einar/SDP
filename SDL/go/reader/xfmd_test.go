package reader

import (
	"context"
	"testing"
)

type capture struct{ received []Delivery }

func (c *capture) Open(_ context.Context, d Delivery) error {
	c.received = append(c.received, d)
	return nil
}
func TestCapturedAddressAndRegistration(t *testing.T) {
	c := &capture{}
	r := Registry{"xfmd": c}
	d := Delivery{Window: "window-one", Pane: "main", Client: "session", Sequence: 4, Entry: "/tmp/view/entry.md", Revision: "r2"}
	if e := r.Open(context.Background(), "xfmd", d); e != nil {
		t.Fatal(e)
	}
	if c.received[0].Window != "window-one" {
		t.Fatal("focus-dependent delivery")
	}
	if e := r.Open(context.Background(), "shell", d); e == nil {
		t.Fatal("unregistered adapter")
	}
	d.Window = "--active"
	d.Pane = "other"
	if e := r.Open(context.Background(), "xfmd", d); e == nil {
		t.Fatal("invalid target")
	}
}
