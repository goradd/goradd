package generator

import (
	"fmt"
	"github.com/goradd/gengen/maps"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
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

func (d Span) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control","fmt"}
}

func (d Span) SupportsColumn(col *generator.ColumnType) bool {
	return true
}

func (d Span) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewSpan(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)

	if generator.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, generator.DefaultWrapper)
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

