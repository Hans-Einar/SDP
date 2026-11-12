package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const ApplicationSchema = "sdp-activities/1"

type InputKind string

const (
	NoInput       InputKind = "none"
	ActivityInput InputKind = "activity"
	RequestInput  InputKind = "request"
	RecoveryInput InputKind = "recovery"
	ViewInput     InputKind = "view"
	TreeInput     InputKind = "tree"
)

type Input struct {
	ActivityID            string
	OperationID           string
	URI, Revision, Output string
	Request               json.RawMessage
}

type OperationContext struct {
	Inspector ExecutionInspector
	View      func(Input) (any, error)
	Tree      func(Input) (any, error)
}

type Response struct {
	Schema    string `json:"schema"`
	Operation string `json:"operation"`
	Data      any    `json:"data"`
}

// Operation binds public metadata and a typed domain call once for both adapters.
type Operation struct {
	Name, ActionID, Label, Description string
	Input                              InputKind
	Write                              bool
	run                                func(*Service, Caller, Input, OperationContext) (any, error)
}

func (o Operation) Execute(ctx context.Context, s *Service, caller Caller, input Input, dependencies OperationContext) (Response, error) {
	result := Response{Schema: ApplicationSchema, Operation: o.Name}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	value, err := o.run(s, caller, input, dependencies)
	result.Data = value
	return result, err // A successful committed operation is never reclassified as canceled.
}

func requestOperation[Q any](name, id, label, description string, call func(*Service, Q, Caller) (Result, error)) Operation {
	return Operation{Name: name, ActionID: id, Label: label, Description: description, Input: RequestInput, Write: true,
		run: func(s *Service, c Caller, in Input, _ OperationContext) (any, error) {
			var q Q
			if err := Decode(in.Request, &q); err != nil {
				return nil, err
			}
			return call(s, q, c)
		}}
}

func Operations() []Operation {
	ops := []Operation{
		requestOperation("create", "activity.create", "Create Activity", "Create a planned Activity and its mandate; no work is dispatched.", (*Service).Create),
		requestOperation("session-create", "activity.session.create", "Create Activity Session", "Create and link one Session document in its Activity; no execution binding is started.", (*Service).CreateSession),
		{Name: "list", ActionID: "activity.list", Label: "Activities", Description: "Read the validated Activity catalog and current canonical history revision.", Input: NoInput, run: func(s *Service, _ Caller, _ Input, _ OperationContext) (any, error) { return s.Read() }},
		{Name: "read", ActionID: "activity.read", Label: "Read Activity", Description: "Read one Activity, current state and declared references.", Input: ActivityInput, run: readOperation},
		{Name: "history", ActionID: "activity.history", Label: "Activity history", Description: "Read original Activity events in causal order.", Input: ActivityInput, run: func(s *Service, c Caller, i Input, d OperationContext) (any, error) {
			a, e := readActivity(s, i.ActivityID)
			return a.Events, e
		}},
		{Name: "validate", ActionID: "activity.validate", Label: "Validate Activities", Description: "Validate canonical Activity history, identities and descriptor snapshots without mutation.", Input: NoInput, run: func(s *Service, _ Caller, _ Input, _ OperationContext) (any, error) { return s.Read() }},
		{Name: "recover", ActionID: "activity.recover", Label: "Recover Activity operation", Description: "Resume the exact retained file/event transaction; refuse changed preimages.", Input: RecoveryInput, Write: true, run: func(s *Service, c Caller, i Input, d OperationContext) (any, error) {
			return s.RecoverLifecycle(i.OperationID, c, d.Inspector)
		}},
	}
	ops = append(ops, Operation{Name: "timeline", ActionID: "activity.timeline", Label: "Activity timeline", Description: "Read one Activity calendar timeline and explicit source markers.", Input: ActivityInput, run: func(s *Service, c Caller, in Input, d OperationContext) (any, error) {
		catalog, e := s.Read()
		if e != nil {
			return nil, e
		}
		p, e := Project(s.Area, catalog, time.Now())
		if e != nil {
			return nil, e
		}
		return p.Render("activity/" + in.ActivityID + "/timeline")
	}})
	ops = append(ops, requestOperation("update", "activity.update", "Update Activity references", "Update reviewed scope and references while retaining descriptor history.", (*Service).Update))
	for _, action := range []string{"activate", "hold", "resume", "switch", "complete", "cancel", "reopen"} {
		action := action
		ops = append(ops, Operation{Name: action, ActionID: "activity." + action, Label: "Activity " + action, Description: "Apply an explicitly requested Activity lifecycle transition.", Input: RequestInput, Write: true, run: func(s *Service, c Caller, in Input, d OperationContext) (any, error) {
			var q TransitionRequest
			if e := Decode(in.Request, &q); e != nil {
				return nil, e
			}
			if q.Action != "" && q.Action != action {
				return nil, fmt.Errorf("request action differs from selected operation")
			}
			q.Action = action
			return s.Transition(q, c, d.Inspector)
		}})
	}
	ops = append(ops, Operation{Name: "adopt", ActionID: "activity.adopt", Label: "Adopt Activity admission", Description: "Enable scoped execution after explicit quiescent store inventory.", Input: RequestInput, Write: true, run: func(s *Service, c Caller, in Input, d OperationContext) (any, error) {
		var q AdoptionRequest
		if e := Decode(in.Request, &q); e != nil {
			return nil, e
		}
		return s.AdoptAdmission(q, c, d.Inspector)
	}})
	ops = append(ops, Operation{Name: "view", ActionID: "activity.view", Label: "View Activity", Description: "Render an on-demand Activity overview, phase, history or timeline.", Input: ViewInput, run: func(s *Service, c Caller, in Input, d OperationContext) (any, error) {
		if d.View == nil {
			return nil, fmt.Errorf("view host unavailable")
		}
		return d.View(in)
	}}, Operation{Name: "tree", ActionID: "activity.tree", Label: "Expand Activities", Description: "Read a bounded revision-bound Activity tree page.", Input: TreeInput, run: func(s *Service, c Caller, in Input, d OperationContext) (any, error) {
		if d.Tree == nil {
			return nil, fmt.Errorf("tree host unavailable")
		}
		return d.Tree(in)
	}})
	return ops
}

func readActivity(s *Service, id string) (Activity, error) {
	c, e := s.Read()
	if e != nil {
		return Activity{}, e
	}
	a, ok := c.Find(id)
	if !ok {
		return Activity{}, fmt.Errorf("Activity %q is not registered", id)
	}
	return a, nil
}
func readOperation(s *Service, _ Caller, i Input, _ OperationContext) (any, error) {
	return readActivity(s, i.ActivityID)
}
func FindOperation(name string) (Operation, bool) {
	for _, o := range Operations() {
		if o.Name == name {
			return o, true
		}
	}
	return Operation{}, false
}
