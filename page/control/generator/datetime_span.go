package generator

import (
	"goradd-project/config"
	"github.com/spekary/goradd/codegen/connector"
	"github.com/spekary/goradd/orm/db"
	"fmt"
	"goradd-project/config/codegen"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/control"
)

func init() {
	if !config.Release {
		connector.RegisterGenerator(DateTimeSpan{})
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

func (d DateTimeSpan) Import() string {
	// TODO: Add fmt to the import list
	return "github.com/spekary/goradd/page/control"
}

func (d DateTimeSpan) SupportsColumn(col *db.ColumnDescription) bool {
	return true
}

func (d DateTimeSpan) GenerateCreate(namespace string, col *db.ColumnDescription) (s string) {
	s = fmt.Sprintf(
		`	ctrl = %s.NewDateTimeSpan(c.ParentControl, id)
	ctrl.SetLabel("%s")
`, namespace, col.GoName)

	if codegen.DefaultWrapper != "" {
		s += fmt.Sprintf(`	ctrl.With(page.NewWrapper("%s"))
`, codegen.DefaultWrapper)
	}

	return
}

func (d DateTimeSpan) GenerateGet(ctrlName string, objName string, col *db.ColumnDescription) (s string) {
	s = fmt.Sprintf(`c.%s.SetDateTime(c.%s.%s())`, ctrlName, objName, col.GoName)
	return
}

func (d DateTimeSpan) GeneratePut(ctrlName string, objName string, col *db.ColumnDescription) (s string) {
	return
}


func (d DateTimeSpan) ConnectorParams() *types.OrderedMap {
	paramControls := page.ControlConnectorParams()
	paramSet := types.NewOrderedMap()
	paramSet.Set("Format", connector.ConnectorParam {
		"Format",
		"format string to use to format the DateTime. See time.Time doc for more info.",
		connector.ControlTypeString,
		`{{var}}.SetFormat{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*control.DateTimeSpan).SetFormat(val.(string))
		}})

	return paramControls
}

