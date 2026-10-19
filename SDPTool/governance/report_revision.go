package governance

import "encoding/json"

// A tool read is itself followed by a host observation. Permit an agent report
// across only proven tool/message observations, never arbitrary stale revisions.
func observedReportRevision(state snapshot, r Run, caller Caller, q TransitionRequest) bool {
	if caller.controller || q.Action != "report" || q.ExpectedRevision >= r.Revision || r.State != "running" || len(r.Attempts) == 0 {
		return false
	}
	a := r.Attempts[len(r.Attempts)-1]
	if a.ID != caller.attemptID || a.CapabilityEpoch() != caller.epoch || a.Status != "running" || a.Gap {
		return false
	}
	revisions := map[uint64]Result{}
	for _, receipt := range state.Receipts {
		if receipt.Result.Run.ID == r.ID {
			revisions[receipt.Result.Run.Revision] = receipt.Result
		}
	}
	base, ok := revisions[q.ExpectedRevision]
	if !ok || len(base.Run.Attempts) != len(r.Attempts) {
		return false
	}
	normalize := func(v Run) string {
		v.Revision = 0
		v.Attempts = append([]Attempt(nil), v.Attempts...)
		// Only the active attempt's observation counters may change.
		i := len(v.Attempts) - 1
		v.Attempts[i].EventCount = 0
		v.Attempts[i].LastObservedAt = ""
		b, _ := json.Marshal(v)
		return string(b)
	}
	basis := normalize(base.Run)
	count := base.Run.Attempts[len(base.Run.Attempts)-1].EventCount
	for rev := q.ExpectedRevision; rev < r.Revision; {
		rev++
		next, ok := revisions[rev]
		if !ok || next.Operation != "host.observe" || len(next.Run.Attempts) != len(r.Attempts) || normalize(next.Run) != basis {
			return false
		}
		n := next.Run.Attempts[len(next.Run.Attempts)-1].EventCount
		if n != count+1 || n > len(a.Events) {
			return false
		}
		event := a.Events[n-1]
		if event.Kind != "tool.completed" && event.Kind != "agentMessage" {
			return false
		}
		count = n
	}
	return true
}
