package activities

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

const legacyShapeDescriptor = "schema: sdp-activity/0.1\nid: ACTIVITY-TEST-0001\nprojectId: TEST\ntitle: Manual need\nowner: project-owner\nmandate: 01-Mandate/Mandate.md\nsourceCards: [KB-TEST-001]\nsystems: []\nreferences:\n- role: session\n  id: SESSION-TEST-0007\n  path: Sessions/Prior.md\n"

func legacyShapeReplay(t *testing.T, descriptor string, extension string, value any) (Catalog, error) {
	t.Helper()
	ref := DescriptorReference{Path: "Activities/Activity-0001--Manual/activity.yaml", SHA256: digest([]byte(descriptor))}
	event := Event{SchemaVersion: "1.0", EventID: "EVT-ACTIVITY-TEST-000001", EventType: "x-activity:created", OccurredAt: "2026-11-12T00:00:00Z", Actor: "manual-owner", SubjectID: "ACTIVITY-TEST-0001", Payload: Payload{SchemaVersion: "0.1", ProjectID: "TEST", OperationID: "historical-manual-operation", RecordedAt: "2026-11-12T00:00:00Z", Authority: Authority{Role: "project-owner", Source: "owner.md"}, To: "planned", Reason: "Actual manual registration", Descriptor: ref, Links: []string{}}}
	raw, e := json.Marshal(event)
	if e != nil {
		t.Fatal(e)
	}
	if extension != "" {
		var envelope map[string]any
		json.Unmarshal(raw, &envelope)
		payload := envelope["payload"].(map[string]any)
		payload[extension] = value
		raw, e = json.Marshal(envelope)
		if e != nil {
			t.Fatal(e)
		}
	}
	raw = append(raw, '\n')
	return ReplayRecords("TEST", projecthistory.Snapshot{Bytes: raw, Revision: digest(raw)}, map[string][]byte{ref.Path: []byte(descriptor)})
}
func TestLegacy01PayloadRejectsExtensionPresence(t *testing.T) {
	if c, e := legacyShapeReplay(t, legacyShapeDescriptor, "", nil); e != nil || len(c.Activities) != 1 {
		t.Fatal("valid immutable legacy shape rejected", c, e)
	}
	for _, tc := range []struct {
		key    string
		values []any
	}{{"requestDigest", []any{nil, ""}}, {"session", []any{nil}}, {"batch", []any{nil}}, {"outcome", []any{nil, ""}}, {"evidence", []any{nil, []any{}}}, {"resumeCondition", []any{nil, "", "unsupported legacy field"}}, {"contributionDispositions", []any{nil, []any{}}}} {
		for _, value := range tc.values {
			t.Run(tc.key+"/"+stringValue(value), func(t *testing.T) {
				c, e := legacyShapeReplay(t, legacyShapeDescriptor, tc.key, value)
				if e == nil || !strings.Contains(e.Error(), "legacy payload extension") || len(c.Activities) != 0 {
					t.Fatalf("0.1 extension presence accepted: %s=%v catalog=%+v error=%v", tc.key, value, c, e)
				}
			})
		}
	}
}
func stringValue(v any) string { raw, _ := json.Marshal(v); return string(raw) }
func TestLegacy01DescriptorRejectsExtensionPresence(t *testing.T) {
	for _, key := range []string{"scope", "dependencies"} {
		for _, value := range []string{"null", "[]"} {
			t.Run(key+"/"+value, func(t *testing.T) {
				c, e := legacyShapeReplay(t, legacyShapeDescriptor+key+": "+value+"\n", "", nil)
				if e == nil || !strings.Contains(e.Error(), "legacy descriptor extension") || len(c.Activities) != 0 {
					t.Fatalf("0.1 descriptor extension accepted: %s=%s catalog=%+v error=%v", key, value, c, e)
				}
			})
		}
	}
}
func TestLegacy01ReferenceRejectsExtensionPresence(t *testing.T) {
	for _, key := range []string{"projectId", "membership", "revision", "availability", "unavailableReason"} {
		for _, value := range []string{"null", `""`} {
			t.Run(key+"/"+value, func(t *testing.T) {
				descriptor := legacyShapeDescriptor + "  " + key + ": " + value + "\n"
				c, e := legacyShapeReplay(t, descriptor, "", nil)
				if e == nil || !strings.Contains(e.Error(), "legacy reference extension") || len(c.Activities) != 0 {
					t.Fatalf("0.1 reference extension accepted: %s=%s catalog=%+v error=%v", key, value, c, e)
				}
			})
		}
	}
}
func TestLegacy01DescriptorPreservesDeclaredTypes(t *testing.T) {
	for _, tc := range []struct{ name, from, to string }{{"null-optional-array", "systems: []", "systems: null"}, {"numeric-string", "title: Manual need", "title: 123"}, {"numeric-reference", "id: SESSION-TEST-0007", "id: 123"}} {
		t.Run(tc.name, func(t *testing.T) {
			descriptor := strings.Replace(legacyShapeDescriptor, tc.from, tc.to, 1)
			if _, e := parseDescriptor([]byte(descriptor)); e == nil {
				t.Fatal("immutable legacy field type silently coerced")
			}
		})
	}
}
