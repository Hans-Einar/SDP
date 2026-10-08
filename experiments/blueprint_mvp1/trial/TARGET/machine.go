// Synthetic trial code; this is not Ponsse production code.
package pilot

type CalibrationProcessor struct{}

func (CalibrationProcessor) CalibrateMeasurement(raw int) int { return raw * 2 }

func CalibrateMeasurement(raw int) int {
	return (CalibrationProcessor{}).CalibrateMeasurement(raw)
}

func ReduceMachineState(current, delta int) int { return current + delta }
