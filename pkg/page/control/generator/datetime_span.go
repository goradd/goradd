package generator

import (
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/codegen/generator"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
)

func init() {
	if !config.Release {
		generator.RegisterControlGenerator(DateTimeSpan{})
	}
}

// This structure describes the DateTimeSpan to the connector dialog and code generator
type DateTimeSpan struct {
}

func (d DateTimeSpan) Type() string {
	return "DateTimeSpan"
}

func (d DateTimeSpan) NewFunc() string {
	return "NewDateTimeSpan"
}

func (d DateTimeSpan) Imports() []string {
	return []string{"github.com/goradd/goradd/pkg/page/control"}
}

func (d DateTimeSpan) SupportsColumn(col *generator.ColumnType) bool {
	return true
}

func (d DateTimeSpan) GenerateCreate(namespace string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewDateTimeSpan(c.ParentControl, id)
	ctrl.SetLabel(ctrl.T("%s"))
`, namespace, col.DefaultLabel)

	if generator.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, generator.DefaultWrapper)
	}

	return
}

func (d DateTimeSpan) GenerateGet(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	s = fmt.Sprintf(`c.%s.SetDateTime(c.%s.%s())`, ctrlName, objName, col.GoName)
	return
}

func (d DateTimeSpan) GeneratePut(ctrlName string, objName string, col *generator.ColumnType) (s string) {
	return
}

func (d DateTimeSpan) ConnectorParams() *maps.SliceMap {
	paramControls := page.ControlConnectorParams()
	paramSet := maps.NewSliceMap()
	paramSet.Set("Format", generator.ConnectorParam{
		"Format",
		"format string to use to format the DateTime. See time.Time doc for more info.",
		generator.ControlTypeString,
		`{{var}}.SetFormat{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*control.DateTimeSpan).SetFormat(val.(string))
		}})

	return paramControls
}
