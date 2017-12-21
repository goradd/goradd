package generator

import (
	"fmt"
	"github.com/spekary/goradd/orm/db"
)

// Utilities used by the code generation process and templates

// Returns the value formatted as a constant. Essentially this just surrounds strings in quotes.
func AsConstant(i interface{}, typ db.GoColumnType) string {
	switch i.(type) {
	case string:
		return "\"" + i.(string) + "\""
	case nil:
		return typ.DefaultValue()
	default:
		return fmt.Sprintf("%v", i)
	}
}