package control

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/base"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page"
	html2 "html"
	"strconv"
)

type SortDirection int

const (
	NotSortable    = SortDirection(0)
	SortAscending  = SortDirection(1)
	SortDescending = SortDirection(-1)
	NotSorted      = SortDirection(-100)
)

// SortButtonHtmlGetter is the injected function for getting the html for sort buttons in the column header.
// The default uses FontAwesome to draw the buttons, which means the css for FontAwesome must be loaded
// into the web page. You can change what html is loaded by setting this function.
var SortButtonHtmlGetter func(SortDirection) string

// ColumnI defines the interface that all columns must support. Most of these functions are provided by the
// default behavior of the ColumnBase class.
type ColumnI interface {
	ID() string
	SetID(string) ColumnI
	setParentTable(TableI)
	ParentTable() TableI
	Title() string
	SetTitle(string) ColumnI
	Span() int
	SetSpan(int) ColumnI
	IsHidden() bool
	AsHeader() bool
	SetHidden(bool) ColumnI
	DrawColumnTag(ctx context.Context, buf *bytes.Buffer)
	DrawFooterCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer)
	DrawCell(ctx context.Context, row int, col int, data interface{}, buf *bytes.Buffer)
	CellText(ctx context.Context, row int, col int, data interface{}) string
	HeaderCellHtml(ctx context.Context, row int, col int) string
	FooterCellHtml(ctx context.Context, row int, col int) string
	HeaderAttributes(ctx context.Context, row int, col int) html.Attributes
	FooterAttributes(ctx context.Context, row int, col int) html.Attributes
	ColTagAttributes() html.Attributes
	UpdateFormValues(ctx *page.Context)
	AddActions(ctrl page.ControlI)
	Action(ctx context.Context, params page.ActionParams)
	SetCellTexter(s CellTexter) ColumnI
	SetHeaderTexter(s CellTexter) ColumnI
	SetFooterTexter(s CellTexter) ColumnI
	SetCellStyler(s CellStyler)
	IsSortable() bool
	SortDirection() SortDirection
	SetSortDirection(SortDirection) ColumnI
	SetSortable() ColumnI
	SetIsHtml(columnIsHtml bool) ColumnI
	PreRender()
	MarshalState(m maps.Setter)
	UnmarshalState(m maps.Loader)
}

// CellTexter defines the interface for a structure that provides the content of a table cell.
// If your CellTexter is not a control, you should register it with gob.
type CellTexter interface {
	CellText(ctx context.Context, col ColumnI, rowNum int, colNum int, data interface{}) string
}

type CellStyler interface {
	CellAttributes(ctx context.Context, col ColumnI, row int, colNum int, data interface{}) html.Attributes
}


// ColumnBase is the base implementation of all table columns
type ColumnBase struct {
	base.Base
	id               string
	parentTable      TableI
	title            string
	html.Attributes // These are static attributes that will appear on each cell
	headerAttributes []html.Attributes // static attributes per header row
	footerAttributes []html.Attributes
	colTagAttributes html.Attributes
	span             int
	asHeader         bool
	isHtml           bool
	cellTexter     CellTexter
	headerTexter   CellTexter
	footerTexter   CellTexter
	cellStyler     CellStyler // for dynamically styling cells
	isHidden         bool
	sortDirection    SortDirection
}

func (c *ColumnBase) Init(self ColumnI) {
	c.Base.Init(self)
	c.Attributes = html.NewAttributes()
}

func (c *ColumnBase) this() ColumnI {
	return c.Self.(ColumnI)
}

// ID returns the id of the column
func (c *ColumnBase) ID() string {
	return c.id
}

// SetID sets the id of the column. If you are going to provide your own id, do this as the first thing after you create
// a table, or the new id might not propagate through the system correctly. Note that the id in html will have the table
// id prepended to it. This is required so that actions can be routed to a column.
func (c *ColumnBase) SetID(id string) ColumnI {
	c.id = id
	return c.this()
}

func (c *ColumnBase) setParentTable(t TableI) {
	c.parentTable = t
}

// ParentTable returns the table that is the parent of the column
func (c *ColumnBase) ParentTable() TableI {
	return c.parentTable
}

// Title returns the title text that will appear in the header of the column
func (c *ColumnBase) Title() string {
	return c.title
}

// SetTitle sets the title of the column. It returns a column reference for chaining.
func (c *ColumnBase) SetTitle(title string) ColumnI {
	c.title = title
	return c.this()
}

// Span returns the number of columns this column will span.
func (c *ColumnBase) Span() int {
	return c.span
}

// SetSpan sets the span indicated in the column tag of the column. This is used to create colgroup tags.
func (c *ColumnBase) SetSpan(span int) ColumnI {
	c.span = span
	return c.this()
}

// SetAsHeader will cause the entire column to be output with th instead of td cells.
func (c *ColumnBase) SetAsHeader(r bool) {
	c.asHeader = r
}

func (c *ColumnBase) AsHeader() bool {
	return c.asHeader
}

// SetIsHtml will cause the cell to treat the text it receives as html rather than raw text it should escape.
// Use this with extreme caution. Do not display unescaped text that might come from user input, as it could
// open you up to XSS attacks.
func (c *ColumnBase) SetIsHtml(columnIsHtml bool) ColumnI {
	c.isHtml = columnIsHtml
	return c.this()
}

// SetCellStyler sets the CellStyler for the body cells.
func (c *ColumnBase) SetCellStyler(s CellStyler) {
	c.cellStyler = s
}

// SetCellTexter sets the CellTexter for getting the content of each body cell.
func (c *ColumnBase) SetCellTexter(s CellTexter) ColumnI {
	c.cellTexter = s
	return c.this()
}

// CellTexter returns the cell texter.
func (c *ColumnBase) CellTexter() CellTexter {
	return c.cellTexter
}

// SetHeaderTexter sets the CellTexter that gets the text for header cells.
func (c *ColumnBase) SetHeaderTexter(s CellTexter) ColumnI {
	c.headerTexter = s
	return c.this()
}

// SetFooterTexter sets the CellTexter that gets the text for footer cells.
func (c *ColumnBase) SetFooterTexter(s CellTexter) ColumnI {
	c.footerTexter = s
	return c.this()
}

// IsHidden returns true if the column is hidden.
func (c *ColumnBase) IsHidden() bool {
	return c.isHidden
}

// SetHidden hides the column without removing it completely from the table.
func (c *ColumnBase) SetHidden(h bool) ColumnI {
	c.isHidden = h
	return c.this()
}

// HeaderAttributes returns the attributes to use on the header cell.
// The default version will return an attribute structure which you can use to directly
// manipulate the attributes. If you want something more customized, create your own column and
// implement this function. row and col are zero based.
func (c *ColumnBase) HeaderAttributes(ctx context.Context, row int, col int) html.Attributes {
	if len(c.headerAttributes) < row + 1 {
		// extend the attributes
		c.headerAttributes = append(c.headerAttributes, make([]html.Attributes, row-len(c.headerAttributes)+1)...)
	}
	if c.headerAttributes[row] == nil {
		c.headerAttributes[row] = html.NewAttributes()
		if row == 0 {
			c.headerAttributes[row].Set("scope", "col") // for screen readers
		}
	}
	return c.headerAttributes[row]
}

// FooterAttributes returns the attributes to use for the footer cell.
func (c *ColumnBase) FooterAttributes(ctx context.Context, row int, col int) html.Attributes {
	if len(c.footerAttributes) < row + 1 {
		// extend the attributes
		c.footerAttributes = append(c.footerAttributes, make([]html.Attributes, row-len(c.footerAttributes)+1)...)
	}
	if c.footerAttributes[row] == nil {
		c.footerAttributes[row] = html.NewAttributes()
	}
	return c.footerAttributes[row]
}

// ColTagAttributes specifies attributes that will appear in the table tag. Note that you have to turn on table
// tags in the table object as well for these to appear.
func (c *ColumnBase) ColTagAttributes() html.Attributes {
	if c.colTagAttributes == nil {
		c.colTagAttributes = html.NewAttributes()
	}
	return c.colTagAttributes
}

// DrawColumnTag draws the column tag if one was requested.
func (c *ColumnBase) DrawColumnTag(ctx context.Context, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}
	a := c.this().ColTagAttributes()
	if c.id != "" {
		a.Set("id", c.this().ParentTable().ID()+"_"+c.id) // so that actions can get routed to a column
	}
	if c.span > 1 {
		a.Set("span", strconv.Itoa(c.span))
	}
	buf.WriteString(html.RenderTag("col", a, ""))
}

// HeaderCellHtml returns the text of the indicated header cell. The default will call into the headerTexter if it
// is provided, or just return the Label value. This function can also be overridden by embedding the ColumnBase object
// into another object.
func (c *ColumnBase) HeaderCellHtml(ctx context.Context, row int, col int) (h string) {
	if c.headerTexter != nil {
		h = c.headerTexter.CellText(ctx, c.this(), row, col, nil)
	} else {
		h = html2.EscapeString(c.title)
	}

	if c.IsSortable() {
		h = c.RenderSortButton(h)
	}
	return
}

// DrawFooterCell will draw the footer cells html into the given buffer.
func (c *ColumnBase) DrawFooterCell(ctx context.Context, row int, col int, count int, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}
	cellHtml := c.this().FooterCellHtml(ctx, row, col)

	a := c.this().FooterAttributes(ctx, row, col)
	tag := "td"
	if c.asHeader {
		tag = "th"
	}
	buf.WriteString(html.RenderTag(tag, a, cellHtml))
}

// FooterCellHtml returns the html to use in the given footer cell.
func (c *ColumnBase) FooterCellHtml(ctx context.Context, row int, col int) string {
	if c.footerTexter != nil {
		return c.footerTexter.CellText(ctx, c.this(), row, col, nil) // careful, this does not get escaped
	}

	return ""
}

// DrawCell is the default cell drawing function.
func (c *ColumnBase) DrawCell(ctx context.Context, row int, col int, data interface{}, buf *bytes.Buffer) {
	if c.isHidden {
		return
	}

	cellHtml := c.this().CellText(ctx, row, col, data)
	if !c.isHtml {
		cellHtml = html2.EscapeString(cellHtml)
	}
	a := c.CellAttributes(ctx, row, col, data)

	tag := "td"
	if c.asHeader {
		tag = "th"
	}
	buf.WriteString(html.RenderTag(tag, a, cellHtml))
}

// CellText returns the text in the cell. It will use the CellTexter if one was provided.
func (c *ColumnBase) CellText(ctx context.Context, row int, col int, data interface{}) string {
	if c.cellTexter != nil {
		return c.cellTexter.CellText(ctx, c.this(), row, col, data)
	}
	return ""
}

// CellAttributes returns the attributes of the cell. Column implementations should call this base version first before
// customizing more. It will use the CellStyler if one was provided.
func (c *ColumnBase) CellAttributes(ctx context.Context, row int, col int, data interface{}) html.Attributes {
	if c.cellStyler != nil {
		return c.cellStyler.CellAttributes(ctx, c.this(), row, col, data)
	}
	return nil
}

// MakeSortable indicates that the column should be drawn with sort indicators.
func (c *ColumnBase) SetSortable() ColumnI {
	c.sortDirection = NotSorted
	return c.this()
}

// IsSortable indicates whether the column is sortable, and has a sort indicator in the head.
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
	if SortButtonHtmlGetter != nil {
		labelHtml += SortButtonHtmlGetter(c.sortDirection)
	}
	switch c.sortDirection {
	case NotSorted:
		labelHtml += ` <i class="fa fa-sort fa-lg"></i>`
	case SortAscending:
		labelHtml += ` <i class="fa fa-sort-asc fa-lg"></i>`
	case SortDescending:
		labelHtml += ` <i class="fa fa-sort-desc fa-lg"></i>`
	}

	return fmt.Sprintf(`<button onclick="g$('%s').trigger('grsort', '%s'); return false;">%s</button>`, c.parentTable.ID(), c.ID(), labelHtml)
}

// SortDirection returns the current sort direction.
func (c *ColumnBase) SortDirection() SortDirection {
	return c.sortDirection
}

// SetSortDirection is used internally to set the sort direction indicator.
func (c *ColumnBase) SetSortDirection(d SortDirection) ColumnI {
	c.sortDirection = d
	return c.this()
}

// PreRender is called just before the table is redrawn.
func (c *ColumnBase) PreRender() {}

// MarshalState is an internal function to save the state of the control.
func (c *ColumnBase) MarshalState(m maps.Setter) {}

// UnmarshalState is an internal function to restore the state of the control.
func (c *ColumnBase) UnmarshalState(m maps.Loader) {}

type ColumnCreator interface {
	Create(context.Context, TableI) ColumnI
}

// Columns is a helper to return a group of columns
func Columns(cols ...ColumnCreator) []ColumnCreator {
	return cols
}

// ColumnOptions are settings you can apply to all types of table columns
type ColumnOptions struct {
	// CellAttributes is a static map of attributes to apply to every cell in the column
	CellAttributes   html.AttributeCreator
	// HeaderAttributes is a slice of attributes to apply to each row of the header cells in the column.
	// Each item in the slice corresponds to a row of the header.
	HeaderAttributes []html.AttributeCreator
	// FooterAttributes is a slice of attributes to apply to each row of the footer cells in the column.
	// Each item in the slice corresponds to a row of the footer.
	FooterAttributes []html.AttributeCreator
	// ColTagAttributes applies attributes to the col tag if col tags are on in the table. There are limited uses for
	// this, but in particular, you can style a column and give it an id. Use Span to set the span attribute.
	ColTagAttributes html.AttributeCreator
	// Span is specifically for col tags to specify the width of the styling in the col tag.
	Span             int
	// AsHeader will cause the entire column to output header tags (th) instead of standard cell tags (td).
	// This is useful for columns on the left or right that contain labels for the rows.
	AsHeader   bool
	// IsHtml will cause the text of the cells to NOT be escaped
	IsHtml           bool
	// HeaderTexter is an object that will provide the text of the header cells. This can be either an
	// object that you have set up prior, or a string id of a control
	HeaderTexter   	 interface{}
	// FooterTexter is an object that will provide the text of the footer cells. This can be either an
	// object that you have set up prior, or a string id of a control
	FooterTexter   	 interface{}
	// IsHidden will start the column out in a hidden state so that it will not initially be drawn
	IsHidden         bool
}

func (c *ColumnBase) ApplyOptions(ctx context.Context, parent TableI, opt ColumnOptions) {
	c.Attributes.Merge(opt.CellAttributes)
	if opt.HeaderAttributes != nil {
		for i,row := range opt.HeaderAttributes {
			attr := c.HeaderAttributes(ctx, i, 0)
			attr.Merge(row)
		}
	}
	if opt.FooterAttributes != nil {
		for i,row := range opt.FooterAttributes {
			attr := c.FooterAttributes(ctx, i, 0)
			attr.Merge(row)
		}
	}
	if opt.ColTagAttributes != nil {
		if c.colTagAttributes == nil {
			c.colTagAttributes = html.NewAttributes()
		}
		c.colTagAttributes.Merge(opt.ColTagAttributes)
	}

	c.isHidden = opt.IsHidden

	if opt.Span != 0 {
		c.SetSpan(opt.Span)
	}
	c.asHeader = opt.AsHeader
	if opt.IsHtml {
		c.isHtml = true
	}
	if opt.HeaderTexter != nil {
		if s,ok := opt.HeaderTexter.(string); ok {
			c.SetHeaderTexter(parent.Page().GetControl(s).(CellTexter))
		} else {
			c.SetHeaderTexter(opt.HeaderTexter.(CellTexter))
		}
	}
	if opt.FooterTexter != nil {
		if s,ok := opt.FooterTexter.(string); ok {
			c.SetFooterTexter(parent.Page().GetControl(s).(CellTexter))
		} else {
			c.SetFooterTexter(opt.FooterTexter.(CellTexter))
		}
	}
}
