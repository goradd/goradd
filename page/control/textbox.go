package control

import (
	"github.com/spekary/goradd/page"
	"goradd/override/control_base"
	"goradd/config"
	"github.com/spekary/goradd/util/types"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/orm/query"
	"github.com/spekary/goradd/codegen/connector"
)

const (
	TextboxTypeDefault  = "text"
	TextboxTypePassword = "password"
	TextboxTypeSearch   = "search"
	TextboxTypeNumber   = "number" // Puts little arrows in box, will need to widen it.
	TextboxTypeEmail    = "email"  // see TextEmail. Prevents submission of RFC5322 email addresses (Gogh Fir <gf@example.com>)
	TextboxTypeTel    = "tel"    // not well supported
	TextboxTypeUrl    = "url"
)

// Text is a basic text entry form item.
type Textbox struct {
	control_base.Textbox
}

type TextboxI interface {
	control_base.TextboxI
}

func NewTextbox(parent page.ControlI, id string) *Textbox {
	t := &Textbox{}
	t.Init(t, parent, id)
	return t
}


// This structure describes the textbox to the connector dialog and code generator

func init() {
	if config.Mode == config.AppModeDevelopment {
		connector.RegisterControl(TextboxDescriber{})
	}
}

type TextboxDescriber struct {

}

func (d TextboxDescriber) Type() string {
	return "Textbox"
}

func (d TextboxDescriber) NewFunc() string {
	return "NewTextbox"
}

func (d TextboxDescriber) Import() string {
	return "github.com/spekary/goradd/page/control"
}

func (d TextboxDescriber) SupportsColumn(col db.ColumnDescription) bool {
	if col.GoType == query.ColTypeBytes ||
		col.GoType == query.ColTypeString {
			return true
	}
	return false
}

func (d TextboxDescriber) ConnectorParams() *types.OrderedMap {
	paramControls := page.ControlConnectorParams()
	paramSet := types.NewOrderedMap()

	paramSet.Set("ColumnCount", connector.ConnectorParam {
		"Column Count",
		"Width of field by the number of characters.",
		connector.ControlTypeInteger,
		`{{var}}.SetColumnCount{{val}}`,
		func(c page.ControlI, val interface{}) {
			c.(*Textbox).SetColumnCount(val.(int))
		}})


	paramControls.Set("Textbox", paramSet)

	return paramControls
}
