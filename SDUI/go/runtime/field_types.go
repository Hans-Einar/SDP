package runtime

const (
	Number          ValueKind = "number"
	OptionID        ValueKind = "option-id"
	ValidationState Property  = "validation"
)

func Numeric(v float64) Value { return Value{Kind: Number, Number: v} }
func Choice(v ItemID) Value   { return Value{Kind: OptionID, OptionID: v} }

type FieldValidation struct{ Code, Message string }
type NumericConstraints struct{ Min, Max, Step string }
type OptionTarget struct {
	Handle                          Handle
	ModelRevision, OptionGeneration uint64
	OptionID                        ItemID
}
type ChoiceOption struct {
	ID      ItemID
	Label   string
	Enabled bool
}
type FieldTarget struct {
	Handle                                                     Handle
	ModelRevision, StateRevision, ValueRevision, DraftRevision uint64
	OptionGeneration                                           uint64
}
type FieldState struct {
	Target                    FieldTarget
	InstancePath              string
	Accepted, Proposed        Value
	RawDraft                  *string
	Dirty, ReadOnly, Required bool
	Validation                FieldValidation
	Numeric                   *NumericConstraints
	Options                   []ChoiceOption
}
type FieldChange struct {
	Field  FieldState
	Value  Value
	Option *OptionTarget
}
type ChangeHandler func(FieldChange)
type FieldValidator func(FieldState) FieldValidation
type ControlCommit struct {
	ValueRevision uint64
	RawDraft      *string
	Option        *OptionTarget
}
