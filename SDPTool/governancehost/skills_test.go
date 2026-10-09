package governancehost

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSkillInventoryDistinguishesAvailabilityFromUse(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Skills", "sdp", "SKILL.md")
	os.MkdirAll(filepath.Dir(path), 0700)
	os.WriteFile(path, []byte("---\nname: sdp\n---\nfixture"), 0600)
	outside := filepath.Join(t.TempDir(), "SKILL.md")
	os.WriteFile(outside, []byte("never retain external body"), 0600)
	left, right := net.Pipe()
	client := newClient(left, left)
	defer client.Close()
	defer right.Close()
	captured := make(chan json.RawMessage, 1)
	go func() {
		scanner := bufio.NewScanner(right)
		if !scanner.Scan() {
			return
		}
		var q frame
		if json.Unmarshal(scanner.Bytes(), &q) != nil {
			return
		}
		captured <- q.Params
		result := map[string]any{"data": []any{map[string]any{"cwd": root, "errors": []any{}, "skills": []any{map[string]any{"name": "sdp", "path": path, "enabled": true}, map[string]any{"name": "external", "path": outside, "enabled": false}}}}}
		b, _ := json.Marshal(map[string]any{"id": q.ID, "result": result})
		right.Write(append(b, '\n'))
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	inventory, err := client.InspectSkills(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if inventory.Status != "observed" || inventory.ExternalCount != 1 || inventory.UseEvidence != "not-observed" || len(inventory.Project) != 1 || inventory.Project[0].Path != "Skills/sdp/SKILL.md" || len(inventory.Project[0].Digest) != 64 {
		t.Fatal(inventory)
	}
	params := <-captured
	var q struct {
		Cwds        []string `json:"cwds"`
		ForceReload bool     `json:"forceReload"`
	}
	json.Unmarshal(params, &q)
	if !q.ForceReload || len(q.Cwds) != 1 || q.Cwds[0] != root {
		t.Fatal(q)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.Marshal([]map[string]any{{"definition": "SkillsListParams", "value": json.RawMessage(params)}})
		if err = os.WriteFile(filepath.Join(dir, "skills-rpc-cases.json"), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSkillInventoryRefusesAmbiguousAndUnboundedMetadata(t *testing.T) {
	for _, scenario := range []string{"cwd", "relative", "duplicate", "limit", "missing-enabled", "missing-skills"} {
		t.Run(scenario, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, "SKILL.md")
			os.WriteFile(path, []byte("fixture"), 0600)
			skill := map[string]any{"name": "sdp", "path": path, "enabled": true}
			skills := []any{skill}
			cwd := root
			switch scenario {
			case "missing-enabled":
				delete(skill, "enabled")
			case "missing-skills":
				skills = nil
			case "cwd":
				cwd = "/another-project"
			case "relative":
				skill["path"] = "SKILL.md"
			case "duplicate":
				skills = append(skills, skill)
			case "limit":
				skills = make([]any, 4097)
			}
			left, right := net.Pipe()
			client := newClient(left, left)
			defer client.Close()
			defer right.Close()
			go func() {
				scanner := bufio.NewScanner(right)
				if !scanner.Scan() {
					return
				}
				var q frame
				json.Unmarshal(scanner.Bytes(), &q)
				b, _ := json.Marshal(map[string]any{"id": q.ID, "result": map[string]any{"data": []any{map[string]any{"cwd": cwd, "errors": []any{}, "skills": skills}}}})
				right.Write(append(b, '\n'))
			}()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if _, err := client.InspectSkills(ctx, root); err == nil {
				t.Fatal("untrustworthy inventory accepted")
			}
		})
	}
}
