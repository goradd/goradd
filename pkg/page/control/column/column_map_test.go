package column

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"github.com/goradd/goradd/pkg/page"
	"github.com/goradd/goradd/pkg/page/control"
	"github.com/goradd/html5tag"
	"github.com/stretchr/testify/assert"
)

func TestMapColumn_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	f := page.NewMockForm()

	f.AddControls(context.Background(),
		control.TableCreator{
			ID: "table",
			Columns: control.Columns(
				MapColumnCreator{
					ID:       "a1",
					Index:    5,
					Title:    "StdMap",
					Sortable: true,
					ColumnOptions: control.ColumnOptions{
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

	c := control.GetTable(f, "table")
	col := c.GetColumnByID("a1")
	assert.Equal(t, "StdMap", col.Title())
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))

	c.Serialize(enc)

	c2 := control.Table{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	col = c2.GetColumnByID("a1")

	assert.True(t, col.IsHidden())
	assert.Equal(t, "StdMap", col.Title())
	assert.Equal(t, "b", col.ColTagAttributes().Get("a"))
	assert.Equal(t, 5, col.(*MapColumn).key)
}
