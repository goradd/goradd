package generator

import (
	"github.com/spekary/gengen/maps"
	"goradd-project/config"
	"fmt"
	"goradd-project/config/codegen"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/codegen/generator"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(Span{})
	}
}

// This structure describes the Span to the connector dialog and code generator
type Span struct {

}

func (d Span) Type() string {
	return "Span"
}

func (d Span) NewFunc() string {
	return "NewSpan"
}

func (d Span) Import() string {
	// TODO: Add fmt to the import list
	return "github.com/spekary/goradd/page/control"
}

func (d Span) SupportsColumn(col *generator.ColumnType) bool {
	return true
}

func (d Span) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewSpan(c.ParentControl, id)
	ctrl.SetLabel("%s")
`, namespace, col.DefaultLabel)

	if codegen.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, codegen.DefaultWrapper)
	}

	return
}

func (d Span) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetText(fmt.Sprintf("%%v", c.%s.%s()))`, ctrlName, objName, col.GoName)
	return
}

func (d Span) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	return
}


func (d Span) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()

	return paramControls
}

