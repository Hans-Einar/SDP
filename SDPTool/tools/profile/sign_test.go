package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
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

func TestProfileRetainsExplicitReleasePredecessor(t *testing.T) {
	// Exercise the real authored payload and output validator, not a copied fixture.
	out := filepath.Join(t.TempDir(), "release.json")
	if err := build("../../..", out, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "development-predecessor-test", ""); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var descriptor struct {
		UpgradesFrom []string `json:"upgradesFrom"`
	}
	if err := json.Unmarshal(b, &descriptor); err != nil {
		t.Fatal(err)
	}
	if len(descriptor.UpgradesFrom) != 1 || descriptor.UpgradesFrom[0] != "edc0c72101a437c6e12c40a081ef59ae41cf0db32bcaedb73824ec48495aaee5" {
		t.Fatalf("published predecessor missing: %v", descriptor.UpgradesFrom)
	}
}
