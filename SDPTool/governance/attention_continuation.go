package governance

import "slices"

const attentionContinuationStateSchema = "sdp-project-governance-state/0.7"

type AttentionReplyRef struct {
	QuestionID       string `json:"questionId"`
	ReplyOperationID string `json:"replyOperationId"`
	Basis            string `json:"basis"`
}

// Later store profiles retain attention capabilities even without questions.
func attentionProfile(schema string) bool {
	return schema == attentionStateSchema || schema == attentionContinuationStateSchema || reportProfile(schema)
}

// Verify the original decision state, before this particular launch input made
// the run ready. Only that exact prepared input may be excluded during claim.
func checkAttentionContinuation(r Run, v OwnerInput, claiming bool) error {
	if len(v.AttentionReplies) == 0 || len(v.AttentionReplies) > 32 || v.Kind != "launch" || len(r.Attempts) == 0 {
		return fail("invalid_request", "bounded attention launch references required")
	}
	a := r.Attempts[len(r.Attempts)-1]
	if a.Status != "completed" || a.Gap {
		return fail("reconciliation_required", "continuation needs a clean completed attempt")
	}
	if hasBlockingAttention(r) {
		return fail("attention_required", "blocking questions remain open")
	}
	if claiming {
		if r.State != "ready" || len(r.Inputs) == 0 || r.Inputs[len(r.Inputs)-1].ID != v.ID {
			return fail("conflict", "prepared continuation is no longer current")
		}
		r.Inputs = r.Inputs[:len(r.Inputs)-1]
		r.State = "running"
	} else if r.State != "running" {
		return fail("conflict", "attention continuation requires a running assignment")
	}
	var original *OwnerInput
	for i := range r.Inputs {
		v := &r.Inputs[i]
		if v.Kind == "launch" && len(v.AttentionReplies) == 0 {
			original = v
			break
		}
	}
	if original == nil || original.Text != v.Text || !slices.Equal(original.Skills, v.Skills) {
		return fail("conflict", "continuation must preserve original owner input and skills")
	}
	seen := map[string]bool{}
	for _, ref := range v.AttentionReplies {
		if !uuidPattern.MatchString(ref.QuestionID) || !uuidPattern.MatchString(ref.ReplyOperationID) || !digestPattern.MatchString(ref.Basis) || seen[ref.QuestionID] {
			return fail("invalid_request", "invalid or duplicate reply reference")
		}
		seen[ref.QuestionID] = true
		found := false
		for _, q := range r.Attention {
			if q.ID == ref.QuestionID && q.Reply != nil && q.Reply.OperationID == ref.ReplyOperationID && q.Basis == ref.Basis && q.Basis == attentionBasis(r) {
				found = true
			}
		}
		if !found {
			return fail("stale_revision", "attention reply no longer applies")
		}
		for _, input := range r.Inputs {
			for _, used := range input.AttentionReplies {
				if used.QuestionID == ref.QuestionID {
					return fail("conflict", "reply already selected for dispatch")
				}
			}
		}
	}
	return nil
}
func attentionDispatch(r Run, questionID string) string {
	for _, v := range r.Inputs {
		for _, ref := range v.AttentionReplies {
			if ref.QuestionID == questionID {
				return "dispatch-" + v.State
			}
		}
	}
	return "not-tracked"
}
func validateAttentionContinuations(state snapshot) error {
	count := 0
	for _, r := range state.Runs {
		used := map[string]bool{}
		for _, v := range r.Inputs {
			if len(v.AttentionReplies) == 0 {
				continue
			}
			count++
			if (state.Schema != attentionContinuationStateSchema && !reportProfile(state.Schema)) || v.Kind != "launch" || len(v.AttentionReplies) > 32 || !slices.Contains([]string{"prepared", "sending", "uncertain", "acknowledged"}, v.State) {
				return fail("reconciliation_required", "invalid attention dispatch input")
			}
			for _, ref := range v.AttentionReplies {
				found := false
				for _, q := range r.Attention {
					if q.ID == ref.QuestionID && q.Reply != nil && q.Reply.OperationID == ref.ReplyOperationID && q.Basis == ref.Basis {
						found = true
					}
				}
				if !found || used[ref.QuestionID] {
					return fail("reconciliation_required", "invalid or repeated attention dispatch reference")
				}
				used[ref.QuestionID] = true
			}
		}
	}
	if state.Schema == attentionContinuationStateSchema && count == 0 {
		return fail("reconciliation_required", "continuation profile has no dispatch references")
	}
	return nil
}

func reportProfile(schema string) bool {
	return schema == reportAcknowledgementStateSchema || schema == leaderStateSchema
}
