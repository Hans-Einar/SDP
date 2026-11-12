package sdptool

import (
	"context"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func activityService(selected string) (*activities.Service, error) {
	p, e := discover(selected, false)
	if e != nil {
		return nil, e
	}
	namespace, e := activityProjectID(p)
	if e != nil {
		return nil, e
	}
	return activities.New(p.Area, namespace), nil
}

func activityProjectID(p Project) (string, error) {
	return activities.ResolveProjectID(p.Area, strings.ToUpper(p.Inventory.ProjectID))
}

func activityOperationContext(ctx context.Context, selected string, s *activities.Service) activities.OperationContext {
	return activities.OperationContext{Inspector: activityInspector{selected: selected}, View: func(in activities.Input) (any, error) {
		p, e := discover(selected, false)
		if e != nil {
			return nil, e
		}
		return ActivityView(ctx, p, in.URI, in.Revision, in.Output)
	}, Tree: func(in activities.Input) (any, error) {
		p, e := discover(selected, false)
		if e != nil {
			return nil, e
		}
		return ExpandActivity(p, in.URI, in.Revision)
	}}
}

// Governance's checkout identity uses discovery's project slug; Activity ledger
// namespaces are separate and must never be passed as an execution-store identity.
type activityInspector struct{ selected string }

func (i activityInspector) InspectExecutions(stores []activities.ExecutionStore) (activities.ExecutionSnapshot, error) {
	p, e := discover(i.selected, false)
	if e != nil {
		return activities.ExecutionSnapshot{}, e
	}
	return governance.ActivityExecutionInspector(p.Root, p.Area, p.Inventory.ProjectID).InspectExecutions(stores)
}
