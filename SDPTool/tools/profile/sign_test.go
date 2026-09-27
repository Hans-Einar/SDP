package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestSigningRejectsUntrustedAndExposedKeys(t *testing.T) {
	dir := t.TempDir()
	desc := filepath.Join(dir, "release.json")
	key := filepath.Join(dir, "key")
	os.WriteFile(desc, []byte("{}"), 0600)
	_, private, _ := ed25519.GenerateKey(nil)
	os.WriteFile(key, []byte(base64.StdEncoding.EncodeToString(private)), 0600)
	if sign(desc, key) == nil {
		t.Fatal("untrusted publisher accepted")
	}
	os.Chmod(key, 0644)
	if sign(desc, key) == nil {
		t.Fatal("publicly readable private key accepted")
	}
	if _, err := os.Stat(desc + ".sig"); !os.IsNotExist(err) {
		t.Fatal("published rejected signature")
	}
}

func TestProductionRequiresExplicitSigner(t *testing.T) {
	if run(t.TempDir(), "unused", "", "0.2.0", "binary", "") == nil {
		t.Fatal("unsigned production descriptor accepted")
	}
}
