package pilot

import "testing"

func TestPreservedBehavior(t *testing.T) {
	if got := CalibrateMeasurement(4); got != 8 {
		t.Fatalf("calibrate = %d; want 8", got)
	}
	if got := ReduceMachineState(10, 3); got != 13 {
		t.Fatalf("reduce = %d; want 13", got)
	}
}
