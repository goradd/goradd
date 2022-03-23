package control

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"github.com/goradd/goradd/pkg/page"

	"github.com/stretchr/testify/assert"
)

func TestCheckbox_Serialize(t *testing.T) {
	p := page.NewMockForm()

	p.AddControls(context.Background(),
		CheckboxCreator{
			ID: "c1",
		},
		CheckboxCreator{
			ID: "c2",
		},
	)
	c1 := GetCheckbox(p, "c1")
	c1.checked = true
	c2 := GetCheckbox(p, "c2")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	c1.Serialize(enc)
	c2.Serialize(enc)

	dec := gob.NewDecoder(&buf)
	c3 := Checkbox{}
	c4 := Checkbox{}
	c3.Deserialize(dec)
	c4.Deserialize(dec)

	assert.Equal(t, true, c3.Checked())
	assert.Equal(t, false, c4.Checked())
}
