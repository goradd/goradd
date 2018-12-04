package generator

import (
	"github.com/spekary/goradd/pkg/orm/db"
	"github.com/spekary/goradd/pkg/orm/query"
)

type ControlCreationInfo struct {
	Typ string
	CreateFunc string
	ImportName string
}

// DefaultControlTyper is the interface that describes the service that determines the default control type that corresponds to a particular column type.
type DefaultControlTyper interface {
	DefaultControlType(col *db.ColumnDescription) ControlCreationInfo
}

// ControlTyper is the global variable that controls how default control types are determined. The default sets up basic html controls
// for various types. Insert your own to change it.

var ControlTyper DefaultControlTyper = new(defaultControlTyper)

// DefaultWrapper defines what wrapper will be used for generated controls. It should correspond to the string the wrapper was registered with.
var DefaultWrapper = "page.Label"

// GenerateControlIDs will determine if the code generator will assign ids to the controls based on table and column names
var GenerateControlIDs = true

//const FormSubDirectory = "/pkg/gen/form"


type defaultControlTyper struct {}

func (t defaultControlTyper) DefaultControlType(col *db.ColumnDescription) ControlCreationInfo {
	if col.IsPk {
		return ControlCreationInfo{"Span", "NewSpan", "github.com/spekary/goradd/pkg/page/control"} // primary keys are not editable
	}

	if col.IsReference() {
		return ControlCreationInfo{"SelectList", "NewSelectList", "github.com/spekary/goradd/pkg/page/control"}
	}

	// default control types for columns
	switch col.GoType {
	case query.ColTypeBytes: return ControlCreationInfo{"", "", ""}
	case query.ColTypeString: return ControlCreationInfo{"Textbox", "NewTextbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeInteger: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeUnsigned: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeInteger64: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeUnsigned64: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeDateTime: return ControlCreationInfo{"DateTimeSpan", "NewDateTimeSpan", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeFloat: return ControlCreationInfo{"FloatTextbox", "NewFloatTextbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeBool: return ControlCreationInfo{"Checkbox", "NewCheckbox", "github.com/spekary/goradd/pkg/page/control"}
	case query.ColTypeUnknown: return ControlCreationInfo{"", "", ""}
	default: return ControlCreationInfo{"", "", ""}
	}
}

