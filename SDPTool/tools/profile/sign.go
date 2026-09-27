package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/bootstrap"
	"github.com/Hans-Einar/SDP/SDPTool/install"
	"os"
	"strings"
)

// Private key bytes are never returned, logged, embedded or published.
func sign(descriptor, keyPath string) error {
	st, err := os.Lstat(keyPath)
	if err != nil {
		return err
	}
	if !st.Mode().IsRegular() || st.Mode().Perm()&0077 != 0 {
		return fmt.Errorf("publisher key must be a private regular file")
	}
	raw, err := install.Read(keyPath, 1024)
	if err != nil {
		return err
	}
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(raw)))
	if err != nil || len(key) != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid publisher key")
	}
	private := ed25519.PrivateKey(key)
	public := private.Public().(ed25519.PublicKey)
	data, err := install.Read(descriptor, install.MetadataLimit)
	if err != nil {
		return err
	}
	signature := ed25519.Sign(private, data)
	id := install.Hash(public)
	if err = bootstrap.Verify(data, signature, id, "signed"); err != nil {
		return err
	}
	b, err := json.Marshal(map[string]any{"keyId": id, "signature": signature})
	if err != nil {
		return err
	}
	f, err := os.OpenFile(descriptor+".sig", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	ce := f.Close()
	if err != nil {
		return err
	}
	return ce
}
