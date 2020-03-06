package generator

import (
	"fmt"
	"github.com/goradd/goradd/codegen/generator"
)

func init() {
	generator.RegisterControlGenerator(Panel{}, "github.com/goradd/goradd/pkg/page/control/Panel")
}

// This structure describes the Panel to the connector dialog and code generator
type Panel struct {
}

func (d Panel) Imports() []string {
	return []string{"fmt"}
}

func (d Panel) SupportsColumn(ref interface{}) bool {
	return true
}

func (d Panel) GenerateCreator(ref interface{}, desc *generator.ControlDescription) (s string) {
	s = fmt.Sprintf(
`%s.PanelCreator{
	ID:        p.ID() + "-%s",
	ControlOptions: page.ControlOptions{
		DataConnector: %s{},
	},
}`, desc.Package, desc.ControlID, desc.Connector)
	return
}


func (d Panel) GenerateRefresh(ref interface{}, desc *generator.ControlDescription) (s string) {
	return `ctrl.SetText(fmt.Sprint(val))`
}

func (d Panel) GenerateUpdate(ref interface{}, desc *generator.ControlDescription) (s string) {
	return ""
}

