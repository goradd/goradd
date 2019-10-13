package control

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataPager_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	p := NewMockForm()

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

	c := GetDataPager(p, "dp")

	c.Serialize(enc)

	c2 := DataPager{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	assert.Equal(t, 11, c2.maxPageButtons)
	assert.Equal(t, "Thing", c2.ObjectName)
}
