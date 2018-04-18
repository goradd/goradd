package control

import (
	"github.com/spekary/goradd/page"
	"bytes"
	"context"
	"github.com/spekary/goradd/html"
	"fmt"
	html2 "golang.org/x/net/html"
	"github.com/spekary/goradd/page/control/column"
	"strconv"
)

type TableI interface {
	DrawCaption(context.Context, *bytes.Buffer) error
	GetHeaderRowAttributes(row int) *html.Attributes
	GetFooterRowAttributes(row int) *html.Attributes
	GetRowAttributes(row int, data interface{}) *html.Attributes
}

type Table struct {
	page.Control
	DataManager

	columns []column.ColumnI
	showHeader bool
	showFooter bool
	renderColumnTags bool
	caption string
	hideIfEmpty bool
	headerRowCount int
	footerRowCount int
	currentHeaderRowIndex int //??
	currentRowIndex int			//??
	rowStyler html.Attributer
	headerRowStyler html.Attributer
	footerRowStyler html.Attributer
	columnIdCounter int
}

func NewTable(parent page.ControlI) *Table {
	t := &Table{}
	t.Init(t, parent)
	return t
}


func (t *Table) Init(self page.ControlI, parent page.ControlI) {
	t.Control.Init(self, parent)
	t.Tag = "table"
	t.columns = []column.ColumnI{}
}

func (t *Table) DrawTag(ctx context.Context) string {
	t.GetData(t)
	defer t.ResetData()
	return t.Control.DrawTag(ctx)
}


func (t *Table) DrawingAttributes() *html.Attributes {
	a := t.Control.DrawingAttributes()
	t.SetDataAttribute("grctl", "table")
	if t.Data == nil {
		a.SetStyle("display", "none")
	}
	return a
}



func (t *Table) DrawInnerHtml(ctx context.Context, buf *bytes.Buffer) (err error) {
	var t2 = t.This().(TableI)	// Get the sub class so we call into its hooks for drawing

	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf1)
	buf2 := page.GetBuffer()
	defer page.PutBuffer(buf2)
	defer func() {buf.WriteString(buf1.String())}()	// Make sure we write out the content of buf 1 even on an error

	if err = t2.DrawCaption(ctx, buf1); err != nil {return}


	if t.renderColumnTags {
		if err = t.drawColumnTags(ctx, buf1); err != nil {return}
	}

	if t.showHeader {
		err = t.drawHeaderRows(ctx, buf2)
		buf1.WriteString(buf2.String())
		buf1.WriteString(html.RenderTag("thead", nil, buf2.String()))
		if err != nil {
			return
		}
		buf2.Reset()
	}

	if t.showFooter {
		err = t.drawFooterRows(ctx, buf2)
		buf1.WriteString(buf2.String())
		buf1.WriteString(html.RenderTag("tfoot", nil, buf2.String()))
		if err != nil {
			return
		}
		buf2.Reset()
	}

	if t.Data != nil && len (t.Data) > 0 {
		for i,row := range t.Data {
			err = t.drawRow(ctx, i, row, buf2)
			if err != nil {return}
		}
	}
	buf1.WriteString(html.RenderTag("tbody", nil, buf2.String()))
	return nil
}

func (t *Table) DrawCaption(ctx context.Context, buf *bytes.Buffer) (err error) {
	if t.caption != "" {
		buf.WriteString(fmt.Sprintf("<caption>%s</caption>\n", html2.EscapeString(t.caption)))
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
	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf)
	for j := 0; j < t.headerRowCount; j++ {
		for i,col := range t.columns {
			col.DrawHeaderCell(ctx, j, i, buf1)
		}
		buf.WriteString(html.RenderTag("tr", t.GetHeaderRowAttributes(j), buf1.String()))
	}
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
	defer page.PutBuffer(buf)
	for j := 0; j < t.footerRowCount; j++ {
		for i,col := range t.columns {
			col.DrawFooterCell(ctx, j, i, buf1)
		}
		buf.WriteString(html.RenderTag("tr", t.GetFooterRowAttributes(j), buf1.String()))
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
	var t2 = t.This().(TableI)	// Get the sub class so we call into its hooks for drawing
	buf1 := page.GetBuffer()
	defer page.PutBuffer(buf1)
	for i,col := range t.columns {
		col.DrawCell(ctx, row, i, data, buf1)
	}
	buf.WriteString(html.RenderTag("tr", t2.GetRowAttributes(row, data), buf1.String()))
	return
}

func (t *Table) GetRowAttributes(row int, data interface{}) *html.Attributes {
	if t.rowStyler != nil {
		return t.rowStyler.Attributes(row)
	}
	return nil
}

func (t *Table) AddColumnAt(column column.ColumnI, loc int) {
	t.columnIdCounter ++
	if column.Id() == "" {
		column.SetId(t.Id() + "_" + strconv.Itoa(t.columnIdCounter))
	}
	if loc < 0 || loc >= len(t.columns) {
		t.columns = append(t.columns, column)
	} else {
		t.columns = append(t.columns, nil)
		copy(t.columns[loc+1:], t.columns[loc:])
		t.columns[loc] = column
	}
	t.Refresh()
}

func (t *Table) AddColumn(column column.ColumnI) {
	t.AddColumnAt(column, -1)
}

func (t *Table) GetColumn(loc int) column.ColumnI {
	return t.columns[loc]
}

func (t *Table) GetColumnById(id string) column.ColumnI {
	for _,col := range t.columns {
		if col.Id() == id {
			return col
		}
	}
	return nil
}

func (t *Table) GetColumnByLabel(label string) column.ColumnI {
	for _,col := range t.columns {
		if col.Label() == label {
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

func (t *Table) RemoveColumnById(id string) {
	for i,col := range t.columns {
		if col.Id() == id {
			t.RemoveColumn(i)
			t.Refresh()
			return
		}
	}
}

func (t *Table) RemoveColumnByLabel(label string) {
	for i,col := range t.columns {
		if col.Label() == label {
			t.RemoveColumn(i)
			t.Refresh()
			return
		}
	}
}

func (t *Table) ClearColumns() {
	if len(t.columns) > 0 {
		t.columns = []column.ColumnI{}
		t.Refresh()
	}
}

func (t *Table) HideColumns() {
	for _,col := range t.columns {
		col.SetHidden(true)
	}
	t.Refresh()
}

func (t *Table) ShowColumns() {
	for _,col := range t.columns {
		col.SetHidden(false)
	}
	t.Refresh()
}

func (t *Table) SetRowStyler(a html.Attributer) {
	t.rowStyler = a
}

func (t *Table) SetHeaderRowStyler(a html.Attributer) {
	t.rowStyler = a
}

func (t *Table) SetFooterRowStyler(a html.Attributer) {
	t.rowStyler = a
}


// UpdateFormValues is called by the system whenever values are sent by client controls. We forward that to the columns.
func (t *Table) UpdateFormValues(ctx *page.Context) {
	for _,col := range t.columns {
		col.UpdateFormValues(ctx)
	}
}
