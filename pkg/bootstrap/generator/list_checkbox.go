package generator

import (
	generator2 "github.com/goradd/goradd/codegen/generator"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	generator2.RegisterControlGenerator(CheckboxList{}, "github.com/goradd/goradd/pkg/bootstrap/control/CheckboxList")
}

// CheckboxList describes the checkbox list to the connector dialog and code generator
type CheckboxList struct {
	generator3.CheckboxList // base it on the built-in generator
}
