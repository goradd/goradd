package table

import (
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"strconv"
	html2 "html"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/page"
)

const ColumnAction = 1000

type ColumnI interface {
	Id() string
	SetId(string)
	Label() string
	Span() int
	IsHidden() bool
	SetHidden(bool)
	DrawColumnTag(ctx context.Context, buf *bytes.Buffer)
	DrawHeaderCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer)
	DrawFooterCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer)
	DrawCell(ctx context.Context, row int, col int, data interface{}, buf *bytes.Buffer)
	CellText(ctx context.Context, row int, col int, data interface{}) string
	HeaderCellText(ctx context.Context, row int, col int) string
	FooterCellText(ctx context.Context, row int, col int) string
	HeaderAttributes() *html.Attributes
	FooterAttributes() *html.Attributes
	ColTagAttributes() *html.Attributes
	UpdateFormValues(ctx *page.Context)
	AddActions(ctrl page.ControlI)
	Action(ctx context.Context, params page.ActionParams)
}

type CellTexter interface {
	CellText(ctx context.Context, row int, col int, data interface{}) string
}

type ColumnBase struct {
	goradd.Base
	id               string
	label            string
	*html.Attributes					// These are attributes that will appear on the cell
	headerAttributes *html.Attributes
	footerAttributes *html.Attributes
	colTagAttributes *html.Attributes
	span             int
	renderAsHeader   bool
	dontEscape       bool
	cellTexter		 CellTexter
	cellStyler       html.Attributer	// for individually styling cells
	headerTexter	 CellTexter
	footerTexter	 CellTexter
	isHidden         bool
	orderByObj		interface{}			// Indicates we are sorting. Is any kind of data that is saved in the table and
										// then fed back to the table to determine how to sort.
	reverseOrderByObj	interface{}
}

func (c *ColumnBase) Init(self ColumnI) {
	c.Base.Init(self)
	c.span = 1
	c.Attributes = html.NewAttributes()
}

func (c *ColumnBase) This() ColumnI {
	return c.Self.(ColumnI)
}

func (c *ColumnBase) Id() string {
	return c.id
}

// SetId sets the id of the table. If you are going to provide your own id, do this as the first thing after you create
// a table, or the new id might not propogate through the system correctly.
func (c *ColumnBase) SetId(id string) {
	c.id = id
}

func (c *ColumnBase) Label() string {
	return c.label
}

func (c *ColumnBase) SetLabel(label string) ColumnI {
	c.label = label
	return c.This()
}

func (c *ColumnBase) Span() int {
	return c.span
}

func (c *ColumnBase) SetSpan(span int) {
	c.span = span
}

func (c *ColumnBase) SetRenderAsHeader(r bool) {
	c.renderAsHeader = r
}

func (c *ColumnBase) SetHtmlEscape(e bool) {
	c.dontEscape = !e
}

func (c *ColumnBase) SetCellStyler(s html.Attributer) {
	c.cellStyler = s
}

func (c *ColumnBase) SetCellTexter(s CellTexter) {
	c.cellTexter = s
}

func (c *ColumnBase) SetHeaderTexter(s CellTexter) {
	c.headerTexter = s
}

func (c *ColumnBase) SetFooterTexter(s CellTexter) {
	c.footerTexter = s
}


func (c *ColumnBase) IsHidden() bool {
	return c.isHidden
}

func (c *ColumnBase) SetHidden(h bool) {
	c.isHidden = h
}

func (c *ColumnBase) HeaderAttributes() *html.Attributes {
	if c.headerAttributes == nil {
		c.headerAttributes = html.NewAttributes()
	}
	return c.headerAttributes
}

func (c *ColumnBase) FooterAttributes() *html.Attributes {
	if c.footerAttributes == nil {
		c.footerAttributes = html.NewAttributes()
	}
	return c.footerAttributes
}

// ColTagAttributes specifies attributes that will appear in the table tag. Note that you have to turn on table
// tags in the table object as well for these to appear.
func (c *ColumnBase) ColTagAttributes() *html.Attributes {
	if c.colTagAttributes == nil {
		c.colTagAttributes = html.NewAttributes()
	}
	return c.colTagAttributes
}

func (c *ColumnBase) DrawColumnTag(ctx context.Context, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}
	a := c.This().ColTagAttributes()
	if c.id != "" {
		a.Set("id", c.id)
	}
	if c.span != 1 {
		a.Set("span", strconv.Itoa(c.span))
	}
	buf.WriteString(html.RenderTag("col", a, ""))
}

func (c *ColumnBase) DrawHeaderCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}
	cellHtml := c.This().HeaderCellText(ctx, row, col)
	if !c.dontEscape {
		cellHtml = html2.EscapeString(cellHtml)
	}
	a := c.This().HeaderAttributes()
	a.Set("scope", "col")
	buf.WriteString(html.RenderTag("th", a, cellHtml))
}

// HeaderCellText returns the text of the indicated header cell. The default will call into the headerTexter if it
// is provided, or just return the Label value. This function can also be overridden by embedding the ColumnBase object
// into another object.
func (c *ColumnBase) HeaderCellText(ctx context.Context, row int, col int) string {
	if c.headerTexter != nil {
		return c.headerTexter.CellText(ctx, row, col, nil)
	}
	return c.label
}

func (c *ColumnBase) DrawFooterCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}
	cellHtml := c.This().FooterCellText(ctx, row, col)
	if !c.dontEscape {
		cellHtml = html2.EscapeString(cellHtml)
	}
	a := c.This().FooterAttributes()
	buf.WriteString(html.RenderTag("td", a, cellHtml))
}

func (c *ColumnBase) FooterCellText(ctx context.Context, row int, col int) string {
	if c.footerTexter != nil {
		return c.footerTexter.CellText(ctx, row, col, nil)
	}

	return ""
}

func (c *ColumnBase) DrawCell(ctx context.Context, row int, col int, data interface{}, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}

	cellHtml := c.This().CellText(ctx, row, col, data)
	if !c.dontEscape {
		cellHtml = html2.EscapeString(cellHtml)
	}
	a := c.CellAttributes(ctx, row, col, data)
	buf.WriteString(html.RenderTag("td", a, cellHtml))
}

func (c *ColumnBase) CellText(ctx context.Context, row int, col int, data interface{}) string {
	if c.cellTexter != nil {
		return c.cellTexter.CellText(ctx, row, col, data)
	}
	return ""
}

func (c *ColumnBase) CellAttributes(ctx context.Context, row int, col int, data interface{}) *html.Attributes {
	if c.cellStyler != nil {
		return c.cellStyler.Attributes(ctx, row, col, data)
	}
	return nil
}

func (c *ColumnBase) OrderByObj() interface{} {
	return c.orderByObj
}

func (c *ColumnBase) SetOrderByObj(o interface{}) {
	c.orderByObj = o
}

func (c *ColumnBase) ReverseOrderByObj() interface{} {
	return c.reverseOrderByObj
}

func (c *ColumnBase) SetReverseOrderByObj(o interface{}) {
	c.reverseOrderByObj = o
}

// UpdateFormValues is called by the system whenever values are sent by client controls.
// This default version does nothing. Columns that need to record information (checkbox columns for example), should
// implement this.
func (c *ColumnBase) UpdateFormValues(ctx *page.Context) {}

func (c *ColumnBase) AddActions(ctrl page.ControlI) {}

// Do a table action that is directed at this table
// Column implementations can implement this method to receive private actions that they have added using AddActions
func (c *ColumnBase) Action(ctx context.Context, params page.ActionParams) {}
