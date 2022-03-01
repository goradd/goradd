package control

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/goradd/goradd/pkg/page"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListSelectString(t *testing.T) {
	p := page.NewMockForm()

	d := NewSelectList(p, "")

	d.AddItem("A", "A")
	d.AddItem("D")
	d.AddItemAt(1, "B", "B")
	d.AddItemAt(-1, "C")
	d.AddItemAt(-10, "- Select a Value -", "")

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

	assert.Equal(t, "C", d.GetItemAt(3).Value())
	assert.Equal(t, "D", d.GetItemAt(4).Value())

	d.SetIsRequired(true)
	valid = d.MockFormValue("")
	assert.False(t, valid)
}

func TestListSelectInt(t *testing.T) {
	p := page.NewMockForm()

	d := NewSelectList(p, "")

	d.AddItem("- Select a Value -", "")
	d.AddItem("A", "1")
	d.AddItem("C", "3")
	d.AddItemAt(2, "B", "2")

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

// This exercises more of the ItemList mixin.
func TestListSelectData(t *testing.T) {
	p := page.NewMockForm()

	d := NewSelectList(p, "")

	d.SetData([]ListValue{
		{"- Select a Value -", nil},
		{"A", 1},
		{"B", 2},
		{"D", 4},
	})
	d.AddItemAt(3, "C", "3")

	d.SetValue(2)
	assert.Equal(t, "B", d.SelectedLabel())
	assert.Equal(t, 2, d.IntValue())
	assert.Equal(t, 2, d.SelectedItem().IntValue())
	assert.Equal(t, 1, d.GetItemAt(1).IntValue())
	assert.Nil(t, d.GetItemAt(7))
	assert.Equal(t, 4, d.ListItems()[4].IntValue())

	valid := d.MockFormValue("3")
	assert.True(t, valid)
	assert.Equal(t, "C", d.SelectedLabel())
	assert.Equal(t, "3", d.Value())
	assert.Equal(t, 3, d.IntValue())

	d.SetIsRequired(true)
	valid = d.MockFormValue("")
	assert.False(t, valid)

	assert.Panics(t, func() { d.RemoveItemAt(7) })
	d.RemoveItemAt(2)
	assert.Equal(t, 4, d.Len())
	assert.Equal(t, 4, d.GetItemAt(3).IntValue())

}

// Test the list item mixin through the SelectList
func TestListItem(t *testing.T) {

}

func TestSelectList_Serialize(t *testing.T) {
	p := page.NewMockForm()

	p.AddControls(context.Background(),
		SelectListCreator{
			ID: "c",
			Items: []ListValue {
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

