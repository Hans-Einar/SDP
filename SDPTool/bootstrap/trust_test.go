package bootstrap

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestPublisherAnchor(t *testing.T) {
	if len(trustedKeys) != 1 {
		t.Fatal("expected selected publisher")
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
