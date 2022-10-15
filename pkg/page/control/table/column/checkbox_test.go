package column

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/page/control/table"
	"testing"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/html5tag"
	"github.com/stretchr/testify/assert"
)

func TestCheckbox_Attributes(t *testing.T) {
	col := NewCheckboxColumn(DefaultCheckboxProvider{})
	col.ColTagAttributes().Set("a", "b")
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))
}

func TestCheckboxColumn_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	f := page.NewMockForm()

	f.AddControls(context.Background(),
		table.TableCreator{
			ID: "table",
			Columns: table.Columns(
				CheckboxColumnCreator{
					ID:               "a1",
					CheckboxProvider: DefaultCheckboxProvider{},
					Title:            "CheckboxList",
					Sortable:         true,
					ColumnOptions: table.ColumnOptions{
						CellAttributes:   nil,
						HeaderAttributes: nil,
						FooterAttributes: nil,
						ColTagAttributes: html5tag.Attributes{
							"a": "b",
						},
						Span:         0,
						AsHeader:     false,
						IsHtml:       false,
						HeaderTexter: nil,
						FooterTexter: nil,
						IsHidden:     true,
						Format:       "",
						TimeFormat:   "",
					},
				},
			),
		},
	)

	c := table.GetTable(f, "table")
	col := c.GetColumnByID("a1")
	assert.Equal(t, "CheckboxList", col.Title())
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))

	c.Serialize(enc)

	c2 := table.Table{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	col = c2.GetColumnByID("a1")

	assert.True(t, col.IsHidden())
	assert.Equal(t, "CheckboxList", col.Title())
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))
}
