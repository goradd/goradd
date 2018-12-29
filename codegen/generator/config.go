package generator

import (
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

type ControlCreationInfo struct {
	Typ string
	CreateFunc string
	ImportName string
}

// DefaultControlTypeFunc is the injected function that determines the default control type for a particular type of database column.
// It gets initialized here, so that if you want to replace it, you can first call the default function
var DefaultControlTypeFunc func(*db.ColumnDescription) ControlCreationInfo = DefaultControlType

// DefaultWrapper defines what wrapper will be used for generated controls. It should correspond to the string the wrapper was registered with.
var DefaultWrapper = "page.Label"

// GenerateControlIDs will determine if the code generator will assign ids to the controls based on table and column names
var GenerateControlIDs = true

func DefaultControlType(col *db.ColumnDescription) ControlCreationInfo {
	if col.IsPk {
		return ControlCreationInfo{"Span", "NewSpan", "github.com/goradd/goradd/pkg/page/control"} // primary keys are not editable
	}

	if col.IsReference() {
		return ControlCreationInfo{"SelectList", "NewSelectList", "github.com/goradd/goradd/pkg/page/control"}
	}

	// default control types for columns
	switch col.GoType {
	case query.ColTypeBytes: return ControlCreationInfo{"", "", ""}
	case query.ColTypeString: return ControlCreationInfo{"Textbox", "NewTextbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeInteger: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeUnsigned: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeInteger64: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeUnsigned64: return ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeDateTime: return ControlCreationInfo{"DateTimeSpan", "NewDateTimeSpan", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeFloat: return ControlCreationInfo{"FloatTextbox", "NewFloatTextbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeBool: return ControlCreationInfo{"Checkbox", "NewCheckbox", "github.com/goradd/goradd/pkg/page/control"}
	case query.ColTypeUnknown: return ControlCreationInfo{"", "", ""}
	default: return ControlCreationInfo{"", "", ""}
	}
}

