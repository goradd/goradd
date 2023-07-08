package generator

import (
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

// BootstrapCodegenSetup sets up the default code generator to generate bootstrap controls when possible.
func BootstrapCodegenSetup() {
	generator.DefaultFormFieldWrapperType = "github.com/goradd/goradd/pkg/bootstrap/control/FormGroup"

	generator.DefaultButtonType = "github.com/goradd/goradd/pkg/bootstrap/control/Button"
	generator.DefaultDataPagerType = "github.com/goradd/goradd/pkg/bootstrap/control/DataPager"

	generator.DefaultControlTypeFunc = func(ref interface{}) (path string) {
		path = generator.DefaultControlType(ref)
		switch col := ref.(type) {
		case *db.ReverseReference:
			if col.IsUnique() {
				return // select list instead
			} else if col.IsNullable() {
				return "github.com/goradd/goradd/pkg/bootstrap/control/CheckboxList"
			}
			return
		case *db.ManyManyReference:
			return "github.com/goradd/goradd/pkg/bootstrap/control/CheckboxList"
		case *db.Column:
			if col.IsPk {
				return
			}

			if col.IsReference() || col.IsEnum() {
				return "github.com/goradd/goradd/pkg/bootstrap/control/SelectList"
			}

			// default control types for columns
			switch col.ColumnType {
			case query.ColTypeString:
				return "github.com/goradd/goradd/pkg/bootstrap/control/Textbox"
			case query.ColTypeInteger:
				fallthrough
			case query.ColTypeUnsigned:
				fallthrough
			case query.ColTypeInteger64:
				fallthrough
			case query.ColTypeUnsigned64:
				return "github.com/goradd/goradd/pkg/bootstrap/control/IntegerTextbox"
			case query.ColTypeFloat32:
				return "github.com/goradd/goradd/pkg/bootstrap/control/FloatTextbox"
			case query.ColTypeFloat64:
				return "github.com/goradd/goradd/pkg/bootstrap/control/FloatTextbox"
			case query.ColTypeBool:
				return "github.com/goradd/goradd/pkg/bootstrap/control/Checkbox"
			case query.ColTypeTime:
				if col.IsTimestamp {
					return "github.com/goradd/goradd/pkg/page/control/DateTimeSpan"
				} else {
					return "github.com/goradd/goradd/pkg/bootstrap/control/DateTextbox"
				}
			default:
				return
			}
		default:
			panic("Unknown column reference type")
		}
	}
}
