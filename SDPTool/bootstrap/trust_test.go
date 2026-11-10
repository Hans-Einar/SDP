package bootstrap

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestPublisherAnchor(t *testing.T) {
	if len(trustedKeys) != 2 {
		t.Fatal("expected overlapping publisher keys")
	}
	for id, key := range trustedKeys {
		if len(key) != ed25519.PublicKeySize || id != fmt.Sprintf("%x", sha256.Sum256(key)) {
			t.Fatal("invalid publisher identity")
		}
		if Verify([]byte("tampered"), make([]byte, ed25519.SignatureSize), id, "signed") == nil {
			t.Fatal("accepted invalid production signature")
		}
	}
}

// Public proof generated on Linux. No private key or release descriptor is a fixture.
func TestLinuxPublisherProof(t *testing.T) {
	const oldID = "bfc4838a7fbe7e7e3e84f86b9b435e69e6e503f4caf01a92319147e9ab701618"
	const linuxID = "2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03"
	if _, ok := trustedKeys[oldID]; !ok {
		t.Fatal("previous publisher removed")
	}
	data := []byte("SDP publisher trust-anchor probe: Linux build host, 2026-11-11. Not a release descriptor.\n")
	sig, err := base64.StdEncoding.DecodeString("adARfG5LbAdInSYkB7oC7n3W6xGVeib2sQlUyc6QeHYVRLkxMeAqk/4uqLrA5uMwXQ5UAIiK21cBJ627MExMAA==")
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(data, sig, linuxID, "signed"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name            string
		data, signature []byte
		id, provenance  string
	}{
		{"tampered data", append(append([]byte{}, data...), '!'), sig, linuxID, "signed"},
		{"truncated signature", data, sig[:len(sig)-1], linuxID, "signed"},
		{"wrong known key", data, sig, oldID, "signed"},
		{"unknown key", data, sig, "unknown", "signed"},
		{"unverified provenance", data, sig, linuxID, "unverified"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if Verify(tc.data, tc.signature, tc.id, tc.provenance) == nil {
				t.Fatal("invalid signature accepted")
			}
		})
	}
}
