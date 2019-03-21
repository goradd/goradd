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
	html2 "html"
	"strconv"
	"github.com/goradd/goradd/pkg/page/event"
	"github.com/goradd/goradd/pkg/page/control/data"
)

const (
	ColumnAction = iota + 2000
	SortClick
)

type TableI interface {
	page.ControlI
	SetCaption(interface{})
	DrawCaption(context.Context, *bytes.Buffer) error
	GetHeaderRowAttributes(row int) *html.Attributes
	GetFooterRowAttributes(row int) *html.Attributes
	GetRowAttributes(row int, data interface{}) *html.Attributes
	HeaderCellDrawingInfo(ctx context.Context, col ColumnI, rowNum int, colNum int) (cellHtml string, cellAttributes *html.Attributes)
}

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
	rowStyler             html.Attributer
	headerRowStyler       html.Attributer
	footerRowStyler       html.Attributer
	columnIdCounter       int

	// Sort info. Sorting is difficult enough, and intertwined with tables enough, that we just make it built in to every column
	sortColumns      []string // keeps a historical list of columns sorted on
	sortHistoryLimit int      // how far back to go
}

func NewTable(parent page.ControlI, id string) *Table {
	t := &Table{}
	t.Init(t, parent, id)
	return t
}

func (t *Table) Init(self page.ControlI, parent page.ControlI, id string) {
	t.Control.Init(self, parent, id)
	t.Tag = "table"
	t.columns = []ColumnI{}
	t.sortHistoryLimit = 1
}

func (t *Table) this() TableI {
	return t.Self.(TableI)
}

// Call Sortable to make a table sortable. It will attach sortable events and show the header if its not shown.
func (t *Table) Sortable() TableI {
	t.On(event.TableSort(), action.Ajax(t.ID(), SortClick), action.PrivateAction{})
	if t.headerRowCount == 0 {
		t.headerRowCount = 1
	}
	return t.this()
}

func (t *Table) SetCaption(caption interface{}) {
	t.caption = caption
}

func (t *Table) SetHeaderRowCount(count int) TableI {
	t.headerRowCount = count
	return t.this()
}

func (t *Table) SetFooterRowCount(count int) TableI {
	t.footerRowCount = count
	return t.this()
}

func (t *Table) ΩDrawTag(ctx context.Context) string {
	log.FrameworkDebug("Drawing table tag")
	if t.HasDataProvider() {
		log.FrameworkDebug("Getting table data")
		t.GetData(ctx, t)
		defer t.ResetData()
	}
	for _,c := range t.columns {
		c.PreRender()
	}
	return t.Control.ΩDrawTag(ctx)
}

func (t *Table) ΩDrawingAttributes() *html.Attributes {
	a := t.Control.ΩDrawingAttributes()
	a.SetDataAttribute("grctl", "table")
	if t.Data == nil {
		a.SetStyle("display", "none")
	}
	return a
}

func (t *Table) ΩDrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	var t2 = t.this().(TableI) // Get the sub class so we call into its hooks for drawing

	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf1)
	buf2 := page.GetBuffer()
	defer page.PutBuffer(buf2)
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

	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf1)
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

func (t *Table) HeaderCellDrawingInfo(ctx context.Context, col ColumnI, rowNum int, colNum int) (cellHtml string, cellAttributes *html.Attributes) {
	cellHtml = col.HeaderCellHtml(ctx, rowNum, colNum)
	cellAttributes = col.HeaderAttributes(rowNum, colNum)
	return
}

func (t *Table) GetHeaderRowAttributes(row int) *html.Attributes {
	if t.headerRowStyler != nil {
		return t.headerRowStyler.Attributes(row)
	}
	return nil
}

func (t *Table) drawFooterRows(ctx context.Context, buf *bytes.Buffer) (err error) {
	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf1)
	for j := 0; j < t.footerRowCount; j++ {
		for i, col := range t.columns {
			col.DrawFooterCell(ctx, j, i, t.footerRowCount, buf1)
		}
		buf.WriteString(html.RenderTag("tr", t.GetFooterRowAttributes(j), buf1.String()))
		buf1.Reset()
	}
	return
}

func (t *Table) GetFooterRowAttributes(row int) *html.Attributes {
	if t.footerRowStyler != nil {
		return t.footerRowStyler.Attributes(row)
	}
	return nil
}

func (t *Table) drawRow(ctx context.Context, row int, data interface{}, buf *bytes.Buffer) (err error) {
	var this = t.this().(TableI) // Get the sub class so we call into its hooks for drawing
	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf1)
	for i, col := range t.columns {
		col.DrawCell(ctx, row, i, data, buf1)
	}
	buf.WriteString(html.RenderTag("tr", this.GetRowAttributes(row, data), buf1.String()))
	return
}

func (t *Table) GetRowAttributes(row int, data interface{}) *html.Attributes {
	if t.rowStyler != nil {
		return t.rowStyler.Attributes(row, data)
	}
	return nil
}

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

func (t *Table) AddColumn(column ColumnI) ColumnI {
	t.AddColumnAt(column, -1)
	return column
}

func (t *Table) GetColumn(loc int) ColumnI {
	return t.columns[loc]
}

func (t *Table) GetColumnByID(id string) ColumnI {
	for _, col := range t.columns {
		if col.ID() == id {
			return col
		}
	}
	return nil
}

func (t *Table) GetColumnByTitle(title string) ColumnI {
	for _, col := range t.columns {
		if col.Title() == title {
			return col
		}
	}
	return nil
}

func (t *Table) RemoveColumn(loc int) {
	copy(t.columns[loc:], t.columns[loc+1:])
	t.columns[len(t.columns)-1] = nil // or the zero value of T
	t.columns = t.columns[:len(t.columns)-1]
	t.Refresh()
}

func (t *Table) RemoveColumnByID(id string) {
	for i, col := range t.columns {
		if col.ID() == id {
			t.RemoveColumn(i)
			t.Refresh()
			return
		}
	}
}

func (t *Table) RemoveColumnByTitle(title string) {
	for i, col := range t.columns {
		if col.Title() == title {
			t.RemoveColumn(i)
			t.Refresh()
			return
		}
	}
}

func (t *Table) ClearColumns() {
	if len(t.columns) > 0 {
		t.columns = []ColumnI{}
		t.Refresh()
	}
}

func (t *Table) HideColumns() {
	for _, col := range t.columns {
		col.SetHidden(true)
	}
	t.Refresh()
}

func (t *Table) ShowColumns() {
	for _, col := range t.columns {
		col.SetHidden(false)
	}
	t.Refresh()
}

func (t *Table) SetRowStyler(a html.Attributer) {
	t.rowStyler = a
}

func (t *Table) RowStyler() html.Attributer {
	return t.rowStyler
}


func (t *Table) SetHeaderRowStyler(a html.Attributer) {
	t.rowStyler = a
}

func (t *Table) SetFooterRowStyler(a html.Attributer) {
	t.rowStyler = a
}

// ΩUpdateFormValues is called by the system whenever values are sent by client controls. We forward that to the columns.
func (t *Table) ΩUpdateFormValues(ctx *page.Context) {
	for _, col := range t.columns {
		col.UpdateFormValues(ctx)
	}
}

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
func (t *Table) SetSortHistoryLimit(n int) {
	t.sortHistoryLimit = n
	t.Refresh()
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
	for _,col := range t.columns {
		col.MarshalState(m)
	}
}

// ΩUnmarshalState is an internal function to restore the state of the control
func (t *Table) ΩUnmarshalState(m maps.Loader) {
	if v,ok := m.Load("sortColumns"); ok {
		if s, ok := v.([]string); ok {
			t.sortColumns = s
		}
	}
	for _,col := range t.columns {
		col.UnmarshalState(m)
	}
}

