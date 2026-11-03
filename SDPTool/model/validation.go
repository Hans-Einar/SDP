package model

import "reflect"

// Validation is a read-only observation of a captured model. It neither commits
// WORK nor grants design acceptance. Frozen artifacts use the same validator.
type Validation struct {
	Schema    string   `json:"schema"`
	Operation string   `json:"operation"`
	Source    Result   `json:"source"`
	Targets   []string `json:"targets"`
}

func Validate(area, ref string) (Validation, error) {
	return validateSource(area, ref, nil)
}

func validateSource(area, ref string, beforeRecheck func()) (Validation, error) {
	if err := validationRecovery(area); err != nil {
		return Validation{}, err
	}
	first, err := Snapshot(area, ref)
	if err != nil {
		return Validation{}, err
	}
	if first.Artifact == nil || first.Status == "conflicted" || len(first.Artifact.PendingParents) > 0 {
		return Validation{}, fail("conflict", "unresolved model source")
	}
	targets, err := validateSources(first.Files)
	if err != nil {
		return Validation{}, err
	}
	if beforeRecheck != nil {
		beforeRecheck()
	}
	if err := validationRecovery(area); err != nil {
		return Validation{}, err
	}
	second, err := Snapshot(area, ref)
	if err != nil {
		return Validation{}, err
	}
	if !reflect.DeepEqual(first, second) {
		return Validation{}, fail("stale", "model changed during validation")
	}
	return Validation{Schema: Schema, Operation: "validate", Source: first.Result, Targets: targets}, nil
}

// Observe journals without acquiring a mutating model lock or repairing staging.
func validationRecovery(area string) error {
	pending, err := journals(area)
	if err != nil {
		return err
	}
	if len(pending) > 0 {
		return fail("recovery", "pending model operation requires recovery")
	}
	return nil
}
