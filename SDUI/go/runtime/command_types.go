package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

type ContextTarget struct {
	Widget        Handle
	ModelRevision uint64
	Item          *CollectionTarget
}
type SurfaceTarget struct {
	Handle                        Handle
	ModelRevision, OpenGeneration uint64
}

// MenuScope is a semantic opening identity, not a native widget/input scope.
type MenuScope struct {
	Handle                       Handle
	ModelRevision, StateRevision uint64
}
type DraftField struct {
	RawDraft                     *string
	FieldValidation              FieldValidation
	OptionTarget                 *OptionTarget
	Handle                       Handle
	Value                        Value
	ValueRevision, DraftRevision uint64
}
type CommandInvocation struct {
	Origin  Handle
	Via     string
	Surface *SurfaceTarget
	Context *ContextTarget
	Checked *bool
}
type DialogRequest struct {
	Surface SurfaceTarget
	Fields  []DraftField
}
type AcceptDecision struct {
	Accepted bool
	Message  string
}
type DialogResult struct {
	Surface                  SurfaceTarget
	Sequence, AcceptSequence uint64
	Kind, Reason             string
	Domain                   DomainOutcome
	Fields                   []DraftField
}
type CommandState struct {
	Handle                                  Handle
	InstancePath, Label, Icon, Tooltip, Key string
	Context, Effect                         string
	Target                                  Handle
	Exclusive, ExclusiveScope               string
	Toggle, Checked, Enabled, Visible       bool
	Binding                                 parser.Reference
}
type CommandPresentation struct {
	Handle, Command                    Handle
	InstancePath, Label, Icon, Tooltip string
	Enabled, Visible                   bool
}
type MenuState struct {
	Handle          Handle
	InstancePath    string
	Open            bool
	Context         *ContextTarget
	Surface         *SurfaceTarget
	Opener          Handle
	CaptureRevision uint64
	Scope           MenuScope
}
type SurfaceState struct {
	Handle              Handle
	InstancePath, Label string
	Modal, Open         bool
	Target              SurfaceTarget
	Parent              Handle
	ParentSurface       *SurfaceTarget
	Opener              Handle
	Context             *ContextTarget
	Focused             string
	AcceptSequence      uint64
	Domain              DomainOutcome
	AcceptBlocked       bool
	Message             string
}

const (
	InvokeCommand EventKind = "invoke-command"
	Accept        EventKind = "accept"
	Cancel        EventKind = "cancel"
	Close         EventKind = "close"
	Checked       Property  = "checked"
)
