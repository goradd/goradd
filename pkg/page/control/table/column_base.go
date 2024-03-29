package table

import (
	"context"
	"fmt"
	"github.com/goradd/goradd/pkg/any"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/event"
	time2 "github.com/goradd/goradd/pkg/time"
	"html"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/goradd/goradd/pkg/base"
	"github.com/goradd/goradd/pkg/config"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
)

// SortDirection indicates the current direction of the sort arrows for the column.
//   - NotSortable indicates the column cannot be sorted.
//   - NotSorted indicates that the column CAN be sorted, but the direction is currently not determined
//   - SortAscending and SortDescending indicate the column is currently sorted in the corresponding direction.
type SortDirection int

const (
	NotSortable    = SortDirection(0)
	SortAscending  = SortDirection(1)
	SortDescending = SortDirection(-1)
	NotSorted      = SortDirection(-100)
)

const ColumnSortEventName = "gr-sort"

// ColumnSortEvent indicates a column's sort button has been clicked.
func ColumnSortEvent() *event.Event {
	return event.NewEvent(ColumnSortEventName)
}

// SortButtonHtmlGetter is the injected function for getting the html for sort buttons in the column header.
var SortButtonHtmlGetter func(SortDirection) string

// ColumnI defines the interface that all columns must support. Most of these functions are provided by the
// default behavior of the ColumnBase class.
type ColumnI interface {
	ID() string
	init(self ColumnI)
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
	DrawColumnTag(ctx context.Context, w io.Writer)
	DrawFooterCell(ctx context.Context, row int, col int, count int, w io.Writer)
	DrawCell(ctx context.Context, row int, col int, data interface{}, w io.Writer)
	CellText(ctx context.Context, row int, col int, data interface{}) string
	CellData(ctx context.Context, row int, col int, data interface{}) interface{}
	HeaderCellHtml(ctx context.Context, row int, col int) string
	FooterCellHtml(ctx context.Context, row int, col int) string
	HeaderAttributes(ctx context.Context, row int, col int) html5tag.Attributes
	FooterAttributes(ctx context.Context, row int, col int) html5tag.Attributes
	ColTagAttributes() html5tag.Attributes
	UpdateFormValues(ctx context.Context)
	AddActions(ctrl page.ControlI)
	DoAction(ctx context.Context, params action.Params)
	SetCellTexter(s CellTexter) ColumnI
	SetHeaderTexter(s CellTexter) ColumnI
	SetFooterTexter(s CellTexter) ColumnI
	SetCellStyler(s CellStyler)
	IsSortable() bool
	SortDirection() SortDirection
	SetSortDirection(SortDirection) ColumnI
	SetSortable() ColumnI
	RenderSortButton(labelHtml string) string
	SetIsHtml(columnIsHtml bool) ColumnI
	PreRender()
	MarshalState(m page.SavedState)
	UnmarshalState(m page.SavedState)
	Serialize(e page.Encoder)
	Deserialize(dec page.Decoder)
	Restore(parentTable TableI)
}

// CellInfo is provided to the cell texter so the cell texter knows how to draw.
// It is a struct here so that the info can grow without the CellTexter signature having to change.
type CellInfo struct {
	RowNum       int
	ColNum       int
	Data         interface{}
	isHeaderCell bool
	isFooterCell bool
}

// CellTexter defines the interface for a structure that provides the content of a table cell.
// If your CellTexter is not a control, you should register it with gob.
type CellTexter interface {
	CellText(ctx context.Context, col ColumnI, info CellInfo) string
}

// CellStyler defines the interface that returns attributes specific to a particular cell.
type CellStyler interface {
	CellAttributes(ctx context.Context, col ColumnI, info CellInfo) html5tag.Attributes
}

// ColumnBase is the base implementation of all table columns
type ColumnBase struct {
	base.Base
	id                  string
	parentTable         TableI
	parentTableID       string // for deserializing
	title               string
	html5tag.Attributes                       // These are static attributes that will appear on each cell
	headerAttributes    []html5tag.Attributes // static attributes per header row
	footerAttributes    []html5tag.Attributes
	colTagAttributes    html5tag.Attributes
	span                int
	asHeader            bool
	isHtml              bool
	cellTexter          CellTexter
	cellTexterID        string // for deserialization
	headerTexter        CellTexter
	headerTexterID      string // for deserialization
	footerTexter        CellTexter
	footerTexterID      string     // for deserialization
	cellStyler          CellStyler // for dynamically styling cells
	cellStylerID        string     // for deserialization
	isHidden            bool
	sortDirection       SortDirection
	// format is a format string. It will be applied using fmt.Sprintf. If you don't provide a Format string, standard
	// string conversion operations will be used.
	format string
	// timeFormat is applied to the data using time.Format. You can have both a Format and TimeFormat, and the Format
	// will be applied using fmt.Sprintf after the TimeFormat is applied using time.Format.
	timeFormat string
	// showLocalTime will convert a time to client local time for display.
	showLocalTime bool
}

func (c *ColumnBase) Init(self ColumnI) {
	c.Base.Init(self)
	c.Attributes = html5tag.NewAttributes()
}

func (c *ColumnBase) init(self ColumnI) {
	c.Init(self)
}

func (c *ColumnBase) this() ColumnI {
	return c.Self().(ColumnI)
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

// Span returns the number of columns this column will span. If the span is not set, it will return 1.
func (c *ColumnBase) Span() int {
	if c.span < 2 {
		return 1
	} else {
		return c.span
	}
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

// AsHeader returns true if the cells in the column will be drawn as header cells (th).
func (c *ColumnBase) AsHeader() bool {
	return c.asHeader
}

// SetFormat sets string that will be passed to fmt.Sprintf with the data of the cell to
// determine how to draw the cell.
func (c *ColumnBase) SetFormat(format string) ColumnI {
	c.format = format
	return c.this()
}

// SetTimeFormat is the format string sent to time.Format and that is used specifically to format
// time values in the cell. If both a time format and cell format are defined, the time format will be
// used first to convert the time.Time into a string, and then the result will be sent to fmt.Sprintf
// using the format string.
func (c *ColumnBase) SetTimeFormat(timeFormat string) ColumnI {
	c.timeFormat = timeFormat
	return c.this()
}

// SetShowLocalTime will convert a time into client local time for display.
func (c *ColumnBase) SetShowLocalTime(showLocalTime bool) ColumnI {
	c.showLocalTime = showLocalTime
	return c.this()
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
func (c *ColumnBase) HeaderAttributes(ctx context.Context, row int, col int) html5tag.Attributes {
	if len(c.headerAttributes) < row+1 {
		// extend the attributes
		c.headerAttributes = append(c.headerAttributes, make([]html5tag.Attributes, row-len(c.headerAttributes)+1)...)
	}
	if c.headerAttributes[row] == nil {
		c.headerAttributes[row] = html5tag.NewAttributes()
	}
	if row == 0 {
		// for screen readers
		c.headerAttributes[0].Set("scope", "col")
		if c.IsSortable() {
			switch c.SortDirection() {
			case SortAscending:
				c.headerAttributes[0].Set("aria-sort", "ascending")
			case SortDescending:
				c.headerAttributes[0].Set("aria-sort", "descending")
			default:
				c.headerAttributes[0].RemoveAttribute("aria-sort")
			}
		}
	}

	return c.headerAttributes[row]
}

// FooterAttributes returns the attributes to use for the footer cell.
func (c *ColumnBase) FooterAttributes(ctx context.Context, row int, col int) html5tag.Attributes {
	if len(c.footerAttributes) < row+1 {
		// extend the attributes
		c.footerAttributes = append(c.footerAttributes, make([]html5tag.Attributes, row-len(c.footerAttributes)+1)...)
	}
	if c.footerAttributes[row] == nil {
		c.footerAttributes[row] = html5tag.NewAttributes()
	}
	return c.footerAttributes[row]
}

// ColTagAttributes specifies attributes that will appear in the column tag. Note that you have to turn on column
// tags in the table object as well for these to appear.
func (c *ColumnBase) ColTagAttributes() html5tag.Attributes {
	if c.colTagAttributes == nil {
		c.colTagAttributes = html5tag.NewAttributes()
	}
	return c.colTagAttributes
}

// DrawColumnTag draws the column tag in the colgroup.
func (c *ColumnBase) DrawColumnTag(ctx context.Context, w io.Writer) {
	if c.isHidden {
		return
	}
	a := c.this().ColTagAttributes()
	a.SetID(c.parentTable.ID() + "-" + c.ID())
	if c.span > 1 {
		a.Set("span", strconv.Itoa(c.span))
	}
	page.WriteString(w, html5tag.RenderTag("col", a, ""))
	return
}

// HeaderCellHtml returns the text of the indicated header cell. The default will call into the headerTexter if it
// is provided, or just return the Label value. This function can also be overridden by embedding the ColumnBase object
// into another object.
func (c *ColumnBase) HeaderCellHtml(ctx context.Context, row int, col int) (h string) {
	if c.headerTexter != nil {
		info := CellInfo{RowNum: row, ColNum: col, isHeaderCell: true}
		h = c.headerTexter.CellText(ctx, c.this(), info)
	} else if row == 0 {
		h = html.EscapeString(c.title)
		if c.IsSortable() {
			h = c.this().RenderSortButton(h)
		}
	}

	return
}

// DrawFooterCell will draw the footer cells html into the given buffer.
func (c *ColumnBase) DrawFooterCell(ctx context.Context, row int, col int, count int, w io.Writer) {
	if c.isHidden {
		return
	}
	cellHtml := c.this().FooterCellHtml(ctx, row, col)

	a := c.this().FooterAttributes(ctx, row, col)
	tag := "td"
	if c.asHeader {
		tag = "th"
	}
	page.WriteString(w, html5tag.RenderTag(tag, a, cellHtml))
	return
}

// FooterCellHtml returns the html to use in the given footer cell.
func (c *ColumnBase) FooterCellHtml(ctx context.Context, row int, col int) string {
	if c.footerTexter != nil {
		info := CellInfo{RowNum: row, ColNum: col, isFooterCell: true}
		return c.footerTexter.CellText(ctx, c.this(), info) // careful, this does not get escaped
	}

	return ""
}

// DrawCell is the default cell drawing function.
func (c *ColumnBase) DrawCell(ctx context.Context, row int, col int, data interface{}, w io.Writer) {
	if c.isHidden {
		return
	}

	cellHtml := c.this().CellText(ctx, row, col, data)
	if !c.isHtml {
		cellHtml = html.EscapeString(cellHtml)
	}
	a := c.CellAttributes(ctx, row, col, data)

	tag := "td"
	if c.asHeader {
		tag = "th"
	}
	page.WriteString(w, html5tag.RenderTag(tag, a, cellHtml))
	return
}

// CellText returns the text in the cell. It will use the CellTexter if one was provided.
func (c *ColumnBase) CellText(ctx context.Context, row int, col int, data interface{}) string {
	if c.cellTexter != nil {
		info := CellInfo{RowNum: row, ColNum: col, Data: data}
		return c.cellTexter.CellText(ctx, c.this(), info)
	}
	d := c.this().CellData(ctx, row, col, data)
	if t, ok := d.(time.Time); c.showLocalTime && ok {
		d = time2.AtGMTOffset(t, page.GetContext(ctx).ClientTimezoneOffset())
	}
	return c.ApplyFormat(d)
}

// CellData returns the data in the cell.
func (c *ColumnBase) CellData(ctx context.Context, row int, col int, data interface{}) interface{} {
	return ""
}

// CellAttributes returns the attributes of the cell. Column implementations should call this base version first before
// customizing more. It will use the CellStyler if one was provided.
func (c *ColumnBase) CellAttributes(ctx context.Context, row int, col int, data interface{}) html5tag.Attributes {
	if c.Attributes == nil && c.cellStyler == nil {
		return nil
	}
	a := c.Attributes.Copy()
	if c.cellStyler != nil {
		info := CellInfo{RowNum: row, ColNum: col, Data: data}
		a2 := c.cellStyler.CellAttributes(ctx, c.this(), info)
		a.Merge(a2)
	}
	return a
}

// SetSortable indicates that the column should be drawn with sort indicators.
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
func (c *ColumnBase) UpdateFormValues(ctx context.Context) {}

func (c *ColumnBase) AddActions(ctrl page.ControlI) {}

// DoAction does an action is directed at this column.
// Column implementations can implement this method to receive private actions that they have added using AddActions.
// By default, the action will be passed to the parent table if not handled.
func (c *ColumnBase) DoAction(ctx context.Context, params action.Params) {
	if t := c.ParentTable(); t != nil {
		t.DoAction(ctx, params)
	}
}

// RenderSortButton returns the HTML that draws the sort button.
func (c *ColumnBase) RenderSortButton(labelHtml string) string {
	labelHtml += ` ` + c.ParentTable().SortIconHtml(c)

	return fmt.Sprintf(
		`<button class="gr-transparent-btn" onclick="g$('%s').trigger('%s', '%s'); return false;">%s</button>`,
		c.parentTable.ID(), ColumnSortEventName, c.ID(), labelHtml)
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
func (c *ColumnBase) MarshalState(m page.SavedState) {}

// UnmarshalState is an internal function to restore the state of the control.
func (c *ColumnBase) UnmarshalState(m page.SavedState) {}

type columnBaseEncoded struct {
	ID               string
	Title            string
	Attributes       html5tag.Attributes
	HeaderAttributes []html5tag.Attributes
	FooterAttributes []html5tag.Attributes
	ColTagAttributes html5tag.Attributes
	Span             int
	AsHeader         bool
	IsHtml           bool
	IsHidden         bool
	SortDirection    SortDirection
	CellTexter       interface{}
	HeaderTexter     interface{}
	FooterTexter     interface{}
	CellStyler       interface{}
	Format           string
	TimeFormat       string
	ShowLocalTime    bool
}

func (c *ColumnBase) Serialize(e page.Encoder) {
	s := columnBaseEncoded{
		ID:               c.id,
		Title:            c.title,
		Attributes:       c.Attributes,
		HeaderAttributes: c.headerAttributes,
		FooterAttributes: c.footerAttributes,
		ColTagAttributes: c.colTagAttributes,
		Span:             c.span,
		AsHeader:         c.asHeader,
		IsHtml:           c.isHtml,
		IsHidden:         c.isHidden,
		SortDirection:    c.sortDirection,
		CellTexter:       c.cellTexter,
		HeaderTexter:     c.headerTexter,
		FooterTexter:     c.footerTexter,
		CellStyler:       c.cellStyler,
		Format:           c.format,
		TimeFormat:       c.timeFormat,
	}

	if ctrl, ok := c.cellTexter.(page.ControlI); ok {
		s.CellTexter = ctrl.ID()
	}
	if ctrl, ok := c.headerTexter.(page.ControlI); ok {
		s.HeaderTexter = ctrl.ID()
	}
	if ctrl, ok := c.footerTexter.(page.ControlI); ok {
		s.FooterTexter = ctrl.ID()
	}
	if ctrl, ok := c.cellStyler.(page.ControlI); ok {
		s.CellStyler = ctrl.ID()
	}
	s.ShowLocalTime = c.showLocalTime

	if err := e.Encode(s); err != nil {
		panic(err)
	}
}

func (c *ColumnBase) Deserialize(dec page.Decoder) {
	var s columnBaseEncoded
	if err := dec.Decode(&s); err != nil {
		panic(err)
	}

	c.id = s.ID
	c.title = s.Title
	c.Attributes = s.Attributes
	c.headerAttributes = s.HeaderAttributes
	c.footerAttributes = s.FooterAttributes
	c.colTagAttributes = s.ColTagAttributes
	c.span = s.Span
	c.asHeader = s.AsHeader
	c.isHtml = s.IsHtml
	c.isHidden = s.IsHidden
	c.sortDirection = s.SortDirection
	c.format = s.Format
	c.timeFormat = s.TimeFormat

	if s.CellTexter != nil {
		if v, ok := s.CellTexter.(string); ok {
			c.cellTexterID = v
		} else {
			c.cellTexter = s.CellTexter.(CellTexter)
		}
	}
	if s.HeaderTexter != nil {
		if v, ok := s.HeaderTexter.(string); ok {
			c.headerTexterID = v
		} else {
			c.headerTexter = s.HeaderTexter.(CellTexter)
		}
	}
	if s.FooterTexter != nil {
		if v, ok := s.FooterTexter.(string); ok {
			c.footerTexterID = v
		} else {
			c.footerTexter = s.FooterTexter.(CellTexter)
		}
	}
	if s.CellStyler != nil {
		if v, ok := s.CellStyler.(string); ok {
			c.cellStylerID = v
		} else {
			c.cellStyler = s.CellStyler.(CellStyler)
		}
	}
	c.showLocalTime = s.ShowLocalTime
}

func (c *ColumnBase) Restore(parentTable TableI) {
	if c.cellTexterID != "" {
		c.cellTexter = parentTable.Page().GetControl(c.cellTexterID).(CellTexter)
	}
	if c.headerTexterID != "" {
		c.headerTexter = parentTable.Page().GetControl(c.headerTexterID).(CellTexter)
	}
	if c.footerTexterID != "" {
		c.footerTexter = parentTable.Page().GetControl(c.footerTexterID).(CellTexter)
	}
	if c.cellStylerID != "" {
		c.cellStyler = parentTable.Page().GetControl(c.cellStylerID).(CellStyler)
	}

	return
}

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
	CellAttributes html5tag.Attributes
	// HeaderAttributes is a slice of attributes to apply to each row of the header cells in the column.
	// Each item in the slice corresponds to a row of the header.
	HeaderAttributes []html5tag.Attributes
	// FooterAttributes is a slice of attributes to apply to each row of the footer cells in the column.
	// Each item in the slice corresponds to a row of the footer.
	FooterAttributes []html5tag.Attributes
	// ColTagAttributes applies attributes to the col tag if col tags are on in the table. There are limited uses for
	// this, but in particular, you can style a column and give it an id. Use Span to set the span attribute.
	ColTagAttributes html5tag.Attributes
	// Span is specifically for col tags to specify the width of the styling in the col tag.
	Span int
	// AsHeader will cause the entire column to output header tags (th) instead of standard cell tags (td).
	// This is useful for columns on the left or right that contain labels for the rows.
	AsHeader bool
	// IsHtml will cause the text of the cells to NOT be escaped
	IsHtml bool
	// HeaderTexter is an object that will provide the text of the header cells.
	// This can be either an object that you have set up prior, or a string id of a control.
	HeaderTexter interface{}
	// FooterTexter is an object that will provide the text of the footer cells.
	// This can be either an object that you have set up prior, or a string id of a control.
	FooterTexter interface{}
	// IsHidden will start the column out in a hidden state so that it will not initially be drawn.
	IsHidden bool
	// Format is a format string applied to the data using fmt.Sprintf.
	Format string
	// TimeFormat is a format string applied specifically to time data using time.Format.
	TimeFormat string
	// ShowLocalTime will convert the time to the client's local time.
	ShowLocalTime bool
}

func (c *ColumnBase) ApplyOptions(ctx context.Context, parent TableI, opt ColumnOptions) {
	c.Attributes.Merge(opt.CellAttributes)
	if opt.HeaderAttributes != nil {
		for i, row := range opt.HeaderAttributes {
			attr := c.HeaderAttributes(ctx, i, 0)
			attr.Merge(row)
		}
	}
	if opt.FooterAttributes != nil {
		for i, row := range opt.FooterAttributes {
			attr := c.FooterAttributes(ctx, i, 0)
			attr.Merge(row)
		}
	}
	if opt.ColTagAttributes != nil {
		if c.colTagAttributes == nil {
			c.colTagAttributes = html5tag.NewAttributes()
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
		if s, ok := opt.HeaderTexter.(string); ok {
			c.SetHeaderTexter(parent.Page().GetControl(s).(CellTexter))
		} else {
			c.SetHeaderTexter(opt.HeaderTexter.(CellTexter))
		}
	}
	if opt.FooterTexter != nil {
		if s, ok := opt.FooterTexter.(string); ok {
			c.SetFooterTexter(parent.Page().GetControl(s).(CellTexter))
		} else {
			c.SetFooterTexter(opt.FooterTexter.(CellTexter))
		}
	}
	if opt.Format != "" {
		c.SetFormat(opt.Format)
	}
	if opt.TimeFormat != "" {
		c.SetTimeFormat(opt.TimeFormat)
	}
	c.SetShowLocalTime(opt.ShowLocalTime)
}

// ApplyFormat is used by table columns to apply the given fmt.Sprintf and time.Format strings to the data.
// It is exported to allow custom cell Texter objects to use it.
// For slice data, the format is applied to each item in the slice and each item is separated with a comma.
func (c *ColumnBase) ApplyFormat(data interface{}) string {
	var out string

	switch d := data.(type) {
	case int:
		if c.format == "" {
			out = fmt.Sprintf("%d", d)
		} else {
			out = fmt.Sprintf(c.format, d)
		}
	case float64:
		if c.format == "" {
			out = fmt.Sprintf("%f", d)
		} else {
			out = fmt.Sprintf(c.format, d)
		}
	case float32:
		if c.format == "" {
			out = fmt.Sprintf("%f", d)
		} else {
			out = fmt.Sprintf(c.format, d)
		}

	case time.Time:
		t := d
		timeFormat := c.timeFormat
		if timeFormat == "" {
			timeFormat = config.DefaultDateTimeFormat
		}
		out = t.Format(timeFormat)

		if c.format != "" {
			out = fmt.Sprintf(c.format, out)
		}
	case nil:
		return ""
	default:
		var format = c.format
		if format == "" {
			format = `%v`
		}
		if any.IsSlice(d) {
			s := any.InterfaceSlice(d)
			var items []string
			for _, i := range s {
				items = append(items, fmt.Sprintf(format, i))
			}
			return strings.Join(items, ", ")
		} else {
			out = fmt.Sprintf(format, d)
		}
	}
	return out
}
