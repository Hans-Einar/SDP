package model

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestValidateReadOnlySharedSemantics(t *testing.T) {
	area := t.TempDir()
	path := initial(t, area, "Design")
	must(t, os.WriteFile(filepath.Join(path, "System.design"), []byte("language design-core version 0.6.\nsystem Example.\nincludes \"Parts.design\".\nExample contains Worker.\n"), 0600))
	must(t, os.WriteFile(filepath.Join(path, "Parts.design"), []byte("language design-core version 0.6.\nunit Worker.\n"), 0600))
	before, e := Snapshot(area, "work:Design")
	must(t, e)
	checked, e := Validate(area, "work:Design")
	must(t, e)
	after, e := Snapshot(area, "work:Design")
	must(t, e)
	if !reflect.DeepEqual(before, after) || checked.Source.Status != "dirty" || !checked.Source.Preliminary || !reflect.DeepEqual(checked.Targets, []string{"System.design"}) {
		t.Fatal(checked)
	}
	candidate, e := Freeze(area, "candidate", "Reviewed", "work:Design", "")
	must(t, e)
	frozen, e := Validate(area, "candidate:Reviewed")
	must(t, e)
	if frozen.Source.Artifact.ID != candidate.Artifact.ID || !reflect.DeepEqual(frozen.Targets, checked.Targets) {
		t.Fatal(frozen)
	}
}
func TestValidateRejectsSemanticDriftAndOrphans(t *testing.T) {
	for _, kind := range []string{"missing-member", "orphan", "changed-during-read"} {
		t.Run(kind, func(t *testing.T) {
			area := t.TempDir()
			path := initial(t, area, "Design")
			entry := filepath.Join(path, "System.design")
			good := []byte("language design-core version 0.6.\nsystem Example.\nunit Worker.\nExample contains Worker.\n")
			must(t, os.WriteFile(entry, good, 0600))
			var recheck func()
			switch kind {
			case "missing-member":
				must(t, os.WriteFile(entry, []byte("language design-core version 0.6.\nsystem Example.\nExample contains Missing.\n"), 0600))
			case "orphan":
				must(t, os.WriteFile(filepath.Join(path, "Orphan.design"), []byte("language design-core version 0.6.\nunit Orphan.\n"), 0600))
			case "changed-during-read":
				recheck = func() { must(t, os.WriteFile(entry, append(good, []byte("// Changed during validation.\n")...), 0600)) }
			}
			if _, e := validateSource(area, "work:Design", recheck); e == nil {
				t.Fatal("invalid or changed model accepted")
			}
		})
	}
}

func TestValidatePreservesPendingRecoveryHold(t *testing.T) {
	for _, during := range []bool{false, true} {
		t.Run(map[bool]string{false: "already-pending", true: "pending-during-validation"}[during], func(t *testing.T) {
			area := t.TempDir()
			path := initial(t, area, "Design")
			must(t, os.WriteFile(filepath.Join(path, "System.design"), []byte("language design-core version 0.6.\nsystem Example.\nunit Worker.\nExample contains Worker.\n"), 0600))
			interrupt := func() {
				faultHook = func(phase string) {
					if phase == "prepared" {
						panic("interrupted fixture")
					}
				}
				defer func() { _ = recover(); faultHook = nil }()
				_, _ = Commit(area, "work:Design", "Interrupted fixture commit", false)
			}
			var recheck func()
			if during {
				recheck = interrupt
			} else {
				interrupt()
			}
			if _, err := validateSource(area, "work:Design", recheck); err == nil {
				t.Fatal("pending operation accepted")
			}
			pending, err := journals(area)
			must(t, err)
			if len(pending) != 1 {
				t.Fatal("validation changed pending recovery", pending)
			}
			before, err := os.ReadFile(filepath.Join(area, operations, pending[0].ID, "journal.yaml"))
			must(t, err)
			if _, err = Validate(area, "work:Design"); err == nil {
				t.Fatal("retry accepted pending operation")
			}
			after, err := os.ReadFile(filepath.Join(area, operations, pending[0].ID, "journal.yaml"))
			must(t, err)
			if !reflect.DeepEqual(before, after) {
				t.Fatal("validation rewrote pending journal")
			}
			_, err = Recover(area, pending[0].ID, "resume")
			must(t, err)
			_, err = Validate(area, "work:Design")
			must(t, err)
		})
	}
}
