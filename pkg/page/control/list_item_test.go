package control

import (
	"bytes"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListItem_Serialize(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	li := NewListItem("test", "5")
	li.Serialize(enc)

	li2 := ListItem{}
	dec := gob.NewDecoder(&buf)
	li2.Deserialize(dec)

	assert.Equal(t, 5, li2.IntValue())
	assert.Equal(t, "test", li2.Label())
}
