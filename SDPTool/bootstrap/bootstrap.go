// Package bootstrap verifies distribution inputs and obtains a compatible local
// executable. It contains no project discovery, installation or migration policy.
package bootstrap

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const Protocol = "sdp-install-command/1"
const maxDescriptor = 16 << 20
const maxBinary = 512 << 20

// No production trust key or release is selected yet. Production keys must be
// compiled into a reviewed distribution; a descriptor cannot add a trusted key.
var trustedKeys = map[string]ed25519.PublicKey{}

type Config struct {
	Descriptor string
	TestKey    string
	Offline    bool
	CacheDir   string
}
type Verified struct {
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	Bytes      []byte `json:"bytes"`
	Signature  []byte `json:"signature"`
	KeyID      string `json:"keyId"`
	Provenance string `json:"provenance"`
}
type signature struct {
	KeyID     string `json:"keyId"`
	Signature []byte `json:"signature"`
}
type asset struct {
	Platform string `json:"platform"`
	Path     string `json:"path"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
}

func hash(b []byte) string { v := sha256.Sum256(b); return hex.EncodeToString(v[:]) }
func Verify(data, sig []byte, keyID, provenance string) error {
	key := trustedKeys[keyID]
	if provenance == "test-signed" {
		if !strings.HasPrefix(keyID, "test:") {
			return fmt.Errorf("missing explicit test trust")
		}
		b, e := base64.StdEncoding.DecodeString(strings.TrimPrefix(keyID, "test:"))
		if e != nil {
			return e
		}
		key = ed25519.PublicKey(b)
	} else if provenance != "signed" {
		return fmt.Errorf("unsupported signing provenance")
	}
	if len(key) != ed25519.PublicKeySize || len(sig) != ed25519.SignatureSize || !ed25519.Verify(key, data, sig) {
		return fmt.Errorf("unknown key or invalid descriptor signature")
	}
	return nil
}
func safe(p string) error {
	p, e := filepath.Abs(p)
	if e != nil {
		return e
	}
	for q := p; ; q = filepath.Dir(q) {
		s, e := os.Lstat(q)
		if e == nil && s.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in distribution path: %s", q)
		}
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if filepath.Dir(q) == q {
			break
		}
	}
	return nil
}
func privateCache(c Config) (string, error) {
	p := c.CacheDir
	if p == "" {
		v, e := os.UserCacheDir()
		if e != nil {
			return "", e
		}
		p = filepath.Join(v, "sdptool", "distribution")
	}
	p, e := filepath.Abs(p)
	if e != nil {
		return "", e
	}
	if e = safe(p); e != nil {
		return "", e
	}
	if e = os.MkdirAll(p, 0700); e != nil {
		return "", e
	}
	s, e := os.Stat(p)
	if e != nil {
		return "", e
	}
	if s.Mode().Perm()&0077 != 0 {
		return "", fmt.Errorf("distribution cache must be private: %s", p)
	}
	return p, nil
}
func read(p string, limit int64) ([]byte, error) {
	if e := safe(p); e != nil {
		return nil, e
	}
	before, e := os.Lstat(p)
	if e != nil {
		return nil, e
	}
	if !before.Mode().IsRegular() || before.Size() > limit {
		return nil, fmt.Errorf("not a bounded regular file: %s", p)
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	s, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !s.Mode().IsRegular() || s.Size() > limit {
		return nil, fmt.Errorf("not a bounded regular file: %s", p)
	}
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("distribution input exceeds limit")
	}
	return b, e
}
func fetch(ctx context.Context, source string, limit int64) ([]byte, error) {
	if !strings.HasPrefix(source, "https://") {
		if strings.Contains(source, "://") {
			return nil, fmt.Errorf("HTTPS or explicit local source required")
		}
		return read(source, limit)
	}
	u, e := url.Parse(source)
	if e != nil || u.User != nil || u.Fragment != "" {
		return nil, fmt.Errorf("invalid distribution URL")
	}
	req, e := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if e != nil {
		return nil, e
	}
	client := http.Client{Timeout: 30 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) > 5 || req.URL.Scheme != "https" || req.URL.User != nil {
			return fmt.Errorf("unsafe distribution redirect")
		}
		return nil
	}}
	resp, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("distribution HTTP %d", resp.StatusCode)
	}
	b, e := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("distribution input exceeds limit")
	}
	return b, e
}
func store(dir, name string, b []byte, mode os.FileMode) error {
	p := filepath.Join(dir, name)
	if e := safe(p); e != nil {
		return e
	}
	if old, e := read(p, int64(len(b))); e == nil {
		if !bytes.Equal(old, b) {
			return fmt.Errorf("corrupt/unequal cache entry %s", name)
		}
		return nil
	} else if !os.IsNotExist(e) {
		return e
	}
	f, e := os.CreateTemp(dir, ".pending-")
	if e != nil {
		return e
	}
	temp := f.Name()
	defer os.Remove(temp)
	if e = f.Chmod(mode); e == nil {
		_, e = f.Write(b)
	}
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e == nil {
		e = ce
	}
	if e != nil {
		return e
	}
	// Link publishes create-new; concurrent identical writers can share bytes.
	if e = os.Link(temp, p); e != nil {
		old, re := read(p, int64(len(b)))
		if re != nil || !bytes.Equal(old, b) {
			return fmt.Errorf("cache publication conflict: %w", e)
		}
	}
	return nil
}

// StrictJSON rejects duplicate/case-colliding keys and trailing documents before
// decoding. Shape validation remains owned by the consumer of the descriptor.
func StrictJSON(b []byte, v any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var walk func(int) error
	walk = func(depth int) error {
		if depth > 64 {
			return fmt.Errorf("JSON depth exceeded")
		}
		t, e := d.Token()
		if e != nil {
			return e
		}
		if delim, ok := t.(json.Delim); ok {
			if delim == '{' {
				seen := map[string]bool{}
				for d.More() {
					k, e := d.Token()
					if e != nil {
						return e
					}
					s, ok := k.(string)
					if !ok {
						return fmt.Errorf("string key required")
					}
					s = strings.ToLower(s)
					if seen[s] {
						return fmt.Errorf("duplicate JSON key")
					}
					seen[s] = true
					if e = walk(depth + 1); e != nil {
						return e
					}
				}
			} else if delim == '[' {
				for d.More() {
					if e = walk(depth + 1); e != nil {
						return e
					}
				}
			} else {
				return fmt.Errorf("unexpected JSON delimiter")
			}
			_, e = d.Token()
			return e
		}
		return nil
	}
	if e := walk(0); e != nil {
		return e
	}
	if _, e := d.Token(); e != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	return json.Unmarshal(b, v)
}
func Resolve(ctx context.Context, c Config) (Verified, error) {
	var v Verified
	if c.Descriptor == "" {
		return v, fmt.Errorf("no release selected; configure SDP_RELEASE or an explicit descriptor")
	}
	dir, e := privateCache(c)
	if e != nil {
		return v, e
	}
	selector := c.Descriptor
	if !strings.HasPrefix(selector, "https://") {
		selector, e = filepath.Abs(selector)
		if e != nil {
			return v, e
		}
	}
	alias := "selection-" + hash([]byte(selector)) + ".json"
	if c.Offline {
		b, e := read(filepath.Join(dir, alias), 32<<20)
		if e != nil {
			return v, fmt.Errorf("offline descriptor unavailable: %w", e)
		}
		if e = StrictJSON(b, &v); e != nil {
			return v, e
		}
	} else {
		b, e := fetch(ctx, selector, maxDescriptor)
		if e != nil {
			return v, e
		}
		sigSource := selector + ".sig"
		if strings.HasPrefix(selector, "https://") {
			u, _ := url.Parse(selector)
			u.Path += ".sig"
			sigSource = u.String()
		}
		sb, e := fetch(ctx, sigSource, 16384)
		if e != nil {
			return v, e
		}
		var sig signature
		if e = StrictJSON(sb, &sig); e != nil {
			return v, e
		}
		v = Verified{Bytes: b, SHA256: hash(b), Signature: sig.Signature, KeyID: sig.KeyID, Provenance: "signed"}
		if c.TestKey != "" {
			b, e := read(c.TestKey, 1024)
			if e != nil {
				return v, e
			}
			key, e := base64.StdEncoding.DecodeString(strings.TrimSpace(string(b)))
			if e != nil || len(key) != ed25519.PublicKeySize || hash(key) != sig.KeyID {
				return v, fmt.Errorf("test key does not match descriptor signer")
			}
			v.KeyID = "test:" + base64.StdEncoding.EncodeToString(key)
			v.Provenance = "test-signed"
		}
	}
	if v.Provenance == "test-signed" {
		if c.TestKey == "" {
			return v, fmt.Errorf("cached test distribution requires explicit SDP_TEST_KEY or --test-key")
		}
		b, e := read(c.TestKey, 1024)
		if e != nil {
			return v, e
		}
		if "test:"+strings.TrimSpace(string(b)) != v.KeyID {
			return v, fmt.Errorf("cached signer differs from selected test key")
		}
	}
	if hash(v.Bytes) != v.SHA256 {
		return v, fmt.Errorf("descriptor cache hash mismatch")
	}
	if e = Verify(v.Bytes, v.Signature, v.KeyID, v.Provenance); e != nil {
		return v, e
	}
	if e = store(dir, v.SHA256+".json", v.Bytes, 0600); e != nil {
		return v, e
	}
	v.Path = filepath.Join(dir, v.SHA256+".json")
	b, _ := json.Marshal(v)
	b = append(b, '\n')
	if !c.Offline {
		if e = store(dir, alias, b, 0600); e != nil {
			return v, fmt.Errorf("immutable selector changed; use an exact versioned selector: %w", e)
		}
	}
	return v, nil
}
func Bootstrap(ctx context.Context, c Config) (string, error) {
	v, e := Resolve(ctx, c)
	if e != nil {
		return "", e
	}
	var d struct {
		Schema   string  `json:"schemaVersion"`
		Protocol string  `json:"protocol"`
		Binaries []asset `json:"binaries"`
	}
	if e = StrictJSON(v.Bytes, &d); e != nil {
		return "", e
	}
	if d.Schema != "sdp-release-descriptor/1" || d.Protocol != Protocol {
		return "", fmt.Errorf("unsupported installation protocol")
	}
	platform := runtime.GOOS + "/" + runtime.GOARCH
	var selected *asset
	for i := range d.Binaries {
		if d.Binaries[i].Platform == platform {
			if selected != nil {
				return "", fmt.Errorf("duplicate platform binary")
			}
			selected = &d.Binaries[i]
		}
	}
	if selected == nil {
		return "", fmt.Errorf("no binary for %s", platform)
	}
	a := *selected
	if a.Size < 1 || a.Size > maxBinary || len(a.SHA256) != 64 {
		return "", fmt.Errorf("invalid executable bounds")
	}
	if _, e = hex.DecodeString(a.SHA256); e != nil {
		return "", e
	}
	if a.Path == "" || filepath.ToSlash(filepath.Clean(a.Path)) != a.Path || strings.HasPrefix(a.Path, "/") || a.Path == ".." || strings.HasPrefix(a.Path, "../") || strings.ContainsAny(a.Path, "\\:") {
		return "", fmt.Errorf("unsafe binary path")
	}
	dir, e := privateCache(c)
	if e != nil {
		return "", e
	}
	name := a.SHA256 + "-sdptool"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	p := filepath.Join(dir, name)
	b, e := read(p, a.Size)
	if os.IsNotExist(e) && !c.Offline {
		source := filepath.Join(filepath.Dir(c.Descriptor), a.Path)
		if strings.HasPrefix(c.Descriptor, "https://") {
			u, _ := url.Parse(c.Descriptor)
			ref := &url.URL{Path: a.Path}
			source = u.ResolveReference(ref).String()
		}
		b, e = fetch(ctx, source, a.Size)
		if e == nil && int64(len(b)) == a.Size && hash(b) == a.SHA256 {
			e = store(dir, name, b, 0700)
		} else if e == nil {
			e = fmt.Errorf("downloaded executable digest/size mismatch")
		}
	}
	if e != nil {
		return "", e
	}
	if int64(len(b)) != a.Size || hash(b) != a.SHA256 {
		return "", fmt.Errorf("cached executable digest/size mismatch")
	}
	probe, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(probe, p, "--version")
	cmd.Dir = dir
	var stdout limitedBuffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard
	if e = cmd.Run(); e != nil {
		return "", fmt.Errorf("executable protocol probe: %w", e)
	}
	var version struct {
		Protocol string `json:"installationProtocol"`
	}
	if e = StrictJSON(stdout.Bytes(), &version); e != nil || version.Protocol != Protocol {
		return "", fmt.Errorf("binary lacks %s", Protocol)
	}
	return p, nil
}

type limitedBuffer struct{ bytes.Buffer }

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 1<<20 {
		return 0, fmt.Errorf("protocol output too large")
	}
	return b.Buffer.Write(p)
}
