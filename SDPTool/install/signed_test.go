package install

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/bootstrap"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPackagedSignedInstall(t *testing.T) {
	binary := os.Getenv("SDP_TEST_BINARY")
	if binary == "" {
		t.Skip("set SDP_TEST_BINARY to packaged sdptool for real-child verification")
	}
	b, e := Read(binary, PayloadLimit)
	if e != nil {
		t.Fatal(e)
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "sdptool")
	if e = os.WriteFile(bin, b, 0700); e != nil {
		t.Fatal(e)
	}
	d := fixture()
	d.Binaries = []Asset{{runtime.GOOS + "/" + runtime.GOARCH, "sdptool", Hash(b), int64(len(b))}}
	data, _ := Canonical(d)
	public, private, _ := ed25519.GenerateKey(nil)
	desc := filepath.Join(dir, "release.json")
	os.WriteFile(desc, data, 0600)
	sig, _ := json.Marshal(map[string]any{"keyId": Hash(public), "signature": ed25519.Sign(private, data)})
	os.WriteFile(desc+".sig", sig, 0600)
	key := filepath.Join(dir, "test.pub")
	os.WriteFile(key, []byte(base64.StdEncoding.EncodeToString(public)), 0600)
	cache := filepath.Join(t.TempDir(), "cache")
	t.Setenv("SDP_CACHE_DIR", cache)
	verified, e := bootstrap.Bootstrap(context.Background(), bootstrap.Config{Descriptor: desc, TestKey: key, CacheDir: cache})
	if e != nil {
		t.Fatal(e)
	}
	root := t.TempDir()
	plan := filepath.Join(t.TempDir(), "plan.json")
	cmd := exec.Command(verified, root, "install", "--release", desc, "--test-key", key, "--plan-output", plan, "--json")
	out, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(out))
	}
	cmd = exec.Command(verified, root, "install", "--apply", plan, "--json")
	out, e = cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(out))
	}
	b, _ = Read(filepath.Join(root, ReceiptPath), MetadataLimit)
	var receipt Receipt
	if e = Decode(b, MetadataLimit, &receipt); e != nil || receipt.Provenance != "test-signed" {
		t.Fatal(receipt, e)
	}
}
