package  column

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/html"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlias_Attributes(t *testing.T) {
	col := NewAliasColumn("Alias")
	col.ColTagAttributes().Set("a","b")
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))
}


func TestAliasColumn_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	f := control.NewMockForm()

	f.AddControls(context.Background(),
		control.TableCreator{
			ID:                "table",
			Columns: control.Columns(
				AliasColumnCreator{
					ID:            "a1",
					Alias:         "DataOption",
					Title:         "Alias",
					Sortable:      true,
					ColumnOptions: control.ColumnOptions{
						CellAttributes:   nil,
						HeaderAttributes: nil,
						FooterAttributes: nil,
						ColTagAttributes: html.Attributes{
							"a":"b",
						},
						Span:             0,
						AsHeader:         false,
						IsHtml:           false,
						HeaderTexter:     nil,
						FooterTexter:     nil,
						IsHidden:         true,
						Format:           "",
						TimeFormat:       "",
					},
				},
			),
		},
	)


	c := control.GetTable(f, "table")
	col := c.GetColumnByID("a1")
	assert.Equal(t, "Alias", col.Title())
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))

	c.Serialize(enc)

	c2 := control.Table{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	col = c2.GetColumnByID("a1")

	assert.True(t, col.IsHidden())
	assert.Equal(t, "Alias", col.Title())
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))
}
