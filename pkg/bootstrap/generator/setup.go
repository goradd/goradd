package generator

import (
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

// Setup sets up the default code generator to generate bootstrap controls when possible.
func BootstrapCodegenSetup() {
	generator.DefaultFormFieldCreator = "bootstrapctrl.FormGroupCreator"

	generator.DefaultControlTypeFunc = func(ref interface{}) (info generator.ControlCreationInfo) {
		info = generator.DefaultControlType(ref)
		switch col := ref.(type) {
		case *db.ReverseReference:
			if col.IsUnique() {
				return // select list instead
			} else if col.IsNullable() {
				return generator.ControlCreationInfo{"CheckboxList", "NewCheckboxList", "github.com/goradd/goradd/pkg/bootstrap/control"}
			}
			return
		case *db.ManyManyReference:
			return generator.ControlCreationInfo{"CheckboxList", "NewCheckboxList", "github.com/goradd/goradd/pkg/bootstrap/control"}
		case *db.Column:
			if col.IsPk {
				return
			}

			if col.IsReference() {
				return generator.ControlCreationInfo{"SelectList", "NewSelectList", "github.com/goradd/goradd/pkg/bootstrap/control"}
			}

			// default control types for columns
			switch col.ColumnType {
			case query.ColTypeString:
				return generator.ControlCreationInfo{"Textbox", "NewTextbox", "github.com/goradd/goradd/pkg/bootstrap/control"}
			case query.ColTypeInteger:
				fallthrough
			case query.ColTypeUnsigned:
				fallthrough
			case query.ColTypeInteger64:
				fallthrough
			case query.ColTypeUnsigned64:
				return generator.ControlCreationInfo{"IntegerTextbox", "NewIntegerTextbox", "github.com/goradd/goradd/pkg/bootstrap/control"}
			case query.ColTypeFloat:
				return generator.ControlCreationInfo{"FloatTextbox", "NewFloatTextbox", "github.com/goradd/goradd/pkg/bootstrap/control"}
			case query.ColTypeDouble:
				return generator.ControlCreationInfo{"FloatTextbox", "NewFloatTextbox", "github.com/goradd/goradd/pkg/bootstrap/control"}
			case query.ColTypeBool:
				return generator.ControlCreationInfo{"Checkbox", "NewCheckbox", "github.com/goradd/goradd/pkg/bootstrap/control"}
			default:
				return
			}
		default:
			panic("Unknown column reference type")
		}
	}
}
