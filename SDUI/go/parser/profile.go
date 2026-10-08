package parser

import "fmt"

// ASTFormat selects the source AST envelope format. An empty Document profile
// is invalid; only normalized legacy Instances use the empty profile sentinel.
func ASTFormat(profile string) (string, error) {
	switch profile {
	case "sdui/0.2":
		return "sdui-ast/0.2", nil
	case "sdui/0.3":
		return "sdui-ast/0.3", nil
	default:
		return "", &Diagnostic{Code: "profile", Message: "Unsupported document profile " + profile}
	}
}

// EffectiveProfile validates the profile of a bounded normalized tree, including
// hidden descendants. Empty Instance.Profile preserves the frozen 0.2 encoding.
func EffectiveProfile(root *Instance) (string, error) {
	if root == nil {
		return "", &Diagnostic{Code: "profile", Message: "Missing normalized root"}
	}
	profile := "sdui/0.2"
	if root.Profile == "sdui/0.3" {
		profile = root.Profile
	} else if root.Profile != "" {
		return "", &Diagnostic{Code: "profile", Message: root.Path + ": unsupported instance profile " + root.Profile, Span: root.Span}
	}
	seen := map[*Instance]bool{}
	var visit func(*Instance, int) error
	visit = func(n *Instance, depth int) error {
		if n == nil || seen[n] {
			return &Diagnostic{Code: "instance-tree", Message: "Nil or repeated normalized instance", Span: root.Span}
		}
		seen[n] = true
		if depth > 64 || len(seen) > 8192 {
			return &Diagnostic{Code: "instance-limit", Message: n.Path + ": normalized tree exceeds limits", Span: n.Span}
		}
		if n.Profile != root.Profile {
			return &Diagnostic{Code: "profile", Message: fmt.Sprintf("%s: instance profile %q differs from root %q", n.Path, n.Profile, root.Profile), Span: n.Span}
		}
		for _, region := range n.Regions {
			if err := visit(region.Node, depth+1); err != nil {
				return err
			}
		}
		for _, row := range n.Rows {
			for _, child := range row {
				if err := visit(child, depth+1); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := visit(root, 1); err != nil {
		return "", err
	}
	return profile, nil
}
