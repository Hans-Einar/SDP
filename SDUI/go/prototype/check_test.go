package prototype

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRevisionAndBindings(t *testing.T) {
	p := filepath.Join(t.TempDir(), "source with spaces.sdui")
	os.WriteFile(p, []byte(`sdui 0.2; ref: model "missing.design"; page=[<go=button("Go",callback=model.go.@callback)>];`), 0600)
	_, r, e := Check(p, "page", "")
	if e != nil {
		t.Fatal(e)
	}
	if r.Status != "prototype-unbound" || len(r.Revision) != 64 {
		t.Fatal(r)
	}
	if _, _, e = Check(p, "page", "stale"); e == nil {
		t.Fatal("accepted stale")
	}
	if _, _, e = Check(p, "absent", ""); e == nil {
		t.Fatal("accepted missing entry")
	}
	os.WriteFile(p, []byte("invalid"), 0600)
	if _, _, e = Check(p, "page", ""); e == nil {
		t.Fatal("accepted invalid")
	}
}

func TestCheckUsesActualMarkdownProviderEvenWhenHidden(t *testing.T) {
	p := filepath.Join(t.TempDir(), "provider.sdui")
	for _, source := range []string{
		`sdui 0.2; page=["<script>alert(1)</script>"];`,
		`sdui 0.2; page=["<script>alert(1)</script>" {visible=false}];`,
	} {
		if e := os.WriteFile(p, []byte(source), 0600); e != nil {
			t.Fatal(e)
		}
		if _, _, e := Check(p, "page", ""); e == nil {
			t.Fatal("unsupported provider content passed readiness")
		}
	}
	if e := os.WriteFile(p, []byte("sdui 0.2; page=[\"# Valid Markdown\"];"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, r, e := Check(p, "page", ""); e != nil || r.Status != "prototype" {
		t.Fatal(r, e)
	}
}
