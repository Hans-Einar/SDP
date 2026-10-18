package presentation

import (
	"encoding/json"
	"fmt"
	"io"
)

func actions(w io.Writer, b json.RawMessage) error {
	var r struct {
		Schema   string `json:"schema"`
		Producer struct {
			Version string `json:"version"`
		} `json:"producer"`
		Actions []struct {
			ID           string `json:"id"`
			Label        string `json:"label"`
			Availability struct {
				State   string `json:"state"`
				Message string `json:"message"`
			} `json:"availability"`
		} `json:"actions"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "SDPTool actions (%s, %s)\n", safe(r.Producer.Version), safe(r.Schema)); err != nil {
		return err
	}
	for _, a := range r.Actions {
		if _, err := fmt.Fprintf(w, "%s — %s [%s]\n  %s\n", safe(a.ID), safe(a.Label), safe(a.Availability.State), safe(a.Availability.Message)); err != nil {
			return err
		}
	}
	return nil
}
