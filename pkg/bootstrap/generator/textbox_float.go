package generator

import (
	generator2 "github.com/spekary/goradd/codegen/generator"
	generator3 "github.com/spekary/goradd/pkg/page/control/generator"
	"github.com/spekary/goradd/pkg/config"
)

func init() {
	if !config.Release {
		generator2.RegisterControlGenerator(FloatTextbox{})
	}
}

// This structure describes the textbox to the connector dialog and code generator
type FloatTextbox struct {
	generator3.FloatTextbox	// base it on the built-in generator
}

func (d FloatTextbox) Imports() []string {
	return []string{"github.com/spekary/goradd/pkg/bootstrap/control"}
}

