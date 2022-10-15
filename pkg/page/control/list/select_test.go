package list

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"github.com/goradd/goradd/pkg/page"

	"github.com/stretchr/testify/assert"
)

func TestListSelectString(t *testing.T) {
	p := page.NewMockForm()

	d := NewSelectList(p, "")

	d.Add("A", "A")
	d.Add("D")
	d.AddAt(1, "B", "B")
	d.AddAt(-1, "C")
	d.AddAt(-10, "- Select a Value -", "")

	d.SetValue("B")

	assert.Equal(t, "B", d.SelectedLabel())
	assert.Equal(t, "B", d.Value())
	assert.Equal(t, "B", d.SelectedItem().Value())

	d.SetValue("D")
	assert.Equal(t, "D", d.Value())
	assert.Equal(t, "D", d.SelectedLabel())

	valid := d.MockFormValue("C")
	assert.True(t, valid)
	assert.Equal(t, "C", d.SelectedLabel())
	assert.Equal(t, "C", d.Value())
	assert.Equal(t, "C", d.StringValue())

	assert.Equal(t, "C", d.ItemAt(3).Value())
	assert.Equal(t, "D", d.ItemAt(4).Value())

	d.SetIsRequired(true)
	valid = d.MockFormValue("")
	assert.False(t, valid)
}

func TestListSelectInt(t *testing.T) {
	p := page.NewMockForm()

	d := NewSelectList(p, "")

	d.Add("- Select a Value -", "")
	d.Add("A", "1")
	d.Add("C", "3")
	d.AddAt(2, "B", "2")

	d.SetValue(2)
	assert.Equal(t, "B", d.SelectedLabel())
	assert.Equal(t, 2, d.IntValue())
	assert.Equal(t, 2, d.SelectedItem().IntValue())

	valid := d.MockFormValue("3")
	assert.True(t, valid)
	assert.Equal(t, "C", d.SelectedLabel())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 3, d.IntValue())

	d.SetIsRequired(true)
	valid = d.MockFormValue("")
	assert.False(t, valid)
}

// This exercises more of the List mixin.
func TestListSelectData(t *testing.T) {
	p := page.NewMockForm()

	d := NewSelectList(p, "")

	d.SetData([]ListValue{
		{"- Select a Value -", nil},
		{"A", 1},
		{"B", 2},
		{"D", 4},
	})
	d.AddAt(3, "C", "3")

	d.SetValue(2)
	assert.Equal(t, "B", d.SelectedLabel())
	assert.Equal(t, 2, d.IntValue())
	assert.Equal(t, 2, d.SelectedItem().IntValue())
	assert.Equal(t, 1, d.ItemAt(1).IntValue())
	assert.Nil(t, d.ItemAt(7))
	assert.Equal(t, 4, d.Items()[4].IntValue())

	valid := d.MockFormValue("3")
	assert.True(t, valid)
	assert.Equal(t, "C", d.SelectedLabel())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 3, d.IntValue())

	d.SetIsRequired(true)
	valid = d.MockFormValue("")
	assert.False(t, valid)

	assert.Panics(t, func() { d.RemoveAt(7) })
	d.RemoveAt(2)
	assert.Equal(t, 4, d.Len())
	assert.Equal(t, 4, d.ItemAt(3).IntValue())

}

// Test the list item mixin through the SelectList
func TestListItem(t *testing.T) {

}

func TestSelectList_Serialize(t *testing.T) {
	p := page.NewMockForm()

	p.AddControls(context.Background(),
		SelectListCreator{
			ID: "c",
			Items: []ListValue{
				{"a", 1},
				{"b", 2},
			},
			Value: "2",
		},
	)
	c := GetSelectList(p, "c")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	c.Serialize(enc)

	c2 := SelectList{}
	dec := gob.NewDecoder(&buf)
	c2.Deserialize(dec)

	assert.Equal(t, "2", c2.Value())
}
