// Package reader addresses an explicitly registered local document consumer.
// A URI cannot supply a program, shell, window or arbitrary after-command.
package reader

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Delivery struct {
	Window, Pane, Client, Entry, Revision, Lease, Broker string
	Sequence                                             uint64
}
type Adapter interface {
	Open(context.Context, Delivery) error
}
type Registry map[string]Adapter

func (r Registry) Open(ctx context.Context, name string, d Delivery) error {
	a, ok := r[name]
	if !ok {
		return fmt.Errorf("unregistered-reader: %s", name)
	}
	if e := d.Validate(); e != nil {
		return e
	}
	return a.Open(ctx, d)
}

var id = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func (d Delivery) Validate() error {
	if !id.MatchString(d.Window) || !id.MatchString(d.Client) || d.Sequence == 0 || (d.Pane != "main" && d.Pane != "navigation") || !filepath.IsAbs(d.Entry) || strings.ContainsAny(d.Entry, "\x00\t\n\r") {
		return fmt.Errorf("invalid document delivery")
	}
	return nil
}

type XFMD struct{ Program string }

func (x XFMD) Open(ctx context.Context, d Delivery) error {
	if e := d.Validate(); e != nil {
		return e
	}
	if !filepath.IsAbs(x.Program) {
		return fmt.Errorf("register absolute XFMD executable")
	}
	args := []string{"--window", d.Window, "--pane", d.Pane, "--client", d.Client, "--request", strconv.FormatUint(d.Sequence, 10)}
	if d.Lease != "" {
		if !id.MatchString(d.Lease) || !filepath.IsAbs(d.Broker) || strings.ContainsAny(d.Broker, "\x00\t\r\n") {
			return fmt.Errorf("invalid lease registration")
		}
		args = append(args, "--lease", d.Lease, "--broker", d.Broker)
	}
	args = append(args, d.Entry)
	c := exec.CommandContext(ctx, x.Program, args...)
	out, e := c.CombinedOutput()
	if e != nil {
		return fmt.Errorf("xfmd-open: %w: %.1000s", e, out)
	}
	if string(out) != "OK\t"+strconv.FormatUint(d.Sequence, 10)+"\n" {
		return fmt.Errorf("xfmd-open: invalid acknowledgement")
	}
	return nil
}
