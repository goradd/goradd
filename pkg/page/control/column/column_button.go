package column

import (
	"context"

	"github.com/goradd/goradd/pkg/javascript"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/html5tag"
)

// ButtonColumnClick returns an event that detects a click on the icon in the column.
// The EventValue will be the row value clicked on.
// If you need to also know the column clicked on, you can set the Action's action value to:
//   javascript.NewJsCode(g$(event.goradd.match).columnId())
// and then get the value from the ActionValue.
func ButtonColumnClick() *page.Event {
	e := page.NewEvent("click").
		Selector("[data-gr-btn-col]").
		ActionValue(javascript.JsCode(
			`g$(event.goradd.match).closest("tr").data("value")`,
		))
	return e
}

// ButtonColumn is a column that draws a button that fires the ButtonColumnClick event.
type ButtonColumn struct {
	control.ColumnBase
	buttonHtml       string
	buttonAttributes html5tag.Attributes
}

// NewButtonColumn creates a new button column.
func NewButtonColumn() *ButtonColumn {
	i := ButtonColumn{}
	i.Self = &i
	i.Init()
	return &i
}

func (c *ButtonColumn) Init() {
	c.ColumnBase.Init(c)
	c.buttonHtml = `&#9998` // default to a standard pencil icon
	c.buttonAttributes = html5tag.NewAttributes().
		AddClass("gr-transparent-btn") // style so button part does not show, only the icon
	c.SetIsHtml(true)
}

// ButtonAttributes returns the attributes of the button. You can directly manipulate those.
func (c *ButtonColumn) ButtonAttributes() html5tag.Attributes {
	return c.buttonAttributes
}

func (c *ButtonColumn) CellData(ctx context.Context, row int, col int, data interface{}) interface{} {
	c.buttonAttributes.SetData("grBtnCol", "1") // make sure we are tagged so event will fire
	return html5tag.RenderTag("button", c.buttonAttributes, c.buttonHtml)
}

func (c *ButtonColumn) Serialize(e page.Encoder) {
	c.ColumnBase.Serialize(e)
	if err := e.Encode(c.buttonHtml); err != nil {
		panic(err)
	}
	if err := e.Encode(c.buttonAttributes); err != nil {
		panic(err)
	}
}

// SetButtonHtml sets the html to use inside the button. By default a pencil is drawn.
// Use `<i class="fas fa-edit" aria-hidden="true"></i>` for a font awesome edit icon
func (c *ButtonColumn) SetButtonHtml(h string) {
	c.buttonHtml = h
}

func (c *ButtonColumn) Deserialize(dec page.Decoder) {
	c.ColumnBase.Deserialize(dec)

	if err := dec.Decode(&c.buttonHtml); err != nil {
		panic(err)
	}
	if err := dec.Decode(&c.buttonAttributes); err != nil {
		panic(err)
	}
}

// ButtonColumnCreator creates a column that displays a clickable icon.
type ButtonColumnCreator struct {
	// ID will assign the given id to the column. If you do not specify it, an id will be given it by the framework.
	ID string
	// Title is the static title string to use in the header row
	Title string
	// IconHtml specifies the html to put inside the button
	ButtonHtml string
	// ButtonAttributes lets you set the button attributes how you want.
	ButtonAttributes html5tag.Attributes
	control.ColumnOptions
}

func (c ButtonColumnCreator) Create(ctx context.Context, parent control.TableI) control.ColumnI {
	col := NewButtonColumn()
	if c.ID != "" {
		col.SetID(c.ID)
	}
	col.SetTitle(c.Title)
	if c.ButtonHtml != "" {
		col.SetButtonHtml(c.ButtonHtml)
	}
	if c.ButtonAttributes != nil {
		col.buttonAttributes = c.ButtonAttributes
	}
	col.ApplyOptions(ctx, parent, c.ColumnOptions)
	return col
}

func init() {
	control.RegisterColumn(ButtonColumn{})
}
