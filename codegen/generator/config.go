package generator

import (
	"github.com/goradd/goradd/pkg/orm/db"
	"github.com/goradd/goradd/pkg/orm/query"
)

// DefaultControlTypeFunc is the injected function that determines the default control type for a particular type of database column.
// It gets initialized here, so that if you want to replace it, you can first call the default function
var DefaultControlTypeFunc = DefaultControlType

// DefaultFormFieldWrapperType defines what form control wrapper will be used for generated controls.
// The wrapper should also have a creator that has the type name DefaultFormFieldWrapperType + "Creator".
var DefaultFormFieldWrapperType = "github.com/goradd/goradd/pkg/page/control/FormFieldWrapper"

// DefaultButtonType defines what buttons will be used for generated forms.
// The button should also have a creator that has the type name DefaultButtonType + "Creator"
var DefaultButtonType = "github.com/goradd/goradd/pkg/page/control/button/Button"

// DefaultDataPagerType defines what pager will be used for generated forms.
// The pager should also have a creator that has the type name DefaultDataPagerType + "Creator"
var DefaultDataPagerType = "github.com/goradd/goradd/pkg/page/control/DataPager"

// DefaultStaticTextType is the type of control to create to display content as static text rather than something editable.
var DefaultStaticTextType = "github.com/goradd/goradd/pkg/page/control/Panel"

// Verbose controls whether to output the list of files being written
var Verbose = false

// DefaultControlType returns the default control type for the given database column
// These types are module paths to the control, and the generator will resolve those to figure out the import paths
// and package names
func DefaultControlType(ref interface{}) string {
	switch col := ref.(type) {
	case *db.ReverseReference:
		if col.IsUnique() {
			return "" // select list I think instead
		} else if col.IsNullable() {
			return "github.com/goradd/goradd/pkg/page/control/list/CheckboxList"
		}
		return ""
	case *db.ManyManyReference:
		return "github.com/goradd/goradd/pkg/page/control/list/CheckboxList"
	case *db.Column:
		if col.IsPk && col.IsId {
			return "github.com/goradd/goradd/pkg/page/control/Span" // primary keys are not editable
		}

		if col.IsReference() || col.IsEnum() {
			return "github.com/goradd/goradd/pkg/page/control/list/SelectList"
		}

		// default control types for columns
		switch col.ColumnType {
		case query.ColTypeBytes:
			return ""
		case query.ColTypeString:
			return "github.com/goradd/goradd/pkg/page/control/textbox/Textbox"
		case query.ColTypeInteger:
			return "github.com/goradd/goradd/pkg/page/control/textbox/IntegerTextbox"
		case query.ColTypeUnsigned:
			return "github.com/goradd/goradd/pkg/page/control/textbox/IntegerTextbox"
		case query.ColTypeInteger64:
			return "github.com/goradd/goradd/pkg/page/control/textbox/IntegerTextbox"
		case query.ColTypeUnsigned64:
			return "github.com/goradd/goradd/pkg/page/control/textbox/IntegerTextbox"
		case query.ColTypeTime:
			return "github.com/goradd/goradd/pkg/page/control/DateTimeSpan"
		case query.ColTypeFloat32:
			return "github.com/goradd/goradd/pkg/page/control/textbox/FloatTextbox"
		case query.ColTypeFloat64:
			return "github.com/goradd/goradd/pkg/page/control/textbox/FloatTextbox"
		case query.ColTypeBool:
			return "github.com/goradd/goradd/pkg/page/control/button/Checkbox"
		case query.ColTypeUnknown:
			return ""
		default:
			return ""
		}
	default:
		panic("Unkown reference type")
	}
}
