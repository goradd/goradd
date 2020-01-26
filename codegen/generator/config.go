package generator

import (
	"fmt"
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

// TODO: Remove this. We just need a module path to the control I think.
// Also, explore making a default generator so that you don't need one for all controls, just ones that are different than normal.
type ControlCreationInfo string

// DefaultControlTypeFunc is the injected function that determines the default control type for a particular type of database column.
// It gets initialized here, so that if you want to replace it, you can first call the default function
var DefaultControlTypeFunc = DefaultControlType

// DefaultFormFieldCreator defines what form field wrapper will be used for generated controls.
var DefaultFormFieldCreator = "github.com/goradd/goradd/pkg/page/control/FormFieldWrapperCreator"

// DefaultButtonCreator defines what buttons will be used for generated forms.
var DefaultButtonCreator = "github.com/goradd/goradd/pkg/page/control/ButtonCreator"

// DefaultControlType returns the default control type for the given database column
// These types are module paths to the control, and the generator will resolve those to figure out the import paths
// and package names
func DefaultControlType(ref interface{}) string {
	switch col := ref.(type) {
	case *db.ReverseReference:
		if col.IsUnique() {
			return "" // select list I think instead
		} else if col.IsNullable() {
			return "github.com/goradd/goradd/pkg/page/control/CheckboxList"
		}
		return ""
	case *db.ManyManyReference:
		return "github.com/goradd/goradd/pkg/page/control/CheckboxList"
	case *db.Column:
		if col.IsPk {
			return "github.com/goradd/goradd/pkg/page/control/Span" // primary keys are not editable
		}

		if col.IsReference() {
			return "github.com/goradd/goradd/pkg/page/control/SelectList"
		}

		// default control types for columns
		switch col.ColumnType {
		case query.ColTypeBytes:
			return ""
		case query.ColTypeString:
			return "github.com/goradd/goradd/pkg/page/control/Textbox"
		case query.ColTypeInteger:
			return "github.com/goradd/goradd/pkg/page/control/IntegerTextbox"
		case query.ColTypeUnsigned:
			return "github.com/goradd/goradd/pkg/page/control/IntegerTextbox"
		case query.ColTypeInteger64:
			return "github.com/goradd/goradd/pkg/page/control/IntegerTextbox"
		case query.ColTypeUnsigned64:
			return "github.com/goradd/goradd/pkg/page/control/IntegerTextbox"
		case query.ColTypeDateTime:
			return "github.com/goradd/goradd/pkg/page/control/DateTimeSpan"
		case query.ColTypeFloat:
			return "github.com/goradd/goradd/pkg/page/control/FloatTextbox"
		case query.ColTypeDouble:
			return "github.com/goradd/goradd/pkg/page/control/FloatTextbox"
		case query.ColTypeBool:
			return "github.com/goradd/goradd/pkg/page/control/Checkbox"
		case query.ColTypeUnknown:
			return ""
		default:
			return ""
		}
	default:
		panic("Unkown reference type")
	}
}


func WrapFormField(wrapper string, label string, forId string, child string) string {
	return fmt.Sprintf(
`%s{
	ID: "%s",
	For: "%s",
	Label: "%s",
	Child: %s,
}
`, wrapper, forId + "-ff", forId, label, child)

}