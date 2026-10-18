package sdptool

import (
	"fmt"
	"path/filepath"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

type SourceLink struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Snapshot string `json:"snapshot"`
	Revision string `json:"revision"`
}

// Read exactly the checked input set and compare per-file hashes; Preview also
// rechecks the entire graph immediately before publication.
func bindSourceLinks(v *viewpoint.Views, entry string) ([]SourceLink, map[string][]byte, error) {
	links := []SourceLink{}
	copies := map[string][]byte{}
	v.SourceBindings = map[string]string{}
	for i, s := range v.Sources {
		id := fmt.Sprintf("s%d", i)
		snapshot := "source-inputs/" + id + ".design"
		path := filepath.Join(filepath.Dir(entry), filepath.FromSlash(s.Path))
		data, e := readSource(path)
		if e != nil {
			return nil, nil, e
		}
		if documents.Hash(data) != s.SHA256 {
			return nil, nil, failure("stale", fmt.Errorf("source changed before snapshot: %s", s.Path))
		}
		links = append(links, SourceLink{id, path, snapshot, s.SHA256})
		copies[snapshot] = data
		v.SourceBindings[s.Path] = id
		if len(v.Sources) == 1 && v.Profile == "design-core/0.5" {
			v.SourceBindings[""] = id
		}
	}
	return links, copies, nil
}
