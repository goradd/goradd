package control

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/log"
	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/action"
	"github.com/goradd/goradd/pkg/page/control/data"
	"github.com/goradd/goradd/pkg/page/event"
	buf3 "github.com/goradd/goradd/pkg/pool"
	html2 "html"
	"strconv"
)

const (
	ColumnAction = iota + 2000
	SortClick
)

// TableI is the table interface that lets you create a "subclass" of the Table object.
// The functions defined here are hooks that you can implement in your subclass.
type TableI interface {
	page.ControlI
	TableEmbedder
	data.DataManagerEmbedder
}

// TableEmbedder is a workaround for a problem in GO interfaces. See https://github.com/golang/go/issues/6977
type TableEmbedder interface {
	SetCaption(interface{}) TableI
	DrawCaption(context.Context, *bytes.Buffer) error
	GetHeaderRowAttributes(row int) *html.Attributes
	GetFooterRowAttributes(row int) *html.Attributes
	GetRowAttributes(row int, data interface{}) *html.Attributes
	HeaderCellDrawingInfo(ctx context.Context, col ColumnI, rowNum int, colNum int) (cellHtml string, cellAttributes *html.Attributes)
	FooterCellDrawingInfo(ctx context.Context, col ColumnI, rowNum int, colNum int) (cellHtml string, cellAttributes *html.Attributes)
	SetRenderColumnTags(r bool) TableI
	SetHideIfEmpty(h bool) TableI
	SetHeaderRowCount(count int) TableI
	SetFooterRowCount(count int) TableI
	SetRowStyler(a TableRowAttributer) TableI
	SetHeaderRowStyler(a TableHeaderRowAttributer) TableI
	SetFooterRowStyler(a TableFooterRowAttributer) TableI
	AddColumnAt(column ColumnI, loc int)
	AddColumn(column ColumnI) ColumnI
	GetColumn(loc int) ColumnI
	GetColumnByID(id string) ColumnI
	GetColumnByTitle(title string) ColumnI
	RemoveColumn(loc int)
	RemoveColumnByID(id string)
	RemoveColumnByTitle(title string)
	ClearColumns()
	HideColumns()
	ShowColumns()
	MakeSortable() TableI
	SetSortHistoryLimit(n int) TableI
}

// TableRowAttributer is used to style particular table rows.
type TableRowAttributer interface {
	// TableRowAttributes returns attributes that should be used on the particular row indicated.
	// Data is the data for that row.
	TableRowAttributes(row int, data interface{}) *html.Attributes
}

type TableHeaderRowAttributer interface {
	// TableHeaderRowAttributes returns attributes to use for the particular header row indicated
	TableHeaderRowAttributes(row int) *html.Attributes
}

type TableFooterRowAttributer interface {
	// TableFooterRowAttributes returns attributes to use for the particular footer row indicated
	TableFooterRowAttributes(row int) *html.Attributes
}

// Table is a goradd control that outputs a dynamic HTML table object, with table, tr, th and td tags,
// as well as optional col, thead, and tfoot tags.
//
// To use a Table, call NewTable and then add column objects to it. The columns use a CellTexter to draw the contents of
// a cell in the table. There are a number of predefined columns to draw text coming from slices of maps, slices, database objects,
// as well as custom functions you define. See the examples directory for examples of using a Table object.
// See also the PagedTable for a table that works with a Pager object to page through a large data set.
//
// Call MakeSortable() to make a table sortable, in which case the user can click in the header of a column to sort
// by that column. The Table maintains a history of what columns have been sorted by what row, so that you can
// implement multi-level sorting if you so desire. This is particularly helpful when some columns have duplicate
// data, that then get further identified by another column.
//
// You can provide data to the table when it is first created by calling SetData,
// or you can provide the data just before drawing by calling SetDataProvider at initialization of the table, and
// then providing the table data in response to the BindData function. The advantage of this approach is that if
// you are showing a lot of data, the data is only loaded in memory during drawing, and not kept in the pagestate
// after drawing.
type Table struct {
	page.Control
	data.DataManager

	columns               []ColumnI
	renderColumnTags      bool
	caption               interface{}
	hideIfEmpty           bool
	headerRowCount        int
	footerRowCount        int
	currentHeaderRowIndex int //??
	currentRowIndex       int //??
	rowStyler             TableRowAttributer
	headerRowStyler       TableHeaderRowAttributer
	footerRowStyler       TableFooterRowAttributer
	columnIdCounter       int

	// Sort info. Sorting is difficult enough, and intertwined with tables enough, that we just make it built in to every column
	sortColumns      []string // keeps a historical list of columns sorted on
	sortHistoryLimit int      // how far back to go
}

// NewTable creates a new table
func NewTable(parent page.ControlI, id string) *Table {
	t := &Table{}
	t.Init(t, parent, id)
	return t
}

// Init is an internal function that enables the object-oriented pattern of calling virtual functions used by the
// goradd controls. You would only call this if you were implementing a "subclass" of the Table. Call it immediately after
// creating your Table structure, passing the newly created table as "self".
func (t *Table) Init(self page.ControlI, parent page.ControlI, id string) {
	t.Control.Init(self, parent, id)
	t.Tag = "table"
	t.columns = []ColumnI{}
	t.sortHistoryLimit = 1
}

// this returns the TableI interface for calling into "virtual" functions. This allows us to call functions defined
// by a subclass.
func (t *Table) this() TableI {
	return t.Self.(TableI)
}

func (t *Table) SetRenderColumnTags(r bool) TableI {
	if t.renderColumnTags != r {
		t.renderColumnTags = r
		t.Refresh()
	}
	return t.this()
}

func (t *Table) SetHideIfEmpty(h bool) TableI {
	if t.hideIfEmpty != h {
		t.hideIfEmpty = h
		t.Refresh()
	}
	return t.this()
}


// MakeSortable makes a table sortable. It will attach sortable events and show the header if its not shown.
func (t *Table) MakeSortable() TableI {
	t.On(event.TableSort(), action.Ajax(t.ID(), SortClick), action.PrivateAction{})
	if t.headerRowCount == 0 {
		t.headerRowCount = 1
	}
	return t.this()
}

// SetCaption sets the caption of the table. The default Table permits a caption to be either a string, or
// a goradd control.
func (t *Table) SetCaption(caption interface{}) TableI {
	t.caption = caption
	return t.this()
}

// SetHeaderRowCount sets the number of header rows shown. Each column will be asked to draw this number of header rows.
func (t *Table) SetHeaderRowCount(count int) TableI {
	t.headerRowCount = count
	return t.this()
}

// SetFooterRowCount sets the number of footer rows shown. Each column will be asked to draw this number of footer rows.
func (t *Table) SetFooterRowCount(count int) TableI {
	t.footerRowCount = count
	return t.this()
}

// ΩDrawTag is called by the framework to draw the table. The Table overrides this to call into the DataProvider
// to load the table's data into memory just before drawing. The data will be unloaded after drawing.
func (t *Table) ΩDrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing table tag")
	if t.HasDataProvider() {
		log.FrameworkDebug("Getting table data")
		t.LoadData(ctx, t.this())
		defer t.ResetData()
	}
	if t.hideIfEmpty && !t.HasData() {
		return ""
	}
	for _, c := range t.columns {
		c.PreRender()
	}
	return t.Control.ΩDrawTag(ctx)
}

// ΩDrawingAttributes is an override to add attributes to the table, including not showing the table at all if there
// is no data to show. This will hide header and footer cells and potentially the outline of the table when there is no
// data in the table.
func (t *Table) ΩDrawingAttributes() *html.Attributes {
	a := t.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "table")
	if !t.HasData() && t.hideIfEmpty {
		a.SetStyle("display", "none")
	}
	return a
}

// ΩDrawInnerHtml is an override to draw the meat of the table.
func (t *Table) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	var t2 = t.this().(TableI) // Get the sub class so we call into its hooks for drawing

	buf1 := buf3.GetBuffer()
	defer buf3.PutBuffer(buf1)
	buf2 := buf3.GetBuffer()
	defer buf3.PutBuffer(buf2)
	defer func() { buf.WriteString(buf1.String()) }() // Make sure we write out the content of buf 1 even on an error

	if err = t2.DrawCaption(ctx, buf1); err != nil {
		return
	}

	if t.renderColumnTags {
		if err = t.drawColumnTags(ctx, buf1); err != nil {
			return
		}
	}

	if t.headerRowCount > 0 {
		err = t.drawHeaderRows(ctx, buf2)
		buf1.WriteString(html.RenderTag("thead", nil, buf2.String()))
		if err != nil {
			return
		}
		buf2.Reset()
	}

	if t.footerRowCount > 0 {
		err = t.drawFooterRows(ctx, buf2)
		buf1.WriteString(html.RenderTag("tfoot", nil, buf2.String()))
		if err != nil {
			return
		}
		buf2.Reset()
	}

	t.RangeData(func(index int, value interface{}) bool {
		err = t.drawRow(ctx, index, value, buf2)
		if err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return
	}

	buf1.WriteString(html.RenderTag("tbody", nil, buf2.String()))
	return nil
}

// DrawCaption is called internally to draw the caption. Subclasses can override this to draw a custom caption.
func (t *Table) DrawCaption(ctx context.Context, buf *bytes.Buffer) (err error) {
	switch obj := t.caption.(type) {
	case string:
		buf.WriteString(fmt.Sprintf("<caption>%s</caption>\n", html2.EscapeString(obj)))
	case page.ControlI:
		buf.WriteString("<caption>")
		err = obj.Draw(ctx, buf)
		if err != nil {
			buf.WriteString("</caption>\n")
		}
	}
	return
}

func (t *Table) drawColumnTags(ctx context.Context, buf *bytes.Buffer) (err error) {
	var colNum int
	var colCount = len(t.columns)

	for colNum < colCount {
		col := t.columns[colNum]
		if !col.IsHidden() {
			col.DrawColumnTag(ctx, buf)
		}
		colNum += col.Span()
	}
	return
}

func (t *Table) drawHeaderRows(ctx context.Context, buf *bytes.Buffer) (err error) {
	var t2 = t.this().(TableI) // Get the sub class so we call into its hooks for drawing

	buf1 := buf3.GetBuffer()
	defer buf3.PutBuffer(buf1)
	for rowNum := 0; rowNum < t.headerRowCount; rowNum++ {
		for colNum, col := range t.columns {
			if !col.IsHidden() {
				cellHtml, attr := t2.HeaderCellDrawingInfo(ctx, col, rowNum, colNum)
				buf1.WriteString(html.RenderTag("th", attr, cellHtml))
			}
		}
		buf.WriteString(html.RenderTag("tr", t.GetHeaderRowAttributes(rowNum), buf1.String()))
		buf1.Reset()
	}
	return
}

// HeaderCellDrawingInfo is called internally to provide the info for each header cell drawn. Subclasses can
// override this.
func (t *Table) HeaderCellDrawingInfo(ctx context.Context, col ColumnI, rowNum int, colNum int) (cellHtml string, cellAttributes *html.Attributes) {
	cellHtml = col.HeaderCellHtml(ctx, rowNum, colNum)
	cellAttributes = col.HeaderAttributes(ctx, rowNum, colNum)
	return
}

// FooterCellDrawingInfo is called internally to provide the info for each header cell drawn. Subclasses can
// override this.
func (t *Table) FooterCellDrawingInfo(ctx context.Context, col ColumnI, rowNum int, colNum int) (cellHtml string, cellAttributes *html.Attributes) {
	cellHtml = col.FooterCellHtml(ctx, rowNum, colNum)
	cellAttributes = col.FooterAttributes(ctx, rowNum, colNum)
	return
}


// GetHeaderRowAttributes is called internally to get the attributes for the tr tags in header rows.
func (t *Table) GetHeaderRowAttributes(row int) *html.Attributes {
	if t.headerRowStyler != nil {
		return t.headerRowStyler.TableHeaderRowAttributes(row)
	}
	return nil
}

func (t *Table) drawFooterRows(ctx context.Context, buf *bytes.Buffer) (err error) {
	var t2 = t.this().(TableI) // Get the sub class so we call into its hooks for drawing

	buf1 := buf3.GetBuffer()
	defer buf3.PutBuffer(buf1)
	for rowNum := 0; rowNum < t.footerRowCount; rowNum++ {
		for colNum, col := range t.columns {
			if !col.IsHidden() {
				cellHtml, attr := t2.FooterCellDrawingInfo(ctx, col, rowNum, colNum)
				tag := "td"
				if col.AsHeader() {
					tag = "th"
				}
				buf1.WriteString(html.RenderTag(tag, attr, cellHtml))
			}
		}
		buf.WriteString(html.RenderTag("tr", t.GetFooterRowAttributes(rowNum), buf1.String()))
		buf1.Reset()
	}
	return
}

// GetFooterRowAttributes is called internally to get the attributes for the tr tags in footer rows.
func (t *Table) GetFooterRowAttributes(row int) *html.Attributes {
	if t.footerRowStyler != nil {
		return t.footerRowStyler.TableFooterRowAttributes(row)
	}
	return nil
}

func (t *Table) drawRow(ctx context.Context, row int, data interface{}, buf *bytes.Buffer) (err error) {
	var this = t.this().(TableI) // Get the sub class so we call into its hooks for drawing
	buf1 := buf3.GetBuffer()
	defer buf3.PutBuffer(buf1)
	for i, col := range t.columns {
		col.DrawCell(ctx, row, i, data, buf1)
	}
	buf.WriteString(html.RenderTag("tr", this.GetRowAttributes(row, data), buf1.String()))
	return
}

// GetRowAttributes is used internally to return the attributes for the tr tag of a data row.
func (t *Table) GetRowAttributes(row int, data interface{}) *html.Attributes {
	if t.rowStyler != nil {
		return t.rowStyler.TableRowAttributes(row, data)
	}
	return nil
}

// AddColumnAt adds the given column at the column offset given. 0 is the first column location.
// Pass a negative number or a number bigger than the current number of columns to put it at the end.
func (t *Table) AddColumnAt(column ColumnI, loc int) {
	t.columnIdCounter++
	column.setParentTable(t)
	if column.ID() == "" {
		column.SetID(strconv.Itoa(t.columnIdCounter))
	}
	if loc < 0 || loc >= len(t.columns) {
		t.columns = append(t.columns, column)
	} else {
		t.columns = append(t.columns, nil)
		copy(t.columns[loc+1:], t.columns[loc:])
		t.columns[loc] = column
	}
	column.AddActions(t)

	t.Refresh()
}

// AddColumn adds the given column to the end of the column list.
func (t *Table) AddColumn(column ColumnI) ColumnI {
	t.AddColumnAt(column, -1)
	return column
}

// GetColumn returns the column at the given location
func (t *Table) GetColumn(loc int) ColumnI {
	return t.columns[loc]
}

// GetColumnByID returns the column with the given id.
func (t *Table) GetColumnByID(id string) ColumnI {
	for _, col := range t.columns {
		if col.ID() == id {
			return col
		}
	}
	return nil
}

// GetColumnByTitle returns the column with the given title.
func (t *Table) GetColumnByTitle(title string) ColumnI {
	for _, col := range t.columns {
		if col.Title() == title {
			return col
		}
	}
	return nil
}

// RemoveColumn removes the column at the given location.
// If you attempt to remove a column that does not exist, it will panic.
func (t *Table) RemoveColumn(loc int) {
	copy(t.columns[loc:], t.columns[loc+1:])
	t.columns[len(t.columns)-1] = nil // or the zero value of T
	t.columns = t.columns[:len(t.columns)-1]
	t.Refresh()
}

// RemoveColumnByID removes the column with the given id. If the column does not exist, nothing will change.
func (t *Table) RemoveColumnByID(id string) {
	for i, col := range t.columns {
		if col.ID() == id {
			t.RemoveColumn(i)
			t.Refresh()
			return
		}
	}
}

// RemoveColumnByTitle removes the given column with the given title. If the column does not exist, nothing will change.
func (t *Table) RemoveColumnByTitle(title string) {
	for i, col := range t.columns {
		if col.Title() == title {
			t.RemoveColumn(i)
			t.Refresh()
			return
		}
	}
}

// ClearColumns removes all of the columns.
func (t *Table) ClearColumns() {
	if len(t.columns) > 0 {
		t.columns = []ColumnI{}
		t.Refresh()
	}
}

// HideColumns hides all of the columns, keeping them in the column list, but causing them not to draw.
func (t *Table) HideColumns() {
	for _, col := range t.columns {
		col.SetHidden(true)
	}
	t.Refresh()
}

// ShowColumns sets all of the columns to be shown.
func (t *Table) ShowColumns() {
	for _, col := range t.columns {
		col.SetHidden(false)
	}
	t.Refresh()
}

// SetRowStyler sets a styler that returns attributes to be used on a particular row.
func (t *Table) SetRowStyler(a TableRowAttributer) TableI {
	t.rowStyler = a
	return t.this()
}

func (t *Table) RowStyler() TableRowAttributer {
	return t.rowStyler
}

// SetHeaderRowStyler sets a styler that returns attributes to be used on a particular header row.
func (t *Table) SetHeaderRowStyler(a TableHeaderRowAttributer) TableI {
	t.headerRowStyler = a
	return t.this()
}

// SetFooterRowStyler sets a styler that returns attributes to be used on a particular footer row.
func (t *Table) SetFooterRowStyler(a TableFooterRowAttributer) TableI {
	t.footerRowStyler = a
	return t.this()
}

// ΩUpdateFormValues is called by the system whenever values are sent by client controls. We forward that to the columns.
func (t *Table) ΩUpdateFormValues(ctx *page.Context) {
	for _, col := range t.columns {
		col.UpdateFormValues(ctx)
	}
}

// PrivateAction is called by the framework to allow controls to process actions internal to themselves.
func (t *Table) PrivateAction(ctx context.Context, p page.ActionParams) {
	switch p.ID {
	case ColumnAction:
		var subId string
		var a action.CallbackActionI
		var ok bool
		if a, ok = p.Action.(action.CallbackActionI); !ok {
			panic("Column actions must be a callback action")
		}
		if subId = a.GetDestinationControlSubID(); subId == "" {
			panic("Column actions must be a callback action")
		}
		c := t.GetColumnByID(subId)
		if c != nil {
			c.Action(ctx, p)
		}
	case SortClick:
		t.sortClick(p.EventValueString())
		t.Refresh()
	}

}

// SetSortHistoryLimit sets the number of columns that the table will remember for the sort history. It defaults to 1,
// meaning it will remember only the current column. Setting it more than 1 will let the system report back on secondary
// sort columns that the user chose. For example, if the user clicks to sort a first name column, and then a last name column,
// it will let you know to sort by last name, and then first name.
func (t *Table) SetSortHistoryLimit(n int) TableI {
	t.sortHistoryLimit = n
	t.Refresh()
	return t.this()
}

func (t *Table) sortClick(id string) {
	var foundLoc = -1
	var firstCol ColumnI

	if t.sortColumns != nil {
		firstCol = t.GetColumnByID(t.sortColumns[0])
		if firstCol.SortDirection() == NotSortable {
			return
		}
	}

	if t.sortColumns != nil {
		// If the column clicked is already the first one in the list, just change direction
		if t.sortColumns[0] == id {
			firstCol.SetSortDirection(firstCol.SortDirection() * -1)
			return
		}

		firstCol.SetSortDirection(NotSorted) // tell the first one in the list to not be sorted

		// remove the column from the sort list if it is there
		for i := 0; i < len(t.sortColumns); i++ {
			if t.sortColumns[i] == id {
				foundLoc = i
				break
			}
		}

		if foundLoc != -1 {
			t.sortColumns = append(t.sortColumns[:foundLoc], t.sortColumns[foundLoc+1:]...)
		}
	}

	//push front
	t.sortColumns = append([]string{id}, t.sortColumns...)
	col := t.GetColumnByID(id)
	col.SetSortDirection(SortAscending) // start out ascending

	//remove back
	if len(t.sortColumns) > t.sortHistoryLimit {
		t.sortColumns = t.sortColumns[:len(t.sortColumns)-1]
	}
}

// SortColumns returns a slice of columns in sort order
func (t *Table) SortColumns() (ret []ColumnI) {
	for _, id := range t.sortColumns {
		if col := t.GetColumnByID(id); col != nil {
			ret = append(ret, col)
		}
	}
	return ret
}

// ΩMarshalState is an internal function to save the state of the control
func (t *Table) ΩMarshalState(m maps.Setter) {
	m.Set("sortColumns", t.sortColumns)
	for _, col := range t.columns {
		col.MarshalState(m)
	}
}

// ΩUnmarshalState is an internal function to restore the state of the control
func (t *Table) ΩUnmarshalState(m maps.Loader) {
	if v, ok := m.Load("sortColumns"); ok {
		if s, ok := v.([]string); ok {
			t.sortColumns = s
		}
	}
	for _, col := range t.columns {
		col.UnmarshalState(m)
	}
}


// TableCreator is the initialization structure for declarative creation of tables
type TableCreator struct {
	// ID is the control id
	ID               string
	// HasColTags will make the table render <col> tags
	HasColTags       bool
	// Caption is the content of the caption tag, and can either be a string, or a data pager
	Caption          interface{}
	// HideIfEmpty will hide the table completely if it has no data. Otherwise, the table and headers will be shown, but no data rows
	HideIfEmpty      bool
	// HeaderRowCount is the number of header rows. You must set this to at least 1 to show header rows.
	HeaderRowCount   int
	// FooterRowCount is the number of footer rows.
	FooterRowCount   int
	// RowStyler returns the attributes to be used in a cell. It can be either a control id or a TableRowAttributer.
	RowStyler        interface{}
	// HeaderRowStyler returns the attributes to be used in a header cell. It can be either a control id or a TableHeaderRowAttributer.
	HeaderRowStyler  interface{}
	// FooterRowStyler returns the attributes to be used in a footer cell. It can be either a control id or a TableFooterRowAttributer.
	FooterRowStyler  interface{}
	// Columns are the column creators that will add columns to the table
	Columns          []ColumnCreator
	// DataProvider is the control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProvider data.DataBinder
	// DataProviderID is the id of a control that will dynamically provide the data for the list and that implements the DataBinder interface.
	DataProviderID string
	// Data is the actual data for the table, and should be a slice of objects
	Data             interface{}
	// Sortable will make the table sortable
	Sortable         bool
	// SortHistoryLimit will set how many columns deep we will remember the sorting for multi-level sorts
	SortHistoryLimit int
	page.ControlOptions
}



// Create is called by the framework to create a new control from the Creator. You
// do not normally need to call this.
func (c TableCreator) Create(ctx context.Context, parent page.ControlI) page.ControlI {
	ctrl := NewTable(parent, c.ID)
	c.Init(ctx, ctrl)
	return ctrl
}

// Init is called by implementations of Buttons to initialize a control with the
// creator. You do not normally need to call this.
func (c TableCreator) Init(ctx context.Context, ctrl TableI) {
	ctrl.SetRenderColumnTags(c.HasColTags)
	ctrl.SetHideIfEmpty(c.HideIfEmpty)
	if c.Caption != nil {
		if ctrl2, ok := c.Caption.(page.Creator); ok {
			ctrl.SetCaption(ctrl2.Create(ctx, ctrl))
		} else {
			ctrl.SetCaption(c.Caption)
		}
	}
	if c.HeaderRowCount > 0 {
		ctrl.SetHeaderRowCount(c.HeaderRowCount)
	}
	if c.FooterRowCount > 0 {
		ctrl.SetFooterRowCount(c.FooterRowCount)
	}

	if c.RowStyler != nil {
		if s,ok := c.RowStyler.(string); ok {
			ctrl.SetRowStyler(ctrl.Page().GetControl(s).(TableRowAttributer))
		} else {
			ctrl.SetRowStyler(c.RowStyler.(TableRowAttributer))
		}
	}
	if c.HeaderRowStyler != nil {
		if s,ok := c.HeaderRowStyler.(string); ok {
			ctrl.SetHeaderRowStyler(ctrl.Page().GetControl(s).(TableHeaderRowAttributer))
		} else {
			ctrl.SetHeaderRowStyler(c.HeaderRowStyler.(TableHeaderRowAttributer))
		}
	}

	if c.FooterRowStyler != nil {
		if s,ok := c.FooterRowStyler.(string); ok {
			ctrl.SetFooterRowStyler(ctrl.Page().GetControl(s).(TableFooterRowAttributer))
		} else {
			ctrl.SetFooterRowStyler(c.FooterRowStyler.(TableFooterRowAttributer))
		}
	}

	if c.DataProvider != nil {
		ctrl.SetDataProvider(c.DataProvider)
	} else if c.DataProviderID != "" {
		provider := ctrl.Page().GetControl(c.DataProviderID).(data.DataBinder)
		ctrl.SetDataProvider(provider)
	}

	if c.Data != nil {
		ctrl.SetData(c.Data)
	}
	if c.Columns != nil {
		for _,colCreator := range c.Columns {
			ctrl.AddColumn(colCreator.Create(ctx, ctrl))
		}
	}
	if c.Sortable {
		ctrl.MakeSortable()
	}
	if c.SortHistoryLimit > 0 {
		ctrl.SetSortHistoryLimit(c.SortHistoryLimit)
	}

	ctrl.ApplyOptions(c.ControlOptions)
}

// GetTable is a convenience method to return the button with the given id from the page.
func GetTable(c page.ControlI, id string) *Table {
	return c.Page().GetControl(id).(*Table)
}
