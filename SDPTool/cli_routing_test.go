package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/install"
)

func TestCLIRoutingErrorPriority(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		args                       []string
		exit                       int
		schema, operation, message string
		stdout                     bool
	}{
		{[]string{"invoke", "--json=broken"}, 2, InvokeSchema, "", "", true},
		{[]string{root, "invoke", "--json"}, 2, InvokeSchema, "", "invoke must be the first argument", true},
		{[]string{"--json", "invoke"}, 2, InvokeSchema, "", "invoke must be the first argument", true},
		{[]string{"--json", "install", "--json=broken"}, 2, install.Protocol, "install", "invalid --json value", true},
		{[]string{"--json", root, "upgrade", "--json=broken"}, 2, install.Protocol, "upgrade", "invalid --json value", true},
		{[]string{"--json", "--version", "--json=broken"}, 1, Version, "", "invalid --json value", false},
		{[]string{"--json", "preview"}, 1, Version, "", "usage: sdptool preview FILE", false},
		{[]string{"unknown", "--json"}, 1, Version, "", "usage: sdptool preview FILE", false},
	}
	for _, tc := range cases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			var out, errs bytes.Buffer
			code := RunWithInput(context.Background(), tc.args, strings.NewReader(""), &out, &errs)
			if code != tc.exit {
				t.Fatalf("exit=%d want=%d: %s %s", code, tc.exit, &out, &errs)
			}
			data := errs.Bytes()
			if tc.stdout {
				data = out.Bytes()
			}
			var envelope struct {
				Schema        string `json:"schema"`
				SchemaVersion string `json:"schemaVersion"`
				Operation     string `json:"operation"`
				Error         struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(data, &envelope); err != nil {
				t.Fatalf("envelope: %v %s", err, data)
			}
			schema := envelope.Schema
			if schema == "" {
				schema = envelope.SchemaVersion
			}
			if schema != tc.schema || envelope.Operation != tc.operation || !strings.Contains(envelope.Error.Message, tc.message) {
				t.Fatalf("unexpected envelope: %s", data)
			}
		})
	}
}
