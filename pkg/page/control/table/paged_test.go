package table

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/page/control"
	"testing"

	"github.com/goradd/html5tag"
	"github.com/stretchr/testify/assert"
)

type pagedTableTestForm struct {
	control.FormBase
}

func (f *pagedTableTestForm) Init(ctx context.Context, id string) {
	f.FormBase.Init(f, ctx, id)
}

func (*pagedTableTestForm) RowAttributes(row int, data interface{}) html5tag.Attributes {
	return html5tag.NewAttributes().AddValues("a", "b")
}

func (*pagedTableTestForm) HeaderRowAttributes(row int) html5tag.Attributes {
	return html5tag.NewAttributes().AddValues("c", "d")
}

func (*pagedTableTestForm) FooterRowAttributes(row int) html5tag.Attributes {
	return html5tag.NewAttributes().AddValues("e", "f")
}

func (*pagedTableTestForm) BindData(ctx context.Context, s control.DataManagerI) {

}

func TestPagedTable_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	f := new(pagedTableTestForm)
	f.Init(context.Background(), "MockFormId")

	f.AddControls(context.Background(),
		PagedTableCreator{
			ID:               "table",
			Caption:          "This is a table",
			HideIfEmpty:      true,
			HeaderRowCount:   2,
			FooterRowCount:   3,
			RowStylerID:      f.ID(),
			HeaderRowStyler:  f,
			FooterRowStyler:  f,
			DataProvider:     f,
			Sortable:         true,
			SortHistoryLimit: 3,
			OnCellClick:      nil,
			PageSize:         7,
			SaveState:        false, // must have a session to test
			Columns:          nil,   // testing columns here will cause circular import
		},
		control.DataPagerCreator{
			ID:             "dp",
			PagedControlID: "table",
		},
	)

	c := GetPagedTable(f, "table")

	c.Serialize(enc)

	c2 := PagedTable{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	assert.Equal(t, "This is a table", c2.caption)
	assert.Equal(t, 3, c2.sortHistoryLimit)
}
