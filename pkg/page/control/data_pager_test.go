package control

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/page"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataPager_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	p := page.NewMockForm()

	p.AddControls(context.Background(),
		PagedTableCreator{
			ID:"table",
		},
		DataPagerCreator{
			ID:               "dp",
			MaxPageButtons:   11,
			ObjectName:       "Thing",
			ObjectPluralName: "Things",
			LabelForNext:     "Next Thing",
			LabelForPrevious: "Previous Thing",
			PagedControl:     "table",
		},
	)

	enc.Encode(p.Page())

	p2 := page.Page{}
	dec := gob.NewDecoder(&buf)
	dec.Decode(&p2)
	c2 := GetDataPager(p2.Form(), "dp").(*DataPager)

	assert.Equal(t, 11, c2.maxPageButtons)
	assert.Equal(t, "Thing", c2.ObjectName)
	assert.Equal(t, "Previous Thing", c2.LabelForPrevious)
}
