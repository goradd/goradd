package control

import (
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDateTextbox(t *testing.T) {
	p := NewMockForm()

	d := NewDateTextbox(p, "")
	d.SetText("2/19/2019 3:04 pm")
	assert.Equal(t, datetime.February, d.Date().Month())
	assert.Equal(t, 19, d.Date().Day())
	assert.Equal(t, 15, d.Date().Hour())
	assert.Equal(t, 4, d.Date().Minute())

	d.SetText("")
	assert.True(t, d.Date().IsZero())

	d.SetText("asdf")
	assert.True(t, d.Date().IsZero())
	assert.Equal(t, "", d.Text())
}
