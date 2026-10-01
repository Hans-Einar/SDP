package bootstrap

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func fixture(t *testing.T) (Config, string) {
	t.Helper()
	dir := t.TempDir()
	public, private, e := ed25519.GenerateKey(nil)
	if e != nil {
		t.Fatal(e)
	}
	binary := []byte("#!/bin/sh\nprintf '%s\\n' '{\"installationProtocol\":\"sdp-install-command/1\"}'\n")
	bin := filepath.Join(dir, "engine")
	os.WriteFile(bin, binary, 0700)
	d := map[string]any{"schemaVersion": "sdp-release-descriptor/1", "protocol": Protocol, "binaries": []asset{{runtime.GOOS + "/" + runtime.GOARCH, "engine", hash(binary), int64(len(binary))}}}
	b, _ := json.Marshal(d)
	p := filepath.Join(dir, "release.json")
	os.WriteFile(p, b, 0600)
	sig, _ := json.Marshal(signature{hash(public), ed25519.Sign(private, b)})
	os.WriteFile(p+".sig", sig, 0600)
	key := filepath.Join(dir, "test.pub")
	os.WriteFile(key, []byte(base64.StdEncoding.EncodeToString(public)), 0600)
	return Config{Descriptor: p, TestKey: key, CacheDir: filepath.Join(t.TempDir(), "cache")}, bin
}
func TestSignedCacheOfflineAndCorruption(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("native test binary is a POSIX fixture")
	}
	c, source := fixture(t)
	p, e := Bootstrap(context.Background(), c)
	if e != nil {
		t.Fatal(e)
	}
	os.Remove(source)
	os.Remove(c.Descriptor)
	os.Remove(c.Descriptor + ".sig")
	c.Offline = true
	if q, e := Bootstrap(context.Background(), c); e != nil || q != p {
		t.Fatal(q, e)
	}
	os.WriteFile(p, []byte("corrupt"), 0700)
	if _, e = Bootstrap(context.Background(), c); e == nil {
		t.Fatal("corrupt cache accepted")
	}
}
func TestTrustAndProtocolFailures(t *testing.T) {
	c, _ := fixture(t)
	missing := c
	missing.Offline = true
	if _, e := Resolve(context.Background(), missing); e == nil {
		t.Fatal("missing offline cache")
	}
	wrong := c
	wrong.TestKey = ""
	if _, e := Resolve(context.Background(), wrong); e == nil {
		t.Fatal("untrusted signer")
	}
	os.WriteFile(c.Descriptor, []byte("{}"), 0600)
	if _, e := Resolve(context.Background(), c); e == nil {
		t.Fatal("signature mismatch")
	}
	if e := StrictJSON([]byte(`{"x":1,"X":2}`), &map[string]any{}); e == nil {
		t.Fatal("ambiguous JSON accepted")
	}
}

func TestSignedDescriptorCompatibilityAndBounds(t *testing.T) {
	for _, kind := range []string{"protocol", "platform", "size", "path", "binary-digest", "probe"} {
		t.Run(kind, func(t *testing.T) {
			c, _ := fixture(t)
			public, private, _ := ed25519.GenerateKey(nil)
			os.WriteFile(c.TestKey, []byte(base64.StdEncoding.EncodeToString(public)), 0600)
			data, _ := os.ReadFile(c.Descriptor)
			var d map[string]any
			json.Unmarshal(data, &d)
			a := d["binaries"].([]any)[0].(map[string]any)
			switch kind {
			case "protocol":
				d["protocol"] = "unsupported"
			case "platform":
				a["platform"] = "unknown/architecture"
			case "size":
				a["size"] = maxBinary + 1
			case "path":
				a["path"] = "../escape"
			case "binary-digest":
				a["sha256"] = hash([]byte("wrong"))
			case "probe":
				b := []byte("#!/bin/sh\nprintf '%s\\n' '{}'\n")
				os.WriteFile(filepath.Join(filepath.Dir(c.Descriptor), "engine"), b, 0700)
				a["sha256"] = hash(b)
				a["size"] = len(b)
			}
			data, _ = json.Marshal(d)
			os.WriteFile(c.Descriptor, data, 0600)
			sig, _ := json.Marshal(signature{hash(public), ed25519.Sign(private, data)})
			os.WriteFile(c.Descriptor+".sig", sig, 0600)
			if _, e := Bootstrap(context.Background(), c); e == nil {
				t.Fatal("accepted incompatible input")
			}
		})
	}
}
func TestImmutableSelectorAndReadBound(t *testing.T) {
	c, _ := fixture(t)
	if _, e := Resolve(context.Background(), c); e != nil {
		t.Fatal(e)
	}
	public, private, _ := ed25519.GenerateKey(nil)
	data := []byte(`{"protocol":"sdp-install-command/1","revision":"changed"}`)
	os.WriteFile(c.Descriptor, data, 0600)
	sig, _ := json.Marshal(signature{hash(public), ed25519.Sign(private, data)})
	os.WriteFile(c.Descriptor+".sig", sig, 0600)
	os.WriteFile(c.TestKey, []byte(base64.StdEncoding.EncodeToString(public)), 0600)
	if _, e := Resolve(context.Background(), c); e == nil {
		t.Fatal("mutable selector accepted")
	}
	p := filepath.Join(t.TempDir(), "large")
	os.WriteFile(p, []byte("12345"), 0600)
	if _, e := read(p, 4); e == nil {
		t.Fatal("local read bound ignored")
	}
}

// The POSIX fixture stands in for old/new signed executables, not a runtime
// implementation dependency. Production probing uses exec.CommandContext.
func TestExplicitJSONProbeAndLegacyFallback(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable fixture")
	}
	for _, tc := range []struct {
		name, script string
		ok           bool
	}{
		{"new", `[ "$1" = "--version" ] && [ "$2" = "--json" ] || exit 2
printf '%s\n' '{"installationProtocol":"sdp-install-command/1"}'`, true},
		{"legacy", `[ "$#" = 1 ] && [ "$1" = "--version" ] || exit 2
printf '%s\n' '{"installationProtocol":"sdp-install-command/1"}'`, true},
		{"successful-invalid-is-not-retried", `if [ "$#" = 2 ]; then printf '%s\n' 'not json'; else printf '%s\n' '{"installationProtocol":"sdp-install-command/1"}'; fi`, false},
		{"wrong-protocol", `printf '%s\n' '{"installationProtocol":"other/1"}'`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, bin := fixture(t)
			binary := []byte("#!/bin/sh\n" + tc.script + "\n")
			if e := os.WriteFile(bin, binary, 0700); e != nil {
				t.Fatal(e)
			}
			public, private, e := ed25519.GenerateKey(nil)
			if e != nil {
				t.Fatal(e)
			}
			d := map[string]any{"schemaVersion": "sdp-release-descriptor/1", "protocol": Protocol, "binaries": []asset{{runtime.GOOS + "/" + runtime.GOARCH, "engine", hash(binary), int64(len(binary))}}}
			b, _ := json.Marshal(d)
			sig, _ := json.Marshal(signature{hash(public), ed25519.Sign(private, b)})
			for p, data := range map[string][]byte{c.Descriptor: b, c.Descriptor + ".sig": sig, c.TestKey: []byte(base64.StdEncoding.EncodeToString(public))} {
				if e := os.WriteFile(p, data, 0600); e != nil {
					t.Fatal(e)
				}
			}
			_, e = Bootstrap(context.Background(), c)
			if (e == nil) != tc.ok {
				t.Fatalf("expected success=%v: %v", tc.ok, e)
			}
		})
	}
}
