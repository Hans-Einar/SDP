package runtime

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func validPlain(s string, limit int, empty bool) bool {
	if len(s) > limit || !utf8.ValidString(s) || (!empty && s == "") {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) || r == '\u2028' || r == '\u2029' {
			return false
		}
	}
	return true
}
func collectionFault(path string, err error) error {
	if f, ok := err.(*Fault); ok {
		return fault(f.Code, path+": "+f.Message)
	}
	return err
}
func validateData(kind string, d CollectionData, rootLoaded, loader bool) error {
	if len(d.Items) > 4096 {
		return fault("collection-limit", "Collection exceeds 4096 items")
	}
	if !rootLoaded && (len(d.Items) != 0 || !loader) {
		return fault("collection-data", "Unloaded root requires a loader and no items")
	}
	items := make(map[ItemID]CollectionItem, len(d.Items))
	bad := func(id ItemID, msg string) error {
		return fault("collection-data", fmt.Sprintf("Item %q: %s", id, msg))
	}
	for _, x := range d.Items {
		if !validPlain(string(x.ID), 1024, false) || !validPlain(string(x.Parent), 1024, true) {
			return bad(x.ID, "invalid identity")
		}
		if _, ok := items[x.ID]; ok {
			return bad(x.ID, "duplicate identity")
		}
		if !validPlain(x.Label, 32768, x.Kind == Separator) {
			return bad(x.ID, "invalid single-line label")
		}
		if x.Kind != Row && x.Kind != Group && x.Kind != Separator {
			return bad(x.ID, "unknown kind")
		}
		if x.Kind == Separator && (x.Label != "" || x.HasChildren || !x.ChildrenLoaded) {
			return bad(x.ID, "separator must be an empty leaf")
		}
		if !x.HasChildren && !x.ChildrenLoaded {
			return bad(x.ID, "leaf must be loaded")
		}
		if x.HasChildren && !x.ChildrenLoaded && !loader {
			return bad(x.ID, "unloaded branch requires loader")
		}
		if kind == "list" && (x.Parent != "" || x.HasChildren) {
			return bad(x.ID, "list items must be flat leaves")
		}
		items[x.ID] = x
	}
	for _, x := range d.Items {
		p := x.Parent
		seen := map[ItemID]bool{x.ID: true}
		depth := 1
		for p != "" {
			if seen[p] {
				return bad(x.ID, "parent cycle")
			}
			seen[p] = true
			parent, ok := items[p]
			if !ok {
				return bad(x.ID, "missing parent")
			}
			if !parent.HasChildren || !parent.ChildrenLoaded || parent.Kind == Separator {
				return bad(x.ID, "parent cannot contain children")
			}
			depth++
			if depth > 64 {
				return bad(x.ID, "parent depth exceeds 64")
			}
			p = parent.Parent
		}
	}
	return nil
}
func boundedError(err error) string {
	s := strings.ToValidUTF8(err.Error(), "�")
	if len(s) > 4096 {
		s = s[:4096]
		for !utf8.ValidString(s) {
			s = s[:len(s)-1]
		}
	}
	return s
}
func (c *collection) resetStatus(previous *collection) {
	c.Status = map[ItemID]LoadStatus{}
	phase := Unloaded
	if c.RootLoaded {
		phase = Loaded
	}
	c.Status[""] = LoadStatus{Phase: phase}
	for _, x := range c.Data.Items {
		if !x.HasChildren {
			continue
		}
		phase = Unloaded
		if x.ChildrenLoaded {
			phase = Loaded
		}
		c.Status[x.ID] = LoadStatus{Phase: phase}
		if previous != nil {
			old, ok := previous.item(x.ID)
			if ok && old.Kind == x.Kind && old.HasChildren && !x.ChildrenLoaded {
				status := previous.Status[x.ID]
				if status.Phase == Loading {
					status = LoadStatus{Phase: Canceled}
				}
				if status.Phase == Canceled || status.Phase == LoadError {
					c.Status[x.ID] = status
				}
			}
		}
	}
}
