package any

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFieldValues(t *testing.T) {
	a := struct {
		A int
		B string
		C float32
		d float64
	}{
		1, "a", 3.4, 6.7,
	}
	b := FieldValues(a)
	assert.Equal(t, 1, b["A"])
	assert.Equal(t, "a", b["B"])
	assert.Empty(t, b["d"])
}

func TestSetFieldValues(t *testing.T) {
	a := struct {
		A int
		B string
		C float32
		d float64
	}{
		1, "a", 3.4, 6.7,
	}
	b := FieldValues(a)
	c := a
	c.A = 2
	e := SetFieldValues(&c, b)
	assert.NoError(t, e)
	assert.Equal(t, 1, c.A)
	assert.Empty(t, b["d"])
}
