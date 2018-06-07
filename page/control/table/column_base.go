package table

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spekary/goradd"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	html2 "html"
	"strconv"
)

const (
	NotSortable    = 0
	SortAscending  = 1
	SortDescending = -1
	NotSorted      = 2
)

type ColumnI interface {
	ID() string
	SetID(string) ColumnI
	setParentTable(TableI)
	Title() string
	SetTitle(string) ColumnI
	Span() int
	SetSpan(int) ColumnI
	IsHidden() bool
	SetHidden(bool)
	DrawColumnTag(ctx context.Context, buf *bytes.Buffer)
	DrawFooterCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer)
	DrawCell(ctx context.Context, row int, col int, data interface{}, buf *bytes.Buffer)
	CellText(ctx context.Context, row int, col int, data interface{}) string
	HeaderCellHtml(ctx context.Context, row int, col int) string
	FooterCellHtml(ctx context.Context, row int, col int) string
	HeaderAttributes(row int, col int) *html.Attributes
	FooterAttributes(row int, col int) *html.Attributes
	ColTagAttributes() *html.Attributes
	UpdateFormValues(ctx *page.Context)
	AddActions(ctrl page.ControlI)
	Action(ctx context.Context, params page.ActionParams)
	SetHeaderTexter(s CellTexter)
	SetCellTexter(s CellTexter)
	SetFooterTexter(s CellTexter)
	SetCellStyler(s html.Attributer)
	IsSortable() bool
	SortDirection() int
	SetSortDirection(int)
	Sortable() ColumnI
	SetIsHtml(columnIsHtml bool)
}

type CellTexter interface {
	CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string
}

type ColumnBase struct {
	goradd.Base
	id               string
	parentTable      TableI
	title            string
	*html.Attributes // These are attributes that will appear on the cell
	headerAttributes *html.Attributes
	footerAttributes *html.Attributes
	colTagAttributes *html.Attributes
	span             int
	renderAsHeader   bool
	isHtml           bool
	cellTexter       CellTexter
	cellStyler       html.Attributer // for individually styling cells
	headerTexter     CellTexter
	footerTexter     CellTexter
	isHidden         bool
	sortDirection    int
}

func (c *ColumnBase) Init(self ColumnI) {
	c.Base.Init(self)
	c.span = 1
	c.Attributes = html.NewAttributes()
}

func (c *ColumnBase) This() ColumnI {
	return c.Self.(ColumnI)
}

func (c *ColumnBase) ID() string {
	return c.id
}

// SetId sets the id of the table. If you are going to provide your own id, do this as the first thing after you create
// a table, or the new id might not propogate through the system correctly.
func (c *ColumnBase) SetID(id string) ColumnI {
	c.id = id
	return c.This()
}

func (c *ColumnBase) setParentTable(t TableI) {
	c.parentTable = t
}

func (c *ColumnBase) Title() string {
	return c.title
}

func (c *ColumnBase) SetTitle(title string) ColumnI {
	c.title = title
	return c.This()
}

func (c *ColumnBase) Span() int {
	return c.span
}

func (c *ColumnBase) SetSpan(span int) ColumnI {
	c.span = span
	return c.This()
}

func (c *ColumnBase) SetRenderAsHeader(r bool) {
	c.renderAsHeader = r
}

func (c *ColumnBase) SetIsHtml(columnIsHtml bool) {
	c.isHtml = columnIsHtml
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

func (c *ColumnBase) HeaderAttributes(row int, col int) *html.Attributes {
	if c.headerAttributes == nil {
		c.headerAttributes = html.NewAttributes()
		c.headerAttributes.Set("scope", "col")
	}
	return c.headerAttributes
}

func (c *ColumnBase) FooterAttributes(row int, col int) *html.Attributes {
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

// HeaderCellHtml returns the text of the indicated header cell. The default will call into the headerTexter if it
// is provided, or just return the Label value. This function can also be overridden by embedding the ColumnBase object
// into another object.
func (c *ColumnBase) HeaderCellHtml(ctx context.Context, row int, col int) (h string) {
	if c.headerTexter != nil {
		h = c.headerTexter.CellText(ctx, c.This(), row, col, nil)
	} else {
		h = html2.EscapeString(c.title)
	}

	if c.IsSortable() {
		h = c.RenderSortButton(h)
	}
	return
}

func (c *ColumnBase) DrawFooterCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}
	cellHtml := c.This().FooterCellHtml(ctx, row, col)

	a := c.This().FooterAttributes(row, col)
	buf.WriteString(html.RenderTag("td", a, cellHtml))
}

func (c *ColumnBase) FooterCellHtml(ctx context.Context, row int, col int) string {
	if c.footerTexter != nil {
		return c.footerTexter.CellText(ctx, c.This(), row, col, nil) // careful, this does not get escaped
	}

	return ""
}

func (c *ColumnBase) DrawCell(ctx context.Context, row int, col int, data interface{}, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}

	cellHtml := c.This().CellText(ctx, row, col, data)
	if !c.isHtml {
		cellHtml = html2.EscapeString(cellHtml)
	}
	a := c.CellAttributes(ctx, row, col, data)
	buf.WriteString(html.RenderTag("td", a, cellHtml))
}

func (c *ColumnBase) CellText(ctx context.Context, row int, col int, data interface{}) string {
	if c.cellTexter != nil {
		return c.cellTexter.CellText(ctx, c.This(), row, col, data)
	}
	return ""
}

func (c *ColumnBase) CellAttributes(ctx context.Context, row int, col int, data interface{}) *html.Attributes {
	if c.cellStyler != nil {
		return c.cellStyler.Attributes(ctx, row, col, data)
	}
	return nil
}

// Sortable indicates that the column should be drawn with sort indicators.
func (c *ColumnBase) Sortable() ColumnI {
	c.sortDirection = NotSorted
	return c.This()
}

func (c *ColumnBase) IsSortable() bool {
	return c.sortDirection != NotSortable
}

// UpdateFormValues is called by the system whenever values are sent by client controls.
// This default version does nothing. Columns that need to record information (checkbox columns for example), should
// implement this.
func (c *ColumnBase) UpdateFormValues(ctx *page.Context) {}

func (c *ColumnBase) AddActions(ctrl page.ControlI) {}

// Do a table action that is directed at this table
// Column implementations can implement this method to receive private actions that they have added using AddActions
func (c *ColumnBase) Action(ctx context.Context, params page.ActionParams) {}

func (c *ColumnBase) RenderSortButton(labelHtml string) string {
	switch c.sortDirection {
	case NotSortable: // do nothing
	case NotSorted:
		labelHtml += " " + html.RenderTag("i", html.NewAttributes().SetClass("fa fa-sort fa-lg"), "")
	case SortAscending:
		labelHtml += " " + html.RenderTag("i", html.NewAttributes().SetClass("fa fa-sort-asc fa-lg"), "")
	case SortDescending:
		labelHtml += " " + html.RenderTag("i", html.NewAttributes().SetClass("fa fa-sort-desc fa-lg"), "")
	}
	return fmt.Sprintf(`<button onclick="$j('#%s').trigger('grsort', '%s'); return false;">%s</button>`, c.parentTable.ID(), c.ID(), labelHtml)
}

func (c *ColumnBase) SortDirection() int {
	return c.sortDirection
}

// SetSortDirection is used internally to set the sort direction indicator
func (c *ColumnBase) SetSortDirection(d int) {
	c.sortDirection = d
}
