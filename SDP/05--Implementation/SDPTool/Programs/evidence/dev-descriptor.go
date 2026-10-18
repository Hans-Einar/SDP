package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func hash(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func main() {
	dir := os.Args[1]
	b, e := os.ReadFile(filepath.Join(dir, "sdptool"))
	if e != nil {
		panic(e)
	}
	pub, priv, e := ed25519.GenerateKey(rand.Reader)
	if e != nil {
		panic(e)
	}
	d := map[string]any{"schemaVersion": "sdp-release-descriptor/1", "release": "dev-rsp1", "sourceCommit": os.Args[2], "protocol": "sdp-install-command/1", "processProfile": "sdp-five-phase/0.1", "managementProfile": "sdp-project-management/0.2", "capabilities": []string{"sdp-install-command/1"}, "files": []any{}, "retired": []any{}, "upgradesFrom": []any{}, "binaries": []any{map[string]any{"platform": "linux/amd64", "path": "sdptool", "sha256": hash(b), "size": len(b)}}}
	data, _ := json.MarshalIndent(d, "", "  ")
	data = append(data, '\n')
	sig, _ := json.MarshalIndent(map[string]any{"keyId": hash(pub), "signature": ed25519.Sign(priv, data)}, "", "  ")
	for p, v := range map[string][]byte{"release.json": data, "release.json.sig": sig, "test-key.pub": []byte(base64.StdEncoding.EncodeToString(pub) + "\n")} {
		if e := os.WriteFile(filepath.Join(dir, p), v, 0600); e != nil {
			panic(e)
		}
	}
	fmt.Println("Created explicit test-signed development descriptor; private key discarded.")
}
