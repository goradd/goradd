package event

// ValidationType is used by active controls, like buttons, to determine what other items on the form will get validated
// when the button is pressed. You can set the ValidationType for a control, but you can also set it for individual events
// and override the control's validation setting.
type ValidationType int

const (
	// ValidateDefault is used by events to indicate they are not overriding a control validation. You should not need to use this.
	ValidateDefault ValidationType = iota
	// ValidateNone indicates the control will not validate the form
	ValidateNone
	// ValidateForm is the default validation for buttons, and indicates the entire form and all controls will validate.
	ValidateForm
	// ValidateSiblingsAndChildren will validate the current control, and all siblings of the control and all
	// children of the siblings and current control.
	ValidateSiblingsAndChildren
	// ValidateSiblingsOnly will validate only the siblings of the current control, but not any child controls.
	ValidateSiblingsOnly
	// ValidateChildrenOnly will validate only the children of the current control.
	ValidateChildrenOnly
	// ValidateContainer will use the validation setting of a parent control with ValidateSiblingsAndChildren, ValidateSiblingsOnly,
	// ValidateChildrenOnly, or ValidateTargetsOnly as the stopping point for validation.
	ValidateContainer
	// ValidateTargetsOnly will only validate the specified targets
	ValidateTargetsOnly
)
