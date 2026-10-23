package codexcompat

import "testing"

func TestStableVersionAdmission(t *testing.T) {
	for _, v := range []string{"0.164.0", "0.164.1", "0.165.0", "0.200.0", "1.0.0"} {
		if !Supported(v) {
			t.Errorf("compatible candidate %q refused", v)
		}
	}
	for _, v := range []string{"", "0.0.0", "0.161.0", "0.163.99", "0.164", "00.164.0", "0.0164.0", "0.164.-1", "0.164.0-dev", "0.164.0+local", "0.164.0\n", "9999999999.0.0"} {
		if Supported(v) {
			t.Errorf("unsupported version %q admitted", v)
		}
	}
}
