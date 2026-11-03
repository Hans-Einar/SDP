package governancehost

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
)

// SessionLeaderTools verifies the entire launcher-selected MCP inventory. An
// extra server, partial page or extra/missing tool refuses leader submission.
func (c *Client) SessionLeaderTools(ctx context.Context, thread string, expected map[string][]string) error {
	var out struct {
		Data []struct {
			Name  string                     `json:"name"`
			Tools map[string]json.RawMessage `json:"tools"`
		} `json:"data"`
		NextCursor *string `json:"nextCursor"`
	}
	if err := c.call(ctx, "mcpServerStatus/list", map[string]any{"limit": 64, "threadId": thread}, &out); err != nil {
		return err
	}
	if out.NextCursor != nil || len(out.Data) != len(expected) {
		return fmt.Errorf("incompatible_profile: exact Session leader MCP inventory unavailable")
	}
	seen := map[string]bool{}
	for _, server := range out.Data {
		want, ok := expected[server.Name]
		if !ok || seen[server.Name] {
			return fmt.Errorf("unauthorized: unexpected Session leader MCP server")
		}
		seen[server.Name] = true
		names := []string{}
		for name := range server.Tools {
			names = append(names, name)
		}
		sort.Strings(names)
		if !slices.Equal(names, want) {
			return fmt.Errorf("incompatible_profile: Session leader tools differ for %s", server.Name)
		}
	}
	return nil
}
