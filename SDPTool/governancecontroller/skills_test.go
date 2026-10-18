package governancecontroller

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type selectedHost struct {
	*fakeHost
	inventory governance.SkillInventory
	drift     bool
	sent      []governance.SkillInput
}

func (f *selectedHost) InspectSkills(context.Context, string) (governance.SkillInventory, error) {
	return f.inventory, nil
}
func (f *selectedHost) StartSkillTurn(_ context.Context, thread, id, text string, skills []governance.SkillInput) (governancehost.Turn, error) {
	f.sent = skills
	v, e := f.c.Service.OwnerInput(governance.LocalController(), f.runID, id)
	if e != nil || v.State != "sending" || len(v.Skills) != 1 || v.Text != "Fix the selected fixture" || strings.TrimSpace(strings.SplitN(text, "\n", 2)[0]) != "$sdp" {
		f.t.Fatal("skill was not bound before send", v, e)
	}
	return f.fakeHost.StartTurn(context.Background(), thread, id, text)
}
func (f *selectedHost) MCPTools(ctx context.Context, server, thread string) ([]string, string, error) {
	if f.drift {
		os.WriteFile(filepath.Join(f.c.Root, "Skills/sdp/SKILL.md"), []byte("changed after claim"), 0600)
	}
	return f.fakeHost.MCPTools(ctx, server, thread)
}
func skillFixture(t *testing.T) (*Controller, *selectedHost, string) {
	c, f, id := setup(t)
	r, _ := c.Service.Read(id)
	ctx, _ := c.Service.Context(r.Run.Selection)
	c.Root = ctx.Root
	path := filepath.Join(c.Root, "Skills/sdp/SKILL.md")
	os.MkdirAll(filepath.Dir(path), 0700)
	os.WriteFile(path, []byte("skill"), 0600)
	host := &selectedHost{fakeHost: f, inventory: governance.SkillInventory{Status: "observed", UseEvidence: "not-observed", Project: []governance.SkillAvailability{{Name: "sdp", Path: "Skills/sdp/SKILL.md", Enabled: true, Digest: fmt.Sprintf("%x", sha256.Sum256([]byte("skill")))}}}}
	c.Host = host
	c.SkillPaths = []string{"Skills/sdp/SKILL.md"}
	return c, host, id
}
func TestExplicitSkillLaunchAndUncertainSend(t *testing.T) {
	for _, lost := range []bool{false, true} {
		t.Run(fmt.Sprint(lost), func(t *testing.T) {
			c, f, id := skillFixture(t)
			f.lost = lost
			r, e := c.Launch(context.Background(), id, "Fix the selected fixture")
			if lost {
				if e == nil {
					t.Fatal("lost reply accepted")
				}
				r, _ = c.Service.Read(id)
				if r.Run.Inputs[0].State != "uncertain" {
					t.Fatal(r)
				}
			} else if e != nil {
				t.Fatal(e)
			}
			if f.starts != 1 || len(f.sent) != 1 {
				t.Fatal("missing skill wire call", f.starts, f.sent)
			}
			if _, e = c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil || f.starts != 1 {
				t.Fatal("resent skill turn", e)
			}
			view, e := c.Service.Provenance(id)
			if e != nil || view.Invocation == nil || view.Invocation.Compliance != "not-observed" {
				t.Fatal(view, e)
			}
		})
	}
}
func TestSelectedSkillRefusalsBeforeIntent(t *testing.T) {
	for _, name := range []string{"missing", "disabled", "duplicate", "ambiguous", "external", "digest", "errors", "too-many"} {
		t.Run(name, func(t *testing.T) {
			c, f, id := skillFixture(t)
			switch name {
			case "missing":
				c.SkillPaths = []string{"Skills/missing/SKILL.md"}
			case "disabled":
				f.inventory.Project[0].Enabled = false
			case "duplicate":
				c.SkillPaths = append(c.SkillPaths, c.SkillPaths[0])
			case "ambiguous":
				f.inventory.Project = append(f.inventory.Project, governance.SkillAvailability{Name: "sdp", Path: "Skills/other/SKILL.md"})
			case "external":
				c.SkillPaths = []string{"/tmp/SKILL.md"}
			case "digest":
				f.inventory.Project[0].Digest = strings.Repeat("0", 64)
			case "errors":
				f.inventory.ErrorCount = 1
			case "too-many":
				c.SkillPaths = make([]string, 9)
			}
			if _, e := c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil {
				t.Fatal("invalid selection accepted")
			}
			r, _ := c.Service.Read(id)
			if len(r.Run.Inputs) != 0 || len(r.Run.Attempts) != 0 || f.starts != 0 {
				t.Fatal("invalid selection caused effects", r)
			}
		})
	}
}
func TestSkillDriftBeforeTurnRefuses(t *testing.T) {
	c, f, id := skillFixture(t)
	f.drift = true
	if _, e := c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil {
		t.Fatal("drift accepted")
	}
	r, _ := c.Service.Read(id)
	if f.starts != 0 || r.Run.Inputs[0].State != "uncertain" {
		t.Fatal(f.starts, r)
	}
}
