package governancecontroller

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func TestApprovalBookkeepingWaitsWithoutResending(t *testing.T) {
	c, h, id, _, reply := approvalFixture(t)
	r, _ := c.Service.Read(id)
	coverage, _ := c.Service.Context(r.Run.Selection)
	lockfile, e := os.OpenFile(filepath.Join(filepath.Dir(coverage.Root), "state", "lock"), os.O_RDWR, 0600)
	if e != nil {
		t.Fatal(e)
	}
	defer lockfile.Close()
	released := make(chan struct{})
	h.reply = func(governance.ApprovalReply) error {
		if e := syscall.Flock(int(lockfile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
			t.Fatal(e)
		}
		go func() {
			time.Sleep(100 * time.Millisecond)
			syscall.Flock(int(lockfile.Fd()), syscall.LOCK_UN)
			close(released)
		}()
		return nil
	}
	result, e := c.Approve(id, reply, "explicit fixture acceptance")
	<-released
	if e != nil || h.replies != 1 || result.Run.Inputs[len(result.Run.Inputs)-1].State != "uncertain" {
		t.Fatal(result, e, h.replies)
	}
}
