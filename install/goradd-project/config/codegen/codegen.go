package codegen

import (
	"github.com/spekary/goradd/pkg/orm/db"
	"github.com/spekary/goradd/pkg/orm/query"
)

// ControlType returns the default type of control for a column.
func DefaultControlType(col *db.ColumnDescription) (typ string, createFunc string, importName string) {
	if col.IsPk {
		return "Span", "NewSpan", "github.com/spekary/goradd/pkg/page/control" // primary keys are not editable
	}

	if col.IsReference() {
		return "SelectList", "NewSelectList", "github.com/spekary/goradd/pkg/page/control"
	}

	// default control types for columns
	switch col.GoType {
	case query.ColTypeUnknown: return "", "", ""
	case query.ColTypeBytes: return "", "", ""
	case query.ColTypeString: return "Textbox", "NewTextbox", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeInteger: return "IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeUnsigned: return "IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeInteger64: return "IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeUnsigned64: return "IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeDateTime: return "DateTimeSpan", "NewDateTimeSpan", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeFloat: return "FloatTextbox", "NewFloatTextbox", "github.com/spekary/goradd/pkg/page/control"
	case query.ColTypeBool: return "Checkbox", "NewCheckbox", "github.com/spekary/goradd/pkg/page/control"
	}
	return
}

// DefaultWrapper defines what wrapper will be used for generated controls
const DefaultWrapper = "bootstrap.FormGroup"

// GenerateControlIDs will determine if the code generator will assign ids to the controls based on table and column names
const GenerateControlIDs = true

const FormSubDirectory = "/pkg/gen/form"