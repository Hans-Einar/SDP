package governance

import "testing"

func TestCurrentDispatchedReplyBoundaries(t *testing.T) {
	q := AttentionQuestion{ID: "question", InputRevision: "same-source", Basis: "predecessor", Reply: &AttentionReply{OperationID: "reply"}}
	r := Run{State: "running", InputRevision: "same-source", Attempts: []Attempt{{ID: "current-attempt", Epoch: "epoch", Status: "running"}}, Inputs: []OwnerInput{{Kind: "launch", State: "acknowledged", AttemptID: "current-attempt", Epoch: "epoch", AttentionReplies: []AttentionReplyRef{{QuestionID: q.ID, ReplyOperationID: q.Reply.OperationID, Basis: q.Basis}}}}}
	if !currentDispatchedReply(r, q) {
		t.Fatal("selected reply unavailable to current attempt")
	}
	for _, change := range []string{"prepared", "uncertain", "steer", "old-attempt", "epoch", "gap", "failed", "interrupted", "terminal", "new-source", "wrong-reply", "withdrawn", "new-input"} {
		t.Run(change, func(t *testing.T) {
			copy := r
			copy.Inputs = append([]OwnerInput{}, r.Inputs...)
			copy.Attempts = append([]Attempt{}, r.Attempts...)
			question := q
			switch change {
			case "prepared", "uncertain":
				copy.Inputs[0].State = change
			case "steer":
				copy.Inputs[0].Kind = "steer"
			case "old-attempt":
				copy.Attempts[0].ID = "different"
			case "epoch":
				copy.Inputs[0].Epoch = "different"
			case "gap":
				copy.Attempts[0].Gap = true
			case "failed", "interrupted":
				copy.Attempts[0].Status = change
			case "terminal":
				copy.State = "verifying"
			case "new-source":
				copy.InputRevision = "different"
			case "wrong-reply":
				question.Reply = &AttentionReply{OperationID: "another"}
			case "withdrawn":
				question.Withdrawal = &AttentionWithdrawal{Reason: "obsolete"}
			case "new-input":
				copy.Inputs = append(copy.Inputs, OwnerInput{Kind: "steer"})
			}
			if currentDispatchedReply(copy, question) {
				t.Fatal("obsolete or uncertain reply projected current")
			}
		})
	}
	for _, status := range []string{"idle", "running", "completed"} {
		r.Attempts[0].Status = status
		if !currentDispatchedReply(r, q) {
			t.Fatal("lost current reply", status)
		}
	}
}
