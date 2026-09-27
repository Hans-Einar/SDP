// profile builds descriptors from shared authored payloads. Production signing
// requires a clean exact source and a locally held, compiled-in publisher key.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/install"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	repo := flag.String("repo", "..", "SDP repository")
	output := flag.String("output", "", "new descriptor")
	commit := flag.String("source-commit", "", "exact source commit")
	release := flag.String("release", "development-gip", "development identity or selected release version")
	key := flag.String("sign-key", "", "private publisher key file for a production release")
	binary := flag.String("binary", "", "packaged executable")
	flag.Parse()
	if e := run(*repo, *output, *commit, *release, *binary, *key); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run(repo, output, commit, release, binary, key string) error {
	if key == "" && !strings.HasPrefix(release, "development-") {
		return fmt.Errorf("production descriptor requires --sign-key")
	}
	if key != "" {
		if !regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`).MatchString(release) || binary == "" {
			return fmt.Errorf("signing requires a final version and binary")
		}
		head, err := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
		if err != nil || strings.TrimSpace(string(head)) != commit {
			return fmt.Errorf("signing requires exact source HEAD")
		}
		status, err := exec.Command("git", "-C", repo, "status", "--porcelain").Output()
		if err != nil || len(status) != 0 {
			return fmt.Errorf("signing requires clean source checkout")
		}
	}
	if err := build(repo, output, commit, release, binary); err != nil {
		return err
	}
	if key != "" {
		return sign(output, key)
	}
	return nil
}
func build(repo, output, commit, release, binary string) error {
	var config struct {
		SchemaVersion     string   `json:"schemaVersion"`
		InventorySource   string   `json:"inventorySource"`
		Reuse             string   `json:"reuse"`
		Protocol          string   `json:"protocol"`
		ProcessProfile    string   `json:"processProfile"`
		ManagementProfile string   `json:"managementProfile"`
		Capabilities      []string `json:"capabilities"`
	}
	b, e := install.Read(filepath.Join(repo, "SDPTool/profiles/five-phase.json"), install.MetadataLimit)
	if e != nil {
		return e
	}
	if e = install.Decode(b, install.MetadataLimit, &config); e != nil {
		return e
	}
	if config.SchemaVersion != "sdp-go-profile/1" || config.Reuse != "files-only" {
		return fmt.Errorf("unsupported authoring profile")
	}
	if e = install.Relative(config.InventorySource); e != nil {
		return e
	}
	b, e = install.Read(filepath.Join(repo, config.InventorySource), install.MetadataLimit)
	if e != nil {
		return e
	}
	var source struct {
		Files []struct {
			Source      string
			Destination string
			Ownership   string
		}
	}
	if e = json.Unmarshal(b, &source); e != nil {
		return e
	}
	d := install.Descriptor{SchemaVersion: install.ReleaseSchema, Release: release, SourceCommit: commit, Protocol: config.Protocol, ProcessProfile: config.ProcessProfile, ManagementProfile: config.ManagementProfile, Capabilities: config.Capabilities, Files: []install.File{}, Retired: []string{}, UpgradesFrom: []string{}, Binaries: []install.Asset{}}
	for _, f := range source.Files {
		if e = install.Relative(f.Source); e != nil {
			return e
		}
		b, e = install.Read(filepath.Join(repo, f.Source), install.FileLimit)
		if e != nil {
			return e
		}
		ownership := f.Ownership
		if ownership == "project" {
			ownership = "initialize-if-missing"
		}
		h := install.Hash(b)
		d.Files = append(d.Files, install.File{Path: f.Destination, Type: "file", Ownership: ownership, SHA256: &h, Content: b})
	}
	if binary != "" {
		b, e = install.Read(binary, install.PayloadLimit)
		if e != nil {
			return e
		}
		d.Binaries = append(d.Binaries, install.Asset{Platform: runtime.GOOS + "/" + runtime.GOARCH, Path: filepath.Base(binary), SHA256: install.Hash(b), Size: int64(len(b))})
	}
	if e = install.ValidateDescriptor(d); e != nil {
		return e
	}
	b, e = install.Canonical(d)
	if e != nil {
		return e
	}
	if len(b) > install.MetadataLimit {
		return fmt.Errorf("descriptor too large")
	}
	f, e := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	_, e = f.Write(b)
	ce := f.Close()
	if e != nil {
		return e
	}
	return ce
}
