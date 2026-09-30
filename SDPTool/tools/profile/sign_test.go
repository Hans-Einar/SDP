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
	want := []string{
		"edc0c72101a437c6e12c40a081ef59ae41cf0db32bcaedb73824ec48495aaee5", // published 0.2.0
		"66590e8e967ede6b36d8fa45cdbee1cd80f69505cca698b0e4a9bd960842735a", // published 0.2.1
		"767527e0d7f54866bab95f4642ffffb2a344024c1805f44b3d5ea9c024829125", // published 1.0.0
	}
	if len(descriptor.UpgradesFrom) != len(want) {
		t.Fatalf("unexpected predecessor inventory: %v", descriptor.UpgradesFrom)
	}
	for i, digest := range want {
		if descriptor.UpgradesFrom[i] != digest {
			t.Fatalf("published predecessor missing: %v", descriptor.UpgradesFrom)
		}
	}
}
