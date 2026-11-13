package sdptool

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// nativeHost materializes only bytes already authenticated inside this engine.
func nativeHost() (string, error) {
	if len(bundledNativeHost) == 0 {
		return "", fmt.Errorf("this source build has no native prototype payload; install the bundled SDP release through gh extension install/upgrade sdp")
	}
	z, err := gzip.NewReader(bytes.NewReader(bundledNativeHost))
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(io.LimitReader(z, (128<<20)+1))
	closeErr := z.Close()
	if err != nil {
		return "", err
	}
	if closeErr != nil {
		return "", closeErr
	}
	if len(data) > 128<<20 {
		return "", fmt.Errorf("native payload exceeds size limit")
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(cache) {
		return "", fmt.Errorf("helper cache must be absolute")
	}
	base := cache
	if err = os.MkdirAll(base, 0700); err != nil {
		return "", err
	}
	for _, part := range []string{"sdptool", "helpers", fmt.Sprintf("%x", sha256.Sum256(data))} {
		base = filepath.Join(base, part)
		if err = os.Mkdir(base, 0700); err != nil && !os.IsExist(err) {
			return "", err
		}
		st, e := os.Lstat(base)
		if e != nil || !st.IsDir() || st.Mode()&os.ModeSymlink != 0 || st.Mode().Perm()&0077 != 0 {
			return "", fmt.Errorf("unsafe helper directory: %s", base)
		}
	}
	path := filepath.Join(base, "sdui-fyne")
	verify := func() error {
		st, e := os.Lstat(path)
		if e != nil {
			return e
		}
		if !st.Mode().IsRegular() || st.Mode().Perm() != 0700 {
			return fmt.Errorf("unsafe cached native helper: %s", path)
		}
		if st.Size() != int64(len(data)) {
			return fmt.Errorf("cached native helper size mismatch: %s", path)
		}
		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		if sha256.Sum256(b) != sha256.Sum256(data) {
			return fmt.Errorf("cached native helper hash mismatch: %s", path)
		}
		return nil
	}
	if _, err = os.Lstat(path); err == nil {
		return path, verify()
	} else if !os.IsNotExist(err) {
		return "", err
	}
	f, err := os.CreateTemp(base, ".native-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	_, err = f.Write(data)
	if err == nil {
		err = f.Chmod(0700)
	}
	if err == nil {
		err = f.Sync()
	}
	e := f.Close()
	if err == nil {
		err = e
	}
	if err != nil {
		return "", err
	}
	if err = os.Link(f.Name(), path); err != nil && !os.IsExist(err) {
		return "", err
	}
	return path, verify()
}
func runNativeHost(ctx context.Context, args []string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := nativeHost()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("native SDUI prototype: %w", err)
	}
	return nil
}
