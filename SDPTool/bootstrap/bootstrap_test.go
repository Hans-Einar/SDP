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
