package activities

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// Pin the observed manual bootstrap prefix instead of treating every future
// Activity/Session as part of a two-Activity compatibility fixture.
func manualPilotPrefix(t *testing.T, raw []byte) []byte {
	t.Helper()
	var out bytes.Buffer
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var e struct {
			EventType string `json:"eventType"`
			Payload   struct {
				SchemaVersion string `json:"schemaVersion"`
			} `json:"payload"`
		}
		if err := json.Unmarshal(line, &e); err != nil {
			t.Fatal(err)
		}
		if strings.HasPrefix(e.EventType, "x-activity") && e.Payload.SchemaVersion != "0.1" {
			break
		}
		out.Write(line)
		out.WriteByte('\n')
	}
	return out.Bytes()
}
