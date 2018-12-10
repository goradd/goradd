package generator

import (
	"github.com/spekary/goradd/codegen/generator"
	"github.com/spekary/goradd/pkg/orm/db"
	"github.com/spekary/goradd/pkg/orm/query"
)

// Setup sets up the default code generator to generate bootstrap controls when possible.
func BootstrapCodegenSetup() {
	generator.DefaultControlTypeFunc = func(col *db.ColumnDescription) (info generator.ControlCreationInfo) {
		info = generator.DefaultControlType(col)

		if col.IsPk {
			return
		}

		if col.IsReference() {
			return generator.ControlCreationInfo{"SelectList", "NewSelectList", "github.com/spekary/goradd/pkg/bootstrap/control"}
		}

		// default control types for columns
		switch col.GoType {
		case query.ColTypeString: return generator.ControlCreationInfo{"Textbox", "NewTextbox", "github.com/spekary/goradd/pkg/bootstrap/control"}
		case query.ColTypeInteger: fallthrough
		case query.ColTypeUnsigned: fallthrough
		case query.ColTypeInteger64: fallthrough
		case query.ColTypeUnsigned64: return generator.ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/spekary/goradd/pkg/bootstrap/control"}
		case query.ColTypeFloat: return generator.ControlCreationInfo{"FloatTextbox", "NewFloatTextbox", "github.com/spekary/goradd/pkg/bootstrap/control"}
		case query.ColTypeBool: return generator.ControlCreationInfo{"Checkbox", "NewCheckbox", "github.com/spekary/goradd/pkg/bootstrap/control"}
		default: return
		}
	}
}
