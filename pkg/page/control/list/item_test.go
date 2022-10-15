package list

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListItem_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	li := NewItem("test", "5")
	li.Serialize(enc)

	li2 := Item{}
	dec := gob.NewDecoder(&buf)
	li2.Deserialize(dec)

	assert.Equal(t, 5, li2.IntValue())
	assert.Equal(t, "test", li2.Label())
}
