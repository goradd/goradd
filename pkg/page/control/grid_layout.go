package control

import "github.com/goradd/goradd/pkg/html"

// LayoutDirection controls whether items are layed out in rows or columns.
type LayoutDirection int

const (
	// LayoutRow lays out items in rows
	LayoutRow LayoutDirection = iota
	// LayoutColumn lays out items in columns, computing the number of rows required to make the specified number of columns.
	LayoutColumn
)

// GridLayoutBuilder is a helper that will allow a slice of items to be layed out in a table like
// pattern. It will compute the number of rows required, and then wrap the rows in
// row html, and the cells in cell html. You can have the items flow with the rows, or flow
// across the row axis. You can use this to build a table or a table-like structure.
type GridLayoutBuilder struct {
	items         []string
	columnCount   int
	direction     LayoutDirection
	rowTag        string
	rowAttributes html.Attributes
}

// Items sets the html for each item to display.
func (g *GridLayoutBuilder) Items(items []string) *GridLayoutBuilder {
	g.items = items
	return g
}

// ColumnCount sets the number of columns.
func (g *GridLayoutBuilder) ColumnCount(count int) *GridLayoutBuilder {
	g.columnCount = count
	return g
}

// LayoutDirection indicates how items are placed, whether they should fill up rows first, or fill up columns.
func (g *GridLayoutBuilder) Direction(placement LayoutDirection) *GridLayoutBuilder {
	g.direction = placement
	return g
}

func (g *GridLayoutBuilder) RowTag(t string) *GridLayoutBuilder {
	g.rowTag = t
	return g
}

func (g *GridLayoutBuilder) RowClass(t string) *GridLayoutBuilder {
	g.getRowAttributes().SetClass(t)
	return g
}

func (g *GridLayoutBuilder) getRowAttributes() html.Attributes {
	if g.rowAttributes == nil {
		g.rowAttributes = html.NewAttributes()
	}
	return g.rowAttributes
}

func (g *GridLayoutBuilder) Build() string {
	if len(g.items) == 0 {
		return ""
	}
	if g.rowTag == "" {
		g.rowTag = "div"
	}
	if g.columnCount == 0 {
		g.columnCount = 1
	}

	if g.direction == LayoutRow {
		return g.wrapRows()
	} else {
		return g.wrapColumns()
	}
}

func (g *GridLayoutBuilder) wrapRows() string {
	var rows string
	var row string
	for i := range g.items {
		row += g.items[i]
		if (i+1)%g.columnCount == 0 {
			rows += html.RenderTag(g.rowTag, g.rowAttributes, row)
			row = ""
		}
	}
	if row != "" {
		// partial row
		rows += html.RenderTag(g.rowTag, g.rowAttributes, row)
	}
	return rows
}

func (g *GridLayoutBuilder) wrapColumns() string {
	l := len(g.items)
	rowCount := (l-1)/g.columnCount + 1

	var row string
	var rows string

	for r := 0; r < rowCount; r++ {
		for c := 0; c < g.columnCount; c++ {
			i := c*rowCount + r
			if i < l {
				row += g.items[i]
			}
		}
		rows += html.RenderTag(g.rowTag, g.rowAttributes, row)
		row = ""
	}
	if row != "" {
		// partial row
		rows += html.RenderTag(g.rowTag, g.rowAttributes, row)
	}
	return rows
}
