package generator

import (
	generator2 "github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	generator3 "github.com/goradd/goradd/pkg/page/control/generator"
)

func init() {
	if !config.Release {
		generator2.RegisterControlGenerator(FloatTextbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type FloatTextbox struct {
	generator3.FloatTextbox // base it on the built-in generator
}

func (d FloatTextbox) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/bootstrap/control"}
}
