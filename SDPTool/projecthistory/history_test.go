package projecthistory

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func event(n int) []byte {
	return []byte(fmt.Sprintf(`{"schemaVersion":"1.0","eventId":"EVT-BP-TEST-%06d","eventType":"x-blueprint:trial","occurredAt":"2026-10-09T00:00:00Z","actor":"test","commit":null,"subjectId":"TEST","payload":{"index":%d}}`, n, n))
}
func requireLinux(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("mutation supported on Linux")
	}
}
func TestReadAndEnvelope(t *testing.T) {
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	s, e := Read(p)
	if e != nil || s.Revision != hash(nil) || len(s.Bytes) != 0 {
		t.Fatal(s, e)
	}
	// Exercise the actual mixed historical envelopes, without writing the repository.
	if _, e = Read("../../SDP/ProjectManagement/Ledger.ndjson"); e != nil {
		t.Fatal(e)
	}
	for _, value := range []string{`123`, `""`, `null`, `true`} {
		bad := bytes.Replace(event(1), []byte(`"subjectId":"TEST"`), []byte(`"subjectId":`+value), 1)
		if _, _, err := envelope(bad); err == nil {
			t.Fatalf("accepted subject %s", value)
		}
	}
	for _, value := range []string{`123`, `""`, `null`, `"invalid"`, `"REL-01.2.3"`, `"REL-1.2.3-01"`} {
		bad := bytes.Replace(event(1), []byte(`"subjectId":"TEST"`), []byte(`"releaseId":`+value), 1)
		if _, _, err := envelope(bad); err == nil {
			t.Fatalf("accepted release %s", value)
		}
	}
	for _, value := range []string{"REL-1.2.3", "REL-0.2.0-alpha.1+build.001"} {
		good := bytes.Replace(event(1), []byte(`"subjectId":"TEST"`), []byte(`"releaseId":"`+value+`"`), 1)
		if _, _, err := envelope(good); err != nil {
			t.Fatal(err)
		}
	}
	for _, bad := range [][]byte{
		bytes.Replace(event(1), []byte(`test`), []byte{0xff}, 1),
		append(event(1), 'x'), []byte(`{}`), []byte(`[]`),
		[]byte(strings.Replace(string(event(1)), `"actor":"test"`, `"actor":"test","actor":"other"`, 1)),
		[]byte(strings.Replace(string(event(1)), `"eventId":`, `"EventId":"EVT-OVERRIDE-1","eventId":`, 1)),
		[]byte(strings.Replace(string(event(1)), `"payload":{"index":1}`, `"payload":null`, 1)),
	} {
		if _, _, e = envelope(bad); e == nil {
			t.Fatalf("accepted invalid envelope %s", bad)
		}
	}
}
func TestAppendPreservesHistoryAndIdempotence(t *testing.T) {
	requireLinux(t)
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	original := append([]byte(" \n"), event(1)...)
	original = append(original, '\n', '\n')
	if e := os.WriteFile(p, original, 0640); e != nil {
		t.Fatal(e)
	}
	r, e := Append(p, hash(original), event(2))
	if e != nil || !r.Appended {
		t.Fatal(r, e)
	}
	got, e := os.ReadFile(p)
	if e != nil || !bytes.HasPrefix(got, original) {
		t.Fatal("history rewritten", e)
	}
	st, _ := os.Stat(p)
	if st.Mode().Perm() != 0640 {
		t.Fatal("permissions changed")
	}
	retry, e := Append(p, hash(original), event(2))
	if e != nil || retry.Appended || retry.Revision != r.Revision {
		t.Fatal(retry, e)
	}
	if _, e = Append(p, hash(original), event(3)); !errors.Is(e, ErrStale) {
		t.Fatal("stale accepted", e)
	}
	changed := bytes.Replace(event(2), []byte(`"actor":"test"`), []byte(`"actor":"other"`), 1)
	if _, e = Append(p, r.Revision, changed); !errors.Is(e, ErrConflict) {
		t.Fatal("conflict accepted", e)
	}
	if _, e = Append(p, r.Revision, event(3)); e != nil {
		t.Fatal(e)
	}
	if r, e = Append(p, hash(original), event(2)); e != nil || r.Appended {
		t.Fatal("retry after later event", r, e)
	}
}
func TestAppendFailuresAndObservedInterference(t *testing.T) {
	requireLinux(t)
	for _, point := range []string{"before-publish", "after-publish", "after-sync"} {
		t.Run(point, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "Ledger.ndjson")
			before := hash(nil)
			r, e := appendEvent(p, before, event(1), func(stage string) error {
				if stage == point {
					return errors.New("injected")
				}
				return nil
			})
			if e == nil {
				t.Fatal("fault missed")
			}
			s, e := Read(p)
			if e != nil {
				t.Fatal(e)
			}
			if point == "before-publish" && len(s.Bytes) != 0 {
				t.Fatal("premature publication")
			}
			if point != "before-publish" && (!r.Appended || len(s.Bytes) == 0) {
				t.Fatal("missing committed result")
			}
			if _, e = Append(p, before, event(1)); e != nil {
				t.Fatal("retry failed", e)
			}
			s, _ = Read(p)
			events, e := validate(s.Bytes)
			if e != nil || len(events) != 1 {
				t.Fatal("duplicate retry", e)
			}
		})
	}
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	other := append(event(9), '\n')
	_, e := appendEvent(p, hash(nil), event(1), func(stage string) error {
		if stage == "before-publish" {
			return os.WriteFile(p, other, 0600)
		}
		return nil
	})
	if !errors.Is(e, ErrStale) {
		t.Fatal(e)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got, other) {
		t.Fatal("interfering write lost")
	}
}
func TestMalformedTailLimitsAndSymlinks(t *testing.T) {
	requireLinux(t)
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	for _, value := range []string{`123`, `""`, `null`, `true`} {
		bad := bytes.Replace(event(1), []byte(`"subjectId":"TEST"`), []byte(`"subjectId":`+value), 1)
		if _, _, err := envelope(bad); err == nil {
			t.Fatalf("accepted subject %s", value)
		}
	}
	for _, value := range []string{`123`, `""`, `null`, `"invalid"`, `"REL-01.2.3"`, `"REL-1.2.3-01"`} {
		bad := bytes.Replace(event(1), []byte(`"subjectId":"TEST"`), []byte(`"releaseId":`+value), 1)
		if _, _, err := envelope(bad); err == nil {
			t.Fatalf("accepted release %s", value)
		}
	}
	for _, value := range []string{"REL-1.2.3", "REL-0.2.0-alpha.1+build.001"} {
		good := bytes.Replace(event(1), []byte(`"subjectId":"TEST"`), []byte(`"releaseId":"`+value+`"`), 1)
		if _, _, err := envelope(good); err != nil {
			t.Fatal(err)
		}
	}
	for _, bad := range [][]byte{event(1), append(append(append(event(1), '\n'), event(1)...), '\n'), []byte("{broken}\n")} {
		os.WriteFile(p, bad, 0600)
		if _, e := Read(p); e == nil {
			t.Fatal("bad stream accepted")
		}
		if _, e := Append(p, hash(bad), event(2)); e == nil {
			t.Fatal("bad stream repaired")
		}
		got, _ := os.ReadFile(p)
		if !bytes.Equal(got, bad) {
			t.Fatal("bad stream changed")
		}
	}
	os.Remove(p)
	target := filepath.Join(t.TempDir(), "real")
	os.WriteFile(target, event(1), 0600)
	if e := os.Symlink(target, p); e != nil {
		t.Skip(e)
	}
	if _, e := Read(p); e == nil {
		t.Fatal("symlink read")
	}
	if _, e := Append(p, hash(nil), event(2)); e == nil {
		t.Fatal("symlink write")
	}
	os.Remove(p)
	os.Remove(p + ".lock")
	if e := os.Symlink(target, p+".lock"); e != nil {
		t.Fatal(e)
	}
	if _, e := Append(p, hash(nil), event(2)); e == nil {
		t.Fatal("lock alias followed")
	}
	big := filepath.Join(t.TempDir(), "big")
	f, _ := os.Create(big)
	f.Truncate(MaxBytes + 1)
	f.Close()
	if _, e := Read(big); e == nil {
		t.Fatal("size limit")
	}
}
func TestConcurrentCompareAndAppend(t *testing.T) {
	requireLinux(t)
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := 0
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			r, e := Append(p, hash(nil), event(n))
			mu.Lock()
			defer mu.Unlock()
			if e == nil && r.Appended {
				success++
			} else if !errors.Is(e, ErrBusy) && !errors.Is(e, ErrStale) {
				t.Errorf("unexpected %v", e)
			}
		}(i)
	}
	wg.Wait()
	if success != 1 {
		t.Fatal("lost compare-and-append", success)
	}
	s, e := Read(p)
	if e != nil {
		t.Fatal(e)
	}
	ev, e := validate(s.Bytes)
	if e != nil || len(ev) != 1 {
		t.Fatal(e, len(ev))
	}
}
func TestHistoryCrashHelper(t *testing.T) {
	p := os.Getenv("SDP_HISTORY_TEST_PATH")
	if p == "" {
		return
	}
	stage := os.Getenv("SDP_HISTORY_TEST_STAGE")
	_, e := appendEvent(p, hash(nil), event(1), func(s string) error {
		if s == stage {
			os.Exit(73)
		}
		return nil
	})
	if e != nil {
		os.Exit(74)
	}
	os.Exit(75)
}
func TestProcessExitAndCrossProcessLock(t *testing.T) {
	requireLinux(t)
	for _, stage := range []string{"before-publish", "after-publish", "after-sync"} {
		t.Run(stage, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "Ledger.ndjson")
			cmd := exec.Command(os.Args[0], "-test.run=^TestHistoryCrashHelper$")
			cmd.Env = append(os.Environ(), "SDP_HISTORY_TEST_PATH="+p, "SDP_HISTORY_TEST_STAGE="+stage)
			if e := cmd.Run(); e == nil || e.(*exec.ExitError).ExitCode() != 73 {
				t.Fatal("child did not exit at selected boundary", e)
			}
			s, e := Read(p)
			if e != nil {
				t.Fatal("torn stream", e)
			}
			if stage == "before-publish" && len(s.Bytes) != 0 {
				t.Fatal("unexpected commit")
			}
			if _, e = Append(p, hash(nil), event(1)); e != nil {
				t.Fatal("lock not released or retry failed", e)
			}
			s, _ = Read(p)
			ev, e := validate(s.Bytes)
			if e != nil || len(ev) != 1 {
				t.Fatal("retry duplicated", e)
			}
		})
	}
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	unlock, e := lock(p)
	if e != nil {
		t.Fatal(e)
	}
	defer unlock()
	cmd := exec.Command(os.Args[0], "-test.run=^TestHistoryCrashHelper$")
	cmd.Env = append(os.Environ(), "SDP_HISTORY_TEST_PATH="+p, "SDP_HISTORY_TEST_STAGE=after-sync")
	if e = cmd.Run(); e == nil || e.(*exec.ExitError).ExitCode() != 74 {
		t.Fatal("cross-process lock bypass", e)
	}
}
func TestAppendedEventIsOneJSONLine(t *testing.T) {
	requireLinux(t)
	p := filepath.Join(t.TempDir(), "Ledger.ndjson")
	var buf bytes.Buffer
	json.Indent(&buf, event(1), "", "  ")
	if _, e := Append(p, hash(nil), buf.Bytes()); e != nil {
		t.Fatal(e)
	}
	s, _ := Read(p)
	if bytes.Count(s.Bytes, []byte{'\n'}) != 1 {
		t.Fatal("multiline stream event")
	}
}
